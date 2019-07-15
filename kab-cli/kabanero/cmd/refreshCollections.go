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

// refreshCollectionsCmd represents the refreshCollections command
var refreshCollectionsCmd = &cobra.Command{
	Use:   "refreshCollections",
	Short: "Refresh the cache of data that is maintained in the client side of the command line interface",
	Long: `The kabanero command line interface caches information on collections and other data needed
	to quickly run requests. This cache is created when you logon, if you stay logged on for long periods
	of time, this command allows you to refresh this cache without logging off and on.`,
	Example: `kabanero refreshCollections`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("refreshCollections called")
	},
}

func init() {
	rootCmd.AddCommand(refreshCollectionsCmd)

}
