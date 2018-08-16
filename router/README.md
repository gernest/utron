# router
--
    import "github.com/NlaakStudios/gowaf/router"


## Usage

```go
var (

	// ErrRouteStringFormat is returned when the route string is of the wrong format
	ErrRouteStringFormat = errors.New("wrong route string, example is\" get,post;/hello/world;Hello\"")
)
```

#### type Middleware

```go
type Middleware struct {
	Type MiddlewareType
}
```

Middleware is the gowaf middleware

#### func (*Middleware) ToHandler

```go
func (m *Middleware) ToHandler(ctx *base.Context) func(http.Handler) http.Handler
```
ToHandler returns a func(http.Handler) http.Handler from the Middleware. gowaf
uses alice to chain middleware.

Use this method to get alice compatible middleware.

#### type MiddlewareType

```go
type MiddlewareType int
```

MiddlewareType is the kind of middleware. gowaf support middleware with variary
of signatures.

```go
const (
	//PlainMiddleware is the middleware with signature
	// func(http.Handler)http.Handler
	PlainMiddleware MiddlewareType = iota

	//CtxMiddleware is the middlewate with signature
	// func(*base.Context)error
	CtxMiddleware
)
```

#### type Options

```go
type Options struct {
	Model        *models.Model
	View         view.View
	Config       *config.Config
	Log          logger.Logger
	SessionStore sessions.Store
}
```

Options additional settings for the router.

#### type Router

```go
type Router struct {
	*mux.Router

	Options   *Options
	NumCtrlrs int
}
```

Router registers routes and handlers. It embeds gorilla mux Router

#### func  NewRouter

```go
func NewRouter(app ...*Options) *Router
```
NewRouter returns a new Router, if app is passed then it is used

#### func (*Router) Add

```go
func (r *Router) Add(ctrlfn func() controller.Controller, middlewares ...interface{}) error
```
Add registers ctrl. It takes additional comma separated list of middleware.
middlewares are of type

    func(http.Handler)http.Handler
    or
    func(*base.Context)error

gowaf uses the alice package to chain middlewares, this means all alice
compatible middleware works out of the box

#### func (*Router) Count

```go
func (r *Router) Count() int
```
Count returns the number of registered models

#### func (*Router) LoadRoutes

```go
func (r *Router) LoadRoutes(cfgPath string)
```
LoadRoutes searches for the route file in the cfgPath. The order of file lookup
is as follows.

    * routes.json
    * routes.toml
    * routes.yml
    * routes.hcl

#### func (*Router) LoadRoutesFile

```go
func (r *Router) LoadRoutesFile(file string) error
```
LoadRoutesFile loads routes from a json file. Example of the routes file.

    {
    	"routes": [
    		"get,post;/hello;Sample.Hello",
    		"get,post;/about;Hello.About"
    	]
    }

supported formats are json, toml, yaml and hcl with extension .json, .toml, .yml
and .hcl respectively.

TODO refactor the decoding part to a separate function? This part shares the
same logic as the one found in NewConfig()

#### func (*Router) Static

```go
func (r *Router) Static(prefix string, h http.FileSystem)
```
Static registers static handler for path perfix
