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
	if Current.APIKeys == nil {
		Current.APIKeys = map[string]string{}
	}

	if token, found := os.LookupEnv("CIVO_TOKEN"); found && token != "" {
		Current.APIKeys["tempKey"] = token
		Current.Meta.CurrentAPIKey = "tempKey"
	}

	if Current.Meta.CurrentAPIKey != "" && Current.RegionToFeatures == nil {
		Current.RegionToFeatures, err = regionsToFeature()
		if err != nil {
			fmt.Printf("Error getting supported regions to feature %s \n", err)
			os.Exit(1)
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
	}

	if time.Since(Current.Meta.LatestReleaseCheck) > (24 * time.Hour) {
		Current.Meta.LatestReleaseCheck = time.Now()

		if Current.Meta.CurrentAPIKey != "" {
			Current.RegionToFeatures, err = regionsToFeature()
			if err != nil {
				fmt.Printf("Error getting supported regions to feature %s \n", err)
				os.Exit(1)
			}
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

func checkConfigFile(filename string) error {
	curr := Config{APIKeys: map[string]string{}}
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
	cliClient, err := civogo.NewClientWithURL(DefaultAPIKey(), Current.Meta.URL, Current.Meta.DefaultRegion)
	if err != nil {
		return nil, err
	}

	var version string
	res, skip := common.VersionCheck()
	if !skip {
		version = res.Current
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
