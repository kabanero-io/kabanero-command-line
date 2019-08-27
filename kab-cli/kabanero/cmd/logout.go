/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Disconnect from Kabanero instance",
	Long: `
Disconnect from the instance of Kabanero that you 
have been interacting with.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// url := cliConfig.GetString(KabURLKey) + "/logout"
		url := getRESTEndpoint("logout")
		resp, err := sendHTTPRequest("POST", url, nil)
		if err != nil {
			return errors.New(err.Error())
		}

		defer resp.Body.Close()
		cliConfig.Set("jwt", "")
		cliConfig.WriteConfig()
		println("Logged out of kab instance: " + cliConfig.GetString(KabURLKey))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
