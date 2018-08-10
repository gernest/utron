# logger
--
    import "github.com/NlaakStudios/gowaf/logger"


## Usage

#### type DefaultLogger

```go
type DefaultLogger struct {
	*log.Logger
}
```

DefaultLogger is the default logger

#### func (*DefaultLogger) Errors

```go
func (d *DefaultLogger) Errors(v ...interface{})
```
Errors log error messages

#### func (*DefaultLogger) Info

```go
func (d *DefaultLogger) Info(v ...interface{})
```
Info logs info messages

#### func (*DefaultLogger) Success

```go
func (d *DefaultLogger) Success(v ...interface{})
```
Success logs success messages

#### func (*DefaultLogger) Warn

```go
func (d *DefaultLogger) Warn(v ...interface{})
```
Warn logs warning messages

#### type Logger

```go
type Logger interface {
	Info(v ...interface{})
	Errors(fv ...interface{})
	Warn(v ...interface{})
	Success(v ...interface{})
}
```

Logger is an interface for gowaf logger

#### func  NewDefaultLogger

```go
func NewDefaultLogger(out io.Writer) Logger
```
NewDefaultLogger returns a default logger writing to out
