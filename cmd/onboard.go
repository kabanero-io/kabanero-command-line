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
	"io/ioutil"

	"github.com/spf13/cobra"
)

// onboardCmd represents the onboard command
var onboardCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(1),
	Use:   "onboard github-id repo-url|org-url",
	Short: "Command to onboard a developer to the Kabanero infrastructure",
	Long: `Under development.  In the future this command causes an email to be sent to the user associated
	with the user-id providing the information necessary to get started using 
	Kabanero.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		gituser := args[0]
		var err error
		url := getRESTEndpoint("v1/onboard")
		requestBody, _ := json.Marshal(map[string]string{"gituser": gituser})
		resp, err := sendHTTPRequest("POST", url, requestBody)
		if err != nil {
			return err
		}
		somedata, _ := ioutil.ReadAll(resp.Body)
		err = printPrettyJSON(somedata)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)

}
