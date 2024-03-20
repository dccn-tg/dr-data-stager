package handler

import (
	"testing"

	"github.com/Donders-Institute/dr-data-stager/internal/api-server/config"
)

func TestGetCollections(t *testing.T) {

	c := config.Configuration{
		RdrGateway: config.RdrGatewayConfig{
			ApiEndpoint: "https://dr-gateway.dccn.nl/v1",
		},
	}

	colls, err := getCollections(c, "dac", "3010000.01")

	if err != nil {
		t.Errorf("%s\n", err)
	}

	for _, coll := range colls {
		t.Logf("%+v\n", coll)
	}

}
