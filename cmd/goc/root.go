package main

import (
	"log"

	"github.com/spf13/cobra"
)

var (
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
	cobra.OnInitialize(initConfigs)
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "specify config file (default: ${HOME}/.goc/config.yaml)")
}
