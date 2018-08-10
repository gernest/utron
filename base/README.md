# base
--
    import "github.com/NlaakStudios/gowaf/base"

Package base is the basic building cblock of gowaf. The main structure here is
Context, but for some reasons to avoid confusion since there is a lot of context
packages I decided to name this package base instead.

## Usage

```go
var Content = struct {
	Type        string
	TextPlain   string
	TextHTML    string
	Application struct {
		Form, JSON, MultipartForm string
	}
}{
	"Content-Type", "text/plain", "text/html",
	struct {
		Form, JSON, MultipartForm string
	}{
		"application/x-www-form-urlencoded",
		"application/json",
		"multipart/form-data",
	},
}
```
Content holds http response content type strings

#### type Context

```go
type Context struct {

	// Params are the parameters specified in the url patterns
	// gowaf uses gorilla mux for routing. So basically Params stores results
	// after calling mux.Vars function .
	//
	// e.g. if you have route /hello/{world}
	// when you make request to /hello/gernest , then
	// in the Params, key named world will have value gernest. meaning Params["world"]=="gernest"
	Params map[string]string

	// Data keeps values that are going to be passed to the view as context
	Data map[string]interface{}

	// Template is the name of the template to be rendered by the view
	Template string

	// Cfg is the application configuration
	Cfg *config.Config

	//DB is the database stuff, with all models registered
	DB *models.Model

	Log logger.Logger

	SessionStore sessions.Store
}
```

Context wraps request and response. It provides methods for handling responses

#### func  NewContext

```go
func NewContext(w http.ResponseWriter, r *http.Request) *Context
```
NewContext creates new context for the given w and r

#### func (*Context) Commit

```go
func (c *Context) Commit() error
```
Commit writes the results on the underlying http.ResponseWriter and commits the
changes. This should be called only once, subsequent calls to this will result
in an error.

If there is a view, and the template is specified the the view is rendered and
its output is written to the response, otherwise any data written to the context
is written to the ResponseWriter.

#### func (*Context) CoreDataInit

```go
func (c *Context) CoreDataInit()
```
CoreDataInit adds default common properties to the templates available template
variables

#### func (*Context) GetData

```go
func (c *Context) GetData(key interface{}) interface{}
```
GetData retrievess any data stored in the request using gorilla.Context package

#### func (*Context) GetSession

```go
func (ctx *Context) GetSession(name string) (*sessions.Session, error)
```
GetSession retrieves session with a given name.

#### func (*Context) HTML

```go
func (c *Context) HTML()
```
HTML renders text/html response

#### func (*Context) Init

```go
func (c *Context) Init()
```
Init initializes the context

#### func (*Context) JSON

```go
func (c *Context) JSON()
```
JSON renders JSON response

#### func (*Context) NewSession

```go
func (ctx *Context) NewSession(name string) (*sessions.Session, error)
```
NewSession returns a new browser session whose key is set to name. This only
works when the *Context.SessionStore is not nil.

The session returned is from grorilla/sessions package.

#### func (*Context) Redirect

```go
func (c *Context) Redirect(url string, code int)
```
Redirect redirects request to url using code as status code.

#### func (*Context) Request

```go
func (c *Context) Request() *http.Request
```
Request returns the *http.Request object used by the context

#### func (*Context) Response

```go
func (c *Context) Response() http.ResponseWriter
```
Response returns the http.ResponseWriter object used by the context

#### func (*Context) SaveSession

```go
func (ctx *Context) SaveSession(s *sessions.Session) error
```
SaveSession saves the given session.

#### func (*Context) Set

```go
func (c *Context) Set(value interface{})
```
Set sets value in the context object. You can use this to change the following

    * Request by passing *http.Request
    * ResponseWriter by passing http.ResponseVritter
    * view by passing View
    * response status code by passing an int

#### func (*Context) SetData

```go
func (c *Context) SetData(key, value interface{})
```
SetData stores key value into the request object attached with the context. this
is a helper method, wraping gorilla/context

#### func (*Context) SetError

```go
func (c *Context) SetError(number int, name, message string)
```
SetError defines template veriables for Template Error Page

#### func (*Context) SetHeader

```go
func (c *Context) SetHeader(key, value string)
```
SetHeader sets response header

#### func (*Context) TextPlain

```go
func (c *Context) TextPlain()
```
TextPlain renders text/plain response

#### func (*Context) Write

```go
func (c *Context) Write(data []byte) (int, error)
```
Write writes the data to the context, data is written to the http.ResponseWriter
upon calling Commit().

data will only be used when Template is not specified and there is no View set.
You can use this for creating APIs (which does not depend on views like JSON
APIs)
