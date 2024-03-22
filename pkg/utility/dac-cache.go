package utility

import (
	"context"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	log "github.com/dccn-tg/tg-toolset-golang/pkg/logger"
)

type PpmFormConfig struct {
	AuthURL      string
	ClientID     string
	ClientSecret string
	ApiEndpoint  string
}

type PpmProjectCollection struct {
	CollID  string `json:"id"`
	OU      string `json:"organizationalUnit"`
	Project string `json:"project"`
}

var (
	// DAC cache with project number as key, collName as value
	cdacs map[string]string
	mu    sync.Mutex
)

// GetDacOfProject returns the corresponding DAC of the project `number` from
// the internal cache.
func GetDacOfProject(number string) string {
	defer mu.Unlock()
	mu.Lock()
	if val, ok := cdacs[number]; ok {
		return val
	}
	return ""
}

func SetDacCache(dacs map[string]string) {
	defer mu.Unlock()
	mu.Lock()
	cdacs = dacs
}

// UpdatePpmDacs runs an blocking loop within the lifetime of the context.
// For `every` duration in this loop, it updates the internal cache for
// DACs of the DCC OU registered in the online PPM form.
//
// A list of static DACs `sdacs` will take precedence over the DACs retrieved
// from the online PPM form.
func UpdatePpmDacs(
	ctx context.Context,
	cfg PpmFormConfig,
	every time.Duration,
	sdacs map[string]string,
) {

	// initialize dac cache with static list
	SetDacCache(sdacs)

	// initial timer to very short duration so that the fetch happens as soon as the call is made.
	timer := time.NewTimer(time.Nanosecond)

	for {
		select {
		case <-timer.C:
			// run update
			log.Debugf("retrieving project to DCC collection mapping")

			c, err := OAuth2HttpClient(
				ctx,
				cfg.ClientID,
				cfg.ClientSecret,
				cfg.AuthURL,
				[]string{
					"urn:dccn:project-proposal:collections",
				},
			)

			if err != nil {
				log.Errorf("%s\n", err)
				continue
			}

			// call ppmform api to get DCC collections
			resp, err := c.Get(
				fmt.Sprintf(
					"%s/api/RepositoryCollections?organizationalUnit=DCC",
					cfg.ApiEndpoint,
				),
			)

			if err != nil {
				log.Errorf("%s\n", err)
				continue
			}

			colls := []PpmProjectCollection{}
			err = UnmarshalFromResponseBody(resp, &colls)

			if err != nil {
				log.Errorf("%s\n", err)
				continue
			}

			// static dacs takes precedence
			newColls := make(map[string]string)
			for p, c := range sdacs {
				newColls[p] = c
			}

			// fill DACs retrieved from PPM form
			for _, coll := range colls {
				if _, ok := newColls[coll.Project]; !ok {
					// convert collID to collName
					name, err := collID2Name(coll.CollID)
					if err != nil {
						log.Errorf("%s\n", err)
					}
					newColls[coll.Project] = name
				}
			}

			// update dac cache with newColls
			SetDacCache(newColls)

			timer.Reset(every)
		case <-ctx.Done():
			log.Debugf("terminating UpdatePpmDacs loop")
			return
		}
	}
}

// collID2Name converts RDR collId to collName
func collID2Name(id string) (string, error) {
	// convert collID to collName
	parts := strings.Split(id, `.`)
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid collId: %s", id)
	}

	o := parts[0]
	ou := parts[1]
	ns := strings.Join(parts[2:], `.`)

	return path.Join(
		"/nl.ru.donders",
		o,
		ou,
		ns,
	), nil
}
