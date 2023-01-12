package cmd

import (
	"go-onedrive-cli/graph"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	credential *viper.Viper
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to onedrive account",
	Long:  "login to onedrive account and then save access token current config file",
	Run: func(cmd *cobra.Command, args []string) {
		refreshToken := credential.Get("refresh_token")
		var access, refresh string
		if refreshToken != nil {
			access, refresh = graph.Login(refreshToken.(string))
		} else {
			access, refresh = graph.Login("")
		}

		credential.Set("access_token", access)
		credential.Set("refresh_token", refresh)
		err := credential.WriteConfig()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	credential = viper.New()
	homeDir, _ := os.UserHomeDir()

	credential.AddConfigPath(path.Join(homeDir, ".go-onedrive-cli"))
	credential.SetConfigName("credentials")
	credential.SetConfigType("yaml")
}
