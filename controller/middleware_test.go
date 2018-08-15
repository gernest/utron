package controller

import (
	"testing"
	"github.com/gorilla/sessions"
	"net/http/httptest"
	"net/http"
	"bytes"
	"fmt"
	"time"
	"github.com/gorilla/securecookie"
)

var store *sessions.CookieStore

const port = "127.0.0.1:9091"
const port2 = "127.0.0.1:9092"
const path = "/account/login"

func TestCheckloginBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	var res []byte
	out := bytes.NewReader(res)

	r, err := http.NewRequest("GET", "http://" + port + path, out)
	if err != nil {
		panic(err)
	}

	createStore()
	session, err := store.Get(r, "new")

	if err != nil {
		panic(err)
	}

	session.Values["state"] = "dsfdsf1"
	session.ID = "new id"
	fmt.Println(session.Name())
	err = session.Save(r, w)

	if err != nil {
		panic(err)
	}

	//s, err := store.Get(r, "new")
	//fmt.Println("TJIS", s.Values)

	test := Checklogin(store, session.Name())
	handler := getTestHandler()

	go func() {
		e := http.ListenAndServe(port, test(handler))
		if e != nil {
			panic(e)
		}
	}()

 	time.Sleep(time.Second)
	client := http.DefaultClient

	resp, err := client.Do(r)

	if err != nil {
		fmt.Println(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Error("Expected http.StatusBadRequest, got: ", resp.Status)
	}
}

func TestCheckloginIternalServerError(t *testing.T) {
	var res []byte
	out := bytes.NewReader(res)
	r, err := http.NewRequest("GET", "http://" + port2 + path, out)
	if err != nil {
		panic(err)
	}

	test := Checklogin(nil, "new")
	handler := getTestHandler()

	go func() {
		e := http.ListenAndServe(port2, test(handler))
		if e != nil {
			panic(e)
		}
	}()

	time.Sleep(time.Second)
	client := http.DefaultClient

	resp, err := client.Do(r)

	if resp.StatusCode != http.StatusInternalServerError {
		t.Error("Expected http.StatusBadRequest, got: ", resp.Status)
	}
}

func getTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		panic("test entered test handler, this should not happen")
	}
	return http.HandlerFunc(fn)
}

func createStore() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
}