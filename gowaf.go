// MIT license that can be found in the LICENSE file.
// Go Web Application Framework Library.

package gowaf

import "github.com/NlaakStudios/gowaf/app"

// NewApp creates a new bare-bone gwaf application. To use the MVC components, you should call
// the Init method before serving requests.
func NewApp() *app.App {
	return app.NewApp()
}

// NewMVC creates a new MVC gwaf app. If cfg is passed, it should be a directory to look for
// the configuration files. The App returned is initialized.
func NewMVC(cfg ...string) (*app.App, error) {
	return app.NewMVC(cfg...)
}
