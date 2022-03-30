package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/api"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/config"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/context"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/http2"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/logger"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/util"
)

type SERVICE struct{}

type (
	// Config information.
	Config struct {
		servicecfg string
	}
)

var configVar Config

var serviceCLi = []cli.Flag{
	cli.StringFlag{
		Name:  "servicecfg",
		Usage: "util file",
	},
}

var initLog *logrus.Entry

func init() {
	initLog = logger.InitLog
}

func (*SERVICE) GetCliCmd() (flags []cli.Flag) {
	return serviceCLi
}

func (service *SERVICE) Initialize(c *cli.Context) error {
	configVar = Config{
		servicecfg: c.String("servicecfg"),
	}

	if configVar.servicecfg != "" {
		initLog.Infof("SERVICE util path is : [%s] ", configVar.servicecfg)
		if err := config.InitConfigFactory(configVar.servicecfg); err != nil {
			return err
		}
	} else {
		initLog.Infof("SERVICE util use default path service/util/servicecfg.yaml ")
		DefaultServiceConfigPath := util.ServicePath("service/util/servicecfg.yaml")
		if err := config.InitConfigFactory(DefaultServiceConfigPath); err != nil {
			return err
		}
	}

	service.setLogLevel()

	return nil
}

func (service *SERVICE) setLogLevel() {
	if config.ServiceConfig.Logger == nil {
		initLog.Warnln("SERVICE util without log level setting!!!")
		return
	}

	if config.ServiceConfig.Logger.Service != nil {
		if config.ServiceConfig.Logger.Service.DebugLevel != "" {
			if level, err := logrus.ParseLevel(config.ServiceConfig.Logger.Service.DebugLevel); err != nil {
				initLog.Warnf("SERVICE Log level [%s] is invalid, set to [info] level",
					config.ServiceConfig.Logger.Service.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				initLog.Infof("SERVICE Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			initLog.Infoln("SERVICE Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(config.ServiceConfig.Logger.Service.ReportCaller)
	}

}

func (service *SERVICE) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range service.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (service *SERVICE) Start() {

	initLog.Infoln("Server started")

	router := logger.NewGinWithLogrus(logger.GinLog)

	api.AddService(router)

	serviceLogPath := util.ServiceLogPath
	servicePemPath := util.ServicePemPath
	serviceKeyPath := util.ServiceKeyPath

	self := context.ServiceSelf()
	context.InitContext(self)

	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)
	initLog.Infof("SERVICE bind on : IP[%s] port : [%d]", self.BindingIPv4, self.SBIPort)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		service.Terminate()
		os.Exit(0)
	}()

	server, err := http2.NewServer(addr, serviceLogPath, router)
	if server == nil {
		initLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		initLog.Warnf("Initialize HTTP server: %+v", err)
	}

	serverScheme := config.ServiceConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(servicePemPath, serviceKeyPath)
	}

	if err != nil {
		initLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (service *SERVICE) Terminate() {
	logger.InitLog.Infof("SERVICE terminated")
}
