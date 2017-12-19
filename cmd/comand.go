package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/proshik/gotrew/view"
	"github.com/proshik/gotrew/yandex"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

type Config struct {
	YDictToken       string `json:"y_dict_token" long:"y_dict_token" description:"token for Yandex Dictionary API"`
	YTranslatorToken string `json:"y_translator_token" long:"y_translator_token" description:"token for Yandex Translator API" `
}

type Provider struct {
	Title  string
	Active bool
}

const (
	VERSION = "0.1"
	NAME    = "gotrew"
)

var listProviders = make([]Provider, 0)
var config *Config

func Execute() {

	app := cli.NewApp()
	app.Usage = `Application for translate words. Support english and russian languages.`
	app.Version = VERSION

	app.Commands = []cli.Command{
		{
			Name:    "translate",
			Aliases: []string{"t"},
			Usage:   "transate words",
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
							fmt.Printf("Title: %s - active: %v\n", p.Title, p.Active)
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
							if p.Title == c.Args().First() {
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
							Usage: "yandex dictionary token",
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
						appPath, err := buildAppPath()
						if err != nil {
							panic(err)
						}
						configFilename := buildConfigPath(appPath)

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
		appPath, err := buildAppPath()
		if err != nil {
			panic(err)
		}

		err = initApplicationDir(appPath)
		if err != nil {
			return nil
		}

		initLogging(appPath)

		err = initConfig(appPath)
		if err != nil {
			return err
		}

		listProviders = append(listProviders, Provider{Title: "yandex", Active: true})

		return nil
	}

	app.Run(os.Args)
}

func initConfig(appPath string) error {

	configFilename := buildConfigPath(appPath)

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

func initLogging(appPath string) {
	f, err := os.OpenFile(path.Join(appPath, "log.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)
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
