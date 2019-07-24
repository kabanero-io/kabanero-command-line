package cmd

//import collections "github.com/kabanero-command-line/kab-cli/kabanero/collections/cmd"
//import access "github.com/kabanero-command-line/kab-cli/kabanero/onboard/cmd"

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	// for logging
	"k8s.io/klog"

	//  homedir "github.com/mitchellh/go-homedir"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// VERSION is set during build
	VERSION         string
	cfgFile         string
	cliConfig       *viper.Viper
	APIVersionV1    = "v1"
	dryrun          bool
	verbose         bool
	klogInitialized = false
)

func homeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		Error.log(err)
		os.Exit(1)
	}
	return home
}

var rootCmd = &cobra.Command{
	Use:   "kabanero-mgmt",
	Short: "A command line interface that can be used to manage the environment.",
	Long: `A command line interface that can be used to manage the collections that 
the environment presents, as well as on-board the people and clusters that will be
used in the environment to build applications.

Complete documentation is available at https://kabanero.io`,
}

func init() {
	// Don't run this on help commands
	// TODO - instead of the isHelpCommand() check, we should delay the config init/ensure until we really need the config
	if !isHelpCommand() {
		cobra.OnInitialize(initLogging)
		cobra.OnInitialize(initConfig)
		//	cobra.OnInitialize(ensureConfig)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kabaneromgmt.yaml)")
	// Added for logging
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Turns on debug output and logging to a file in $HOME/.kabaneromgmt/logs")

	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Turns on dry run mode")

	// The subbcommand processor for commands to manage the apphub
	//	var apphubCmd = cobra.Command(collections.GetCollectionsCLI())
	//	apphubCmd.Use = "apphub"
	//	apphubCmd.Short = "Commands manage an application hub (apphub)"
	//	rootCmd.AddCommand(&apphubCmd)

	// The subcommand processoor for commands to grant and revoke access to Git Repos and there relationship to
	// name spaces in Kubes
	//	var accessCmd = cobra.Command(access.GetAccessCLI())
	//	accessCmd.Use = "access"
	//	accessCmd.Short = "Commands to grant and revoke access to Kabanero Repos"
	//	rootCmd.AddCommand(&accessCmd)
}

func isHelpCommand() bool {
	if len(os.Args) <= 1 {
		return true
	}
	for _, arg := range os.Args {
		if arg == "help" || arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func initConfig() {
	Debug.log("Running with command line args: kabanero-mgmt ", strings.Join(os.Args[1:], " "))
	cliConfig = viper.New()

	cliConfig.SetDefault("home", filepath.Join(homeDir(), ".kabaneromgmt"))
	cliConfig.SetDefault("images", "index.docker.io")
	cliConfig.SetDefault("tektonserver", "")
	if cfgFile != "" {
		// Use config file from the flag.
		cliConfig.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".hello-cobra" (without extension).
		cliConfig.AddConfigPath(cliConfig.GetString("home"))
		cliConfig.SetConfigName(".kabaneromgmt")
	}

	cliConfig.SetEnvPrefix("kabaneromgmt")
	cliConfig.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// Ignore errors, if the config isn't found, we will create a default later
	_ = cliConfig.ReadInConfig()
}

func getDefaultConfigFile() string {
	return filepath.Join(cliConfig.GetString("home"), ".kabaneromgmt.yaml")
}

func Execute(version string) {
	VERSION = version

	if err := rootCmd.Execute(); err != nil {
		Error.log(err)
		os.Exit(1)
	}
}

type kabaneromgmtlogger string

// define the logging levels
var (
	Info       kabaneromgmtlogger = "Info"
	Warning    kabaneromgmtlogger = "Warning"
	Error      kabaneromgmtlogger = "Error"
	Debug      kabaneromgmtlogger = "Debug"
	Container  kabaneromgmtlogger = "Container"
	InitScript kabaneromgmtlogger = "InitScript"
)

func (l kabaneromgmtlogger) log(args ...interface{}) {
	msgString := fmt.Sprint(args...)
	l.internalLog(msgString)
}

func (l kabaneromgmtlogger) logf(fmtString string, args ...interface{}) {
	msgString := fmt.Sprintf(fmtString, args...)
	l.internalLog(msgString)
}

func (l kabaneromgmtlogger) internalLog(msgString string) {
	if l == Debug && !verbose {
		return
	}

	if verbose || l != Info {
		msgString = "[" + string(l) + "] " + msgString
	}

	// Print to console
	if l == Info {
		fmt.Fprintln(os.Stdout, msgString)
	} else {
		fmt.Fprintln(os.Stderr, msgString)
	}

	// Print to log file
	if verbose && klogInitialized {
		klog.InfoDepth(2, msgString)
		klog.Flush()
	}
}

func initLogging() {

	if verbose {

		logDir := filepath.Join(homeDir(), ".kabaneromgmt", "logs")

		_, errPath := os.Stat(logDir)
		if errPath != nil {
			Debug.log("Creating log dir ", logDir)
			if err := os.MkdirAll(logDir, 0755); err != nil {
				Error.logf("Could not create %s: %s", logDir, err)
			}
		}

		currentTimeValues := strings.Split(time.Now().Local().String(), " ")
		fileName := strings.ReplaceAll("kabaneromgmt"+currentTimeValues[0]+"T"+currentTimeValues[1]+".log", ":", "-")
		pathString := filepath.Join(homeDir(), ".kabaneromgmt", "logs", fileName)
		klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
		klog.InitFlags(klogFlags)
		_ = klogFlags.Set("v", "4")
		_ = klogFlags.Set("skip_headers", "false")
		_ = klogFlags.Set("skip_log_headers", "true")
		_ = klogFlags.Set("log_file", pathString)
		_ = klogFlags.Set("logtostderr", "false")
		_ = klogFlags.Set("alsologtostderr", "false")
		klogInitialized = true
		Debug.log("Logging to file ", pathString)
	}
}
