package main

import (
	"github.com/proshik/gotrew/cmd"
	"github.com/proshik/gotrew/utils"
	"log"
	"os"
	"path"
	"io"
)

//NAME is a title application
const NAME = "gotrew"

//VERSION is a version of application
var VERSION = "0.1.2"

func main() {
	//logging configure
	appPath, err := utils.BuildAppPath(NAME)
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(path.Join(appPath, "log.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	//run a command line interface
	cmd.Execute(NAME, VERSION)
}
