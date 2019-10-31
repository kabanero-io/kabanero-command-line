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
	return "https://" + cliConfig.GetString(KabURLKey) + "/" + appendValue
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
		return resp, err
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
		return resp, errors.New(err.Error())
	}
	if resp.StatusCode == 401 || resp.StatusCode == 503 {
		data := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return nil, err
		}
		expJWTResp := data["message"].(string)
		return nil, errors.New(expJWTResp)
	}
	Debug.log("RESPONSE ", url, resp.StatusCode, http.StatusText(resp.StatusCode))
	return resp, nil
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh the collections list",
	Long:  `Run the kabanero refresh command to refresh the list of collections from the curated collection, making these collections current with the activated collections across all namespaces in the Kabanero instance. This command can also be used to restore deactivated collections. See kabanero deactivate.`,
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
		if err != nil {
			return err
		}

		Debug.log(data)
		tWriter := new(tabwriter.Writer)
		tWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)
		if len(data.NewColl) == 0 && (len(data.KabColl) == 0) && len(data.ObsoleteColl) == 0 && len(data.CuratedColl) == 0 && len(data.VChangeColl) == 0 && len(data.ActivateColl) == 0 {
			fmt.Println("active collection is already synchronized with the curated collection")
		} else {
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", "Collection Name", "Version", "Status")
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", "----", "----", "----")
			for i := 0; i < len(data.NewColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.NewColl[i].Name, data.NewColl[i].Version, "added to kabanero")
			}
			for i := 0; i < len(data.ActivateColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ActivateColl[i].Name, data.ActivateColl[i].Version, "inactive ==> active")
			}
			for i := 0; i < len(data.KabColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.KabColl[i].Name, data.KabColl[i].Version, data.KabColl[i].Status)
			}
			for i := 0; i < len(data.ObsoleteColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ObsoleteColl[i].Name, data.ObsoleteColl[i].Version, "deactivated")
			}
			for i := 0; i < len(data.VChangeColl); i++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.VChangeColl[i].Name, data.VChangeColl[i].Version, "version changed")
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
