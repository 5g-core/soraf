package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/logger"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/service"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/version"
)

var APPService = &service.SERVICE{}

var appLog *logrus.Entry

func init() {
	appLog = logger.AppLog
}

func main() {
	app := cli.NewApp()
	app.Name = "APP name"
	appLog.Infoln(app.Name)
	appLog.Infoln("APPService version: ", version.GetVersion())
	app.Usage = "- to do "
	app.Action = action
	app.Flags = APPService.GetCliCmd()
	if err := app.Run(os.Args); err != nil {
		appLog.Errorf("APPService Run error: %v", err)
	}
}

func action(c *cli.Context) error {
	if err := APPService.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("Failed to initialize !!")
	}

	APPService.Start()

	return nil
}
