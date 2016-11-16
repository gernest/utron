package utron

import (
	"os"

	"github.com/gernest/utron/app"
	"github.com/gernest/utron/logger"
	"github.com/gernest/utron/models"
	"github.com/gernest/utron/router"
)

// NewApp creates a new bare-bone utron application. To use the MVC components, you should call
// the Init method before serving requests.
func NewApp() *app.App {
	return &app.App{
		Log:    logger.NewDefaultLogger(os.Stdout),
		Router: router.NewRouter(),
		Model:  models.NewModel(),
	}
}

// NewMVC creates a new MVC utron app. If cfg is passed, it should be a directory to look for
// the configuration files. The App returned is initialized.
func NewMVC(cfg ...string) (*app.App, error) {
	app := NewApp()
	if len(cfg) > 0 {
		app.SetConfigPath(cfg[0])
	}
	if err := app.Init(); err != nil {
		return nil, err
	}
	return app, nil
}
