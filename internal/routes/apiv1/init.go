package apiv1

import (
	"regexp"

	"github.com/valyala/fastjson"
	"github.com/zercle/gofiber-skelton/internal/datasources"
)

// RouterResources DB handler
type RouterResources struct {
	regxNum        *regexp.Regexp
	jsonParserPool *fastjson.ParserPool
	mainDB         *datasources.MariadbConfig
}

// InitRouterResources returns a new DBHandler
func InitRouterResources(mainDB *datasources.MariadbConfig) *RouterResources {
	return &RouterResources{
		regxNum:        regexp.MustCompile(`[0-9]+`),
		jsonParserPool: new(fastjson.ParserPool),
		mainDB:         mainDB,
	}
}
