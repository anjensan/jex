package runtime

import (
	"fmt"
	"strings"
)

const Unused = 0

type Exception interface {
	Err() error
	Wrap() ExceptionError

	Suppress(e Exception)
	Suppressed() []Exception
	ClearSuppressed()

	Log(v interface{})
	Logf(f string, a ...interface{})
	Logs() []interface{}

	private()
}

type ExceptionError interface {
	error
	Unwrap() Exception

	private()
}

type exception struct {
	err      error
	suppress []Exception
	logs     []interface{}
}

type exError struct {
	ex *exception
}

func (e exError) private() {}

func (e exError) Error() string {
	return e.ex.String()
}

func (e exError) Cause() error {
	return e.ex.Err()
}

func (e exError) Unwrap() Exception {
	return e.ex
}
func (e *exception) private() {}

func (e *exception) Err() error {
	return e.err
}

func (e *exception) Wrap() ExceptionError {
	return exError{e}
}

func (e *exception) Suppress(ex Exception) {
	if e != ex {
		e.suppress = append(e.suppress, ex)
	}
}

func (e *exception) Suppressed() []Exception {
	return e.suppress
}

func (e *exception) ClearSuppressed() {
	e.suppress = nil
}

func (e *exception) Log(v interface{}) {
	e.logs = append(e.logs, v)
}

func (e *exception) Logf(f string, a ...interface{}) {
	e.logs = append(e.logs, fmt.Sprintf(f, a...))
}

func (e *exception) Logs() []interface{} {
	return e.logs
}

func (e *exception) String() string {
	var b strings.Builder
	writeException(&b, e)
	return b.String()
}

func TryCatch(body func(), catch func(Exception)) {
	defer func() {
		switch e := recover().(type) {
		case nil:
			return
		case Exception:
			catch(e)
		case ExceptionError:
			catch(e.Unwrap())
		default:
			panic(e)
		}
	}()
	body()
}

func Suppress(ex Exception) {
	p := recover()
	if p == nil {
		return
	}
	if p != ex {
		if e, ok := p.(Exception); ok {
			e.Suppress(ex)
		}
	}
	panic(p)
}

func NewException(err error, logs ...interface{}) Exception {
	if ew, ok := err.(ExceptionError); ok {
		ew := ew.Unwrap()
		for _, l := range logs {
			ew.Log(l)
		}
		return ew
	}
	return &exception{err: err, logs: logs}
}

type MultiDefer []func()

func (md *MultiDefer) Defer(df func()) {
	*md = append(*md, df)
}

func (md *MultiDefer) Run() {
	for _, df := range *md {
		defer df()
	}
}

func Must(err error) {
	if err == nil {
		return
	}
	panic(NewException(err))
}

func writeException(b *strings.Builder, exception Exception) {
	b.WriteString(fmt.Sprintf("%v\n", exception.Err()))
	for _, d := range exception.Logs() {
		b.WriteString(fmt.Sprintf(" - %v\n", d))
	}
	for _, s := range exception.Suppressed() {
		b.WriteString("suppress:\n")
		writeException(b, s)
	}
}
