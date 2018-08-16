# flash
--
    import "github.com/NlaakStudios/gowaf/flash"


## Usage

```go
const (
	// FlashSuccess is the context key for success flash messages
	FlashSuccess = "FlashSuccess"

	// FlashWarn is a context key for warning flash messages
	FlashWarn = "FlashWarn"

	// FlashErr is a context key for flash error message
	FlashErr = "FlashError"
)
```

#### func  AddFlashToCtx

```go
func AddFlashToCtx(ctx *base.Context, name, key string) error
```
AddFlashToCtx takes flash messages stored in a cookie which is associated with
the request found in ctx, and puts them inside the ctx object. The flash
messages can then be retrieved by calling ctx.Get( FlashKey).

NOTE When there are no flash messages then nothing is set.

#### type Flash

```go
type Flash struct {
	Kind    string
	Message string
}
```

Flash implements flash messages, like ones in gorilla/sessions

#### type Flasher

```go
type Flasher struct {
}
```

Flasher tracks flash messages

#### func  New

```go
func New() *Flasher
```
New creates new flasher. This alllows accumulation of lash messages. To save the
flash messages the Save method should be called explicitly.

#### func (*Flasher) Add

```go
func (f *Flasher) Add(kind, message string)
```
Add adds the flash message

#### func (*Flasher) Err

```go
func (f *Flasher) Err(msg string)
```
Err adds error flash message

#### func (*Flasher) Save

```go
func (f *Flasher) Save(ctx *base.Context, name, key string) error
```
Save saves flash messages to context

#### func (*Flasher) Success

```go
func (f *Flasher) Success(msg string)
```
Success adds success flash message

#### func (*Flasher) Warn

```go
func (f *Flasher) Warn(msg string)
```
Warn adds warning flash message

#### type Flashes

```go
type Flashes []*Flash
```

Flashes is a collection of flash messages

#### func  GetFlashes

```go
func GetFlashes(ctx *base.Context, name, key string) (Flashes, error)
```
GetFlashes retieves all flash messages found in a cookie session associated with
ctx..

name is the session name which is used to store the flash messages. The flash
messages can be stored in any session, but it is a good idea to separate session
for flash messages from other sessions.

key is the key that is used to identiry which flash messages are of interest.
