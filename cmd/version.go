// Copyright © 2019 IBM Corporation and others.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

type VersionJSON struct {
	Version string
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Kabanero CLI version",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		Info.log("kabanero cli version: ", VERSION)

		if cliConfig.GetString(KabURLKey) == "" {
			return nil
		}

		url := getRESTEndpoint("v1/version")
		resp, err := sendHTTPRequest("GET", url, nil)
		if err != nil {
			return err
		}
		var versionJSON VersionJSON
		err = json.NewDecoder(resp.Body).Decode(&versionJSON)
		if err != nil {
			return err
		}
		Info.log("kabanero command line service version: ", versionJSON.Version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
