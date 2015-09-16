package main

import (
	"github.com/gernest/utron"
	_ "github.com/gernest/utron/fixtures/todo/controllers"
	_ "github.com/gernest/utron/fixtures/todo/models"
)

func main() {
	utron.Run()
}
