package utron

// Controller is an interface for utron controllers
type Controller interface {
	New(*Context)
	Render() error
}

// BaseController implements the Controlller interface, It is recommended all
// user defined Contollers should embed *BaseController.
type BaseController struct {
	Ctx *Context
}

// New sets ctx as the active context
func (b *BaseController) New(ctx *Context) {
	b.Ctx = ctx
}

// Render commits the changes made in the active context.
func (b *BaseController) Render() error {
	return b.Ctx.Commit()
}

// HTML renders text/html with the given code as status code
func (b *BaseController) HTML(code int) {
	b.Ctx.Set(code)
	b.Ctx.HTML()
}

// String renders text/plain with given code as status code
func (b *BaseController) String(code int) {
	b.Ctx.Set(code)
	b.Ctx.TextPlain()
}

// JSON renders application/json with the given code
func (b *BaseController) JSON(code int) {
	b.Ctx.Set(code)
	b.Ctx.JSON()
}
