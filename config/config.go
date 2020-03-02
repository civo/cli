package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/civo/civogo"
	"github.com/mitchellh/go-homedir"
)

// Config describes the configuration for Civo's CLI
type Config struct {
	APIKeys map[string]string `json:"apikeys"`
	Meta    struct {
		Admin              bool      `json:"admin"`
		CurrentAPIKey      string    `json:"current_apikey"`
		DefaultRegion      string    `json:"default_region"`
		LatestReleaseCheck time.Time `json:"latest_release_check"`
		URL                string    `json:"url"`
	} `json:"meta"`
}

// Current contains the parsed ~/.civo.json file
var Current Config

// Filename is set to a full filename if the default config
// file is overriden by a command-line switch
var Filename string

// ReadConfig reads in config file and ENV variables if set.
func ReadConfig() {
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
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&Current)
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

	configJSON, _ := json.Marshal(Current)
	err := ioutil.WriteFile(filename, configJSON, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	return civogo.NewClient(DefaultAPIKey())
}
