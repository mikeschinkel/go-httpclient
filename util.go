package httpclient

import "net/http"

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func HeaderEquals(a, b http.Header) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		println(i,v)
		//if v != b.Get(i) {
			return false
		//}
	}
	return true
}
