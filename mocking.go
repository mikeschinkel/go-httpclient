package httpclient

import (
	"net/http"
)


// DoFuncType defines the type of func used
type DoFuncType = func(req *http.Request) (*http.Response, error)

// GetDoFunc fetches the mock client's `Do` func
var GetDoFunc DoFuncType

// MockClient is the mock client
type MockClient struct {
	DoFunc DoFuncType
}

// Do is the mock client's `Do` func
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}


// RequestDoer interface
type RequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}
