package main

import (
	"os"
	"github.com/urfave/cli"
	"log"
	"os/user"
	"io/ioutil"
	"encoding/json"
	"github.com/proshik/goswe/yandex"
	view "github.com/proshik/goswe/view"
)

var VERSION = "0.1.1"

type Config struct {
	YDictToken       string `json:"y_dict_token"`
	YTranslatorToken string `json:"y_translator_token"`
}

func main() {

	var yDictToken string
	var yTranslatorToken string

	app := cli.NewApp()
	app.Usage = "translate ru/en words"
	app.Version = VERSION

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "yd-token, ydt",
			Usage:       "token for Yandex Dictionary API",
			Destination: &yDictToken,
		},
		cli.StringFlag{
			Name:        "yt-token, ytt",
			Usage:       "token for Yandex Translator API",
			Destination: &yTranslatorToken,
		},
	}

	app.Before = func(c *cli.Context) error {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Error, not found user home directory, %v", err)
		}
		configFilename := usr.HomeDir + "/.gotrew/config.json"

		config := extractConfigValue(configFilename)
		if config != nil {
			yDictToken = config.YDictToken
			yTranslatorToken = config.YTranslatorToken
			return nil
		}

		if yDictToken == "" || yTranslatorToken == "" {
			log.Fatal("GOTREW requires exactly 2 arguments.\n\nSee 'gotrew --help'.")
		}

		config = &Config{YDictToken: yDictToken, YTranslatorToken: yTranslatorToken}

		data, err := json.Marshal(&config)
		if err != nil {
			log.Fatalf("Error on save token data in file by path=%s", configFilename)
		}

		err = os.MkdirAll(usr.HomeDir+"/.gotrew", 0755)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(configFilename)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		file.Chmod(0755)
		_, err = file.Write(data)
		if err != nil {
			log.Fatal(err)
		}
		file.Sync()

		return nil
	}

	app.Run(os.Args)

	yDict := yandex.NewYDictionary(yDictToken)
	ui := view.NewUI(yDict)

	ui.Run()
}
func extractConfigValue(filename string) *Config {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	var config Config
	json.Unmarshal(configFile, &config)

	return &config
}
