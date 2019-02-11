package outlinesdk

import (
	"crypto/tls"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"
)

// A Client for communicating with the shadowbox (Outline Server)
type Client struct {
	ApiUrl     *url.URL
	CertSha256 []byte
}

// Create a new Outline SDK Client with apiUrl and certSha256 as provided by /opt/outline/access.txt
func NewClient(apiUrl string, certSha256 string) (*Client, error) {
	cert, err := hex.DecodeString(certSha256)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
	}
	return &Client{
		ApiUrl:     u,
		CertSha256: cert,
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

// Retrieve information about the server with /access
// Returns ServerInfo if success and an error if failed
func (c *Client) GetServerInfo() (*ServerInfo, error) {
	req, err := c.newReq(http.MethodGet, "/server", nil)
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

// Renames the server with /name
// Returns nil if success and an error if failed
func (c *Client) RenameServer(name string) (err error) {
	reqQbj := NameType{Name: name}
	req, err := c.newReq(http.MethodPut, "/name", reqQbj)
	if err != nil {
		return err
	}

	err = c.do(req, http.StatusNoContent, nil)
	if err != nil {
		return err
	}
	return nil
}
