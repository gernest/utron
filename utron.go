package utron

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

var baseApp *App

func init() {
	baseApp = NewApp()
	if err := baseApp.Init(); err != nil {
		// TODO log this?
		log.Println(err)
	}
}

// App is the main utron application.
type App struct {
	router     *Router
	cfg        *Config
	view       View
	log        Logger
	model      *Model
	configPath string
	isInit     bool
}

// NewApp creates a new bare-bone utron application. To use the MVC components, you should call
// the Init method before serving requests.
func NewApp() *App {
	app := &App{}
	app.Set(logThis)
	r := NewRouter(app)
	app.Set(r)
	app.Set(NewModel())
	return app
}

// NewMVC creates a new MVC utron app. If cfg is passed, it should be a directory to look for
// the configuration files. The App returned is initialized.
func NewMVC(cfg ...string) (*App, error) {
	app := NewApp()
	if len(cfg) > 0 {
		app.SetConfigPath(cfg[0])
	}
	if err := app.Init(); err != nil {
		return nil, err
	}
	return app, nil
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

	views, err := NewSimpleView(appConfig.ViewsDir)
	if err != nil {
		return err
	}
	if a.model != nil && !a.model.IsOpen() {
		oerr := a.model.OpenWithConfig(appConfig)
		if oerr != nil {
			return oerr
		}
	} else {
		model, err := NewModelWithConfig(appConfig)
		if err != nil {
			return err
		}
		a.Set(model)
	}
	a.router.loadRoutes(a.configPath) // Load a routes file if available.
	a.Set(appConfig)
	a.Set(views)
	a.isInit = true

	// In case the StaticDir is specified in the Config file, register
	// a handler serving contents of that directory under the PathPrefix /static/.
	if appConfig.StaticDir != "" {
		static, err := getAbsolutePath(appConfig.StaticDir)
		if err != nil {
			logThis.Errors(err)
		}
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
func loadConfig(cfg ...string) (*Config, error) {
	cfgDir := "config"
	if len(cfg) > 0 {
		cfgDir = cfg[0]
	}

	// Load configurations.
	cfgFile, err := findConfigFile(cfgDir, "app")
	if err != nil {
		return nil, err
	}
	return NewConfig(cfgFile)
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
func (a *App) AddController(ctrlfn func() Controller, middlewares ...interface{}) {
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
	case Logger:
		a.log = value.(Logger)
	case *Router:
		a.router = value.(*Router)
	case View:
		a.view = value.(View)
	case *Config:
		a.cfg = value.(*Config)
	case *Model:
		a.model = value.(*Model)
	}
}

// ServeHTTP serves http requests. It can be used with other http.Handler implementations.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// SetConfigPath sets the path to look for the configuration files in the
// global utron App.
func SetConfigPath(path string) {
	baseApp.SetConfigPath(path)
}

// RegisterModels registers models in the global utron App.
func RegisterModels(models ...interface{}) {
	_ = baseApp.model.Register(models...)
}

// RegisterController registers a controller in the global utron App.
func RegisterController(ctrl Controller, middlewares ...interface{}) {
	_ = baseApp.router.Add(GetCtrlFunc(ctrl), middlewares...)
}

// GetCtrlFunc returns a new copy of the contoller everytime the function is called
func GetCtrlFunc(ctrl Controller) func() Controller {
	v := reflect.ValueOf(ctrl)
	return func() Controller {
		e := v
		if e.Kind() == reflect.Ptr {
			e = e.Elem()
			return e.Addr().Interface().(Controller)
		}
		return e.Interface().(Controller)
	}
}

// ServeHTTP serves request using global utron App.
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !baseApp.isInit {
		if err := baseApp.Init(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	baseApp.ServeHTTP(w, r)
}

// Migrate runs migrations on the global utron app.
func Migrate() {
	baseApp.model.AutoMigrateAll()
}

// Run runs a http server, serving the global utron App.
//
// By using this, you should make sure you followed the MVC pattern.
func Run() {
	if baseApp.cfg.Automigrate {
		Migrate()
	}
	port := baseApp.cfg.Port
	logThis.Info("starting server at ", baseApp.cfg.BaseURL)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), baseApp))
}
