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

// listMyinfoCmd represents the listMyinfo command
var listMyinfoCmd = &cobra.Command{
	Use:   "listMyinfo",
	Short: "List the information I need in order to use Kabanero",
	Long: `Kabanero provides a development flow for building 
applications quickly.  In order for a developer to engage in this 
process, the developer must be added to a team, know where information
is stored in Git,  where their containers are stored,  how to gain access
to the test integration server.  All this information is returned from this
command.`,
	Example: `kabanero listMyinfo`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listMyinfo called")
	},
}

func init() {
	rootCmd.AddCommand(listMyinfoCmd)
}
