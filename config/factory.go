/*
 *  Configuration Factory
 */

package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/logger"
)

var ServiceConfig Config

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		logger.CfgLog.Infof("yaml content : %s", content)
		ServiceConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &ServiceConfig); yamlErr != nil {
			return yamlErr
		}
	}

	return nil
}
