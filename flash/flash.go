package flash

import (
	"encoding/gob"
	"errors"

	"github.com/gernest/utron/base"
)

const (
	// FlashSuccess is the context key for success flash messages
	FlashSuccess = "FlashSuccess"

	// FlashWarn is a context key for warning flash messages
	FlashWarn = "FlashWarn"

	// FlashErr is a context key for flash error message
	FlashErr = "FlashError"
)

func init() {
	gob.Register(&Flash{})
	gob.Register(Flashes{})
}

// Flash implements flash messages, like ones in gorilla/sessions
type Flash struct {
	Kind    string
	Message string
}

// Flashes is a collection of flash messages
type Flashes []*Flash

// GetFlashes retieves all flash messages found in a cookie session associated with ctx..
//
// name is the session name which is used to store the flash messages. The flash
// messages can be stored in any session, but it is a good idea to separate
// session for flash messages from other sessions.
//
// key is the key that is used to identiry which flash messages are of interest.
func GetFlashes(ctx *base.Context, name, key string) (Flashes, error) {
	ss, err := ctx.GetSession(name)
	if err != nil {
		return nil, err
	}
	if v, ok := ss.Values[key]; ok {
		delete(ss.Values, key)
		serr := ss.Save(ctx.Request(), ctx.Response())
		if serr != nil {
			return nil, serr
		}
		return v.(Flashes), nil
	}
	return nil, errors.New("no flashes found")
}

// AddFlashToCtx takes flash messages stored in a cookie which is associated with the
// request found in ctx, and puts them inside the ctx object. The flash messages can then
// be retrieved by calling ctx.Get( FlashKey).
//
// NOTE When there are no flash messages then nothing is set.
func AddFlashToCtx(ctx *base.Context, name, key string) error {
	f, err := GetFlashes(ctx, name, key)
	if err != nil {
		return err
	}
	ctx.SetData(key, f)
	return nil
}

//Flasher tracks flash messages
type Flasher struct {
	f Flashes
}

//New creates new flasher. This alllows accumulation of lash messages. To save the flash messages
//the Save method should be called explicitly.
func New() *Flasher {
	return &Flasher{}
}

// Add adds the flash message
func (f *Flasher) Add(kind, message string) {
	fl := &Flash{kind, message}
	f.f = append(f.f, fl)
}

// Success adds success flash message
func (f *Flasher) Success(msg string) {
	f.Add(FlashSuccess, msg)
}

// Err adds error flash message
func (f *Flasher) Err(msg string) {
	f.Add(FlashErr, msg)
}

// Warn adds warning flash message
func (f *Flasher) Warn(msg string) {
	f.Add(FlashWarn, msg)
}

// Save saves flash messages to context
func (f *Flasher) Save(ctx *base.Context, name, key string) error {
	ss, err := ctx.GetSession(name)
	if err != nil {
		return err
	}
	var flashes Flashes
	if v, ok := ss.Values[key]; ok {
		flashes = v.(Flashes)
	}
	ss.Values[key] = append(flashes, f.f...)
	err = ss.Save(ctx.Request(), ctx.Response())
	if err != nil {
		return err
	}
	return nil
}
