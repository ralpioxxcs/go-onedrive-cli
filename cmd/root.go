package cmd

import (
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "go-onedrive-cli",
		Short: "go-onedrive-cli is command line tool using onedrive REST API",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "specify config file (default: ${HOME}/.go-onedrive-cli/config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)

		defaultCfgDir := path.Join(homeDir, ".go-onedrive-cli")
		if _, err := os.Stat(defaultCfgDir); os.IsExist(err) {
			if err := os.MkdirAll(defaultCfgDir, os.ModePerm); err != nil {
				log.Fatal(err)
			}
		}

		viper.AddConfigPath(path.Join(homeDir, ".go-onedrive-cli"))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		credential.AddConfigPath(path.Join(homeDir, ".go-onedrive-cli"))
		credential.SetConfigName("credentials")
		credential.SetConfigType("yaml")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Println("using config file:", viper.ConfigFileUsed())
	}

	if err := credential.ReadInConfig(); err == nil {
		log.Fatalf("failed to read credential (err: %v)", err.Error())
	}
}
