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
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/kabanero-io/kabanero-command-line/pkg/security"
	"github.com/spf13/cobra"

	"encoding/base64"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	InsecureTLS bool
	clientCert  string
)

type JWTResponse struct {
	JWT     string
	Message string
}

func parseKabURL(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, "/")
	return url
}

func is06Compatible() bool {
	url := getRESTEndpoint("v1/image")
	resp, err := sendHTTPRequest("GET", url, nil)
	if err != nil {
		Info.logf("kabanero command line service image: CLI service unavailable, %s", err.Error())
		return false
	}
	var versionJSON VersionJSON
	err = json.NewDecoder(resp.Body).Decode(&versionJSON)
	if err != nil {
		messageAndExit("Error decoding version response for compatibility check") //this will osexit, not return
		return false
	}
	servicesVersion := strings.Split(versionJSON.Image, ":")[1]
	if servicesVersion == "latest" {
		return true
	}
	if strings.Contains(servicesVersion, ".") { //check if this is a version eg: 0.9.1 vs a sha256
		servicesVersion1stVal, _ := strconv.Atoi(strings.Split(servicesVersion, ".")[0])
		servicesVersion2ndVal, _ := strconv.Atoi(strings.Split(servicesVersion, ".")[1])
		if servicesVersion2ndVal < 6 && servicesVersion1stVal == 0 {

			fmt.Printf("\nYour current CLI version (%s) is incompatible with the command line service image (%s). Please upgrade your command line service to version 0.6.0 or greater, or get a version of the CLI that matches the service image\n", VERSION, servicesVersion)
			return false
		}
	}
	return true
}

func HandleTLSFLag(insecureTLS bool) {
	cliConfig.Set("insecureTLS", insecureTLS)
	err := cliConfig.WriteConfig()
	if err != nil {
		messageAndExit("There was a problem writing to the cli config")
	}

	if clientCert != "" {
		cliConfig.Set(CertKey, clientCert)
		err = cliConfig.WriteConfig()
		if err != nil {
			messageAndExit("There was a problem writing to the cli config")
		}
		return
	}

	if !insecureTLS && clientCert == "" {

		fmt.Print("Are you sure you want to continue with an insecure connection to " + cliConfig.GetString(KabURLKey) + " (y/n): ")

		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println(err)
			//TODO handle incorrect characters or full yes
		}
		fmt.Println()
		switch unicode.ToLower(char) {
		case 'y':
			cliConfig.Set("insecureTLS", true)
			err = cliConfig.WriteConfig()
			if err != nil {
				messageAndExit("There was a problem writing to the cli config")
			}
		case 'n':
			cliConfig.Set("insecureTLS", false)
			err = cliConfig.WriteConfig()
			if err != nil {
				messageAndExit("There was a problem writing to the cli config")
			}

			if cliConfig.GetString(CertKey) == "" {
				messageAndExit("To continue with a secure connection, provide path to certificate with --certificate-authority at login. See login -h for help.")
			}

		default:
			messageAndExit("Please enter y or n")

		}

	}

}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login kabanero-cli-url -u Github userid \n  At the password prompt, enter your GitHub Password/PAT",
	Short: "Will authenticate you to a Kabanero instance",
	Long: `
	Logs in to a Kabanero instance using Github credentials, and stores a temporary access token for subsequent command line calls.
	The temporary authentication token will be stored in your-home-directory/.kabanero/config.yaml.
	Use your Github userid and either password or Personal Access Token (PAT).
	* If you use a GitHub Personal Access Token, make sure it has **read:org** scope. You can select this when creating your PAT in GitHub
	`,
	Example: `
	# login with Github userid and password:
	kabanero login my.kabaneroInstance.io -u myGithubID 
	# login with previously used url Github userid and PAT:
	kabanero login -u myGithubID 
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		Debug.log("login called")
		var err error

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		if username == "" {
			fmt.Printf("Username:")
			bytePwd, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return err
			}
			//ePass = base64.StdEncoding.EncodeToString(bytePwd)
			username = strings.TrimSpace(string(bytePwd))
			fmt.Println()
		}

		if password == "" {
			fmt.Printf("Password:")
			bytePwd, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return err
			}
			//ePass = base64.StdEncoding.EncodeToString(bytePwd)
			password = strings.TrimSpace(string(bytePwd))
			fmt.Println()
		}
		ePass := base64.StdEncoding.EncodeToString([]byte(password))
		eUser := base64.StdEncoding.EncodeToString([]byte(username))

		password = ""
		username = ""
		fmt.Println(username)
		fmt.Println(password)

		var kabLoginURL string

		viper.SetEnvPrefix("KABANERO")

		if len(args) > 0 {
			cliConfig.Set(KabURLKey, parseKabURL(args[0]))
			err = cliConfig.WriteConfig()
			if err != nil {
				return err
			}
		} else {
			if cliConfig.GetString(KabURLKey) == "" {
				messageAndExit("No Kabanero instance url specified")
			}
		}

		HandleTLSFLag(InsecureTLS)

		kabLoginURL = getRESTEndpoint("login")

		requestBody, _ := json.Marshal(map[string]string{"000_ERG_TEN_TWENTY": eUser, "010_BOHM_THIRTY_FIVE": ePass})

		resp, err := sendHTTPRequest("POST", kabLoginURL, requestBody)
		if err != nil {
			messageAndExit("login: Error on sendHTTPRequest:")
		}
		requestBody = nil
		if requestBody == nil {
			fmt.Print()
		}

		Debug.log("RESPONSE ", kabLoginURL, resp.StatusCode, http.StatusText(resp.StatusCode))
		if resp.StatusCode == 404 {
			messageAndExit("The url: " + cliConfig.GetString(KabURLKey) + " is not a valid kabanero url")
		}

		var data JWTResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			Debug.log(err)
			return err
		}
		key := security.Create32BKey((time.Now().String()))
		cliConfig.Set("key", key)

		encryptedJWT := security.EncryptString(data.JWT, key)
		cliConfig.Set("jwt", encryptedJWT)

		err = cliConfig.WriteConfig()
		if err != nil {
			return err
		}
		if cliConfig.GetString("jwt") == "" {
			messageAndExit("Unable to validate user: " + username + " to " + cliConfig.GetString(KabURLKey))
		}
		key = ""
		fmt.Println(key)

		if !is06Compatible() {

			url := getRESTEndpoint("logout")
			resp, err := sendHTTPRequest("POST", url, nil)
			if err != nil {
				return err
			}

			defer resp.Body.Close()
			cliConfig.Set("jwt", "")
			err = cliConfig.WriteConfig()
			if err != nil {
				return err
			}
		} else {

			fmt.Println("Logged in to Kabanero instance: " + cliConfig.GetString(KabURLKey))
			Debug.log("Logged in to Kabanero instance: " + cliConfig.GetString(KabURLKey))
		}
		defer resp.Body.Close()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("username", "u", "", "github username")

	//_ = loginCmd.MarkFlagRequired("username") // possibly comment out to make username flad not required and add promot for username
	loginCmd.Flags().StringP("password", "p", "", "github password/PAT. If no password is provided, prompt will appear")
	loginCmd.Flags().BoolVar(&InsecureTLS, "insecure-skip-tls-verify", false, "If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")
	loginCmd.Flags().StringVar(&clientCert, "certificate-authority", "", "Path to a cert file for the certificate authority")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
