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

// createNewCmd represents the createNew command
var createNewCmd = &cobra.Command{
	Use:   "createNew collection-name",
	Short: "Create a new collection, starting with a template project.",
	Long: `
A collection is a set of meta data and base images that are used by
the developer to build a microservice or application.  This meta data includes
information about the type of application, languages, pipeline definitions,
management artifacts to install, and the base container.

In the createNew case, we provide a template for the metadata. But the 
collection provider must provide the information and containers in the 
local registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createNew called")
	},
}

func init() {
	rootCmd.AddCommand(createNewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createNewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createNewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
