package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client provides an instance of the RequestDoer interface
type Client struct {
	*http.Client
}

func NewClient() *Client {
	return NewClientWithTransport(http.DefaultTransport)
}

func NewClientWithTransport(rt http.RoundTripper) *Client {
	return &Client{
		Client: &http.Client{
			Transport: rt,
			Timeout:   3 * time.Second, // Most testing should fail quickly.
		},
	}
}

// GET requests an HTTP(S) GET of the URL with provided headers.
func (c *Client) GET(url string, headers http.Header) (resp *http.Response, err error) {
	return c.do(http.MethodGet, url, headers, nil)
}

// PUT requests an HTTP(S) PUT of the URL with the body with provided headers
func (c *Client) PUT(url string, body interface{}, headers http.Header) (resp *http.Response, err error) {
	return c.requestWithBody(http.MethodPut, url, body, headers)
}

// POST requests an HTTP(S) POST of the URL with the body and provided headers
func (c *Client) POST(url string, body interface{}, headers http.Header) (resp *http.Response, err error) {
	return c.requestWithBody(http.MethodPost, url, body, headers)
}

// DELETE requests an HTTP(S) DELETE of the URL with provided headers.
func (c *Client) DELETE(url string, headers http.Header) (resp *http.Response, err error) {
	return c.do(http.MethodDelete, url, headers, nil)
}

// requestWithBody requests via HTTP(S) of the URL with the method, body and provided headers
func (c *Client) requestWithBody(method, url string, body interface{}, headers http.Header) (resp *http.Response, err error) {
	var jsonBytes []byte

	jsonBytes, err = json.Marshal(body)
	if err != nil {
		err = fmt.Errorf("unable to marshal JSON body for POST request of %s", url)
		goto end
	}

	resp, err = c.do(method, url, headers, bytes.NewReader(jsonBytes))

end:
	return resp, err
}

// do call a http client to do an HTTP(S) request
//goland:noinspection GoUnusedParameter
func (c *Client) do(method, url string, headers http.Header, body interface{}) (resp *http.Response, err error) {
	var scheme string
	var req *http.Request

	scheme, err = getURLScheme(url)
	if err != nil {
		err = fmt.Errorf("unable to get scheme of URL %s on %s; %w",
			url,
			method,
			err)
		goto end
	}

	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		err = fmt.Errorf("unable to instantiate %s %s request; %w",
			scheme,
			method,
			err)
		goto end
	}

	req.Close = true // Avoid failed calls that send EOF
	req.Header = headers
	resp, err = c.Do(req)
	if err != nil {
		err = fmt.Errorf("unable to perform %s %s request: %s",
			scheme,
			method,
			err,
		)
	}

end:
	return resp, err
}

// getURLScheme returns the scheme from a URL
func getURLScheme(u string) (s string, err error) {
	uo, err := url.Parse(u)
	if err != nil {
		err = fmt.Errorf("unable to parse URL '%s': %s",
			u,
			err,
		)
		goto end
	}
	s = uo.Scheme
end:
	return s, err
}
