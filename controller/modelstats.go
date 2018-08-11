package controller

import (
	"net/http"

	"github.com/NlaakStudios/gowaf/models"
)

// ModelStats is the controller for the ModelStats Model
//Is used to track all registered models and information about them.
//Each models CRUD handlers MUST call
type ModelStats struct {
	BaseController
	Routes []string
}

// Index displays the account ModelStats (index) page
func (c *ModelStats) Index() {
	c.List()
}

// List shows a paginated list of all model items based on filter / search info
func (c *ModelStats) List() {
	ModelStatss := []*models.ModelStats{}
	c.Ctx.DB.Order("created_at desc").Find(&ModelStatss)
	c.Ctx.Data["List"] = ModelStatss
	c.Ctx.Data["use_styles"] = false
	c.Ctx.Data["use_sparkline"] = false
	c.Ctx.Data["use_datatables"] = true
	c.Ctx.Template = "application/modelstats/index"
	c.HTML(http.StatusOK)
}

// Create creates or updates a model in the database
func (c *ModelStats) Create() {
	ModelStats := &models.ModelStats{}
	req := c.Ctx.Request()
	_ = req.ParseForm()
	if err := Decoder.Decode(ModelStats, req.PostForm); err != nil {
		c.Ctx.SetError(400, "Internal Server Errror", err.Error())
		c.HTML(http.StatusInternalServerError)
		return
	}

	c.Ctx.DB.Create(ModelStats)
	c.Ctx.Redirect("/modelstats", http.StatusFound)
}
