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
var SkipAPIInitialization bool

// Config describes the configuration for Civo's CLI
type Config struct {
	APIKeys          map[string]string         `json:"apikeys"`
	Meta             Metadata                  `json:"meta"`
	RegionToFeatures map[string]civogo.Feature `json:"region_to_features"`
}

// Metadata describes the metadata for Civo's CLI
type Metadata struct {
	Admin               bool      `json:"admin"`
	CurrentAPIKey       string    `json:"current_apikey"`
	DefaultRegion       string    `json:"default_region"`
	LatestReleaseCheck  time.Time `json:"latest_release_check"`
	URL                 string    `json:"url"`
	LastCmdExecuted     time.Time `json:"last_command_executed"`
	DisableVersionCheck bool      `json:"disable_version_check"` 
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
		Current.APIKeys = map[string]string{}
	}

	if token, found := os.LookupEnv("CIVO_TOKEN"); found && token != "" {
		Current.APIKeys["tempKey"] = token
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

		
		if !Current.Meta.DisableVersionCheck && time.Since(Current.Meta.LatestReleaseCheck) > (24*time.Hour) {
			Current.Meta.LatestReleaseCheck = time.Now()
			saveUpdatedConfig(filename)
			common.CheckVersionUpdate()
		}
	}
}

func saveUpdatedConfig(filename string) {
	dataBytes, err := json.Marshal(Current)
	if err != nil {
		fmt.Printf("Error serializing configuration to JSON: %s\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(filename, dataBytes, 0600)
	if err != nil {
		fmt.Printf("Error writing configuration to file '%s': %s\n", filename, err)
		os.Exit(1)
	}

	if err := os.Chmod(filename, 0600); err != nil {
		fmt.Printf("Error setting file permissions for '%s': %s\n", filename, err)
		os.Exit(1)
	}
}

// SaveConfig saves the current configuration back to ~/.civo.json
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

	saveUpdatedConfig(filename)
}

func checkConfigFile(filename string) error {
	curr := Config{APIKeys: map[string]string{}}
	curr.Meta = Metadata{
		Admin:               false,
		DefaultRegion:       "NYC1",
		URL:                 "https://api.civo.com",
		LastCmdExecuted:     time.Now(),
		DisableVersionCheck: false,  
	}

	if Current.Meta.CurrentAPIKey != "" {
		var err error
		curr.RegionToFeatures, err = regionsToFeature()
		if err != nil {
			return err
		}
	}

	fileContent, jsonErr := json.Marshal(curr)
	if jsonErr != nil {
		fmt.Printf("Error parsing the JSON")
		os.Exit(1)
	}

	file, err := os.Stat(filename)
	if os.IsNotExist(err) {
		err = os.WriteFile(filename, fileContent, 0600)
		if err != nil {
			return err
		}
	} else if file.Size() == 0 {
		err = os.WriteFile(filename, fileContent, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err := os.Chmod(filename, 0600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}

// regionsToFeature fetches supported features for regions
func regionsToFeature() (map[string]civogo.Feature, error) {
	client, err := CivoAPIClient()
	if err != nil {
		fmt.Printf("Creating the connection to Civo's API failed with %s", err)
		return nil, err
	}

	regions, err := client.ListRegions()
	if err != nil {
		fmt.Printf("Unable to list regions: %s", err)
		return nil, err
	}

	regionFeatures := make(map[string]civogo.Feature)
	for _, region := range regions {
		regionFeatures[region.Code] = region.Features
	}

	return regionFeatures, nil
}

// DefaultAPIKey returns the current API key
func DefaultAPIKey() string {
	if Current.Meta.CurrentAPIKey != "" {
		return Current.APIKeys[Current.Meta.CurrentAPIKey]
	}
	return ""
}

// CivoAPIClient creates a client with the current API key
func CivoAPIClient() (*civogo.Client, error) {
	apiKey := DefaultAPIKey()
	if apiKey == "" {
		fmt.Println("Error: No API key supplied. Please save it using 'civo apikey save YOUR_API_KEY'.")
		return nil, fmt.Errorf("no API Key supplied")
	}

	cliClient, err := civogo.NewClientWithURL(apiKey, Current.Meta.URL, Current.Meta.DefaultRegion)
	if err != nil {
		return nil, err
	}

	var version string
	if !Current.Meta.DisableVersionCheck { 
		res, skip := common.VersionCheck(common.GithubClient())
		if !skip {
			version = *res.TagName
		} else {
			version = "0.0.0"
		}
	} else {
		version = "0.0.0"
	}

	cliClient.SetUserAgent(&civogo.Component{
		Name:    "civo-cli",
		Version: version,
	})

	return cliClient, nil
}

func initializeDefaultConfig(filename string) {
	Current = Config{
		APIKeys: make(map[string]string),
		Meta: Metadata{
			Admin:               false,
			DefaultRegion:       "LON1",
			URL:                 "https://api.civo.com",
			LastCmdExecuted:     time.Now(),
			DisableVersionCheck: false, 
		},
		RegionToFeatures: make(map[string]civogo.Feature),
	}

	saveUpdatedConfig(filename)
}
