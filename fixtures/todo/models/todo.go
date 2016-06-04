package models

import (
	"time"

	"github.com/gernest/utron"
)

//Todo represent an item of todo list
type Todo struct {
	ID        int       `schema:"-"`
	Body      string    `schema:"body"`
	CreatedAt time.Time `schema:"-"`
	UpdatedAt time.Time `schema:"-"`
}

func init() {
	utron.RegisterModels(&Todo{})
}
