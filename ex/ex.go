package ex

import (
	"errors"
	"fmt"

	jex "github.com/anjensan/jex/runtime"
)

func Must_(err error, log ...interface{}) {
	if err != nil {
		panic(jex.NewException(err, log...))
	}
}

func Check_(cond bool, err error, log ...interface{}) {
	if !cond {
		panic(jex.NewException(err, log...))
	}
}

func Assert_(cond bool, f string, a ...interface{}) {
	if !cond {
		panic(jex.NewException(errors.New(fmt.Sprintf(f, a...))))
	}
}

func Log(v interface{}) {
	if p := recover(); p != nil {
		if e, ok := p.(jex.Exception); ok {
			e.Log(v)
		}
		panic(p)
	}
}

func Logf(f string, a ...interface{}) {
	if p := recover(); p != nil {
		if e, ok := p.(jex.Exception); ok {
			e.Log(fmt.Sprintf(f, a...))
		}
		panic(p)
	}
}

func Catch(catch func(error)) {
	switch e := recover().(type) {
	case nil:
		return
	case jex.Exception:
		catch(e.Wrap())
	case jex.ExceptionError:
		catch(e)
	default:
		panic(e)
	}
}
