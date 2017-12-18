package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/proshik/gotrew/view"
	"github.com/proshik/gotrew/yandex"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"github.com/urfave/cli"
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
	////Logging
	//f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()
	//
	//log.SetOutput(f)

	var yDictToken string
	var yTranslatorToken string

	app := cli.NewApp()
	app.Usage = "translate ru/en words"
	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "y_dict_token, ydt",
			Usage:       "token for Yandex Dictionary API",
			Destination: &yDictToken,
		},
		cli.StringFlag{
			Name:        "y_translator_token, ytt",
			Usage:       "token for Yandex Translator API",
			Destination: &yTranslatorToken,
		},
	}
	app.Run(os.Args)

	if yDictToken != "" || yTranslatorToken != "" {
		if err := createConfig(&Config{yDictToken, yTranslatorToken}); err != nil {
			panic(err)
		}
	}

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
