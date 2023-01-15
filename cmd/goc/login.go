package main

import (
	"log"

	"github.com/ralpioxxcs/go-onedrive-cli/graph"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to onedrive account",
	Long:  "login to onedrive account and then save access token current config file",
	Run: func(cmd *cobra.Command, args []string) {
		refreshToken := credential.Get("refresh_token")
		var access, refresh string

		auth := graph.Authform{
			RedirectPort: config.GetString("auth.redirect_port"),
			RedirectPath: config.GetString("auth.redirect_path"),
			Scope:        config.GetString("auth.scope"),
			Tenant:       config.GetString("auth.tenant"),
			ClientId:     config.GetString("auth.client_id"),
			ClientSecret: config.GetString("auth.client_secret"),
		}

		if refreshToken != nil {
			access, refresh = graph.Login(refreshToken.(string), auth)
		} else {
			access, refresh = graph.Login("", auth)
		}

		credential.Set("access_token", access)
		credential.Set("refresh_token", refresh)
		err := credential.WriteConfigAs(credentialFilepath)
		if err != nil {
			log.Fatalf("failed to write credential (err: %v)", err.Error())
		}

		log.Println("success to login")

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
