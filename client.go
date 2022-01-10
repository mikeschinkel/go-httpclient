package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// DefaultTimeout is short so testing fails quickly
const DefaultTimeout = 3 // Seconds

// Client provides an instance of the RequestDoer interface
type Client struct {
	*http.Client
	Logger Logger
}

func NewClient() *Client {
	return NewClientWithTransport(http.DefaultTransport)
}

func NewClientWithTransport(rt http.RoundTripper) *Client {
	c := defaultHttpClient()
	c.Transport = rt
	return &Client{
		Logger: DefaultLogger(),
		Client: c,
	}
}

func NewClientWithLogger(l Logger) *Client {
	return &Client{
		Logger: l,
		Client: defaultHttpClient(),
	}
}

func defaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: DefaultTimeout * time.Second,
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
	logger := c.Logger
	youarehere := "[testclient.do()]"
	logger.Logf("DEBUG: %s Requesting %s %s", youarehere, method, url)

	scheme, err = getURLScheme(url)
	if err != nil {
		err = fmt.Errorf("unable to get scheme of URL; %w", err)
		goto end
	}
	logger.Logf("DEBUG: %s Scheme is %s", youarehere, scheme)

	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		err = fmt.Errorf("unable to instantiate HTTP request; %w", err)
		goto end
	}
	logger.Logf("DEBUG: %s Request is %#v", youarehere, req)

	req.Header = headers

	if len(headers) == 0 {
		logger.Log("DEBUG: %s No headers provided")
	} else {
		logger.Logf("DEBUG: %s Headers are %#v", headers)
	}
	req.Close = true // Avoid failed calls that send EOF
	resp, err = c.Do(req)
	if err != nil {
		err = fmt.Errorf("unable to perform HTTP request; %w", err)
	}
	logger.Logf("DEBUG: %s Response is %#v", youarehere, resp)

end:
	if err != nil {
		err = fmt.Errorf("%s unable to do %s %s; %w",
			youarehere, scheme, method, err)
		logger.Logf("ERROR: %s", err.Error())
	}
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
