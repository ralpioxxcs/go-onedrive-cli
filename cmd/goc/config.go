package main

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/ralpioxxcs/go-onedrive-cli/internal/util"
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

	credentialFilepath = path.Join(homeDir, ".go-onedrive-cli", "credentials.yaml")

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

	//
	if util.IsFileExists(credentialFilepath) {
		if err := credential.ReadInConfig(); err != nil {
			log.Fatalf("failed to read credential (err: %v)", err.Error())
		}
	}

}

func CheckCredentials() bool {
	return util.IsFileExists(credentialFilepath)
}
