package ecot

import "github.com/labstack/echo"

type Route struct {
	Method string // HTTP method
	Path string // route path
	Handler func(ctx echo.Context)error
	MiddlewareFunc []echo.MiddlewareFunc
	Version string // API version e.g. v1
}

type RouteGroup struct {
	Routes []Route
	MiddlewareFunc []echo.MiddlewareFunc
}