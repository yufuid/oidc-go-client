package utils

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

var (
	once  sync.Once
	httpc *http.Client
)

func NewHTTPClient() *http.Client {
	once.Do(func() {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		httpc = &http.Client{
			Transport: tr,
			Timeout:   10 * time.Second,
		}
	})
	return httpc
}
