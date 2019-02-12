package outlinesdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Helper method that returns a http.Request with optional JSON body
func (c *Client) newReq(method, path string, obj interface{}) (*http.Request, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	realPath := c.APIURL.ResolveReference(u)
	if obj == nil {
		return http.NewRequest(method, realPath.String(), nil)
	}
	body, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, realPath.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Helper method to run a http.Request
func (c *Client) do(req *http.Request, supposedStatusCode int, obj interface{}) error {
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	if !CheckTLSCert(res.TLS, c.CertSHA256) {
		return errors.New("invalid https certificate")
	}
	if supposedStatusCode != res.StatusCode {
		return errors.New("http status code not match")
	}

	// No decode if not required
	if obj == nil {
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(obj)

	if err != nil {
		return err
	}
	return nil
}
