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

	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(2),
	Use:   "deactivate stack-name version",
	Short: "Remove the specified stack from the list of available application types, without deleting it from the Kabanero instance.",
	Long: `
Run the deactivate command to remove the specified stack from the list of available application types, without deleting it from the Kabanero instance.

This command is useful in a case where you have cloned a stack and customized it for your business needs. Deactivation keeps the base stack in the app hub. The base stack continues to be updated and the updates percolate up to your cloned stack. To restore a deactivated stack, run the kabanero refresh command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// fmt.Println("deactivate called")
		stackName := args[0]
		version := args[1]
		url := getRESTEndpoint("v1/stacks/" + stackName + "/versions/" + version)
		resp, err := sendHTTPRequest("DELETE", url, nil)
		if err != nil {
			Debug.log("deactivate: Error on sendHTTPRequest:")
			return err
		}
		data := make(map[string]interface{})

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			Debug.log("deactivate: Error on Decode:")
			return err
		}
		deactivateResponse := data["status"]
		if deactivateResponse == nil {
			return errors.New("no status with deactivate response")
		}
		// if _, found := data["exception message"]; found {
		fmt.Println(deactivateResponse)
		// }
		// Debug.log(deactivateResponse)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deactivateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deactivateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
