package main

import (
	"encoding/json"
	"github.com/proshik/goswe/view"
	"github.com/proshik/goswe/yandex"
	"io/ioutil"
	"os/user"
	"os"
	"path"
	"bufio"
	"fmt"
)

const (
	VERSION = "0.1.0"
	NAME    = "gotrew"
)

type Config struct {
	YDictToken       string `json:"y_dict_token" long:"y_dict_token" description:"token for Yandex Dictionary API"`
	YTranslatorToken string `json:"y_translator_token" long:"y_translator_token" description:"token for Yandex Translator API" `
}

func main() {
	config, err := readConfigFromFS()
	if err != nil {
		var yDictToken string
		fmt.Printf("Please, enter Yandex dictionary token:\n")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			yDictToken = scanner.Text()
			if yDictToken == "" {
				fmt.Printf("Yandex Dictionary API token not may be empty!\n")
			} else {
				break
			}
		}

		config = &Config{yDictToken, ""}
		if err := createConfig(config); err != nil {
			panic(err)
		}
	}

	yDict := yandex.NewYDictionary(config.YDictToken)
	//yTr := yandex.NewYTranslator(yTranslatorToken)
	ui := view.NewUI(yDict)

	ui.Run()
}

//May throw *PathError
func readConfigFromFS() (*Config, error) {
	appPath, err := buildAppPath()
	if err != nil {
		return nil, err
	}

	filename := buildConfigPath(appPath)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func createConfig(config *Config) error {
	appPath, err := buildAppPath()
	if err != nil {
		return err
	}

	filename := buildConfigPath(appPath)
	//check file and if not exists then create directory
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err = os.MkdirAll(appPath, 0755)
		if err != nil {
			return err
		}
	}
	//create file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	file.Chmod(0755)

	defer file.Close()

	//transform config to []byte
	data, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	file.Sync()

	return nil
}

func buildAppPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(usr.HomeDir, "."+NAME), nil
}

func buildConfigPath(appPath string) string {
	return path.Join(appPath, "config.json")
}
