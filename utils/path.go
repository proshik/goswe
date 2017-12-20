package utils

import (
	"os/user"
	"path"
)

func BuildAppPath(appName string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(usr.HomeDir, "."+appName), nil
}

func BuildConfigPath(appPath string) string {
	return path.Join(appPath, "config.json")
}
