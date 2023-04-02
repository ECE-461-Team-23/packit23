package github_apis

import "net/http"

type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func NewHTTPClientWithHeaders(baseRoundTripper http.RoundTripper, headers map[string]string) *http.Client {
	if baseRoundTripper == nil {
		baseRoundTripper = http.DefaultTransport
	}

	return &http.Client{
		Transport: &headerTransport{
			base:    baseRoundTripper,
			headers: headers,
		},
	}
}

func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := CloneRequest(req)
	for key, val := range h.headers {
		req2.Header.Set(key, val)
	}
	return h.base.RoundTrip(req2)
}

// CloneRequest and CloneHeader copied from https://github.com/kubernetes/apimachinery/blob/master/pkg/util/net/http.go#L424

// CloneRequest creates a shallow copy of the request along with a deep copy of the Headers.
func CloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)

	// shallow clone
	*r = *req

	// deep copy headers
	r.Header = CloneHeader(req.Header)

	return r
}

// CloneHeader creates a deep copy of an http.Header.
func CloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}
