package main

import (
	"github.com/proshik/gotrew/cmd"
	"github.com/proshik/gotrew/utils"
	"log"
	"os"
	"path"
)

//NAME is a title application
const NAME = "gotrew"

//VERSION is a version of application
var VERSION = "0.1.5"

func main() {
	//logging configure
	appPath, err := utils.BuildAppPath(NAME)
	if err != nil {
		panic(err)
	}

	logFilePath := path.Join(appPath, "log.log")
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		f, err := os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	//run a command line interface
	cmd.Execute(NAME, VERSION)
}
