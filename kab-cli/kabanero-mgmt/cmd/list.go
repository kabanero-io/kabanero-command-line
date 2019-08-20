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
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type collStruct struct {
	name      string
	version   string
	status    string
	exception string
}

type collectionsResponse struct {
	newColl      []map[string][]string `json:"new collections"`
	kabColl      string                `json:"kabanero collection"`
	obsoleteColl string                `json:"obsolete collections"`
	masterColl   string                `json:"master collection"`
	vChangeColl  string                `json:"version change collection"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [status]",
	Short: "List all the collections in the apphub, and optionally their status",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "http://10.211.54.244:31000/KabCollections-1.0-SNAPSHOT/v1/collections"
		fmt.Println("list called")
		client := &http.Client{
			Timeout: time.Second * 30,
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Print("Problem with the new request")
			return errors.New(err.Error())
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", string(cliConfig.GetString("jwt")))
		if cliConfig.GetString("jwt") == "" {
			return errors.New("Login to your kabanero instance")
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print("Unable to retrieve collections")
			return errors.New(err.Error())
		}
		defer resp.Body.Close()

		somedata, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(somedata))

		var data collectionsResponse
		json.NewDecoder(resp.Body).Decode(&data)
		fmt.Println("**********************************")
		fmt.Println(data.kabColl)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
