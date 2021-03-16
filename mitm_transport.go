package httpclient

import (
	"github.com/google/uuid"
	"net/http"
)

func NewMITMTransport() *MITMTransport {
	return &MITMTransport{
		sessionId: uuid.New(),
	}
}

// MockTransport provides a mock for RoundTrip() for testing
type MITMTransport struct{
	counter int
	sessionId uuid.UUID
}

// RoundTrip mocks the HTTP request-response process.
func (t *MITMTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	req.Header = http.Header{}
	req.Header.Add("X-TestSession","ABC")
	req.Header.Add("X-TestSession-Counter","1")
	return http.DefaultTransport.RoundTrip(req)
}


