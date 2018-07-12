package models

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/NlaakStudios/gowaf"
	c "github.com/NlaakStudios/gowaf/controller"
	"github.com/NlaakStudios/gowaf/models"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func TestModels(t *testing.T) {
	// Start the MVC App
	app, err := gowaf.NewMVC()
	if err != nil {
		log.Fatal(err)
	}

	// Register Models
	app.Model.Register(&models.TestModel{})

	// CReate Models tables if they dont exist yet
	app.Model.AutoMigrateAll()

	// Register Controller
	app.AddController(c.NewTestModel)

	// Start the server
	port := fmt.Sprintf(":%d", app.Config.Port)
	log.Fatal(http.ListenAndServe(port, app))
}
