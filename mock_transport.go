package httpclient

import (
	"net/http"
)


// MockTransport provides a mock for RoundTrip() for testing
type MockTransport struct{}

// RoundTrip mocks the HTTP request-response process.
func (t *MockTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	// TODO Expand this

	res = &http.Response{
		Header: req.Header,
	}
	return res,err
}
