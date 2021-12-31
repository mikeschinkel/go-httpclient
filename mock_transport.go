package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const ContentTypeHeader = "Content-Type"

// MockTransport provides a mock for RoundTrip() for testing
type MockTransport struct {
	ExpectedResponse *ExpectedResponse
}

// RoundTrip mocks the HTTP request-response process.
//goland:noinspection GoUnusedParameter
func (mt *MockTransport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	var body io.ReadCloser

	header := make(http.Header, 0)
	header.Add(ContentTypeHeader, mt.ExpectedResponse.ContentType)
	body, err = mt.LoadBody()
	if err != nil {
		goto end
	}

	res = &http.Response{
		StatusCode: mt.ExpectedResponse.StatusCode,
		Header:     header,
		// See https://gist.github.com/crgimenes/92d851b944ca2e459da7daa5c44801bf
		Body: body,
	}

end:
	return res, err
}

func (mt *MockTransport) LoadBody() (rc io.ReadCloser, err error) {
	var action string
	var body []byte

	fp := mt.ExpectedResponse.Filepath()
	r := struct {
		Body interface{} `json:"body"`
	}{}

	var b []byte
	b, err = ioutil.ReadFile(fp)
	if err != nil {
		action = "read"
		goto end
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		action = "load JSON"
		goto end
	}
	body, err = json.Marshal(r.Body)
	if err != nil {
		action = "marshal body to JSON"
		goto end
	}
	rc = ioutil.NopCloser(bytes.NewReader(body))
end:
	if err != nil {
		err = fmt.Errorf("unable to %s from %s: %s",
			action,
			mt.ExpectedResponse.Filepath(),
			err)

	}
	return rc, err
}
