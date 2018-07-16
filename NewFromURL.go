package paperfishGo

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

//func NewFromURL(url string, client *http.Client) (map[string]Operation, error) {
func NewFromURL(uri string, client *http.Client) ([]WSClientT, error) {
	var err error
	var resp *http.Response

	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
					RootCAs:            x509.NewCertPool(),
					Certificates:       []tls.Certificate{tls.Certificate{}},
				},
				DisableCompression: true,
			},
		}
	}

	resp, err = client.Get(uri)
	if err != nil {
		Goose.New.Logf(1, "%s (%s)", ErrFetchingContract, err)
		return nil, err
	}

	return NewFromReader(resp.Body, client)
}
