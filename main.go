package main

import (
	"github.com/proshik/gotrew/cmd"
	"github.com/proshik/gotrew/utils"
	"io/ioutil"
	"log"
	"os"
	"path"
)

//NAME is a title application
const NAME = "gotrew"

//VERSION is a version of application
var VERSION string

func main() {
	//read VERSION file
	versionFile, err := ioutil.ReadFile("VERSION")
	if err != nil {
		panic(err)
	}
	VERSION = string(versionFile)

	//logging configure
	appPath, err := utils.BuildAppPath(NAME)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(path.Join(appPath, "log.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(f)

	//run a command line interface
	cmd.Execute(NAME, VERSION)

}
