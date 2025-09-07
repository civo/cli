package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/civo/civogo"
	"github.com/civo/cli/common"
	"github.com/mitchellh/go-homedir"
)

// SkipAPIInitialization can be set to true to bypass API-related initialization
var (
	SkipAPIInitialization bool
	DefaultAPIURL         = "https://api.civo.com"
)

type APIKey struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	APIURL string `json:"url"`
}

// Config describes the configuration for Civo's CLI
type Config struct {
	APIKeys          []APIKey                  `json:"apikeys"`
	Meta             Metadata                  `json:"meta"`
	RegionToFeatures map[string]civogo.Feature `json:"region_to_features"`
}

// Metadata describes the metadata for Civo's CLI
type Metadata struct {
	Admin              bool      `json:"admin"`
	CurrentAPIKey      string    `json:"current_apikey"`
	DefaultRegion      string    `json:"default_region"`
	LatestReleaseCheck time.Time `json:"latest_release_check"`
	URL                string    `json:"url"`
	LastCmdExecuted    time.Time `json:"last_command_executed"`
}

// Current contains the parsed ~/.civo.json file
var Current Config

// Filename is set to a full filename if the default config
// file is overridden by a command-line switch
var Filename string

// ReadConfig reads in config file and ENV variables if set.
func ReadConfig() {
	filename, found := os.LookupEnv("CIVO_CONFIG")
	if found {
		Filename = filename
	}

	if Filename != "" {
		loadConfig(Filename)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		loadConfig(fmt.Sprintf("%s/%s", home, ".civo.json"))
	}
}

func loadConfig(filename string) {
	var err error
	err = checkConfigFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	configFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Current)
	if err != nil {
		initializeDefaultConfig(filename)
	}
	if Current.APIKeys == nil {
		Current.APIKeys = []APIKey{}
	}

	if token, found := os.LookupEnv("CIVO_TOKEN"); found && token != "" {
		Current.APIKeys = append(Current.APIKeys, APIKey{Name: "tempKey", Value: "token", APIURL: DefaultAPIURL})
		Current.Meta.CurrentAPIKey = "tempKey"
	}

	// Skip API initialization if the flag is set
	if !SkipAPIInitialization {
		if Current.Meta.CurrentAPIKey != "" && Current.RegionToFeatures == nil {
			Current.RegionToFeatures, err = regionsToFeature()
			if err != nil {
				fmt.Printf("Error getting supported regions to feature: %s\n", err)
				os.Exit(1)
			}

			saveUpdatedConfig(filename)
		}

		if time.Since(Current.Meta.LatestReleaseCheck) > (24 * time.Hour) {
			Current.Meta.LatestReleaseCheck = time.Now()

			if Current.Meta.CurrentAPIKey != "" {
				Current.RegionToFeatures, err = regionsToFeature()
				if err != nil {
					fmt.Printf("Error getting supported regions to feature: %s\n", err)
					os.Exit(1)
				}
			}

			saveUpdatedConfig(filename)
			common.CheckVersionUpdate()
		}
	}
}

func saveUpdatedConfig(filename string) {
	// Marshal the Current configuration into JSON
	dataBytes, err := json.Marshal(Current)
	if err != nil {
		fmt.Printf("Error serializing configuration to JSON: %s\n", err)
		os.Exit(1)
	}

	// Write the JSON data to the specified configuration file
	err = os.WriteFile(filename, dataBytes, 0o600)
	if err != nil {
		fmt.Printf("Error writing configuration to file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	// Set file permissions to be read-write for the owner only
	if err := os.Chmod(filename, 0o600); err != nil {
		fmt.Printf("Error setting file permissions for '%s': %s\n", filename, err)
		os.Exit(1)
	}

	fmt.Println("Configuration successfully updated.")
}

// SaveConfig saves the current configuration back out to a JSON file in
// either ~/.civo.json or Filename if one was set
func SaveConfig() {
	var filename string

	if Filename != "" {
		filename = Filename
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filename = fmt.Sprintf("%s/%s", home, ".civo.json")
	}

	dataBytes, err := json.Marshal(Current)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.WriteFile(filename, dataBytes, 0o600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Chmod(filename, 0o600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkConfigFile(filename string) error {
	curr := Config{APIKeys: []APIKey{}}
	curr.Meta = Metadata{
		Admin:           false,
		DefaultRegion:   "NYC1",
		URL:             "https://api.civo.com",
		LastCmdExecuted: time.Now(),
	}

	if Current.Meta.CurrentAPIKey != "" {
		var err error
		curr.RegionToFeatures, err = regionsToFeature()
		if err != nil {
			return err
		}
	}

	fileContend, jsonErr := json.Marshal(curr)
	if jsonErr != nil {
		fmt.Printf("Error parsing the JSON")
		os.Exit(1)
	}

	file, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
		err = os.WriteFile(filename, fileContend, 0o600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		size := file.Size()
		if size == 0 {
			err = os.WriteFile(filename, fileContend, 0o600)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if err := os.Chmod(filename, 0o600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

// regionsToFeature get the region to supported features map
func regionsToFeature() (map[string]civogo.Feature, error) {
	regionsToFeature := map[string]civogo.Feature{}
	client, err := CivoAPIClient()
	if err != nil {
		fmt.Printf("Creating the connection to Civo's API failed with %s", err)
		return regionsToFeature, err
	}

	regions, err := client.ListRegions()
	if err != nil {
		fmt.Printf("Unable to list regions: %s", err)
		return regionsToFeature, err
	}

	for _, region := range regions {
		regionsToFeature[region.Code] = region.Features
	}

	return regionsToFeature, nil
}

// DefaultAPIKey returns the current default API key
func DefaultAPIKey() string {
	if Current.Meta.CurrentAPIKey != "" {
		for _, apiKey := range Current.APIKeys {
			if apiKey.Name == Current.Meta.CurrentAPIKey {
				return apiKey.Value
			}
		}
	}
	return ""
}

// CivoAPIClient returns a civogo client using the current default API key
func CivoAPIClient() (*civogo.Client, error) {
	apiKey := DefaultAPIKey()
	if apiKey == "" {
		fmt.Printf("Error: Creating the connection to Civo's API failed because no API key is supplied. This is required to authenticate requests. Please go to https://dashboard.civo.com/security to obtain your API key, then save it using the command 'civo apikey save YOUR_API_KEY'.\n")
		return nil, fmt.Errorf("no API Key supplied, this is required")
	}
	cliClient, err := civogo.NewClientWithURL(apiKey, Current.Meta.URL, Current.Meta.DefaultRegion)
	if err != nil {
		return nil, err
	}

	var version string
	res, skip := common.VersionCheck(common.GithubClient())
	if !skip {
		version = *res.TagName
	} else {
		version = "0.0.0"
	}

	// Update the user agent to include the version of the CLI
	cliComponent := &civogo.Component{
		Name:    "civo-cli",
		Version: version,
	}
	cliClient.SetUserAgent(cliComponent)

	return cliClient, nil
}

func initializeDefaultConfig(filename string) {
	// Set up a default configuration
	Current = Config{
		APIKeys: []APIKey{}, // Initialize an empty API keys map
		Meta: Metadata{
			Admin:           false,
			DefaultRegion:   "LON1", // Set a default region
			URL:             "https://api.civo.com",
			LastCmdExecuted: time.Now(), // Set the current time for the last executed command
		},
		RegionToFeatures: make(map[string]civogo.Feature), // Initialize an empty map for regions to features
	}

	// Marshal the default configuration to JSON
	dataBytes, err := json.Marshal(Current)
	if err != nil {
		fmt.Printf("Error creating default configuration: %s\n", err)
		os.Exit(1)
	}

	// Write the default configuration to the file
	err = os.WriteFile(filename, dataBytes, 0o600)
	if err != nil {
		fmt.Printf("Error saving default configuration to file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	// Set secure file permissions
	if err := os.Chmod(filename, 0o600); err != nil {
		fmt.Printf("Error setting file permissions for '%s': %s\n", filename, err)
		os.Exit(1)
	}

	fmt.Println("Default configuration initialized and saved.")
}
