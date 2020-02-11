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
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
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

	if verboseHTTP {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		Info.log("requestDump: " + string(requestDump))
	}

	resp, err = client.Do(req)
	if err != nil {
		return resp, errors.New(err.Error())
	}
	if resp.StatusCode == 503 {
		data := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return resp, errors.New(cliConfig.GetString(KabURLKey) + " is unreachable")
		}

	}
	if resp.StatusCode == 429 {
		return resp, errors.New("Github retry Limited Exceeded, please try again in 2 minutes")
	}
	if resp.StatusCode == 401 {
		return nil, errors.New("Your session may have expired or the credentials entered may be invalid")
	}
	if resp.StatusCode == 539 || resp.StatusCode == 424 || resp.StatusCode == 500 {
		message := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&message)
		if err != nil {
			Debug.log("sync: Decode error for 500/539/424")
			return nil, err
		}
		fmt.Println("HTTP Status " + string(resp.StatusCode) + ": " + message["message"].(string))
		return resp, errors.New("Invalid Response")
	}

	if verboseHTTP {
		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		Info.log("responseDump: " + string(responseDump))
	}
	Debug.log("RESPONSE ", url, " ", resp.StatusCode, " ", http.StatusText(resp.StatusCode))
	return resp, nil
}

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync the stack list",
	Long: `Run the kabanero sync command to synchronize the list of kabanero instance stacks with the curated stacks from github. This will activate/deactivate as well as update versions of the kabanero stacks to reflect the state of the curated stacks.
	See also kabanero deactivate.
	Modifications to the curated stacks may be slow to replicate in git hub and therefore may not be reflected immediately in KABANERO LIST or SYNC display output
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url := getRESTEndpoint("v1/stacks")
		resp, err := sendHTTPRequest("PUT", url, nil)
		if err != nil {
			Debug.log("sync: Error on sendHTTPRequest:")
			return errors.New(err.Error())
		}
		Debug.log("RESPONSE ", url, resp.StatusCode, http.StatusText(resp.StatusCode))
		defer resp.Body.Close()
		//Decode the response into data
		decoder := json.NewDecoder(resp.Body)
		var data StacksResponse
		err = decoder.Decode(&data)
		if err != nil {
			Debug.log("sync: Error on Decode:")
			return err
		}

		Debug.log(data)
		tWriter := new(tabwriter.Writer)
		tWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)
		if len(data.NewStack) == 0 && (len(data.KabStack) == 0) && len(data.ObsoleteStack) == 0 && len(data.CuratedStack) == 0 && len(data.ActivateStack) == 0 {
			syncedOutput := KabStacksHeader + " is already synchronized with the " + GHStacksHeader
			fmt.Println(strings.ToLower(syncedOutput))
		} else {
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", KabStacksHeader, "Version", "Status")
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", strings.Repeat("-", len(KabStacksHeader)), "-------", "------")

			for i := 0; i < len(data.NewStack); i++ {
				for j := 0; j < len(data.NewStack[i].Versions); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.NewStack[i].Name, data.NewStack[i].Versions[j].Version, "added to Kabanero")
				}
			}
			for i := 0; i < len(data.ActivateStack); i++ {
				for j := 0; j < len(data.ActivateStack[i].Versions); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ActivateStack[i].Name, data.ActivateStack[i].Versions[j].Version, "inactive ==> active")
				}
			}
			for i := 0; i < len(data.KabStack); i++ {
				for j := 0; j < len(data.KabStack[i].Status); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.KabStack[i].Name, data.KabStack[i].Status[j].Version, data.KabStack[i].Status[j].Status)
				}
			}
			for i := 0; i < len(data.ObsoleteStack); i++ {
				for j := 0; j < len(data.ObsoleteStack[i].Versions); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ObsoleteStack[i].Name, data.ObsoleteStack[i].Versions[j].Version, "deleted")
				}
			}

			fmt.Fprintln(tWriter)
			tWriter.Flush()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
