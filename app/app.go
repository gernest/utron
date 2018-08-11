package app

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NlaakStudios/gowaf/config"
	"github.com/NlaakStudios/gowaf/controller"
	"github.com/NlaakStudios/gowaf/logger"
	"github.com/NlaakStudios/gowaf/models"
	"github.com/NlaakStudios/gowaf/router"
	"github.com/NlaakStudios/gowaf/view"
	_ "github.com/cznic/ql/driver"
	"github.com/gernest/qlstore"
	"github.com/gorilla/sessions"
)

//StaticServerFunc is a function that returns the static assetsfiles server.
//
// The first argument returned is the path prefix for the static assets. If strp
// is set to true then the prefix is going to be stripped.
type StaticServerFunc func(*config.Config) (prefix string, strip bool, h http.Handler)

// App is the main gowaf application.
type App struct {
	Version       string
	GoWAFVersion  string
	Router        *router.Router
	Config        *config.Config
	View          view.View
	Log           logger.Logger
	Model         *models.Model
	FixtureFolder string
	ConfigName    string
	ConfigFolder  string
	ViewFolder    string
	StaticFolder  string
	StaticServer  StaticServerFunc
	SessionStore  sessions.Store
	isInit        bool
}

// ResiterFunc is used to pass in a reference to the actual user defined function for registering
// the webapp's Models and Controllers
type RegisterFunc func(*App)

// VersionFunc is used to pass in a reference to the webapps version() func which returns a string
// representing the version.
type VersionFunc func(*App)

// NewApp creates a new bare-bone gowaf application. To use the MVC components, you should call
// the Init method before serving requests.
func NewApp() *App {
	return &App{
		//Version:      "0.0.0-notset", ///This is set in the webapp not framework
		GoWAFVersion: Version(),
		Log:          logger.NewDefaultLogger(os.Stdout),
		Router:       router.NewRouter(),
		Model:        models.NewModel(),
	}
}

// NewMVC creates a new MVC gowaf app. If dir is passed, it should be a directory to look for
// all project folders (config, static, views, models, controllers, etc). The App returned is initialized.
func NewMVC(ver, cfgDir, cfgName string) (*App, error) {
	app := NewApp()
	app.Version = ver

	//Prepare Config Folder
	if len(cfgDir) > 0 {
		app.SetFixturePath(cfgDir)
	} else {
		app.SetFixturePath("./fixtures")
	}

	//Prepare Config Name (without extension)
	if len(cfgName) > 0 {
		app.ConfigName = cfgName
	} else {
		app.ConfigName = "webapp"
	}

	app.SetConfigFolder(fmt.Sprintf("%s/config", app.FixtureFolder))

	if err := app.Init(); err != nil {
		return nil, err
	}

	fmt.Printf(
		"%s v%s\nGoWAf Framework v%s\nHost: %s\n%s\n",
		app.Config.AppName, app.Version,
		app.GoWAFVersion,
		app.Config.BaseURL,
		"------------------------------------------------------",
	)
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
	if a.ConfigFolder == "" {
		a.SetConfigFolder("config")
	}
	return a.init()
}

// SetFixturePath sets the directory path as a base to all other folders (config, views, etc).
func (a *App) SetFixturePath(dir string) {
	a.FixtureFolder = dir
}

// SetConfigFolder sets the directory path to search for the config files.
func (a *App) SetConfigFolder(dir string) {
	a.ConfigFolder = dir
}

// SetViewPath sets the directory path to search for the view files.
func (a *App) SetViewPath(dir string) {
	if dir == "" {
		dir = "views"
	}
	a.ViewFolder = fmt.Sprintf("%s/%s", a.FixtureFolder, dir)
}

// SetStaticPath sets the directory path to search for the static asset files being served
func (a *App) SetStaticPath(dir string) {
	if dir == "" {
		dir = "assets"
	}
	a.StaticFolder = fmt.Sprintf("%s/%s", a.FixtureFolder, dir)
}

// SetNoModel sets the Config.NoModel value manually in the config.
func (a *App) SetNoModel(no bool) {
	a.Config.NoModel = no
}

// init initializes values to the app components.
func (a *App) init() error {
	appConfig, err := loadConfig(a.ConfigName, a.ConfigFolder)
	if err != nil {
		return err
	}
	a.Config = appConfig

	a.SetViewPath("views")
	viewsabs, _ := getAbsolutePath(a.ViewFolder)
	if viewsabs == "" {
		return err
	}
	a.ViewFolder = viewsabs
	if _, err := os.Stat(a.ViewFolder); os.IsNotExist(err) {
		// path/to/view folder does not exist use default fixtures/view
		a.ViewFolder = fmt.Sprintf("%s/%s", a.FixtureFolder, "/views")
	}

	views, err := view.NewSimpleView(a.ViewFolder)
	if err != nil {
		//TODO: Coverage - Need Failure here
		return err
	}
	a.View = views

	// only when model is allowed
	if !appConfig.NoModel {
		model := models.NewModel()
		err = model.OpenWithConfig(appConfig)
		if err != nil {
			return err
		}

		//TODO: Coverage - Need good config here with No_Model = true
		a.Model = model
	}

	// The sessionistore is really not critical. The application can just run
	// without session set
	store, err := getSesionStore(appConfig)
	if err == nil {
		a.SessionStore = store
	}

	a.Router.Options = a.options()
	a.Router.LoadRoutes(a.ConfigFolder) // Load a routes file if available.
	a.isInit = true

	// In case the StaticDir is specified in the Config file, register
	// a handler serving contents of that directory under the PathPrefix /static/.
	//if appConfig.StaticDir != "" {
	a.SetStaticPath("assets")
	static, _ := getAbsolutePath(a.StaticFolder)
	if static != "" {
		//TODO: Coverage -  Need to hit here
		a.SetStaticPath(static)
		a.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))
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
		//TODO: Coverage -  Need to hit here
		return nil, err
	}

	err = qlstore.Migrate(db)
	if err != nil {
		//TODO: Coverage -  Need to hit here
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
		//TODO: Coverage -  Need to hit here
		return "", fmt.Errorf("gowaf: %s is not a directory", dir)
	}

	if filepath.IsAbs(dir) { // If dir is already absolute, return it.
		return dir, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		//TODO: Coverage -  Need to hit here
		return "", err
	}
	absDir := filepath.Join(wd, dir)
	_, err = os.Stat(absDir)
	if err != nil {
		//TODO: Coverage -  Need to hit here
		return "", err
	}
	return absDir, nil
}

// loadConfig loads the configuration file. If cfg is provided, then it is used as the directory
// for searching the configuration files. It defaults to the directory named config in the current
// working directory.
func loadConfig(name string, cfg ...string) (*config.Config, error) {
	cfgDir := "config"
	if len(cfg) > 0 {
		cfgDir = cfg[0]
	}

	// Load configurations.
	cfgFile, err := findConfigFile(cfgDir, name)
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
			//TODO: Coverage -  Need to hit here
			return
		}
		file = file + ext
		if info, serr := os.Stat(file); serr == nil && !info.IsDir() {
			return
		}
	}
	return "", fmt.Errorf("gowaf: can't find configuration file %s in %s", name, dir)
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
	return errors.New("gowaf: application router is not set")
}

//************COMMAND LINE STUFF ***************/
func (a *App) printHeader() {
	fmt.Printf("%s v%s Daemon\n", a.Config.AppName, a.Version)
	fmt.Println("-----------------------------------------------------------------------------------------")
}

//printUsage diplay commandline usage information to the user.
func (a *App) printUsage() {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("	usefolder -path/to/fixtures/folder - Defines the target fixture folder to use")
	fmt.Println("	startnode - Start a node")
	fmt.Println("	version - Display version")
	fmt.Println()
}

// Run parses command line arguments and processes commands
func (a *App) Run(f RegisterFunc) {

	if len(os.Args) < 2 {
		a.printUsage()
		os.Exit(1)
	}

	userFolderCmd := flag.NewFlagSet("userfolder", flag.ExitOnError)
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
	useFolder := userFolderCmd.String("path", "", "The path to the fixtures folder to use")

	switch os.Args[1] {
	case "userfolder":
		err := userFolderCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "version":
		err := versionCmd.Parse(os.Args[1:])
		if err != nil {
			log.Panic(err)
		}
	default:
		a.printUsage()
		os.Exit(1)
	}

	if userFolderCmd.Parsed() {
		if *useFolder == "" {
			userFolderCmd.Usage()
			os.Exit(1)
		}
		a.FixtureFolder = *useFolder
	}

	if startNodeCmd.Parsed() {
		if a.Config.Verbose {
			a.Log.Info("Using base fixture folder at ", a.FixtureFolder)
			a.Log.Info("Using static assets at ", a.Config.StaticDir)
			a.Log.Info("Using views at ", a.ViewFolder)

			if a.Config.LoadTestData {
				a.Log.Warn("Load Test Data is enabled in config, please turn off for production.")
			}

			if a.Config.GoogleID != "" {
				a.Log.Success("Google Analytics enabled with ID: ", a.Config.GoogleID)
			}
		}

		// Call the users Register Function
		a.Log.Info("Registering Models & Controllers...")
		f(a)
		a.Log.Success("Done. ", a.Model.Count(), " Models and ", a.Router.Count(), " Controllers registered.")

		port := fmt.Sprintf(":%d", a.Config.Port)
		a.Log.Info("Starting server, listening on port", port)
		log.Fatal(http.ListenAndServe(port, a))

	}

	if versionCmd.Parsed() {
		fmt.Println(a.Version)
	}

}
