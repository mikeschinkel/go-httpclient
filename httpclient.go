package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mikeschinkel/go-only"
	"net/http"
	"net/url"
)

// Client provides an instance of the RequestDoer interface
type Client struct {
	RequestDoer RequestDoer
}

func NewClient() *Client {
	return &Client{
		RequestDoer: &http.Client{
			Transport: http.DefaultTransport,
		},
	}
}

// SetTransport allows setting a RoundTripper
// as httpclient.Transport before GET or POST
func (c *Client) SetTransport(t http.RoundTripper) {
	c.RequestDoer.(*http.Client).Transport = t
}


// Get sends a post request to the URL with provided headers.
func (c *Client) Get(url string, headers http.Header) (resp *http.Response, req *http.Request, err error) {
	return c.do(http.MethodGet, url, headers, nil)
}

// Post sends a post request to the URL with the body with provided headers
func (c *Client) Post(url string, body interface{}, headers http.Header) (resp *http.Response, req *http.Request, err error) {
	var jsonBytes []byte
	for range only.Once {

		var scheme string
		scheme,err = getURLScheme(url)
		if err != nil {
			break
		}

		jsonBytes, err = json.Marshal(body)
		if err != nil {
 			err = fmt.Errorf("unable to marshal JSON body for %s POST request", scheme )
			break
		}

		resp, req, err = c.do(http.MethodPost, url, headers, bytes.NewReader(jsonBytes))

	}
	return resp, req,err
}

// do calls on http client to do an HTTP(S) request
func (c *Client) do(method, url string, headers http.Header, body interface{}) (resp *http.Response, req *http.Request, err error) {
	for range only.Once {

		scheme,err := getURLScheme(url)
		if err != nil {
			break
		}

		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			err = fmt.Errorf("unable to instantiate %s %s request",
				scheme,
				method,
			)
			break
		}

		req.Header = headers
		resp, err = c.RequestDoer.Do(req)
		if err != nil {
			err = fmt.Errorf("unable to perform an %s %s request",
				scheme,
				method,
			)
		}

	}
	return resp, req, err
}

// getURLScheme returns the scheme from a URL
func getURLScheme(u string) (s string, err error) {
	for range only.Once {
		uo, err := url.Parse(u)
		if err != nil {
			err = fmt.Errorf("unable to parse URL '%s': %s",
				u,
				err,
			)
		}
		s = uo.Scheme
	}
	return s,err
}
