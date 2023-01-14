package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "go-onedrive-cli",
		Short: "go-onedrive-cli is command line tool using onedrive REST API",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("test!", cmd)

		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfigs)
}
