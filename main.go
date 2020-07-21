package main

import (
	"fmt"
	"os"
	"path"
)

func GetConfigDir() string {
	configPath := os.Getenv("STM_CONFIG_PATH")
	if configPath != "" {
		return configPath
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = path.Join(os.Getenv("HOME"), ".config")
	}

	return path.Join(configHome, "stm")
}

func SetupConfigPath() {
	configPath := GetConfigDir()
	os.Setenv("STM_CONFIG_PATH", configPath)
	os.MkdirAll(configPath, 0755)
}

func main() {
	fmt.Println(GetConfigDir())
}
