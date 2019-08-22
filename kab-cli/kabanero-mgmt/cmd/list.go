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
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type CollStruct struct {
	OriginalName string
	Name         string
	Version      string
}

type CollectionsResponse struct {
	NewColl      []CollStruct `json:"new collections"`
	KabColl      []CollStruct `json:"kabanero collection"`
	ObsoleteColl []CollStruct `json:"obsolete collections"`
	MasterColl   []CollStruct `json:"master collection"`
	VChangeColl  []CollStruct `json:"version change collections"`
}

func printPrettyJSON(jsonData []byte) {
	var testBuffer bytes.Buffer
	json.Indent(&testBuffer, jsonData, "", "\t")
	fmt.Println(string(testBuffer.Bytes()))

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
		url := cliConfig.GetString(KabURLKey) + "/v1/collections"
		resp, err := sendHTTPRequest("GET", url, nil)
		if err != nil {
			return errors.New(err.Error())
		}

		Debug.log("RESPONSE ", url, resp.StatusCode, http.StatusText(resp.StatusCode))

		//Decode the response into data
		decoder := json.NewDecoder(resp.Body)
		var data CollectionsResponse
		err = decoder.Decode(&data)
		if err != nil {
			return err
		}
		//
		Debug.log(data)
		tWriter := new(tabwriter.Writer)
		tWriter.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Println("Masters Collections")
		fmt.Fprintf(tWriter, "\n%s\t%s\t", "Name", "Version")
		fmt.Fprintf(tWriter, "\n%s\t%s\t", "----", "----")

		for i := 0; i < len(data.MasterColl); i++ {
			fmt.Fprintf(tWriter, "\n %s\t%s", data.MasterColl[i].Name, data.MasterColl[i].Version)
		}
		tWriter.Flush()
		fmt.Println("")

		// fmt.Println(data.KabColl[0].OriginalName)

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
