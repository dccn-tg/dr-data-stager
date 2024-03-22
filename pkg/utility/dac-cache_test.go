package utility

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestGetDacOfProject(t *testing.T) {

	var (
		authURL      = os.Getenv("TEST_AUTH_URL")
		clientID     = os.Getenv("TEST_CLIENT_ID")
		clientSecret = os.Getenv("TEST_CLIENT_SECRET")
		apiEndpoint  = os.Getenv("TEST_PPM_API_ENDPOINT")
	)

	ctx, cancel := context.WithCancel(context.Background())

	cfg := PpmFormConfig{
		ApiEndpoint:  apiEndpoint,
		AuthURL:      authURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	sdacs := map[string]string{
		"3010000.01": "/nl.ru.donders/di/dccn/DAC_3010000.01_173",
	}

	t.Logf("run UpdatePpmDacs in go routine")

	go UpdatePpmDacs(ctx, cfg, 2*time.Second, sdacs)

	t.Logf("wait for 5 seconds ...")
	time.Sleep(5 * time.Second)

	// cancel the context to terminate the go routine
	t.Logf("stop go routine")
	cancel()

	for p, collNameExpected := range map[string]string{
		"3010000.01": sdacs["3010000.01"],
		"3017079.01": "/nl.ru.donders/di/dcc/DAC_2022.00149_912",
	} {
		collName := GetDacOfProject(p)
		if collName != collNameExpected {
			t.Errorf("%s: collName mismatch: %s != %s", p, collName, collNameExpected)
		} else {
			t.Logf("%s: %s == %s", p, collName, collNameExpected)
		}
	}
}
