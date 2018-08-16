package controller

import (
	"github.com/NlaakStudios/gowaf/models"
	"github.com/casbin/casbin"
	"github.com/gorilla/sessions"
	"net/http"
)

var role byte

//Authorizer - check for user access
func Authorizer(e *casbin.Enforcer, sessionStore sessions.Store, sessionName string, model *models.Model) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			session, err := sessionStore.Get(r, sessionName)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			getRole(session)

			if role > 3 || role < 0 {
				role = 0
			}

			if role > 0 {
				uid := session.Values["uid"].(int)

				account := &models.Account{ID: uid}
				rows := model.Find(account)

				if rows.RowsAffected == 0 {
					w.WriteHeader(http.StatusForbidden)
					return
				}

				role = account.Access
			}

			res, err := e.EnforceSafe(string(role), r.URL.Path, r.Method)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			if res {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func getRole(session *sessions.Session) {
	roleI := session.Values["role"]

	if roleI == nil {
		role = 0
	} else {
		role = roleI.(byte)
	}
}