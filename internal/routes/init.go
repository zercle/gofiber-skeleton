package routes

import (
	"github.com/zercle/gofiber-skelton/internal/datasources"
)

// RouterResources DB handler
type RouterResources struct {
	*datasources.Resources
}

// NewRouterResources returns a new DBHandler
func NewRouterResources(resources *datasources.Resources) *RouterResources {
	return &RouterResources{resources}
}
