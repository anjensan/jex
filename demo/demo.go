//+build jex
//go:generate jex

package demo

import . "github.com/anjensan/jex"

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/anjensan/jex/ex"
)

var zeroError = errors.New("zero error")


// All "exceptional" functions must ends with "_"
func checkZero_() {

	x := rand.Intn(5)
	switch x {
	case 0:
		return
	case 1:
		// Just throw
		THROW(zeroError)
	case 2:
		// Throw & attach comment (not part of the error!)
		THROW(zeroError, "with comment")
	case 3:
		// Create new instance of error & throw
		THROW(errors.New("ZERO error"))
	case 4:
		// Special macro, expands to call of `ex.Must_`
		ERR = zeroError
	default:
		THROW(zeroError, fmt.Sprintf("n = %d", x))
	}
}

func catchErr() {
	// Special macro to catch exceptions
	if TRY() {
		checkZero_()
	} else {
		// EX() returns current exception
		// EX().Err() - original error
		// EX().Wrap() - modified error (original error + debug messages)
		// EX().Logs() - comments, logs etc
		// EX().Suppressed() - list of suppressed errors
		fmt.Println(EX().Err())
	}
}

func rethrowError_() {
	if TRY() {
		checkZero_()
	} else {
		// Rethrow error
		THROW()
	}
}

func rethrowNewError_() {
	if TRY() {
		checkZero_()
	} else {
		// Rethrow error
		THROW(errors.New("another error"))
	}
}

func fail_() {
	THROW(errors.New("fail"))
}

func suppressError_() {
	if TRY() {
		checkZero_()
	} else {
		// Implicitly keep original exception as `suppressed`
		fail_()
	}
}

func causeOfException() {
	if TRY() {
		suppressError_()
	} else {
		fmt.Println(EX().Suppressed())
	}
}

func supportsDefer() {
	// Defers are supported.
	// However `recover` is NOT supported!

	defer fmt.Println("before try")
	if TRY() {
		defer fmt.Println("before check")
		checkZero_()
		defer fmt.Println("after check")
	} else {
		defer fmt.Println("handle err")
	}
	defer fmt.Println("after try")
}

func supportsReturn() int {
	// Returns are supported.
	if TRY() {
		checkZero_()
		return 1
	} else {
		return 0
	}
}

func addDebugInfo_() {
	// Special helper to attach debug info to exception.
	defer ex.Log("call addDebugInfo")
	checkZero_()
}

func checkZeroLegacy() (err error) {
	// How to adapt exceptional functions to old api
	if TRY() {
		checkZero_()
		return nil
	} else {
		return EX().Wrap()
	}
}

func checkZeroAsync_() {
	// Catch & rethrow exceptions from goroutine
	e := make(chan error)
	go func() {
		defer close(e)
		if TRY() {
			checkZero_()
		} else {
			e <- EX().Wrap()
		}
	}()
	// Read wrapped exception from channel, maybe unwrap & rethow it
	ex.Must_(<-e)
}
