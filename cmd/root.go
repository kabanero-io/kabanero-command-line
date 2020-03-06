package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	// for logging
	"k8s.io/klog"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// VERSION is set during build
	VERSION      string
	cfgFile      string
	cliConfig    *viper.Viper
	APIVersionV1 = "v1"
	//dryrun          bool
	verbose         bool
	verboseHTTP     bool
	developerMode   bool
	klogInitialized = false
	KabURLKey       = "KABURL"
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
	Use:   "kabanero",
	Short: "This repo defines a command line interface used by the enterprise, solution, or application architect who defines and manages the kabanero stacks that are used by developers to create governed applications for their business.",
	Long: `**kabanero** is a command line interface for managing the stacks in a Kabanero 
environment, as well as to on-board the people that will use 
the environment to build applications.

Before using the cli please configure the github authorization for the cli service. Steps can be found in the following documentation: https://kabanero.io/docs/ref/general/configuration/github-authorization.html


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

	// we will only allow default config file name/location for now.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kabanero.yaml)")
	// Added for logging
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Turns on debug output and logging to a file in $HOME/.kabanero/logs. If installed with brew the file is in ~/Library/Logs/kabanero/")
	rootCmd.PersistentFlags().BoolVarP(&verboseHTTP, "debug http", "x", false, "Turns on debug output for http request/responses")
	rootCmd.PersistentFlags().BoolVarP(&developerMode, "developer mode", "d", false, "Bypasses certain blocks in the CLI")
	err := rootCmd.PersistentFlags().MarkHidden("debug http")
	if err != nil {
		fmt.Fprintln(os.Stdout, "err with MarkHidden")
	}
	//rootCmd.Execute()
	// not implemented: rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Turns on dry run mode")

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
	//unsafe to do this cause we might log pw:  Debug.log("Running with command line args: kabanero ", strings.Join(os.Args[1:], " "))
	// handle user supplied config file:
	//if cfgFile != "" {
	//	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
	//		Info.log("user supplied --config file does not exist.  Creating: " + cfgFile)
	//		_, err := os.OpenFile(cfgFile, os.O_RDWR|os.O_CREATE, 0755)
	//		if err != nil {
	//			Error.log("ERROR opening user supplied config file: "+cfgFile, err)
	//			os.Exit(1)
	//		}
	//	} else {
	//		Debug.log("using --config file: " + cfgFile)
	//	}
	//}

	// verify the config directory and file:
	cfgDir := filepath.Join(homeDir(), ".kabanero")
	Debug.log("Kabanero config directory: " + cfgDir)
	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		if err := os.Mkdir(cfgDir, os.ModePerm); err != nil {
			Error.log("failed to create config dir: " + cfgDir)
			os.Exit(1)
		}
	}

	// setup Viper  and some defaults
	cliConfig = viper.New()
	cliConfig.SetDefault("home", cfgDir)
	cliConfig.SetDefault("images", "index.docker.io")
	cliConfig.SetDefault("tektonserver", "")

	if cfgFile == "" {
		//viper needs cfgFile to NOT include the file type
		cfgFile := filepath.Join(cfgDir, "config")
		f, err := os.OpenFile(cfgFile+".yaml", os.O_RDWR|os.O_CREATE, 0700)
		if err != nil {
			Error.log("ERROR creating config file config.yaml", err)
			os.Exit(1)
		}
		Debug.log("Config file name: " + cfgFile + ".yaml")
		f.Close()
	}
	cliConfig.SetConfigName("config") // name of config file without extension
	cliConfig.AddConfigPath(cfgDir)

	cliConfig.SetEnvPrefix("KABANERO") // will expect all env vars to be prefixed with "KABANERO_"
	cliConfig.AutomaticEnv()           // read in environment variables that match

	cliConfig.SetConfigType("yaml")
	// If a config file is found, read it in.
	if err := cliConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// no config file
			Debug.log("Config file not found: " + cfgFile)
		} else {
			Error.log("ERROR: config file error: ", err)
		}
	}
	Debug.log("config file used: " + cliConfig.ConfigFileUsed())
}

//func getDefaultConfigFile() string {
//	return filepath.Join(cliConfig.GetString("home"), ".kabanero.yaml")
//}

func Execute(version string) {
	VERSION = version

	if err := rootCmd.Execute(); err != nil {
		Error.log(err)
		os.Exit(1)
	}
}

type kabanerologger string

// define the logging levels
var (
	Info       kabanerologger = "Info"
	Warning    kabanerologger = "Warning"
	Error      kabanerologger = "Error"
	Debug      kabanerologger = "Debug"
	Container  kabanerologger = "Container"
	InitScript kabanerologger = "InitScript"
)

func (l kabanerologger) log(args ...interface{}) {
	msgString := fmt.Sprint(args...)
	l.internalLog(msgString)
}

func (l kabanerologger) logf(fmtString string, args ...interface{}) {
	msgString := fmt.Sprintf(fmtString, args...)
	l.internalLog(msgString)
}

func (l kabanerologger) internalLog(msgString string) {
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

		logDir := filepath.Join(homeDir(), ".kabanero", "logs")

		_, errPath := os.Stat(logDir)
		if errPath != nil {
			Debug.log("Creating log dir ", logDir)
			if err := os.MkdirAll(logDir, 0755); err != nil {
				Error.logf("Could not create %s: %s", logDir, err)
			}
		}

		currentTimeValues := strings.Split(time.Now().Local().String(), " ")
		fileName := strings.ReplaceAll("kabanero"+currentTimeValues[0]+"T"+currentTimeValues[1]+".log", ":", "-")
		pathString := filepath.Join(homeDir(), ".kabanero", "logs", fileName)
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
