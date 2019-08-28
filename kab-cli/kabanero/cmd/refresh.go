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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func getRESTEndpoint(appendValue string) string {
	return cliConfig.GetString(KabURLKey) + "/" + KabURLContext + "/" + appendValue
}

func sendHTTPRequest(method string, url string, jsonBody []byte) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Timeout:   time.Second * 30,
		Transport: tr,
	}

	var resp *http.Response
	var requestBody *bytes.Buffer
	var req *http.Request
	var err error
	if jsonBody != nil {
		requestBody = bytes.NewBuffer(jsonBody)
		req, err = http.NewRequest(method, url, requestBody)

	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		fmt.Print("Problem with the new request")
		return resp, errors.New(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	if !strings.Contains(url, "login") {
		req.Header.Set("Authorization", string(cliConfig.GetString("jwt")))
		if cliConfig.GetString("jwt") == "" {
			return resp, errors.New("Login to your kabanero instance")
		}
	}

	resp, err = client.Do(req)
	if err != nil {
		fmt.Print("Unable to retrieve collections")
		return resp, errors.New(err.Error())
	}
	Debug.log("RESPONSE ", url, resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, nil
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the collections list",
	Long:  `Refresh reconciles the list of collections from master to make them current with the activated collections across all namespace in the kabanero instance`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := getRESTEndpoint("v1/collections")
		resp, err := sendHTTPRequest("PUT", url, nil)
		if err != nil {
			return errors.New(err.Error())
		}
		Debug.log("RESPONSE ", url, resp.StatusCode, http.StatusText(resp.StatusCode))
		defer resp.Body.Close()
		//Decode the response into data
		decoder := json.NewDecoder(resp.Body)
		var data CollectionsResponse
		err = decoder.Decode(&data)
		//

		Debug.log(data)
		tWriter := new(tabwriter.Writer)
		tWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)
		if len(data.NewColl) == 0 && (len(data.ActiveColl) == 0) && len(data.ObsoleteColl) == 0 && len(data.MasterColl) == 0 && len(data.VChangeColl) == 0 {
			fmt.Println("active collections synchronized with master")
		} else {
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", "Name", "Version", "Collection")
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", "----", "----", "----")
			for i := 0; i < len(data.NewColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.NewColl[i].Name, data.NewColl[i].Version, "new collection")
			}
			for i := 0; i < len(data.ActiveColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ActiveColl[i].Name, data.ActiveColl[i].Version, "active collections")
			}
			for i := 0; i < len(data.ObsoleteColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ObsoleteColl[i].Name, data.ObsoleteColl[i].Version, "obsolete collections")
			}
			for i := 0; i < len(data.MasterColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.MasterColl[i].Name, data.MasterColl[i].Version, "master collection")
			}
			for i := 0; i < len(data.VChangeColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.VChangeColl[i].Name, data.VChangeColl[i].Version, "changed collection")
			}
			fmt.Fprintln(tWriter)
			tWriter.Flush()
		}
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
