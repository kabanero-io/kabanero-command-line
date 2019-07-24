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

// listCollectionsCmd represents the listCollections command
var listCollectionsCmd = &cobra.Command{
	Use:   "listCollections",
	Short: "Show the list of named collections that can be used to create applications.",
	Long: `A Kabanero collection is a governed set of meta-data
that provides pre-configured application structure, and deployment
information that allows the developer to focus on developing code. 
And once the code is stored in the source control system, the 
meta-data directs out the application to be built and deployed into
the develop downstream chain that gets it deployed into production
eventually. 

The critical piece of meta-data inside the kabanero collection is the 
Appsody stack. Appsody gives you pre-configured stacks and templates for a 
growing set of popular open source runtimes and frameworks, providing a 
foundation on which to build applications for Kubernetes and Knative deployments. 
This allows developers to focus on their code, reducing the learning curve for 
cloud-native development and enabling rapid development for these cloud-native
applications.`,
	Example: "listCollectionsCmd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listCollections called")
	},
}

func init() {
	rootCmd.AddCommand(listCollectionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCollectionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCollectionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
