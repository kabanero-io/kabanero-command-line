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
	"fmt"

	// "net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// The list command gets a set of "stacks" back from the CLI Service.

// KabStruct : represents the JSON returned for the Kabanero stacks.
type KabStruct struct {
	Name   string
	Status []StatusStruct `json:"status"`
}
type StatusStruct struct {
	Version string
	Status  string
}

// CommonStackStruct : represents the JSON returned for all the other stacks.
type CommonStackStruct struct {
	Name     string
	Versions []VersionStruct `json:"versions"`
}

type VersionStruct struct {
	Images   []string `json:"image"`
	Reponame string
	Version  string
}

// StacksResponse : all the stacks from GH
type StacksResponse struct {
	NewStack      []CommonStackStruct `json:"new curated stacks"`
	ActivateStack []CommonStackStruct `json:"activate stacks"`
	KabStack      []KabStruct         `json:"kabanero stacks"`
	ObsoleteStack []CommonStackStruct `json:"obsolete stacks"`
	CuratedStack  []CommonStackStruct `json:"curated stacks"`
	Repos         []ReposStruct       `json:"repositories"`
}

// Repos : yaml sources that the stacks are coming form
type ReposStruct struct {
	Name string
	URL  string
}

// KabStacksHeader for all references to what we call the "Kab stacks"
var KabStacksHeader = "Kabanero Instance Stacks "

// GHStacksHeader for all references to the "curated stacks"
var GHStacksHeader = "GitHub Curated Stacks"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list ",
	Short: "List all the stacks in the kabanero instance, and their status",
	Long: `List all the stacks in the kabanero instance, and their status. 
	Modifications to the curated stack may be slow to replicate in git hub and therefore may not be reflected immediately in KABANERO LIST or SYNC display output`,
	RunE: func(cmd *cobra.Command, args []string) error {
		Debug.log("List called...")
		url := getRESTEndpoint("v1/stacks")
		resp, err := sendHTTPRequest("GET", url, nil)
		if err != nil {
			Debug.log("list: Error on sendHTTPRequest:")
			return err
		}
		// cannot reference resp here.  May not be fully formed and cause nil pointer deref: Debug.log("RESPONSE ", url, resp.StatusCode, http.StatusText(resp.StatusCode))
		//Decode the response into data
		decoder := json.NewDecoder(resp.Body)
		var data StacksResponse
		err = decoder.Decode(&data)
		if err != nil {
			Debug.log("list: Error on Decode:")
			return err
		}

		Debug.log(data)
		fmt.Println()
		fmt.Println("Kabanero CLI service url: ", cliConfig.GetString(KabURLKey))

		tWriter := new(tabwriter.Writer)
		tWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)

		//Kabenero Stacks
		fmt.Fprintf(tWriter, "\n%s\t%s\t%s", KabStacksHeader, "Version", "Status")
		fmt.Fprintf(tWriter, "\n%s\t%s\t%s", strings.Repeat("-", len(KabStacksHeader)), "-------", "------")

		for i := 0; i < len(data.KabStack); i++ {
			for j := 0; j < len(data.KabStack[i].Status); j++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.KabStack[i].Name, data.KabStack[i].Status[j].Version, data.KabStack[i].Status[j].Status)
			}
		}
		for i := 0; i < len(data.ObsoleteStack); i++ {
			for j := 0; j < len(data.ObsoleteStack[i].Versions); j++ {
				fmt.Fprintf(tWriter, "\n%s\t%s\t%s", data.ObsoleteStack[i].Name, data.ObsoleteStack[i].Versions[j].Version, "obsolete")
			}
		}
		fmt.Fprintln(tWriter)
		tWriter.Flush()

		// put new stacks name/version into a map to compare to curated because the new stacks overlap with the regular list of stacks
		mNewStack := make(map[string]string)
		for i := 0; i < len(data.NewStack); i++ {
			for j := 0; j < len(data.NewStack[i].Versions); j++ {
				mNewStack[data.NewStack[i].Name] = data.NewStack[i].Name + data.NewStack[i].Versions[j].Version
			}
		}

		fmt.Println()
		fmt.Println()
		fmt.Println("GitHub Curated Stacks (repo name - url):")
		for i := 0; i < len(data.Repos); i++ {
			fmt.Printf("   %s - %s", data.Repos[0].Name, data.Repos[0].URL)
		}
		fmt.Println()

		fmt.Fprintf(tWriter, "\n%s\t%s\t%s", GHStacksHeader, "Version", "Repo")
		fmt.Fprintf(tWriter, "\n%s\t%s\t%s", strings.Repeat("-", len(GHStacksHeader)), "-------", "----")
		for i := 0; i < len(data.CuratedStack); i++ {
			name := data.CuratedStack[i].Name
			for j := 0; j < len(data.CuratedStack[i].Versions); j++ {
				version := data.CuratedStack[i].Versions[j].Version
				nameAndVersion := name + version

				//fmt.Fprintf(tWriter, "\n%s", name)
				_, found := mNewStack[nameAndVersion]
				if found {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", name+" (new)", version, data.CuratedStack[i].Versions[j].Reponame)
				} else {
					fmt.Fprintf(tWriter, "\n%s\t%s\t%s", name, version, data.CuratedStack[i].Versions[j].Reponame)
				}
			}
		}
		fmt.Fprintln(tWriter)
		tWriter.Flush()
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
