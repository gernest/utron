package controller

import (
	"net/http"
	"github.com/gorilla/sessions"
	"log"
)

//Checklogin - check if user login
func Checklogin(sessionStore sessions.Store, sessionName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if sessionStore == nil {
				log.Print("ERR")
				http.Redirect(w,r,"/account/login", http.StatusInternalServerError)
				return
			}

			session, _ := sessionStore.Get(r, sessionName)
			state := session.Values["state"]

			if state == nil {
				http.Redirect(w,r,"/account/login", http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
