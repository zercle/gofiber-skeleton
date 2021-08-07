package apiv1

import (
	"github.com/valyala/fastjson"
	"github.com/zercle/gofiber-skelton/internal/datasources"
)

// RouterResources DB handler
type RouterResources struct {
	jsonParserPool *fastjson.ParserPool
	mainDB         *datasources.MariadbConfig
}

// InitRouterResources returns a new DBHandler
func InitRouterResources(mainDB *datasources.MariadbConfig) *RouterResources {
	return &RouterResources{
		jsonParserPool: new(fastjson.ParserPool),
		mainDB:         mainDB,
	}
}
