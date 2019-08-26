package ecot

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	err2 "github.com/uauteam/ecot/err"
	"github.com/uauteam/ecot/repo"
	"os"
	"path/filepath"
	"strings"
)

type Ecot struct {
	*echo.Echo
	routeGroup map[string]*echo.Group
}

type EcotValidator struct {
	validator *validator.Validate
}

func (ev *EcotValidator) Validate(i interface{}) error {
	return ev.validator.Struct(i)
}

func New() (e *Ecot) {
	e = &Ecot{
		echo.New(),
		make(map[string]*echo.Group),
	}
	e.HideBanner = true

	e.Validator = &EcotValidator{validator:validator.New()}

	return
}

// Register service with custom config
// first param configFuncHandler is a handler within returning	 a function accepting a Config parameter and returning a Config
// second param config is the param that configFuncHandler will invoke with
// the second param is a slice that only for variadic param, only the first item in the slice effects, other items will be discarded
func (ecot *Ecot) Register(configFuncHandler func(Config) func() Config, config ...Config) (err error) {
	cfg := Config{}
	if len(config) > 0 {
		cfg = config[0]
	}
	configFunc := configFuncHandler(cfg)
	c := configFunc()

	ecot.Logger.Infof("init %s database", c.Name)

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
			e := os.MkdirAll(dbPath, 0755)
			if e != nil {
				return e
			}
		}
	}

	database, err := gorm.Open(c.DBDialect, c.DBArgs...)
	if err != nil {
		return
	}

	if ecot.Logger.Level() != log.OFF {
		database.LogMode(true)
	}

	if e := repo.RegisterDB("db_"+c.Name, database); e != nil {
		ecot.Logger.Printf(e.Error())
	}

	if c.AutoMigrateEntityRegister != nil {
		entities := c.AutoMigrateEntityRegister()
		database.AutoMigrate(entities...)
	}

	// register service api
	ecot.Logger.Printf("registering %s api", c.Name)
	if c.ApiRegister == nil {
		return err2.NoAPIRegistered
	}

	routeGroups := c.ApiRegister()
	for prefix, routeGroup := range routeGroups {
		for _, route := range routeGroup.Routes {
			routeGroupPrefix := route.Version + prefix
			if route.Version != "" {
				routeGroupPrefix = "/" + routeGroupPrefix
			}

			if g, ok := ecot.routeGroup[routeGroupPrefix]; !ok {
				g = ecot.Group(routeGroupPrefix, routeGroup.MiddlewareFunc...)
				ecot.routeGroup[routeGroupPrefix] = g
			}

			g := ecot.routeGroup[routeGroupPrefix]

			ecot.Logger.Printf("mapping %s %s%s", route.Method, routeGroupPrefix, route.Path)
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
				ecot.Logger.Printf("HTTP method not supported: %s", route.Method)
			}
		}
	}

	return
}
