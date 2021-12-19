package cmd

import (
	"fmt"
	"go-onedrive-cli/graph"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to onedrive account",
	Long:  "login to onedrive account, save access token",
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := graph.Login()

		viper.Set("access_token", accessToken)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var (
	tenant string
)

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&tenant, "tenant", "common", "specify tenant of account")
}
