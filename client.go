package outlinesdk

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Implements Outline Server APIs listed in
// https://rebilly.github.io/ReDoc/?url=https://raw.githubusercontent.com/Jigsaw-Code/outline-server/master/src/shadowbox/server/api.yml

// A Client for communicating with the shadowbox (Outline Server)
type Client struct {
	APIURL     *url.URL
	CertSHA256 []byte
}

// NewClient creates a new Outline SDK Client with apiUrl and certSha256 as provided by /opt/outline/access.txt
func NewClient(apiURL string, certSha256 string) (*Client, error) {
	if apiURL[len(apiURL)-1:] != "/" {
		apiURL += "/"
	}
	cert, err := hex.DecodeString(certSha256)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		APIURL:     u,
		CertSHA256: cert,
	}, nil
}

// Custom http.Client used in this code
// Skip TLS verification to cope with shadowbox's special HTTPS cert
// TODO(michaellee8): Implement custom TLS transport to prevent the security loophole mentioned in util.go
var httpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
	Timeout: time.Second * 10,
}

// GetServerInfo retrieves information about the server with /access.
// Returns ServerInfo if success and an error if failed
func (c *Client) GetServerInfo() (*ServerInfo, error) {
	req, err := c.newReq(http.MethodGet, "server", nil)
	if err != nil {
		return nil, err
	}
	var serverInfo ServerInfo
	err = c.do(req, http.StatusOK, &serverInfo)
	if err != nil {
		return nil, err
	}
	return &serverInfo, nil
}

// RenameServer renames the server with /name.
// Returns nil if success and an error if failed
func (c *Client) RenameServer(name string) (err error) {
	reqQbj := NameType{Name: name}
	req, err := c.newReq(http.MethodPut, "name", reqQbj)
	if err != nil {
		return err
	}

	err = c.do(req, http.StatusNoContent, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetMetricsSetting get metrics sharing settings of the server with /metrics/enabled
func (c *Client) GetMetricsSetting() (*bool, error) {
	req, err := c.newReq(http.MethodGet, "metrics/enabled", nil)
	if err != nil {
		return nil, err
	}
	resObj := &MetricsSetting{}
	err = c.do(req, http.StatusOK, resObj)
	if err != nil {
		return nil, err
	}
	return &resObj.MetricsEnabled, nil
}

// SetMetricsSetting set metric sharing settings of the server with /metrics/enabled
func (c *Client) SetMetricsSetting(opt bool) error {
	req, err := c.newReq(http.MethodPut, "metrics/enabled", &MetricsSetting{MetricsEnabled: opt})
	if err != nil {
		return err
	}
	err = c.do(req, http.StatusNoContent, nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateAccessKey creates a new access key entry
func (c *Client) CreateAccessKey() (*AccessKey, error) {
	req, err := c.newReq(http.MethodPost, "access-keys", nil)
	if err != nil {
		return nil, err
	}
	resObj := &AccessKey{}
	err = c.do(req, http.StatusCreated, resObj)
	if err != nil {
		return nil, err
	}
	return resObj, nil
}

// GetAccessKeys retrieve the list of access keys on the server
func (c *Client) GetAccessKeys() (*AccessKeyList, error) {
	req, err := c.newReq(http.MethodGet, "access-keys", nil)
	if err != nil {
		return nil, err
	}
	resObj := &getAccessKeysResponse{}
	err = c.do(req, http.StatusOK, resObj)
	if err != nil {
		return nil, err
	}
	return &resObj.AccessKeys, nil

}

// DeleteAccessKey deletes the access key with given id
func (c *Client) DeleteAccessKey(id string) error {
	req, err := c.newReq(http.MethodDelete, fmt.Sprintf("access-keys/%s", url.PathEscape(id)), nil)
	if err != nil {
		return err
	}
	err = c.do(req, http.StatusNoContent, nil)
	if err != nil {
		return err
	}
	return nil
}

// RenameAccessKey rename the access key with given id to given name
func (c *Client) RenameAccessKey(id, name string) error {
	req, err := c.newReq(http.MethodPut, fmt.Sprintf("access-keys/%s/name", url.PathEscape(id)), NameType{Name: name})
	if err != nil {
		return err
	}
	err = c.do(req, http.StatusNoContent, nil)
	if err != nil {
		return err
	}
	return nil
}

// GetUsageMetrics get the usage statistics of different account in the server
func (c *Client) GetUsageMetrics() (*map[string]int64, error) {
	req, err := c.newReq(http.MethodGet, "metrics/transfer", nil)
	if err != nil {
		return nil, err
	}
	resObj := &UsageInfo{}
	err = c.do(req, http.StatusOK, resObj)
	if err != nil {
		return nil, err
	}
	return &resObj.BytesTransferredByUserID, err
}
