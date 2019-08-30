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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type JWTResponse struct {
	JWT     string
	Message string
}

func parseKabURL(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimSuffix(url, "/")
	return url
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(2),
	Use:   "login userid Github-PAT|Github-password kabanero-url",
	Short: "Will authentic you to a Kabanero instance",
	Long: `
	Login to a Kabanero instance using github credentials, and store a temporary access token for subsequent command line calls.
	The temporary authentication token will be stored in your-home-directory/.kabanero/config.yaml.
	Use your github userid and either password or Personal Access Token (PAT).
	`,
	Example: `
	# login with Github userid and password:
	kabanero login myGithubID myGithubPassword my.kabaneroInstance.io

	# login with Github userid and PAT:
	kabanero login myGithubID myGithubPAT my.kabaneroInstance.io

	# You can also include https.   If you omit it as above, the command will add it automatically
	kabanero login myGithubID myGithubPassword https://my.kabaneroInstance.io
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		Debug.log("login called")

		username := args[0]
		password := args[1]

		var kabLoginURL string

		viper.SetEnvPrefix("KABANERO")

		if len(args) > 2 {
			cliConfig.Set(KabURLKey, parseKabURL(args[2]))
			cliConfig.WriteConfig()
		} else {
			if cliConfig.GetString(KabURLKey) == "" {
				return errors.New("No Kabanero instance url specified")
			}
		}
		kabLoginURL = getRESTEndpoint("login")

		requestBody, _ := json.Marshal(map[string]string{"gituser": username, "gitpat": password})

		resp, err := sendHTTPRequest("POST", kabLoginURL, requestBody)
		if err != nil {
			return err
		}
		Debug.log("RESPONSE ", kabLoginURL, resp.StatusCode, http.StatusText(resp.StatusCode))

		var data JWTResponse
		json.NewDecoder(resp.Body).Decode(&data)
		cliConfig.Set("jwt", data.JWT)
		cliConfig.WriteConfig()
		if cliConfig.GetString("jwt") == "" {
			Debug.log("Unable to validate user: " + username + " to " + cliConfig.GetString(KabURLKey))
			return errors.New("Unable to validate user: " + username + " to " + cliConfig.GetString(KabURLKey))
		}
		fmt.Println("Logged in to Kabanero instance: " + cliConfig.GetString(KabURLKey))
		Debug.log("Logged in to Kabanero instance: " + cliConfig.GetString(KabURLKey))
		defer resp.Body.Close()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
