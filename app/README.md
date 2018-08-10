# app
--
    import "github.com/NlaakStudios/gowaf/app"


## Usage

#### func  StaticServer

```go
func StaticServer(cfg *config.Config) (string, bool, http.Handler)
```
StaticServer implements StaticServerFunc.

This uses the http.Fileserver to handle static assets. The routes prefixed with
/static/ are static asset routes by default.

#### func  Version

```go
func Version() string
```
version returns the application version as a properly formed string per the
semantic versioning 2.0.0 spec (http://semver.org/).

#### type App

```go
type App struct {
	Version       string
	GoWAFVersion  string
	Router        *router.Router
	Config        *config.Config
	View          view.View
	Log           logger.Logger
	Model         *models.Model
	FixtureFolder string
	ConfigPath    string
	StaticServer  StaticServerFunc
	SessionStore  sessions.Store
}
```

App is the main gowaf application.

#### func  NewApp

```go
func NewApp() *App
```
NewApp creates a new bare-bone gowaf application. To use the MVC components, you
should call the Init method before serving requests.

#### func  NewMVC

```go
func NewMVC(ver string, dir ...string) (*App, error)
```
NewMVC creates a new MVC gowaf app. If dir is passed, it should be a directory
to look for all project folders (config, static, views, models, controllers,
etc). The App returned is initialized.

#### func (*App) AddController

```go
func (a *App) AddController(ctrlfn func() controller.Controller, middlewares ...interface{})
```
AddController registers a controller, and middlewares if any is provided.

#### func (*App) Init

```go
func (a *App) Init() error
```
Init initializes the MVC App.

#### func (*App) Run

```go
func (a *App) Run(f RegisterFunc)
```
Run parses command line arguments and processes commands

#### func (*App) ServeHTTP

```go
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request)
```
ServeHTTP serves http requests. It can be used with other http.Handler
implementations.

#### func (*App) SetConfigPath

```go
func (a *App) SetConfigPath(dir string)
```
SetConfigPath sets the directory path to search for the config files.

#### func (*App) SetFixturePath

```go
func (a *App) SetFixturePath(dir string)
```
SetFixturePath sets the directory path as a base to all other folders (config,
views, etc).

#### func (*App) SetNoModel

```go
func (a *App) SetNoModel(no bool)
```
SetNoModel sets the Config.NoModel value manually in the config.

#### func (*App) SetNotFoundHandler

```go
func (a *App) SetNotFoundHandler(h http.Handler) error
```
SetNotFoundHandler this sets the hadler that is will execute when the route is
not found.

#### func (*App) SetStaticPath

```go
func (a *App) SetStaticPath(dir string)
```
SetStaticPath sets the directory path to search for the static asset files being
served

#### func (*App) SetViewPath

```go
func (a *App) SetViewPath(dir string)
```
SetViewPath sets the directory path to search for the view files.

#### type RegisterFunc

```go
type RegisterFunc func(*App)
```

ResiterFunc is used to pass in a reference to the actual user defined function
for registering the webapp's Models and Controllers

#### type StaticServerFunc

```go
type StaticServerFunc func(*config.Config) (prefix string, strip bool, h http.Handler)
```

StaticServerFunc is a function that returns the static assetsfiles server.

The first argument returned is the path prefix for the static assets. If strp is
set to true then the prefix is going to be stripped.

#### type VersionFunc

```go
type VersionFunc func(*App)
```

VersionFunc is used to pass in a reference to the webapps version() func which
returns a string representing the version.
