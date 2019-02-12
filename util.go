package outlinesdk

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"regexp"
)

// CheckTLSCert checks the custom tls cert from shadowbox.
// These cert are indeed strange given that they have no SANs, neither IP or domain.
// Hence I can only implement an after connection check for it.
// However, it may means that a MITM attack would be able to get the access key
// since the access key is transferred to the server in the request
//
// TODO(michaellee8): Implement custom TLS transport to prevent this security loophole
func CheckTLSCert(con *tls.ConnectionState, fp []byte) bool {
	for _, cert := range con.PeerCertificates {
		cs := sha256.Sum256(cert.Raw)
		if bytes.Compare(cs[:], fp) == 0 {
			return true
		}
	}
	return false
}

// ParseAccessTxt parse the /opt/outline/access.txt
// so they can get the connection info without grabbing the JSON in
// installation script.
func ParseAccessTxt(accessTxt string) (apiURL, certSHA256 string, err error) {
	result := regexp.MustCompile(`certSha256:(.+)`).FindStringSubmatch(accessTxt)
	if result == nil {
		return "", "", errors.New("invalid certSha256")
	}
	certSHA256 = result[1]
	result = regexp.MustCompile(`apiURL:(.+)`).FindStringSubmatch(accessTxt)
	if result == nil {
		return "", "", errors.New("invalid apiURL")
	}
	apiURL = result[1]
	return apiURL, certSHA256, nil
}
