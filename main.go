package main

import (
	//"bufio"
	//"encoding/json"
	//"fmt"
	//"github.com/proshik/gotrew/view"
	//"github.com/proshik/gotrew/yandex"
	//"io/ioutil"
	//"os"
	//"os/user"
	//"path"
	//"github.com/spf13/cobra/cobra/cmd"
	//"github.com/spf13/cobra"
	//"strings"
	"github.com/proshik/gotrew/commands"
)

const (
	VERSION = "0.1.0"
	NAME    = "gotrew"
)

func main() {
	////Logging
	//f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()
	//
	//log.SetOutput(f)

	commands.Execute()

	//OLD

	//var yDictToken string
	//var yTranslatorToken string
	//
	//if yDictToken != "" || yTranslatorToken != "" {
	//	if err := createConfig(&Config{yDictToken, yTranslatorToken}); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//config, err := readConfigFromFS()
	//if err != nil {
	//	var yDictToken string
	//	fmt.Printf("Please, enter Yandex dictionary token:\n")
	//	scanner := bufio.NewScanner(os.Stdin)
	//	for scanner.Scan() {
	//		yDictToken = scanner.Text()
	//		if yDictToken == "" {
	//			fmt.Printf("Yandex Dictionary API token not may be empty!\n")
	//		} else {
	//			break
	//		}
	//	}
	//
	//	config = &Config{yDictToken, ""}
	//	if err := createConfig(config); err != nil {
	//		panic(err)
	//	}
	//}
	//
	//yDict := yandex.NewYDictionary(config.YDictToken)
	//_ = view.NewUI(yDict)

	//ui.Run()
}

//May throw *PathError
//func readConfigFromFS() (*Config, error) {
//	appPath, err := buildAppPath()
//	if err != nil {
//		return nil, err
//	}
//
//	filename := buildConfigPath(appPath)
//
//	if _, err := os.Stat(filename); os.IsNotExist(err) {
//		return nil, err
//	}
//
//	configFile, err := ioutil.ReadFile(filename)
//	if err != nil {
//		return nil, err
//	}
//
//	var config Config
//	err = json.Unmarshal(configFile, &config)
//	if err != nil {
//		return nil, err
//	}
//	return &config, nil
//}
//
//func createConfig(config *Config) error {
//	appPath, err := buildAppPath()
//	if err != nil {
//		return err
//	}
//
//	filename := buildConfigPath(appPath)
//	//check file and if not exists then create directory
//	if _, err := os.Stat(filename); os.IsNotExist(err) {
//		err = os.MkdirAll(appPath, 0755)
//		if err != nil {
//			return err
//		}
//	}
//	//create file
//	file, err := os.Create(filename)
//	if err != nil {
//		return err
//	}
//	file.Chmod(0755)
//
//	defer file.Close()
//
//	//transform config to []byte
//	data, err := json.Marshal(&config)
//	if err != nil {
//		return err
//	}
//
//	_, err = file.Write(data)
//	if err != nil {
//		return err
//	}
//
//	file.Sync()
//
//	return nil
//}
//
//func buildAppPath() (string, error) {
//	usr, err := user.Current()
//	if err != nil {
//		return "", err
//	}
//
//	return path.Join(usr.HomeDir, "."+NAME), nil
//}
//
//func buildConfigPath(appPath string) string {
//	return path.Join(appPath, "config.json")
//}
