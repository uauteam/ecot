package ecot


type Config struct {
	Name string // Use same name for same database of projects
	DBDialect string
	DBArgs []interface{}

	AutoMigrateEntityRegister func()[]interface{}

	ApiRegister func()map[string]RouteGroup
}