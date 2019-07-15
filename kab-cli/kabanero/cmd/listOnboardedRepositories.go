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
	"fmt"

	"github.com/spf13/cobra"
)

// listOnboardedRepositoriesCmd represents the listOnboardedRepositories command
var listOnboardedRepositoriesCmd = &cobra.Command{
	Use:   "listOnboardedRepositories",
	Short: "As a developer, you have been assigned a team, and each team has a set of repositories into which to put code. This provides the list that you are able to access.",
	Long: `As a developer, you have been assigned a team, and each team has a set of repositories into which to put code. 
	This provides the list of existing repositories that you are able to access. It will also provide the organization and name prefix for 
	any new repositories that you will create.
	`,
	Example: `
	kabanero listOnBoardedRepositories
	
	07142019 11:30:52 Existing Repositories:
	  kabanario.mytest-testrepo
	  kabanerio.mytest-test2repo
	07142019 11:30:53 New repoositories use the following naming convention:
	  kabaneroio.mytest-*`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listOnboardedRepositories called")
	},
}

func init() {
	rootCmd.AddCommand(listOnboardedRepositoriesCmd)

}
