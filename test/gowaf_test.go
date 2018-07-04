package test

import (
	"log"
	"testing"

	"github.com/NlaakStudios/gowaf"
)

//TestNewApp test the root gowaf NewApp() function
func TestNewApp(t *testing.T) {
	app := gowaf.NewApp()
	app.SetConfigPath("./fixtures/config")
	if err := app.Init(); err != nil {
		log.Fatal(err)
	}

	log.Printf("App Name: %s", app.Config.AppName)
}

//TestNewMVC test the root gowaf NewMVC() function
func TestNewMVC(t *testing.T) {
	app, err := gowaf.NewMVC("fixtures/config")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("App Name: %s", app.Config.AppName)
}
