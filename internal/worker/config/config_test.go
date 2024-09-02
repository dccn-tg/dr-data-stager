package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	cfg, err := LoadConfig(os.Getenv("TEST_CONFIG_FILE"))

	if err != nil {
		t.Errorf("%s\n", err)
	} else {
		t.Logf("%+v\n", cfg)
	}

}
