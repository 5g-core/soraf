/*
 *  Configuration
 */

package config

import (
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/logger"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *logger.Logger `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version,omitempty"`
	Description string `yaml:"description,omitempty"`
}

const (
	SORAF_DEFAULT_IPV4     = "127.0.0.4"
	SORAF_DEFAULT_PORT     = "8000"
	SORAF_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	Sbi *Sbi `yaml:"sbi"`
}

type Sbi struct {
	Scheme      string `yaml:"scheme"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty"` // IP used to run the server in the node.
	Port        int    `yaml:"port"`
	Tls         *Tls   `yaml:"tls,omitempty"`
}

type Tls struct {
	Log string `yaml:"log"`
	Pem string `yaml:"pem"`
	Key string `yaml:"key"`
}
