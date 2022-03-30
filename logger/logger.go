package logger

import (
	"os"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log         *logrus.Logger
	AppLog      *logrus.Entry
	InitLog     *logrus.Entry
	CfgLog      *logrus.Entry
	HandlerLog  *logrus.Entry
	UtilLog     *logrus.Entry
	HttpLog     *logrus.Entry
	ConsumerLog *logrus.Entry
	GinLog      *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	ServiceLogHook, err := NewFileHook(ServiceLogfle, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err == nil {
		log.Hooks.Add(ServiceLogHook)
	}
	// to do  different component different logger
	AppLog = log.WithFields(logrus.Fields{"component": "SERVICE", "category": "App"})
	InitLog = log.WithFields(logrus.Fields{"component": "INIT", "category": "Init"})
	CfgLog = log.WithFields(logrus.Fields{"component": "CONFIG", "category": "CFG"})
	HandlerLog = log.WithFields(logrus.Fields{"component": "HANDLER", "category": "HDLR"})
	UtilLog = log.WithFields(logrus.Fields{"component": "UTIL", "category": "Util"})
	HttpLog = log.WithFields(logrus.Fields{"component": "HTTP", "category": "HTTP"})
	ConsumerLog = log.WithFields(logrus.Fields{"component": "CONSUMER", "category": "Consumer"})
	GinLog = log.WithFields(logrus.Fields{"component": "GIN", "category": "GIN"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(set bool) {
	log.SetReportCaller(set)
}
