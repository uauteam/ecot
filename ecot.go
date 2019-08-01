package ecot

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	err2 "github.com/uauteam/ecot/err"
	"github.com/uauteam/ecot/log"
	"github.com/uauteam/ecot/repo"
	"os"
	"path/filepath"
	"strings"
)

type Ecot struct {
	*echo.Echo
}

func New() (e *Ecot) {
	e = &Ecot{
		echo.New(),
	}
	e.HideBanner = true

	return
}

// Register service with custom config
// first param configFuncHandler is a handler within returning	 a function accepting a Config parameter and returning a Config
// second param config is the param that configFuncHandler will invoke with
// the second param is a slice that only for variadic param, only the first item in the slice effects, other items will be discarded
func (e *Ecot) Register(configFuncHandler func(Config) func() Config, config ...Config) (err error) {
	cfg := Config{}
	if len(config) > 0 {
		cfg = config[0]
	}
	configFunc := configFuncHandler(cfg)
	c := configFunc()

	e.Logger.Infof("init %s database", c.Name)

	if strings.TrimSpace(c.DBDialect) == "" {
		return err2.DBDialectNotSet
	}

	if len(c.DBArgs) <= 0 {
		return err2.DBArgsNotSet
	}

	// init database
	if c.DBDialect == "sqlite3" {
		dbPath := fmt.Sprintf("%v", c.DBArgs[0])
		dbPath = filepath.Dir(dbPath)
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			err = os.MkdirAll(dbPath, 0755)
			if err != nil {
				return
			}
		}
	}

	var database *gorm.DB
	database, err = gorm.Open(c.DBDialect, c.DBArgs...)
	if err != nil {
		return
	}
	if c.LogLevel != log.OFF {
		database.LogMode(true)
	}

	if err := repo.RegisterDB("db_"+c.Name, database); err != nil {
		e.Logger.Printf(err.Error())
	}

	if c.AutoMigrateEntityRegister != nil {
		entities := c.AutoMigrateEntityRegister()
		database.AutoMigrate(entities...)
	}

	// register service api
	e.Logger.Printf("registering %s api", c.Name)
	if c.ApiRegister == nil {
		return err2.NoAPIRegistered
	}

	routeGroups := c.ApiRegister()
	for prefix, routeGroup := range routeGroups {
		g := e.Group(prefix)
		g.Use(routeGroup.MiddlewareFunc...)

		for _, route := range routeGroup.Routes {
			e.Logger.Printf("mapping %s %s%s", route.Method, prefix, route.Path)
			switch route.Method {
			case echo.POST:
				g.POST(route.Path, route.Handler, route.MiddlewareFunc...)
			case echo.GET:
				g.GET(route.Path, route.Handler, route.MiddlewareFunc...)
			case echo.PUT:
				g.PUT(route.Path, route.Handler, route.MiddlewareFunc...)
			case echo.DELETE:
				g.DELETE(route.Path, route.Handler, route.MiddlewareFunc...)
			case echo.PATCH:
				g.PATCH(route.Path, route.Handler, route.MiddlewareFunc...)
			default:
				e.Logger.Printf("HTTP method not supported: %s", route.Method)
			}
		}
	}

	return
}
