package base

import (
	"errors"

	"github.com/gorilla/sessions"
)

var errNoStore = errors.New("no session store was found")

//NewSession returns a new browser session whose key is set to name. This only
//works when the *Context.SessionStore is not nil.
//
// The session returned is from grorilla/sessions package.
func (ctx *Context) NewSession(name string) (*sessions.Session, error) {
	if ctx.SessionStore != nil {
		return ctx.SessionStore.New(ctx.Request(), name)
	}
	return nil, errNoStore
}

//GetSession retrieves session with a given name.
func (ctx *Context) GetSession(name string) (*sessions.Session, error) {
	if ctx.SessionStore != nil {
		return ctx.SessionStore.New(ctx.Request(), name)
	}
	return nil, errNoStore
}

//SaveSession saves the given session.
func (ctx *Context) SaveSession(s *sessions.Session) error {
	if ctx.SessionStore != nil {
		return ctx.SessionStore.Save(ctx.Request(), ctx.Response(), s)
	}
	return errNoStore
}
