package services

import (
	"github.com/valyala/fastjson"
	"github.com/zercle/gofiber-skelton/internal/datasources"
)

// ServiceResources DB handler
type ServiceResources struct {
	jsonParserPool *fastjson.ParserPool
	mainDB         *datasources.MariadbConfig
}

// InitServiceResources returns a new DBHandler
func InitServiceResources(mainDB *datasources.MariadbConfig) *ServiceResources {
	return &ServiceResources{
		jsonParserPool: new(fastjson.ParserPool),
		mainDB:         mainDB,
	}
}
