package app

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gernest/qlstore"
	"github.com/gernest/utron/config"
	"github.com/gernest/utron/controller"
	"github.com/gernest/utron/logger"
	"github.com/gernest/utron/models"
	"github.com/gernest/utron/router"
	"github.com/gernest/utron/view"
	"github.com/gorilla/sessions"
	// load ql drier
	_ "github.com/cznic/ql/driver"
)

//StaticServerFunc is a function that returns the static assetsfiles server.
//
// The first argument retrued is the path prefix for the static assets. If strp
// is set to true then the prefix is going to be stripped.
type StaticServerFunc func(*config.Config) (prefix string, strip bool, h http.Handler)

// App is the main utron application.
type App struct {
	Router       *router.Router
	Config       *config.Config
	View         view.View
	Log          logger.Logger
	Model        *models.Model
	ConfigPath   string
	StaticServer StaticServerFunc
	SessionStore sessions.Store
	isInit       bool
}

// NewApp creates a new bare-bone utron application. To use the MVC components, you should call
// the Init method before serving requests.
func NewApp() *App {
	return &App{
		Log:    logger.NewDefaultLogger(os.Stdout),
		Router: router.NewRouter(),
		Model:  models.NewModel(),
	}
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

//StaticServer implements StaticServerFunc.
//
// This uses the http.Fileserver to handle static assets. The routes prefixed
// with /static/ are static asset routes by default.
func StaticServer(cfg *config.Config) (string, bool, http.Handler) {
	static, _ := getAbsolutePath(cfg.StaticDir)
	if static != "" {
		return "/static/", true, http.FileServer(http.Dir(static))
	}
	return "", false, nil
}

func (a *App) options() *router.Options {
	return &router.Options{
		Model:        a.Model,
		View:         a.View,
		Config:       a.Config,
		Log:          a.Log,
		SessionStore: a.SessionStore,
	}
}

// Init initializes the MVC App.
func (a *App) Init() error {
	if a.ConfigPath == "" {
		a.SetConfigPath("config")
	}
	return a.init()
}

// SetConfigPath sets the directory path to search for the config files.
func (a *App) SetConfigPath(dir string) {
	a.ConfigPath = dir
}

// init initializes values to the app components.
func (a *App) init() error {
	appConfig, err := loadConfig(a.ConfigPath)
	if err != nil {
		return err
	}
	a.Config = appConfig

	views, err := view.NewSimpleView(appConfig.ViewsDir)
	if err != nil {
		return err
	}
	a.View = views

	// only when mode is allowed
	if !appConfig.NoModel {
		model := models.NewModel()
		err = model.OpenWithConfig(appConfig)
		if err != nil {
			return err
		}
		a.Model = model
	}

	// The sessionistore s really not critical. The application can just run
	// without session set
	store, err := getSesionStore(appConfig)
	if err == nil {
		a.SessionStore = store
	}

	a.Router.Options = a.options()
	a.Router.LoadRoutes(a.ConfigPath) // Load a routes file if available.	
	a.isInit = true

	// In case the StaticDir is specified in the Config file, register
	// a handler serving contents of that directory under the PathPrefix /static/.
	if appConfig.StaticDir != "" {
		static, _ := getAbsolutePath(appConfig.StaticDir)
		if static != "" {
			a.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))
		}

	}
	return nil
}

func getSesionStore(cfg *config.Config) (sessions.Store, error) {
	opts := &sessions.Options{
		Path:     cfg.SessionPath,
		Domain:   cfg.SessionDomain,
		MaxAge:   cfg.SessionMaxAge,
		Secure:   cfg.SessionSecure,
		HttpOnly: cfg.SessionSecure,
	}
	db, err := sql.Open("ql-mem", "session.db")
	if err != nil {
		return nil, err
	}
	err = qlstore.Migrate(db)
	if err != nil {
		return nil, err
	}

	store := qlstore.NewQLStore(db, "/", 2592000, keyPairs(cfg.SessionKeyPair)...)
	store.Options = opts
	return store, nil
}

func keyPairs(src []string) [][]byte {
	var pairs [][]byte
	for _, v := range src {
		pairs = append(pairs, []byte(v))
	}
	return pairs
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
	_ = a.Router.Add(ctrlfn, middlewares...)
}

// ServeHTTP serves http requests. It can be used with other http.Handler implementations.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

//SetNotFoundHandler this sets the hadler that is will execute when the route is
//not found.
func (a *App) SetNotFoundHandler(h http.Handler) error {
	if a.Router != nil {
		a.Router.NotFoundHandler = h
		return nil
	}
	return errors.New("untron: application router is not set")
}
