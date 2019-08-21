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
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func sendHTTPRequest(method string, url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	var resp *http.Response

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Print("Problem with the new request")
		return resp, errors.New(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", string(cliConfig.GetString("jwt")))
	if cliConfig.GetString("jwt") == "" {
		return resp, errors.New("Login to your kabanero instance")
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Print("Unable to retrieve collections")
		return resp, errors.New(err.Error())
	}
	return resp, nil
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the collections list",
	Long:  `Refresh reconciles the list of collections from master to make them current with the activated collections across all namespace in the kabanero instance`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := cliConfig.GetString(KabURLKey) + "/v1/collections"
		resp, err := sendHTTPRequest("PUT", url)
		if err != nil {
			return errors.New(err.Error())
		}
		defer resp.Body.Close()

		somedata, _ := ioutil.ReadAll(resp.Body)
		Debug.log(string(somedata))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// refreshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// refreshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
