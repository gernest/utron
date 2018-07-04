package gowaf

import (
	"errors"
	"testing"

	"fmt"
	"log"

	"github.com/NlaakStudios/gowaf"
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

var ISTESTING = true

type Account struct {
	gorm.Model
	Username       string `valid:"required,length(6|16)" schema:username`
	Password       string `gorm:"-" valid:"required,length(6|24)" schema:"password"`
	Email          string `valid:"required,email" schema:"email"`
	VerifyPass     string `gorm:"-" schema:"verifypass"`
	HashedPassword string
}

func (u *Account) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}
	if u.Password != u.VerifyPass {
		return errors.New("password missmatch")
	}
	return err
}

func Test_GoWAF(t *testing.T) {
	app, err := gowaf.NewMVC("./fixtures/config")
	if err != nil {
		log.Fatal(err)
	}
	// register models (not for testing)
	//app.Model.Register(&Account{})
	//app.Model.LogMode(true)
	//app.Model.AutoMigrateAll()

	// Register Controllers
	//app.AddController(controllers.NewAccount)

	// Start the server
	port := fmt.Sprintf(":%d", app.Config.Port)
	app.Log.Info("staring server on port", port)
	//log.Fatal(http.ListenAndServe(port, app))
}
