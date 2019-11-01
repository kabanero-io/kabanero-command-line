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
	"fmt"

	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:    "activate collection-name",
	Hidden: true,
	Short:  "Activate a collection for use by the development team",
	Long: `
A collection can be available to a development team
to use for building applications or not. Activate
will cause the collection to be shown to 
the develoopment team when they list the types of
application they can build.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("activate called")
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
