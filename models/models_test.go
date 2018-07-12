package models

import (
	"testing"

	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
)

type TestData struct {
	gorm.Model
	AString  string `valid:"required,length(6|16)" schema:"astring"`
	AInteger int    `schema:"ainteger"`
	ABoolean bool   `schema:"aboolean"`
}

var myTestModel TestData
var decoder = schema.NewDecoder()

func TestModel(t *testing.T) {

	myTestModel := NewModel()
	if myTestModel == nil {
		t.Fail()
	}

	myTestModel.Model.Register(&models.TestData{})
}
