package httpclient

import (
	"github.com/mikeschinkel/go-only"
	"testing"
)

const URLToTest = "https://api.github.com/users/google/repos"

func TestGet(t *testing.T) {
	for range only.Once {
		client := NewClient()
		client.SetTransport(&MockTransport{})
		got,want,err := client.Get(URLToTest,nil)
		if err != nil {
			t.Fatalf("Failed to retrieve %s: %s",
				URLToTest,
				err)
		}
		if HeaderEquals(got.Header, want.Header) {
			t.Fatalf("Failed to retrieve webpage for google.com")
		}

	}
}










