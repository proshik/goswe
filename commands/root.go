package commands

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"errors"
	"github.com/spf13/viper"
)

type RootConfig struct {
	App             App
	Provider        Provider
	DictionaryToken string `yaml:"dictionary_token"`
	TranslatorToken string `yaml:"translator_token"`
}

type App struct {
	Name string
}

type Provider struct {
	Active string
}

var RootCfgFile string
var Config RootConfig

func init() {
	RootCmd.PersistentFlags().StringVar(&RootCfgFile, "config", "", "config file (default is C:\\Users\\Prokhor_Krylov\\.gotrew)")
}

var RootCmd = &cobra.Command{
	Use:   "gotrew",
	Short: "Application for translate words",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gotrew application for translate words")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg RootConfig

		viper.AddConfigPath(RootCfgFile) // call multiple times to add many search paths

		err := viper.Unmarshal(&cfg)
		if err != nil {
			panic(err)
		}

		return errors.New("Error on init config")
	},
}

func Execute() {

	RootCmd.SilenceUsage = true

	RootCmd.AddCommand(translate)
	RootCmd.AddCommand()

	//RootCmd.SilenceUsage = false

	RootCmd.Usage()

	if c, err := RootCmd.ExecuteC(); err != nil {
		//if isUserError(err) {
		c.Println("")
		c.Println(c.UsageString())
		//}

		os.Exit(-1)
	}
}
