package cmd

import (
	"context"
	"fmt"
	"github.com/cclose/dnsmasq-api/constant/envvar"
	"github.com/cclose/dnsmasq-api/constant/key"
	"github.com/cclose/dnsmasq-api/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

const dnsMasqRootCmdName = "dnsMasqAPI"

var (
	// pflag targets
	// configFile The location of the Configuration File
	configFile = "config.yaml"
	// verbose Print more better
	verbose bool
	// versionMode Captures version mode flag
	versionMode bool

	// build-time injection targets

	// Version The Version of the code
	Version string
	// Commit the SHA of the tip Commit when the code was built
	Commit string
	// BuildTimeStr when the code was built string representation
	BuildTimeStr string
	// BuildTime time object for BuildTimeStr
	BuildTime time.Time

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   dnsMasqRootCmdName,
		Short: "An API Companion for DNSMasq",
		Long: `DNSMasq API

A Companion API for DNSMasq that provides the ability to manage DNS entries remotely over REST.
Provides the ability to GET, DELETE, and POST entries into the DNSMasq config, as well as reload the configuration.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// This is a cobra thing
func Execute(buildTime, commit, version string) {
	BuildTimeStr = buildTime
	Commit = commit
	Version = version

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init Setup the Root command's configuration and options
func init() {
	// Set default BuildTime if it's not provided
	if BuildTimeStr == "" {
		BuildTime = time.Now().UTC()
	} else {
		var err error
		BuildTime, err = time.Parse(time.RFC3339, BuildTimeStr)
		if err != nil {
			fmt.Printf("Invalid build time format: %v. Defaulting to current time.\n", err)
			BuildTime = time.Now().UTC()
		}
	}

	viper.SetDefault("author", "Cory Close <pulsar2612@hotmail.com>")
	viper.SetDefault("license", "bsd-3-clause")
	cobra.OnInitialize(initViper)

	// Add -c flag and bind to viper
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yaml",
		"config file (default is config.yaml)",
	)
	err := viper.BindPFlag(envvar.ViperConfig, rootCmd.PersistentFlags().Lookup(envvar.ViperConfig))
	if err != nil {
		log.Fatal(err)
	}

	// Add -v flag and bind to viper
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Display more verbose output in console output. (default: false)",
	)
	err = viper.BindPFlag(envvar.ViperVerbose, rootCmd.PersistentFlags().Lookup(envvar.ViperVerbose))
	if err != nil {
		log.Fatal(err)
	}

	// Add -v flag and bind to viper
	rootCmd.PersistentFlags().BoolVarP(&versionMode, "version", "V", false,
		"Display version information and exit",
	)
	err = viper.BindPFlag(envvar.ViperVersion, rootCmd.PersistentFlags().Lookup(envvar.ViperVersion))
	if err != nil {
		log.Fatal(err)
	}
}

// initViper Sets the Bindings and envvars for viper to read
func initViper() {
	// Bind the environment variable CONFIG to the config key
	err := viper.BindEnv(envvar.ViperConfig, envvar.Config)
	if err != nil {
		log.Fatal(err)
	}
}

// initConfig Reads the configuration file, build data, and builds the application context
func initConfig(cmd *cobra.Command) error {
	cFile := viper.GetString(envvar.ViperConfig)
	viper.SetConfigFile(cFile)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Build Config from configuration file
	config := &model.Config{}
	if err := viper.Unmarshal(config); err != nil {
		return fmt.Errorf("unable to decode config, %v", err)
	}

	// Load app build data
	buildInfo := model.BuildInfo{
		BuildTime: BuildTime,
		Commit:    Commit,
		Version:   Version,
	}

	// Package information in one struct
	appConfig := model.AppConfig{
		BuildInfo: buildInfo,
		Config:    *config,
	}

	// Load a context and store the AppConfig
	baseCtx := cmd.Context()
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	ctx := context.WithValue(baseCtx, key.ContextConfig, appConfig)
	cmd.SetContext(ctx)

	// TODO find a better way to do this
	if viper.GetBool(envvar.ViperVersion) {
		fmt.Printf("%s version %s (%s) - BuildTime %s\n", dnsMasqRootCmdName, Version, Commit, BuildTimeStr)
		os.Exit(0)
	}

	return nil
}
