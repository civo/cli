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

// Config describes the configuration for Civo's CLI
type Config struct {
	APIKeys          map[string]string         `json:"apikeys"`
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
	filename := GetConfigFilename()
	if filename != "" {
		LoadConfig(filename)
		ProcessConfig(filename)
	} else {
		fmt.Println("No configuration file found")
		os.Exit(1)
	}
}

// Function to retrieve the config filename from environment variable or default
func GetConfigFilename() string {
	if filename, found := os.LookupEnv("CIVO_CONFIG"); found {
		return filename
	}

	homeDir := getHomeDir()
	if homeDir != "" {
		return fmt.Sprintf("%s/%s", homeDir, ".civo.json")
	}

	return ""
}

// Function to get the home directory of the current user
func getHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Error retrieving home directory:", err)
		os.Exit(1)
	}

	return home
}

func LoadConfig(filename string) {
	var err error
	err = CheckConfigFile(filename)
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
		Current.Meta.Admin = false
		Current.Meta.DefaultRegion = "LON1"
		Current.Meta.URL = "https://api.civo.com"
		Current.Meta.LastCmdExecuted = time.Now()

		fileContend, jsonErr := json.Marshal(Current)
		if jsonErr != nil {
			fmt.Printf("Error parsing the JSON")
			os.Exit(1)
		}
		err = os.WriteFile(filename, fileContend, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func ProcessConfig(filename string) {
	if Current.APIKeys == nil {
		Current.APIKeys = map[string]string{}
	}

	if token, found := os.LookupEnv("CIVO_TOKEN"); found && token != "" {
		Current.APIKeys["tempKey"] = token
		Current.Meta.CurrentAPIKey = "tempKey"
	}

	if Current.Meta.CurrentAPIKey != "" && Current.RegionToFeatures == nil {
		regionstoFeature, err := regionsToFeature()
		if err != nil {
			fmt.Printf("Error getting supported regions to feature %s \n", err)
			os.Exit(1)
		}

		Current.RegionToFeatures = regionstoFeature

		dataBytes, err := json.Marshal(Current)
		if err != nil {
			fmt.Printf("Error parsing JSON %s \n", err)
			os.Exit(1)
		}

		err = os.WriteFile(filename, dataBytes, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if time.Since(Current.Meta.LatestReleaseCheck) > (24 * time.Hour) {
		Current.Meta.LatestReleaseCheck = time.Now()

		if Current.Meta.CurrentAPIKey != "" {
			regionFeatures, err := regionsToFeature()
			if err != nil {
				fmt.Printf("Error getting supported regions to feature %s \n", err)
				os.Exit(1)
			}

			Current.RegionToFeatures = regionFeatures
		}

		dataBytes, err := json.Marshal(Current)
		if err != nil {
			fmt.Printf("Error parsing JSON %s \n", err)
			os.Exit(1)
		}

		err = os.WriteFile(filename, dataBytes, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		common.CheckVersionUpdate()
	}
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

	err = os.WriteFile(filename, dataBytes, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Chmod(filename, 0600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func CheckConfigFile(filename string) error {
	curr := Config{APIKeys: map[string]string{}}
	curr.Meta = Metadata{
		Admin:           false,
		DefaultRegion:   "NYC1",
		URL:             "https://api.civo.com",
		LastCmdExecuted: time.Now(),
	}

	currApiKey := Current.Meta.CurrentAPIKey
	if currApiKey != "" {
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
		err = os.WriteFile(filename, fileContend, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		size := file.Size()
		if size == 0 {
			err = os.WriteFile(filename, fileContend, 0600)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if err := os.Chmod(filename, 0600); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Current = curr
	if currApiKey != "" {
		Current.Meta.CurrentAPIKey = currApiKey
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
		return Current.APIKeys[Current.Meta.CurrentAPIKey]
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
