package utility

import (
	"context"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// OAuth2HttpClient returns a HTTP client wrapped with HTTP round tripper for Oauth token.
// It has lifetime within the `ctx`.
func OAuth2HttpClient(ctx context.Context, clientID, clientSecret, authURL string, scopes []string) (*http.Client, error) {

	cfg := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		TokenURL:     strings.Join([]string{authURL, "connect/token"}, "/"),
	}

	// custom http client for the context of retriving/refreshing oauth2 token
	ctxt := context.WithValue(
		context.Background(),
		oauth2.HTTPClient,
		NewHTTPSClient(5*time.Second, false),
	)

	ts := cfg.TokenSource(ctxt)

	return oauth2.NewClient(ctx, ts), nil
}

func GetOAuth2Token(clientID, clientSecret, authURL string, scopes []string) (*oauth2.Token, error) {

	cfg := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		TokenURL:     strings.Join([]string{authURL, "connect/token"}, "/"),
	}

	ctxt := context.WithValue(
		context.Background(),
		oauth2.HTTPClient,
		NewHTTPSClient(5*time.Second, false),
	)

	return cfg.Token(ctxt)
}
