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
	if aConfig.Config.SSLEnabled {
		msg += "  SSL Enabled\n"
	}
	msg += strings.Repeat("#", 73)

	return msg
}

// startServer configures and boots the webservice
func startServer(ctx context.Context, appConfig model.AppConfig) error {
	config := appConfig.Config
	// TODO configure logging from config
	logger := logrus.New()

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	initMetrics(e)

	// Boot our services
	ds, err := service.NewDNSMasqService(config, service.WithLogger(logger))
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
	if config.SSLEnabled {
		err = e.StartTLS(address, config.SSLCertFile, config.SSLKeyFile)
	} else {
		err = e.Start(address)
	}

	return err
}
