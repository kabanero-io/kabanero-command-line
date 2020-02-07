/*
Copyright © 2019 IBM Corporation and others.

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
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
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
	Use:   "login kabanero-url -u Github userid /n PASSWORDPROMPT:GitHub Password|PAT",
	Short: "Will authenticate you to a Kabanero instance",
	Long: `
	Logs in to a Kabanero instance using Github credentials, and stores a temporary access token for subsequent command line calls.
	The temporary authentication token will be stored in your-home-directory/.kabanero/config.yaml.
	Use your Github userid and either password or Personal Access Token (PAT).
	* If you use a GitHub Personal Access Token, make sure it has **read:org** scope. You can select this when creating your PAT in GitHub
	`,
	Example: `
	# login with Github userid and password:
	kabanero login my.kabaneroInstance.io -u myGithubID 
	# login with previously used url Github userid and PAT:
	kabanero login -u myGithubID 
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		Debug.log("login called")
		var err error

		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			fmt.Println("EMPTY USERNAME")
		}
		fmt.Printf("Password:")
		bytePwd, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		password := strings.TrimSpace(string(bytePwd))

		// password, _ := cmd.Flags().GetString("password")
		// if password == "" {
		// 	fmt.Println("EMPTY PASSWORD")
		// }
		var kabLoginURL string

		viper.SetEnvPrefix("KABANERO")

		if len(args) > 0 {
			cliConfig.Set(KabURLKey, parseKabURL(args[0]))
			err = cliConfig.WriteConfig()
			if err != nil {
				return err
			}
		} else {
			if cliConfig.GetString(KabURLKey) == "" {
				return errors.New("No Kabanero instance url specified")
			}
		}
		kabLoginURL = getRESTEndpoint("login")

		requestBody, _ := json.Marshal(map[string]string{"gituser": username, "gitpat": password})

		resp, err := sendHTTPRequest("POST", kabLoginURL, requestBody)
		if err != nil {
			Debug.log("login: Error on sendHTTPRequest:")
			return err
		}

		Debug.log("RESPONSE ", kabLoginURL, resp.StatusCode, http.StatusText(resp.StatusCode))
		if resp.StatusCode == 404 {
			return errors.New("The url: " + cliConfig.GetString(KabURLKey) + " is not a valid kabanero url")
		}

		var data JWTResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return err
		}
		cliConfig.Set("jwt", data.JWT)
		err = cliConfig.WriteConfig()
		if err != nil {
			return err
		}
		if cliConfig.GetString("jwt") == "" {
			Debug.log("Unable to validate user: " + username + " to " + cliConfig.GetString(KabURLKey))
			return errors.New("Unable to validate user: " + username + " to " + cliConfig.GetString(KabURLKey))
		}
		fmt.Println("\nLogged in to Kabanero instance: " + cliConfig.GetString(KabURLKey))
		Debug.log("Logged in to Kabanero instance: " + cliConfig.GetString(KabURLKey))
		defer resp.Body.Close()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("username", "u", "", "github username")
	_ = loginCmd.MarkFlagRequired("username")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
