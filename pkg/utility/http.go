package utility

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"
)

// NewHTTPSClient initiates a new HTTPS client with timeout and skip certificate
// verification when `insecure` flag is set to `true`.
func NewHTTPSClient(timeout time.Duration, insecure bool) (client *http.Client) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
	}

	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client = &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	return
}

// UnmarshalFromResponseBody unmarshals data from the `response.Body` stream
// into a struct `i`.  The `response.Body` is also closed after calling this
// function.
func UnmarshalFromResponseBody(response *http.Response, i interface{}) error {

	defer response.Body.Close()

	d, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	return json.Unmarshal(d, i)
}
