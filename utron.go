package utron

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var baseApp *App

func init() {
	baseApp = NewApp()
}

// App is the main utron application
type App struct {
	router     *Router
	cfg        *Config
	view       View
	log        Logger
	model      *Model
	configPath string
	isInit     bool
}

// NewApp creates a new bare-bone utron application. To use MVC components, you should call
// Init method before serving requests.
func NewApp() *App {
	app := &App{}
	app.Set(logThis)
	r := NewRouter(app)
	app.Set(r)
	app.Set(NewModel())
	return app
}

// NewMVC creates a new MVC utron app, if cfg is passed, it should be a directory to look for
// configuration file. The App returned is initialized.
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

// Init initializes MVC App
func (a *App) Init() error {
	if a.configPath == "" {
		a.SetConfigPath("config")
	}
	return a.init()
}

// SetConfigPath sets dir as a path to search for config files
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
	a.router.loadRoutes(a.configPath) // load routes file if any
	a.Set(appConfig)
	a.Set(views)
	a.isInit = true

	// Case the StaticDir is specified in the Config fille, register
	// a handler serving contents of the directory under the PathPrefix /static/
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

//getAbsolutePath returns absolute path to dir, if dir is relative then we add current working directory.
// Checks are made to ensure the directory exist.In case of any error, and empty string is returned.
func getAbsolutePath(dir string) (string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("untron: %s is not a directory", dir)
	}

	if filepath.IsAbs(dir) { // dir is already absolute, return it
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

// loadConfig loads configuration file, if cfg is provided then it is used as the directory
// for searching configuration file else defaults to directory named config in the current
// working directory.
func loadConfig(cfg ...string) (*Config, error) {
	cfgDir := "config"
	if len(cfg) > 0 {
		cfgDir = cfg[0]
	}

	// load configurations
	cfgFile, err := findConfigFile(cfgDir, "app")
	if err != nil {
		return nil, err
	}
	return NewConfig(cfgFile)
}

//findConfigFile finds the configuration file name, in the directory dir
func findConfigFile(dir string, name string) (file string, err error) {
	extensions := []string{".json", ".toml", ".yml"}

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

// AddController registers ctrl, and middlwares if provided
func (a *App) AddController(ctrl Controller, middlewares ...interface{}) {
	a.router.Add(ctrl, middlewares...)
}

// Set assigns value to *App components. The following can be set
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

// ServeHTTP serves http, it can be used with other http.Handler implementations
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

// SetConfigPath sets the path to look for configurations files in the
// global utron App.
func SetConfigPath(path string) {
	baseApp.SetConfigPath(path)
}

// RegisterModels registers models in the global utron App.
func RegisterModels(models ...interface{}) {
	baseApp.model.Register(models...)
}

// RegisterController register ctrl in the global utron App.
func RegisterController(ctrl Controller, middlewares ...interface{}) {
	baseApp.router.Add(ctrl, middlewares...)
}

// ServeHTTP serves request using global utron App
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
// By using this, you should make sure you followed MVC pattern,
func Run() {
	if err := baseApp.Init(); err != nil {
		logThis.Errors(err)
		os.Exit(1)
	}
	Migrate()
	port := baseApp.cfg.Port
	logThis.Info("starting server at ", baseApp.cfg.BaseURL)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), baseApp))
}
