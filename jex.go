package jex

import r "github.com/anjensan/jex/runtime"

func THROW(args ...interface{}) {
	panic("Unsubstituted macro call!")
}

func TRY() bool {
	panic("Unsubstituted macro call!")
}

func EX() r.Exception {
	panic("Unsubstituted macro call!")
}
