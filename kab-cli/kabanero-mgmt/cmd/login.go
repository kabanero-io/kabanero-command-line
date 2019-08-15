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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type JWTResponse struct {
	JWT     string
	Message string
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(2),
	Use:   "login userid password kabanero-url",
	Short: "Will authentic you to the Kabanero instance",
	Long: `
	The userid and password passed will be used
	to authenticate the user with kabanero instance.
	
	By authenticating with the Kabanero instance, 
	you will be able to manage the instance of kabanero.`,
	Example: `
		kabanero-management champ champpassword https://kabanero1.io
		`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("login called")

		username := args[0]
		password := args[1]

		var kabURL string

		if len(args) > 2 {
			kabURL = args[2]
		} else {
			return errors.New("No Kabanero instance url specified")
		}

		client := &http.Client{
			Timeout: time.Second * 30,
		}

		requestBody, _ := json.Marshal(map[string]string{"gituser": username, "gitpat": password})

		req, err := http.NewRequest("POST", kabURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return errors.New(err.Error())
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			return errors.New("Login failed to endpoint: " + kabURL + " \n")
		}

		// data, _ := ioutil.ReadAll(resp.Body)
		var data JWTResponse
		json.Unmarshal(resp.Body, &data)
		defer resp.Body.Close()

		// fmt.Println(string(data))
		fmt.Println(data.JWT)
		fmt.Printf("USERNAME/PWD/KAB" + username + "-- " + password + "____" + kabURL)
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
