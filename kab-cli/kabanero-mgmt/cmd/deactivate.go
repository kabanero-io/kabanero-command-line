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

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate collection-name",
	Short: "Prevent this collection from being shown to the development team, while not deleting it.",
	Long: `
A collection can be available to a development team
to use for building applications or not. Deactivate
will cause the collection to not be shown to 
the development team when they list the types of
application they can build.

This would be done in the case where you have cloned the collection
and made changes for your business.  This keeps the base collection
in the apphub, and it will continue to be updated, and the 
updates will be percolated up to your cloned collection.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deactivate called")
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
