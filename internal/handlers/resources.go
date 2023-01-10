package handlers

import "github.com/golang-jwt/jwt/v4"

// RouterResources DB handler
type RouterResources struct {
	JwtKeyfunc jwt.Keyfunc
}

// NewRouterResources returns a new DBHandler
func NewRouterResources(jwtKeyfunc jwt.Keyfunc) *RouterResources {
	return &RouterResources{
		JwtKeyfunc: jwtKeyfunc,
	}
}
