package ecot

import "github.com/labstack/echo"

type Route struct {
	Method string
	Path string
	Handler func(ctx echo.Context)error
	MiddlewareFunc []echo.MiddlewareFunc
}

type RouteGroup struct {
	Routes []Route
	MiddlewareFunc []echo.MiddlewareFunc
}