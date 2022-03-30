package logger

type Logger struct {
	Service *LogSetting `yaml:"SERVICE"`
}

type LogSetting struct {
	DebugLevel   string `yaml:"debugLevel"`
	ReportCaller bool   `yaml:"ReportCaller"`
}
