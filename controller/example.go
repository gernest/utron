package controller

// Example is the controller for the Example Model
type Example struct {
	BaseController
	Routes []string
}

// Index displays the account example (index) page
func (a *Example) Index() {
	a.Ctx.Template = "application/example/index"
	//a.Ctx.Template = "layout/webapp"
	a.Ctx.Data["title"] = "Home"
	a.Ctx.Data["route"] = "../application/example/index"
	a.Ctx.Data["use_styles"] = false
	a.Ctx.Data["use_sparkline"] = true
	a.Ctx.Data["use_datatables"] = true

	a.Ctx.Data["model_route"] = "example"
	a.Ctx.Log.Success(a.Ctx.Request().Method, " : ", a.Ctx.Template)
}

// List shows a paginated list of all model items based on filter / search info
func (a *Example) List() {

}

// Create creates a new model in the database
func (a *Example) Create() {

}

// Edit edits an existing model in the database with correct access level
func (a *Example) Edit() {

}

// Delete deletes a model in the database with correct access level
func (a *Example) Delete() {

}

// NewExample returns a new account controller object
func NewExample() Controller {
	return &Example{
		Routes: []string{
			"get;/example;Index",
			"get;/example/list;List",
			"get,post;/example/create;Create",
			"get,post;/example/delete;Delete",
		},
	}
}
