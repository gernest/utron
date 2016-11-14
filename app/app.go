package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gernest/utron/config"
	"github.com/gernest/utron/controller"
	"github.com/gernest/utron/logger"
	"github.com/gernest/utron/models"
	"github.com/gernest/utron/router"
	"github.com/gernest/utron/view"
)

// App is the main utron application.
type App struct {
	router     *router.Router
	cfg        *config.Config
	view       view.View
	log        logger.Logger
	model      *models.Model
	configPath string
	isInit     bool
}

func (a *App) options() *router.Options {
	return &router.Options{
		Model:  a.model,
		View:   a.view,
		Config: a.cfg,
	}
}

// Init initializes the MVC App.
func (a *App) Init() error {
	if a.configPath == "" {
		a.SetConfigPath("config")
	}
	return a.init()
}

// SetConfigPath sets the directory path to search for the config files.
func (a *App) SetConfigPath(dir string) {
	a.configPath = dir
}

// init initializes values to the app components.
func (a *App) init() error {
	appConfig, err := loadConfig(a.configPath)
	if err != nil {
		return err
	}

	views, err := view.NewSimpleView(appConfig.ViewsDir)
	if err != nil {
		return err
	}
	if a.model != nil && !a.model.IsOpen() {
		oerr := a.model.OpenWithConfig(appConfig)
		if oerr != nil {
			return oerr
		}
	} else {
		model, err := models.NewModelWithConfig(appConfig)
		if err != nil {
			return err
		}
		a.Set(model)
	}
	a.router.Options = a.options()
	a.router.LoadRoutes(a.configPath) // Load a routes file if available.
	a.Set(appConfig)
	a.Set(views)
	a.isInit = true

	// In case the StaticDir is specified in the Config file, register
	// a handler serving contents of that directory under the PathPrefix /static/.
	if appConfig.StaticDir != "" {
		static, _ := getAbsolutePath(appConfig.StaticDir)
		if static != "" {
			a.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))
		}

	}
	return nil
}

// getAbsolutePath returns the absolute path to dir. If the dir is relative, then we add
// the current working directory. Checks are made to ensure the directory exist.
// In case of any error, an empty string is returned.
func getAbsolutePath(dir string) (string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("untron: %s is not a directory", dir)
	}

	if filepath.IsAbs(dir) { // If dir is already absolute, return it.
		return dir, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	absDir := filepath.Join(wd, dir)
	_, err = os.Stat(absDir)
	if err != nil {
		return "", err
	}
	return absDir, nil
}

// loadConfig loads the configuration file. If cfg is provided, then it is used as the directory
// for searching the configuration files. It defaults to the directory named config in the current
// working directory.
func loadConfig(cfg ...string) (*config.Config, error) {
	cfgDir := "config"
	if len(cfg) > 0 {
		cfgDir = cfg[0]
	}

	// Load configurations.
	cfgFile, err := findConfigFile(cfgDir, "app")
	if err != nil {
		return nil, err
	}
	return config.NewConfig(cfgFile)
}

// findConfigFile finds the configuration file name in the directory dir.
func findConfigFile(dir string, name string) (file string, err error) {
	extensions := []string{".json", ".toml", ".yml", ".hcl"}

	for _, ext := range extensions {
		file = filepath.Join(dir, name)
		if info, serr := os.Stat(file); serr == nil && !info.IsDir() {
			return
		}
		file = file + ext
		if info, serr := os.Stat(file); serr == nil && !info.IsDir() {
			return
		}
	}
	return "", fmt.Errorf("utron: can't find configuration file %s in %s", name, dir)
}

// AddController registers a controller, and middlewares if any is provided.
func (a *App) AddController(ctrlfn func() controller.Controller, middlewares ...interface{}) {
	_ = a.router.Add(ctrlfn, middlewares...)
}

// Set is for assigning a value to *App components. The following can be set:
//	Logger by passing Logger
//	View by passing View
//	Router by passing *Router
//	Config by passing *Config
//	Model by passing *Model
func (a *App) Set(value interface{}) {
	switch value.(type) {
	case logger.Logger:
		a.log = value.(logger.Logger)
	case *router.Router:
		a.router = value.(*router.Router)
	case view.View:
		a.view = value.(view.View)
	case *config.Config:
		a.cfg = value.(*config.Config)
	case *models.Model:
		a.model = value.(*models.Model)
	}
}

// ServeHTTP serves http requests. It can be used with other http.Handler implementations.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
