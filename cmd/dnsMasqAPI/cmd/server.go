package cmd

import (
	"context"
	"fmt"
	"github.com/cclose/dnsmasq-api/constant/key"
	"github.com/cclose/dnsmasq-api/controller"
	"github.com/cclose/dnsmasq-api/model"
	"github.com/cclose/dnsmasq-api/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

const serverCmdName = "server"

// serverCmd The server subcommand
var serverCmd = &cobra.Command{
	Use:   serverCmdName,
	Short: "run the web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		appConfig, ok := ctx.Value(key.ContextConfig).(model.AppConfig)
		if !ok {
			return fmt.Errorf("app config not found in context. Context failed to load")
		}

		return startServer(ctx, appConfig)
	},
}

// init Register the server subcommand with cobra root cmd
func init() {
	rootCmd.AddCommand(serverCmd)
}

// initMetrics Register Middleware to calculate request metrics
func initMetrics(e *echo.Echo) {
	e.Use(controller.MetricsMiddleware())
}

// startUpMessage creates the boot up message for the service, informing the log of the service's Version and config
func startUpMessage(aConfig model.AppConfig, lAddr string) string {
	msg := strings.Repeat("#", 73)
	msg += fmt.Sprintf(
		"\n%s - %s Version %s (%s)\n"+
			"Starting Server:\n"+
			"  Listening on %s\n"+
			"  Tracking DNSMasq Config: %s\n"+
			"  DNSMasq Service Reloading: %v\n",
		rootCmd.Name(), serverCmdName, aConfig.BuildInfo.Version, aConfig.BuildInfo.Commit,
		//
		lAddr,
		aConfig.Config.DnsmasqConfig,
		!aConfig.Config.SkipDNSMasqReload, // invert the boolean since it is for skipping
	)
	if aConfig.Config.SSL.Enabled {
		msg += "  SSL Enabled\n"
	}
	msg += strings.Repeat("#", 73)

	return msg
}

// configureLogging Configures the logging provider based on the logging config
func configureLogging(lConfig model.LoggingConfig) (*logrus.Logger, error) {
	logger := logrus.New()
	// If we set a file path
	if lConfig.FilePath != "" {
		// try to open the file and set as output
		logFile, err := os.OpenFile(lConfig.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(logFile)
		if lConfig.Output != "" && strings.ToLower(lConfig.Output) != "file" {
			// We defer this warning so we can finish setting up
			defer func() {
				logger.Warnf("ignoring contradictory config log.output = %s because log.file_path was set",
					lConfig.Output)
			}()
		}

		// file_path cancels out output
	} else if lConfig.Output != "" {
		// We're stupidly permissive here... lol
		switch strings.ToLower(lConfig.Output) {
		case "stdout", "standardout", "out":
			logger.SetOutput(os.Stdout)
		case "stderr", "standarderror", "standarderr", "stderrout", "err":
			logger.SetOutput(os.Stderr)
		case "syslog":
			//logger.SetOutput()
			return nil, fmt.Errorf("configuration log.output = syslog is not yet implemented. Use stdout")
		case "file":
			return nil, fmt.Errorf("invalid configuration log.output = file: set log.file_path instead")
		default:
			return nil, fmt.Errorf("unknown configuration log.output option '%s'", lConfig.Output)
		}
	}

	if lConfig.Level != "" {
		loggingLevel, err := logrus.ParseLevel(lConfig.Level)
		if err != nil {
			return nil, err
		}
		logger.SetLevel(loggingLevel)
	}

	return logger, nil
}

// startServer configures and boots the webservice
func startServer(ctx context.Context, appConfig model.AppConfig) error {
	config := appConfig.Config
	logger, err := configureLogging(config.Logging)
	// Make sure we close our log file when the server shuts down
	// check this BEFORE err incase we opened a file but still returned an err
	if logCloser, ok := logger.Out.(io.Closer); ok {
		defer logCloser.Close() // we're just gunna eat this error since we're shtting down
	}
	if err != nil {
		return err
	}

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	initMetrics(e)

	// Boot our services
	ds, err := service.NewDNSMasqService(config, service.WithLogger(logger), service.WithConfig(config.DB))
	if err != nil {
		return err
	}

	// Register our Controllers
	dc := controller.NewDnsController(ds)
	dc.Register(e)
	sc := controller.NewStatusController(appConfig.BuildInfo)
	sc.Register(e)

	// Calculate service address and boot
	address := fmt.Sprintf(":%d", config.Port)
	logger.Print(startUpMessage(appConfig, address))
	if config.SSL.Enabled {
		err = e.StartTLS(address, config.SSL.CertFile, config.SSL.KeyFile)
	} else {
		err = e.Start(address)
	}

	return err
}
