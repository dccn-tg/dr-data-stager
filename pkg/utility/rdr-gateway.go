package utility

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	crdr "github.com/dccn-tg/dr-gateway/pkg/swagger/client/client"
	ordr "github.com/dccn-tg/dr-gateway/pkg/swagger/client/client/operations"
	httptransport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/strfmt"
)

type RdrGatewayConfig struct {
	ApiEndpoint string
}

// GetCollections return a list of collName of the RDR collections with type `ctype`
// and associated with project number `project`.
func GetCollections(cfg RdrGatewayConfig, ctype, project string) ([]string, error) {

	lctype := strings.ToLower(ctype)

	apiURL, err := url.Parse(cfg.ApiEndpoint)

	if err != nil {
		return nil, err
	}

	c := crdr.New(
		httptransport.New(
			apiURL.Host,
			apiURL.Path,
			[]string{apiURL.Scheme},
		),
		strfmt.Default,
	)

	params := ordr.GetCollectionsProjectIDParams{
		ID: project,
	}
	params.WithTimeout(10 * time.Second)

	rslt, err := c.Operations.GetCollectionsProjectID(&params)

	if err != nil {
		return nil, err
	}

	if !rslt.IsSuccess() {
		return nil, fmt.Errorf("%s (%d)", rslt.Error(), rslt.Code())
	}

	colls := []string{}
	for _, m := range rslt.GetPayload().Collections {
		if strings.ToLower(string(*m.Type)) == lctype {
			colls = append(
				colls,
				*m.Path,
			)
		}
	}

	return colls, nil
}
