package datasources

import (
	"crypto/tls"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

func NewFastHTTPClient(insecureSkipVerify bool) (client *fasthttp.Client) {
	client = &fasthttp.Client{
		MaxConnsPerHost: (runtime.NumCPU() * 512) / 2,
		ReadTimeout:     time.Second * 45,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
		},
	}
	return
}
