package controller

import (
	"database/sql"
	"fmt"
	"github.com/NlaakStudios/gowaf/base"
	"github.com/NlaakStudios/gowaf/config"
	"github.com/NlaakStudios/gowaf/logger"
	"github.com/NlaakStudios/gowaf/models"
	_ "github.com/cznic/ql/driver"
	"github.com/gernest/qlstore"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

var (
	account *Account

	accountRoutes = []string{
		"get;/account;Index",
		"get,post;/account/register;Register",
		"get,post;/account/login;Login",
		"get;/account/logout;Logout",
	}

	accountUsername      = "username"
	accountPassword      = "password"
	accountEmail         = accountEmailUsername + "@" + accountEmailDomain
	accountEmailUsername = "mail"
	accountEmailDomain   = "gmail.com"
	aID                  = int64(1)

	findQueryAccount   = "SELECT * FROM `accounts` WHERE (Username = ?) ORDER BY `accounts`.`id` ASC LIMIT 1"
	updateQueryAccount = "UPDATE `accounts` SET `created_at` = ?, `hashed_password` = ?, `id` = ?, `state` = ?, `updated_at` = ?, `username` = ? WHERE `accounts`.`id` = ?"

	accountFields = []string{"id", "username", "email", "hashed_password", "state", "access", "email_address", "email_id", "company_id", "person_id", "created_at", "updated_at"}
)

func TestAccount_Index(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", routeToAccount)
	account, _ = prepareAccount(req, rr)

	account.Index()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

func TestAccount_RegisterGet(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", routeToAccount)
	account, _ = prepareAccount(req, rr)

	account.Register()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

func TestAccount_RegisterWithoutData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", routeToAccount)
	account, _ = prepareAccount(req, rr)

	account.Register()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.BadRequest, got: ", rr.Result().Status)
	}
}

func TestAccount_RegisterInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", routeToAccount)
	account, _ = prepareAccount(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("sdsd", "fdsf")

	account.Register()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestAccount_Register(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", routeToAccount)
	account, _ = prepareAccount(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("username", accountUsername)
	req.PostForm.Add("password", accountPassword)
	req.PostForm.Add("verify_password", accountPassword)
	req.PostForm.Add("email", accountEmail)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `emails`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), accountEmailUsername, accountEmailDomain,
		sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(aID, 1))
	mock.ExpectExec("INSERT INTO `accounts`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), accountUsername, accountEmail,
		sqlmock.AnyArg(), models.UserStateVerifyEmailSent, sqlmock.AnyArg(), aID, 0, 0).WillReturnResult(sqlmock.NewResult(aID, 1))
	mock.ExpectCommit()

	account.Register()
	if rr.Result().StatusCode != http.StatusFound {
		t.Error("Expected http.StatusFound, got: ", rr.Result().Status)
	}
}

//TODO fix bug with session
func TestAccount_LoginGet(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", routeToAccount)
	account, _ = prepareAccount(req, rr)

	session, _ := account.Ctx.GetSession(account.Ctx.Cfg.SessionName)

	session.Values["state"] = "2"

	account.Ctx.SessionStore.Save(account.Ctx.Request(), account.Ctx.Response(), session)

	s, _ := account.Ctx.SessionStore.Get(account.Ctx.Request(), account.Ctx.Cfg.SessionName)
	fmt.Println(s.Values)
	account.Login()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusOK, got: ", rr.Result().Status)
	}
}

func TestAccount_LoginWithoutData(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", routeToAccount)
	account, _ = prepareAccount(req, rr)

	account.Login()
	if rr.Result().StatusCode != http.StatusUnauthorized {
		t.Error("Expected http.StatusUnauthorized, got: ", rr.Result().Status)
	}
}

func TestAccount_LoginInvalidForm(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", routeToAccount)
	account, _ = prepareAccount(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("sdsd", "fdsf")

	account.Login()
	if rr.Result().StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusInternalServerError, got: ", rr.Result().Status)
	}
}

func TestAccount_LoginWithoutPassword(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", routeToAccount)
	account, _ = prepareAccount(req, rr)

	req.PostForm = url.Values{}
	req.PostForm.Add("username", accountUsername)

	rows := sqlmock.NewRows(accountFields)
	mock.ExpectQuery(fixedFullRe(findQueryAccount)).WithArgs(accountUsername).WillReturnRows(
		rows.AddRow(aID, accountUsername, &models.Email{}, "somehash", 2, 0, 0, 0, 0, 0, time.Now(), time.Now()))

	account.Login()
	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestAccount_Login(t *testing.T) {
	req, rr = prepareReqAndRecorder("POST", routeToAccount)
	account, _ = prepareAccount(req, rr)

	hash, err := bcrypt.GenerateFromPassword([]byte(accountPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	req.PostForm = url.Values{}
	req.PostForm.Add("username", accountUsername)
	req.PostForm.Add("password", accountPassword)

	rows := sqlmock.NewRows(accountFields)
	mock.ExpectQuery(fixedFullRe(findQueryAccount)).WithArgs(accountUsername).WillReturnRows(
		rows.AddRow(aID, accountUsername, &models.Email{}, hash, 2, 0, 0, 0, 0, 0, time.Now(), time.Now()))
	mock.ExpectExec(fixedFullRe(updateQueryAccount)).WithArgs(sqlmock.AnyArg(), string(hash), pID, models.UserStateSignedIn,
		sqlmock.AnyArg(), accountUsername, aID).WillReturnResult(sqlmock.NewResult(aID, 1))

	account.Login()
	if rr.Result().StatusCode != http.StatusOK {
		t.Error("Expected http.StatusBadRequest, got: ", rr.Result().Status)
	}
}

func TestAccount_Logout(t *testing.T) {
	req, rr = prepareReqAndRecorder("GET", routeToAccount)
	account, _ = prepareAccount(req, rr)

	session, _ := account.Ctx.GetSession(account.Ctx.Cfg.SessionName)

	session.Values["uid"] = "2"
	session.Save(account.Ctx.Request(), account.Ctx.Response())

	account.Logout()
	if rr.Result().StatusCode != http.StatusUnauthorized {
		t.Error("Expected http.StatusUnauthorized, got: ", rr.Result().Status)
	}
}

func prepareAccount(req *http.Request, rr *httptest.ResponseRecorder) (*Account, *base.Context) {
	model := models.NewModel()

	var db *sql.DB

	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("can't create sqlmock: %s", err)
	}

	gormDB, gerr := gorm.Open("mysql", db)
	if gerr != nil {
		log.Fatalf("can't open gorm connection: %s", err)
	}
	gormDB.LogMode(true)

	model.DB = gormDB

	ctx = base.NewContext(rr, req)
	ctx.DB = model
	ctx.Log = logger.NewDefaultLogger(os.Stdout)

	contr := BaseController{Ctx: ctx, Routes: accountRoutes}
	account := &Account{}

	account.Routes = accountRoutes
	account.BaseController = contr

	conf := config.DefaultConfig()
	account.Ctx.Cfg = conf

	account.Ctx.SessionStore = getStore(conf)

	return account, ctx
}

func getStore(cfg *config.Config) sessions.Store {
	db, err := sql.Open("ql-mem", "session.db")
	if err != nil {
		//TODO: Coverage -  Need to hit here
		panic(err)
	}

	err = qlstore.Migrate(db)
	if err != nil {
		//TODO: Coverage -  Need to hit here
		panic(err)
	}
	res := qlstore.NewQLStore(db, "/", 2592000, keyPairs(cfg.SessionKeyPair)...)
	res.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	return res
}

func keyPairs(src []string) [][]byte {
	var pairs [][]byte
	for _, v := range src {
		pairs = append(pairs, []byte(v))
	}
	return pairs
}
