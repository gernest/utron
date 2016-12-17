package controller

import (
	"encoding/json"
	"reflect"

	"github.com/gernest/utron/base"
)

// Controller is an interface for utron controllers
type Controller interface {
	New(*base.Context)
	Render() error
}

// BaseController implements the Controller interface, It is recommended all
// user defined Controllers should embed *BaseController.
type BaseController struct {
	Ctx    *base.Context
	Routes []string
}

// New sets ctx as the active context
func (b *BaseController) New(ctx *base.Context) {
	b.Ctx = ctx
}

// Render commits the changes made in the active context.
func (b *BaseController) Render() error {
	return b.Ctx.Commit()
}

// HTML renders text/html with the given code as status code
func (b *BaseController) HTML(code int) {
	b.Ctx.HTML()
	b.Ctx.Set(code)
}

// String renders text/plain with given code as status code
func (b *BaseController) String(code int) {
	b.Ctx.TextPlain()
	b.Ctx.Set(code)
}

// JSON renders application/json with the given code
func (b *BaseController) JSON(code int) {
	b.Ctx.JSON()
	b.Ctx.Set(code)
}

// RenderJSON encodes value into json and renders the response as JSON
func (b *BaseController) RenderJSON(value interface{}, code int) {
	_ = json.NewEncoder(b.Ctx).Encode(value)
	b.JSON(code)
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
