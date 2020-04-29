package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type DescribeInfo struct {
	Name           string
	Version        string
	Project        string
	Source         string `json:"git repo url"`
	Image          string `json:"image"`
	Status         string
	DigestCheck    string `json:"digest check"`
	ImageDigest    string `json:"image digest"`
	KabaneroDigest string `json:"kabanero digest"`
}

var describeCmd = &cobra.Command{
	Args:  cobra.MinimumNArgs(2),
	Use:   "describe stack-name version",
	Short: "Get more information about the specified stack",
	Long:  `Get more information about the specified stacks, including the digest values `,
	Run: func(cmd *cobra.Command, args []string) {
		stackName := args[0]
		version := args[1]
		url := getRESTEndpoint("v1/describe/stacks/" + stackName + "/versions/" + version)
		resp, err := sendHTTPRequest("GET", url, nil)
		if err != nil {
			messageAndExit("describe: Error sending HTTP request")
		}

		decoder := json.NewDecoder(resp.Body)
		var describeData DescribeInfo
		err = decoder.Decode(&describeData)
		if err != nil {
			messageAndExit("describe: Error decoding http response")
		}

		Debug.log(describeData)
		fmt.Println("stack name: ", describeData.Name)
		fmt.Println("version: ", describeData.Version)
		fmt.Println("project: ", describeData.Project)
		fmt.Println("source: ", describeData.Source)
		fmt.Println("image: ", describeData.Image)
		fmt.Println("status: ", describeData.Status)
		fmt.Println("digest check: ", describeData.DigestCheck)
		fmt.Println("kabanero digest: ", describeData.KabaneroDigest)
		fmt.Println("image digest: ", describeData.ImageDigest)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
