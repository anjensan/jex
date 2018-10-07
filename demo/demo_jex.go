//line demo.go:1
//+build !jex
//jex:off

package demo

//line demo.go:4
import _jex "github.com/anjensan/jex/runtime"

//line demo.go:8
import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/anjensan/jex/ex"
)

var zeroError = errors.New("zero error")

//line demo.go:19
// All "exceptional" functions must ends with "_"
func checkZero_() {

	x := rand.Intn(5)
	switch x {
	case 0:
		return
	case 1:
//line demo.go:26
		panic(
		// Just throw
//line demo.go:26
		_jex.NewException(zeroError))
//line demo.go:29
	case 2:
//line demo.go:29
		panic(
		// Throw & attach comment (not part of the error!)
//line demo.go:29
		_jex.NewException(zeroError, "with comment"))
//line demo.go:32
	case 3:
//line demo.go:32
		panic(
		// Create new instance of error & throw
//line demo.go:32
		_jex.NewException(errors.New("ZERO error")))
//line demo.go:35
	case 4:
//line demo.go:35
		var _jex_e606 error
		// Special macro, expands to call of `ex.Must_`
//line demo.go:35
		_jex_e606 = zeroError
//line demo.go:37
		_jex.Must(_jex_e606)
	default:
//line demo.go:38
		panic(_jex.NewException(zeroError, fmt.Sprintf("n = %d", x)))
//line demo.go:40
	}
}

func catchErr() {
//line demo.go:43
	_jex.
	// Special macro to catch exceptions
//line demo.go:43
	TryCatch(func() {
//line demo.go:46
		checkZero_()
	}, func(_jex_ex _jex.
	// EX() returns current exception
	// EX().Err() - original error
	// EX().Wrap() - modified error (original error + debug messages)
	// EX().Logs() - comments, logs etc
	// EX().Suppressed() - list of suppressed errors
//line demo.go:47
	Exception) {
//line demo.go:47
		defer _jex.Suppress(_jex_ex)

//line demo.go:53
		fmt.Println(_jex_ex.Err())
	})
}

func rethrowError_() {
//line demo.go:57
	_jex.TryCatch(func() {
//line demo.go:59
		checkZero_()
	}, func(_jex_ex _jex.
	// Rethrow error
//line demo.go:60
	Exception) {
//line demo.go:60
		defer _jex.Suppress(_jex_ex)
//line demo.go:60
		panic(_jex_ex)

//line demo.go:63
	})
}

func rethrowNewError_() {
//line demo.go:66
	_jex.TryCatch(func() {
//line demo.go:68
		checkZero_()
	}, func(_jex_ex _jex.
	// Rethrow error
//line demo.go:69
	Exception) {
//line demo.go:69
		defer _jex.Suppress(_jex_ex)
//line demo.go:69
		panic(_jex.NewException(errors.New("another error")))
//line demo.go:72
	})
}

func fail_() {
//line demo.go:75
	panic(_jex.NewException(errors.New("fail")))
//line demo.go:77
}

func suppressError_() {
//line demo.go:79
	_jex.TryCatch(func() {
//line demo.go:81
		checkZero_()
	}, func(_jex_ex _jex.
	// Implicitly keep original exception as `suppressed`
//line demo.go:82
	Exception) {
//line demo.go:82
		defer _jex.Suppress(_jex_ex)

		fail_()
	})
}

func causeOfException() {
//line demo.go:88
	_jex.TryCatch(func() {
//line demo.go:90
		suppressError_()
	}, func(_jex_ex _jex.Exception) {
//line demo.go:91
		defer _jex.Suppress(_jex_ex)
		fmt.Println(_jex_ex.Suppressed())
	})
}

func supportsDefer() {
	// Defers are supported.
	// However `recover` is NOT supported!

	defer fmt.Println("before try")
//line demo.go:100
	var _jex_md1662 _jex.MultiDefer
//line demo.go:100
	defer _jex_md1662.Run()
//line demo.go:100
	_jex.TryCatch(func() {
		{
//line demo.go:101
			_f, _p0 := fmt.Println, "before check"
			_jex_md1662.Defer(func() {
//line demo.go:102
				_f(_p0)
//line demo.go:102
			})
//line demo.go:102
		}
		checkZero_()
//line demo.go:103
		{
//line demo.go:103
			_f, _p0 := fmt.Println, "after check"
			_jex_md1662.Defer(func() {
//line demo.go:104
				_f(_p0)
//line demo.go:104
			})
//line demo.go:104
		}
	}, func(_jex_ex _jex.Exception) {
//line demo.go:105
		defer _jex.Suppress(_jex_ex)
//line demo.go:105
		{
//line demo.go:105
			_f, _p0 := fmt.Println, "handle err"
			_jex_md1662.Defer(func() {
//line demo.go:106
				_f(_p0)
//line demo.go:106
			})
//line demo.go:106
		}
	})
	defer fmt.Println("after try")
}

func supportsReturn() (_jex_r0 int) {
//line demo.go:111
	_jex.
	// Returns are supported.
//line demo.go:111
	TryCatch(func() {
//line demo.go:114
		checkZero_()
//line demo.go:114
		_jex_r0 = 1
		return
	}, func(_jex_ex _jex.Exception) {
//line demo.go:116
		defer _jex.Suppress(_jex_ex)
//line demo.go:116
		_jex_r0 = 0
		return
	})
//line demo.go:118
	return
}

func addDebugInfo_() {
	// Special helper to attach debug info to exception.
	defer ex.Log("call addDebugInfo")
	checkZero_()
}

func checkZeroLegacy() (err error) {
//line demo.go:127
	_jex.
	// How to adapt exceptional functions to old api
//line demo.go:127
	TryCatch(func() {
//line demo.go:130
		checkZero_()
//line demo.go:127
		err = nil
//line demo.go:131
		return
	}, func(_jex_ex _jex.Exception) {
//line demo.go:132
		defer _jex.Suppress(_jex_ex)
//line demo.go:127
		err = _jex_ex.
//line demo.go:133
			Wrap()
//line demo.go:133
		return
	})
//line demo.go:134
	return
}

func checkZeroAsync_() {
	// Catch & rethrow exceptions from goroutine
	e := make(chan error)
	go func() {
		defer close(e)
//line demo.go:141
		_jex.TryCatch(func() {
//line demo.go:143
			checkZero_()
		}, func(_jex_ex _jex.Exception) {
//line demo.go:144
			defer _jex.Suppress(_jex_ex)
			e <- _jex_ex.Wrap()
		})
	}()
	// Read wrapped exception from channel, maybe unwrap & rethow it
	ex.Must_(<-e)
}

//line demo.go:150
const _ = _jex.Unused
