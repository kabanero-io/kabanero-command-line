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
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/kabanero-io/kabanero-command-line/pkg/security"

	"github.com/spf13/cobra"
)

func messageAndExit(msg string) {
	Debug.log(msg)
	fmt.Println(msg)
	os.Exit(3)
}

func messageandDebugExit(msg string, dbgmsg string) {
	Debug.log(dbgmsg)
	fmt.Println(msg)
	os.Exit(3)
}

func getRESTEndpoint(appendValue string) string {
	return "https://" + cliConfig.GetString(KabURLKey) + "/" + appendValue
}

func sendHTTPRequest(method string, url string, jsonBody []byte) (*http.Response, error) {

	var resp *http.Response
	var requestBody *bytes.Buffer
	var req *http.Request
	var err error

	// commented out codes have their own handling and are here just for error code tracking
	serviceErrorCodes := map[int]string{
		400: "Stack Version not found/version not found/ jwt expired",
		// 401: "Session expired or invalid certs",
		// 404: "Unable to reach services endpoint",
		424: "Kab CR config errors",
		429: "GitHub retry limit exceeded",
		500: "Internal Server Error",
		503: "Operator pod is not fully up",
		539: "CLI has not been configured",
	}

	rootCAPool, _ := x509.SystemCertPool()
	if rootCAPool == nil {
		rootCAPool = x509.NewCertPool()
	}
	if !cliConfig.GetBool("insecureTLS") {
		cert, err := ioutil.ReadFile(cliConfig.GetString(CertKey))
		if err != nil {
			messageAndExit(fmt.Sprintf("Problem with the certificate for %s, provided at %s", cliConfig.GetString(KabURLKey), cliConfig.GetString(CertKey)))
		}
		rootCAPool.AppendCertsFromPEM(cert)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            rootCAPool,
			InsecureSkipVerify: cliConfig.GetBool("insecureTLS")},
	}
	client := &http.Client{
		Timeout:   time.Second * 30,
		Transport: tr,
	}

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
		if cliConfig.GetString("jwt") == "" || cliConfig.GetString("key") == "" {
			messageAndExit("Login to your kabanero instance")
		}

		jwt := security.DecryptString(cliConfig.GetString("jwt"), cliConfig.GetString("key"))
		req.Header.Set("Authorization", "Bearer "+jwt)

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
		msg := "Could not connect to url: " + cliConfig.GetString(KabURLKey) + "\nError: " + err.Error()
		messageAndExit(msg)
	}
	if verboseHTTP {
		responseDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Println(err)
		}
		Info.log("responseDump: " + string(responseDump))
	}
	if resp.StatusCode == 401 {
		messageAndExit("Session expired or your token is invalid. Please try logging in again")
	}

	if _, found := serviceErrorCodes[resp.StatusCode]; found {
		message := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&message)
		if err != nil {
			messageAndExit("No Response, check CLI service status")
		}
		if message["message"] == nil {
			messageAndExit("No message in http response")
		}
		messageandDebugExit(message["message"].(string), fmt.Sprintf("HTTP Status %d : %s", resp.StatusCode, message["message"].(string)))
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
			messageAndExit("sync: Error on sendHTTPRequest:")

		}
		Debug.log("RESPONSE ", url, resp.StatusCode, http.StatusText(resp.StatusCode))
		defer resp.Body.Close()
		//Decode the response into data
		decoder := json.NewDecoder(resp.Body)
		var data StacksResponse
		err = decoder.Decode(&data)
		if err != nil {
			messageAndExit("sync: Error on Decode")
		}

		Debug.log(data)
		tWriter := new(tabwriter.Writer)
		tWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)
		if len(data.NewStack) == 0 && (len(data.KabStack) == 0) && len(data.ObsoleteStack) == 0 && len(data.CuratedStack) == 0 && len(data.ActivateStack) == 0 {
			syncedOutput := KabStacksHeader + " are already synchronized with the " + GHStacksHeader
			fmt.Println(strings.ToLower(syncedOutput))
		} else {
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", KabStacksHeader, "Version", "Status")
			fmt.Fprintf(tWriter, "\n%s\t%s\t%s", strings.Repeat("-", len(KabStacksHeader)), "-------", "------")

			var statusMsg string
			var exceptionMsgs []string

			for i := 0; i < len(data.NewStack); i++ {
				statusMsg = "added to Kabanero"
				if data.NewStack[i].ExceptionMessage != "" {
					statusMsg = data.NewStack[i].Status
					exceptionMsgs = append(exceptionMsgs, data.NewStack[i].ExceptionMessage)
				}
				for j := 0; j < len(data.NewStack[i].Versions); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.NewStack[i].Name, data.NewStack[i].Versions[j].Version, statusMsg)
				}
			}
			for i := 0; i < len(data.ActivateStack); i++ {
				statusMsg = "inactive ==> active"
				if data.ActivateStack[i].ExceptionMessage != "" {
					statusMsg = data.ActivateStack[i].Status
					exceptionMsgs = append(exceptionMsgs, data.ActivateStack[i].ExceptionMessage)
				}

				for j := 0; j < len(data.ActivateStack[i].Versions); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ActivateStack[i].Name, data.ActivateStack[i].Versions[j].Version, statusMsg)
				}
			}
			for i := 0; i < len(data.KabStack); i++ {
				for j := 0; j < len(data.KabStack[i].Status); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.KabStack[i].Name, data.KabStack[i].Status[j].Version, data.KabStack[i].Status[j].Status)
				}
			}
			for i := 0; i < len(data.ObsoleteStack); i++ {
				statusMsg = "deleted"
				if data.ObsoleteStack[i].ExceptionMessage != "" {
					statusMsg = data.ObsoleteStack[i].Status
					exceptionMsgs = append(exceptionMsgs, data.ObsoleteStack[i].ExceptionMessage)
				}
				for j := 0; j < len(data.ObsoleteStack[i].Versions); j++ {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ObsoleteStack[i].Name, data.ObsoleteStack[i].Versions[j].Version, statusMsg)
				}
			}

			fmt.Fprintln(tWriter)
			tWriter.Flush()

			for i := 0; i < len(exceptionMsgs); i++ {
				fmt.Println(exceptionMsgs[i])
			}
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
