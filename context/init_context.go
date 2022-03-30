package context

import (
	"os"

	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/config"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/logger"
)

type UriScheme string

// List of UriScheme
const (
	UriScheme_HTTP  UriScheme = "http"
	UriScheme_HTTPS UriScheme = "https"
)

func InitContext(context *ServiceContext) {
	localConfig := config.ServiceConfig
	logger.UtilLog.Infof("config Info: Description[%s]", localConfig.Info.Version, localConfig.Info.Description)
	configuration := localConfig.Configuration

	context.SBIPort = config.SORAF_DEFAULT_PORT_INT // default port
	if sbi := configuration.Sbi; sbi != nil {
		context.UriScheme = sbi.Scheme

		if sbi.Port != 0 {
			context.SBIPort = sbi.Port
		}

		context.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if context.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			context.BindingIPv4 = sbi.BindingIPv4
			if context.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				context.BindingIPv4 = "0.0.0.0"
			}
		}
	}

}
