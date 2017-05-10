// Package base is the basic building cblock of utron. The main structure here is
// Context, but for some reasons to avoid confusion since there is a lot of
// context packages I decided to name this package base instead.
package base

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/gernest/utron/config"
	"github.com/gernest/utron/logger"
	"github.com/gernest/utron/models"
	"github.com/gernest/utron/view"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Content holds http response content type strings
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

// Context wraps request and response. It provides methods for handling responses
type Context struct {

	// Params are the parameters specified in the url patterns
	// utron uses gorilla mux for routing. So basically Params stores results
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

	request    *http.Request
	response   http.ResponseWriter
	out        io.ReadWriter
	isCommited bool
	view       view.View
}

// NewContext creates new context for the given w and r
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		Params:   make(map[string]string),
		Data:     make(map[string]interface{}),
		request:  r,
		response: w,
		out:      &bytes.Buffer{},
	}
	ctx.Init()
	return ctx
}

// Init initializes the context
func (c *Context) Init() {
	c.Params = mux.Vars(c.request)
}

// Write writes the data to the context, data is written to the http.ResponseWriter
// upon calling Commit().
//
// data will only be used when Template is not specified and there is no View set. You can use
// this for creating APIs (which does not depend on views like JSON APIs)
func (c *Context) Write(data []byte) (int, error) {
	return c.out.Write(data)
}

// TextPlain renders text/plain response
func (c *Context) TextPlain() {
	c.SetHeader(Content.Type, Content.TextPlain)
}

// JSON renders JSON response
func (c *Context) JSON() {
	c.SetHeader(Content.Type, Content.Application.JSON)
}

// HTML renders text/html response
func (c *Context) HTML() {
	c.SetHeader(Content.Type, Content.TextHTML)
}

// Request returns the *http.Request object used by the context
func (c *Context) Request() *http.Request {
	return c.request
}

// Response returns the http.ResponseWriter object used by the context
func (c *Context) Response() http.ResponseWriter {
	return c.response
}

// GetData retrievess any data stored in the request using
// gorilla.Context package
func (c *Context) GetData(key interface{}) interface{} {
	return context.Get(c.Request(), key)
}

//SetData stores key value into the request object attached with the context.
//this is a helper method, wraping gorilla/context
func (c *Context) SetData(key, value interface{}) {
	context.Set(c.Request(), key, value)
}

// Set sets value in the context object. You can use this to change the following
//
//	 * Request by passing *http.Request
//	 * ResponseWriter by passing http.ResponseVritter
//	 * view by passing View
//	 * response status code by passing an int
func (c *Context) Set(value interface{}) {
	switch value.(type) {
	case view.View:
		c.view = value.(view.View)
	case *http.Request:
		c.request = value.(*http.Request)
	case http.ResponseWriter:
		c.response = value.(http.ResponseWriter)
	case int:
		c.response.WriteHeader(value.(int))
	}
}

// SetHeader sets response header
func (c *Context) SetHeader(key, value string) {
	c.response.Header().Set(key, value)
}

// Commit writes the results on the underlying http.ResponseWriter and commits the changes.
// This should be called only once, subsequent calls to this will result in an error.
//
// If there is a view, and the template is specified the the view is rendered and its
// output is written to the response, otherwise any data written to the context is written to the
// ResponseWriter.
func (c *Context) Commit() error {
	if c.isCommited {
		return errors.New("already committed")
	}
	if c.Template != "" && c.view != nil {
		out := &bytes.Buffer{}		
		err := c.view.Render(out, c.Template, c.Data)
		if err != nil {
			return err
		}
		_, _ = io.Copy(c.response, out)
	} else {
		_, _ = io.Copy(c.response, c.out)
	}
	c.isCommited = true
	return nil
}

// Redirect redirects request to url using code as status code.
func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.Response(), c.Request(), url, code)
}
