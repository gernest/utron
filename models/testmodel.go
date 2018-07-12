package models

import (
	"time"
)

// TestModel model used for testing and code coverage
type TestModel struct {
	ID        int       `schema: "-"`
	Body      string    `schema:"body"`
	CreatedAt time.Time `schema:"-"`
	UpdatedAt time.Time `schema:"-"`
}
