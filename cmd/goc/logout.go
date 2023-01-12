package main

import (
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use: "logout",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
