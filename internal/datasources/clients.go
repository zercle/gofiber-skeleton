package datasources

import (
	"crypto/tls"
	"net/http"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	FasthttpClient *fasthttp.Client
	HttpClient     *http.Client
)

func InitFasthttpClient() (client *fasthttp.Client) {
	client = &fasthttp.Client{
		MaxConnsPerHost: (runtime.NumCPU() * 512) / 2,
		ReadTimeout:     time.Second * 45,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return
}

func InitHttpClient() (client *http.Client) {
	client = &http.Client{
		Timeout: time.Second * 45,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return
}
