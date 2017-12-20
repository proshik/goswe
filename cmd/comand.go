package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/proshik/gotrew/utils"
	"github.com/proshik/gotrew/view"
	"github.com/proshik/gotrew/yandex"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	YDictToken       string `json:"y_dict_token"`
	YTranslatorToken string `json:"y_translator_token"`
}

type Provider struct {
	Title  string
	Active bool
	URL    string
}

var listProviders = make([]Provider, 0)
var config *Config

func Execute(appName string, appVersion string) {
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Usage = "Application for translate words. Support english and russian languages."
	app.HideVersion = false

	app.Commands = []cli.Command{
		{
			Name:    "translate",
			Aliases: []string{"t"},
			Usage:   "translate words mode",
			Action: func(c *cli.Context) error {
				yDict := yandex.NewYDictionary(config.YDictToken)
				ui := view.NewUI(yDict)

				ui.Run()
				return nil
			},
		},
		{
			Name:    "provider",
			Aliases: []string{"p"},
			Usage:   "show and select provider for translate",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list available providers",
					Action: func(c *cli.Context) error {
						for _, p := range listProviders {
							fmt.Printf("- %s, active: %v, url: %s\n", strings.ToUpper(p.Title), p.Active, p.URL)
						}
						return nil
					},
				},
				{
					Name:  "select",
					Usage: "select a provider",
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							fmt.Println("Need printed provider title")
							return nil
						}

						for _, p := range listProviders {
							if strings.ToUpper(p.Title) == strings.ToUpper(c.Args().First()) {
								fmt.Printf("Select provider: %s\n", c.Args().First())
								return nil
							}
						}

						fmt.Printf("Unrecognized provider title: %s\n", c.Args().First())
						return nil
					},
				},
				{
					Name:  "config",
					Usage: "set config for providers",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "token, t",
							Usage: "provider dictionary token",
						},
					},
					Action: func(c *cli.Context) error {
						if len(c.Args()) == 0 {
							fmt.Println("Need printed provider title")
							return nil
						}

						token := c.String("token")
						if token == "" {
							fmt.Println("Need printed flag [token]")
							return nil
						}

						var found = false
						for _, p := range listProviders {
							if p.Title == c.Args().First() {
								found = true
							}
						}
						if found == false {
							fmt.Printf("Unrecognized provider title: %s\n", c.Args().First())
							return nil
						}
						appPath, err := utils.BuildAppPath(appName)
						if err != nil {
							panic(err)
						}
						configFilename := utils.BuildConfigPath(appPath)

						err = saveConfig(configFilename, &Config{token, ""})
						if err != nil {
							return err
						}

						return nil
					},
				},
			},
		},
	}

	app.Before = func(context *cli.Context) error {
		appPath, err := utils.BuildAppPath(appName)
		if err != nil {
			panic(err)
		}

		err = initApplicationDir(appPath)
		if err != nil {
			return nil
		}

		err = initConfig(appPath)
		if err != nil {
			return err
		}

		listProviders = append(listProviders, Provider{Title: "yandex", Active: true, URL: "https://tech.yandex.ru/dictionary"})

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

	log.Println("Application was running!")
}

func initConfig(appPath string) error {
	configFilename := utils.BuildConfigPath(appPath)

	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		if err := saveConfig(configFilename, &Config{"", ""}); err != nil {
			return err
		}
	}

	c, err := readConfig(configFilename)
	if err != nil {
		return err
	}
	config = c

	return nil
}

func initApplicationDir(appPath string) error {
	//check file and if not exists then create directory
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		err = os.MkdirAll(appPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveConfig(filename string, config *Config) error {
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

func readConfig(filename string) (*Config, error) {
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
