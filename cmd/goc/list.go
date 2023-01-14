package main

import (
	"fmt"

	"github.com/ralpioxxcs/go-onedrive-cli/graph"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list items",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := credential.GetString("access_token")

		fmt.Println(accessToken)

		res := graph.List(accessToken)
		for _, v := range res {
			fmt.Println(v)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
