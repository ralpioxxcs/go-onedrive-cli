package cmd

import (
	"fmt"
	"go-onedrive-cli/graph"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use: "list",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := viper.Get("access_token")

		res := graph.List(accessToken.(string))
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
