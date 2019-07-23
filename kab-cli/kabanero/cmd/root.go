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
	Use:   "kabanero",
	Short: "kabanero",
	Long: `Kabanero enables developers, architects, and operations to work together, 
faster. In a single solution, architects and operations can include their company’s 
standards for aspects like security and build pipelines into a customized stack that 
developers use. Kabanero gives enterprises the control they need for areas related 
to governance and compliance, while also meeting the developers’ need for agility 
and speed.

The kabanero command line tool is targetted at the developer.  It provides a mechanism
to get a set of governed application stacks from the Kabanero Application Hub, as well 
as learn about your role within the governed organization, what Git repos and docker
registries your code will be in, and where your test cluster is for combining your 
changes with others on your team.

Complete documentation is available at https://kabanero.io`,
}

func init() {
	// Don't run this on help commands
	// TODO - instead of the isHelpCommand() check, we should delay the config init/ensure until we really need the config
	if !isHelpCommand() {
		cobra.OnInitialize(initLogging)
		cobra.OnInitialize(initConfig)

		//cobra.OnInitialize(ensureConfig)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kabanero.yaml)")
	// Added for logging
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Turns on debug output and logging to a file in $HOME/.kabanero/logs")

	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "Turns on dry run mode")

}

func GetkabaneroCLI() cobra.Command {
	return *rootCmd
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
	Debug.log("Running with command line args: kabanero ", strings.Join(os.Args[1:], " "))
	cliConfig = viper.New()

	cliConfig.SetDefault("home", filepath.Join(homeDir(), ".kabanero"))
	cliConfig.SetDefault("images", "index.docker.io")
	cliConfig.SetDefault("tektonserver", "")
	if cfgFile != "" {
		// Use config file from the flag.
		cliConfig.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".hello-cobra" (without extension).
		cliConfig.AddConfigPath(cliConfig.GetString("home"))
		cliConfig.SetConfigName(".kabanero")
	}

	cliConfig.SetEnvPrefix("kabanero")
	cliConfig.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// Ignore errors, if the config isn't found, we will create a default later
	_ = cliConfig.ReadInConfig()
}

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
