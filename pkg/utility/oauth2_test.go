package utility

import (
	"os"
	"testing"
)

func TestGetOAuth2Token(t *testing.T) {

	var (
		authURL      = os.Getenv("TEST_AUTH_URL")
		clientID     = os.Getenv("TEST_CLIENT_ID")
		clientSecret = os.Getenv("TEST_CLIENT_SECRET")
		scopes       = []string{"urn:dccn:project-proposal:collections"}
	)

	t1, err := GetOAuth2Token(clientID, clientSecret, authURL, scopes)

	if err != nil {
		t.Errorf("%s", err)
	}

	t.Logf("token1: %s %s\n", t1.AccessToken, t1.Expiry)

	t2, err := GetOAuth2Token(clientID, clientSecret, authURL, scopes)

	if err != nil {
		t.Errorf("%s", err)
	}

	t.Logf("token2: %s %s\n", t2.AccessToken, t2.Expiry)

}
