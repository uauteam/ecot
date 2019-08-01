package ecot

import "github.com/uauteam/ecot/log"

type Config struct {
	Name string // Use same name for same database of projects
	DBDialect string
	DBArgs []interface{}

	AutoMigrateEntityRegister func()[]interface{}

	ApiRegister func()map[string]RouteGroup

	LogLevel log.Level
}