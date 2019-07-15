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

// createCloneCmd represents the createClone command
var createCloneCmd = &cobra.Command{
	Use:   "createClone new-collectiom-name collection-name",
	Short: `Clone an existing collection and build a new collection from it`,
	Long: `
A collection is a set of meta data and base images that are use by
the developer to build a microservice or application.  This meta data includes
information about the type of appliciation, languages, pipeline definitions,
management aftifacts to install, and the base container.  
	
In the createClone case, we provide clone the metadata, renaming the collection
in the process. The collection provider can then update the information 
and containers in the local registry`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createClone called")
	},
}

func init() {
	rootCmd.AddCommand(createCloneCmd)

}
