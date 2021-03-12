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
// file is overridden by a command-line switch
var Filename string

// ReadConfig reads in config file and ENV variables if set.
func ReadConfig() {
	filename, found := os.LookupEnv("CIVOCONFIG")
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
	}

	configFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Current)
	if err != nil {
		fmt.Println(err.Error())
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

	err = ioutil.WriteFile(filename, dataBytes, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := os.Chmod(filename, 0600); err != nil {
		fmt.Println(err)
	}

}

func checkConfigFile(filename string) error {
	file, err := os.Stat(filename)
	fileContend := []byte(fmt.Sprintf("{\"apikeys\":{},\"meta\":{\"admin\":false,\"current_apikey\":\"\",\"default_region\":\"SVG1\",\"latest_release_check\":\"%s\",\"url\":\"https://api.civo.com\"}}", time.Now().Format(time.RFC3339)))
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, fileContend, 0600)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		size := file.Size()
		if size == 0 {
			err = ioutil.WriteFile(filename, fileContend, 0600)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if err := os.Chmod(filename, 0600); err != nil {
		fmt.Println(err)
	}

	return nil
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
	return civogo.NewClientWithURL(DefaultAPIKey(), Current.Meta.URL, Current.Meta.DefaultRegion)
}
