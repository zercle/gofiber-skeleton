package datasources

import (
	"crypto/tls"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"gorm.io/gorm"
)

var (
	FastHttpClient *fasthttp.Client
	JsonParserPool *fastjson.ParserPool
)

type Resources struct {
	MainDbConn *gorm.DB
}

func InitResources(mainDbConn *gorm.DB) *Resources {
	return &Resources{
		MainDbConn: mainDbConn,
	}
}

func InitFasthttpClient() (client *fasthttp.Client) {
	client = &fasthttp.Client{
		MaxConnsPerHost: (runtime.NumCPU() * 512) / 2,
		ReadTimeout:     time.Second * 45,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	return
}

func InitJsonParserPool() (jsonParserPool *fastjson.ParserPool) {
	jsonParserPool = new(fastjson.ParserPool)
	return
}
