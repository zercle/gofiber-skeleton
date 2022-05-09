package routes

import (
	"github.com/valyala/fastjson"
	"gorm.io/gorm"
)

// RouterResources DB handler
type RouterResources struct {
	jsonParserPool *fastjson.ParserPool
	mainDB         *gorm.DB
}

// InitRouterResources returns a new DBHandler
func InitRouterResources(mainDB *gorm.DB) *RouterResources {
	return &RouterResources{
		jsonParserPool: new(fastjson.ParserPool),
		mainDB:         mainDB,
	}
}
