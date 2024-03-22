package utility

import (
	"testing"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

func init() {
	// setup logger
	log.NewLogger(
		log.Configuration{
			EnableConsole:     true,
			ConsoleJSONFormat: false,
			ConsoleLevel:      log.Debug,
		},
		log.InstanceLogrusLogger,
	)
}

func TestGetCollections(t *testing.T) {

	c := RdrGatewayConfig{
		ApiEndpoint: "https://dr-gateway.dccn.nl/v1",
	}

	colls, err := GetCollections(c, "dac", "3010000.01")

	if err != nil {
		t.Errorf("%s\n", err)
	}

	for _, coll := range colls {
		t.Logf("%+v\n", coll)
	}

}
