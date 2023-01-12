package main

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrNotExistConfig    = errors.New("config file is not exists")
	ErrNotExistDirectory = errors.New("base directory is not exists")
)

var (
	config             *viper.Viper
	credential         *viper.Viper
	configFilepath     string
	credentialFilepath string
)

func initConfigs() {
	config = viper.New()
	credential = viper.New()

	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	defaultConfigDir := path.Join(homeDir, ".go-onedrive-cli")
	if _, err := os.Stat(defaultConfigDir); os.IsExist(err) {
		if err := os.MkdirAll(defaultConfigDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	configFilepath = path.Join(homeDir, ".go-onedrive-cli", "config.yaml")
	credentialFilepath = path.Join(homeDir, ".go-onedrive-cli", "credentials.yaml")

	log.Println("config: ", configFilepath)
	log.Println("credential: ", credentialFilepath)

	config.AddConfigPath(path.Join(homeDir, ".go-onedrive-cli"))
	config.SetConfigName("config")
	config.SetConfigType("yaml")

	credential.AddConfigPath(path.Join(homeDir, ".go-onedrive-cli"))
	credential.SetConfigName("credentials")
	credential.SetConfigType("yaml")

	viper.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config (err: %v)", err.Error())
	}

	/*
		if err := credential.ReadInConfig(); err != nil {
			log.Fatalf("failed to read credential (err: %v)", err.Error())
		}
	*/

}
