package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/Donders-Institute/dr-data-stager/pkg/utility"
)

// Configuration is the data structure for marshaling the
// config.yml file using the viper configuration framework.
type Configuration struct {
	Auth struct {
		Basic  map[string]string
		Oauth2 struct {
			JwksEndpoint     string
			UserInfoEndpoint string
		}
	}
	RdrGateway utility.RdrGatewayConfig
	PpmForm    utility.PpmFormConfig
	Dacs       map[string]string
}

// LoadConfig reads configuration file `cpath` and returns the
// `Configuration` data structure.
func LoadConfig(cpath string) (Configuration, error) {

	var conf Configuration

	// load configuration
	cfg, err := filepath.Abs(cpath)
	if err != nil {
		return conf, fmt.Errorf("cannot resolve config path: %s", cpath)
	}

	if _, err := os.Stat(cfg); err != nil {
		return conf, fmt.Errorf("cannot load config: %s", cfg)
	}

	// new viper with key delimiter set to something else than `.`
	// this allows to unmarshal key containing `.` in configuration file.
	v := viper.NewWithOptions(
		viper.KeyDelimiter(`#`),
	)

	v.SetConfigFile(cfg)
	if err := v.ReadInConfig(); err != nil {
		return conf, fmt.Errorf("error reading config file, %s", err)
	}

	err = v.Unmarshal(&conf)
	if err != nil {
		return conf, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return conf, nil
}
