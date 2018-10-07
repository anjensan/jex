//+build jex
//go:generate jex

package must

import . "github.com/anjensan/jex"

func Interface_(v interface{}, err error) interface{} { ERR := err; return v }
func String_(v string, err error) string              { ERR := err; return v }
func Uintptr_(v uintptr, err error) uintptr           { ERR := err; return v }

func Bool_(v bool, err error) bool { ERR := err; return v }
func Byte_(v byte, err error) byte { ERR := err; return v }
func Rune_(v rune, err error) rune { ERR := err; return v }

func Int_(v int, err error) int       { ERR := err; return v }
func Int8_(v int8, err error) int8    { ERR := err; return v }
func Int16_(v int16, err error) int16 { ERR := err; return v }
func Int32_(v int32, err error) int32 { ERR := err; return v }
func Int64_(v int64, err error) int64 { ERR := err; return v }

func Uint_(v uint, err error) uint       { ERR := err; return v }
func Uint8_(v uint8, err error) uint8    { ERR := err; return v }
func Uint16_(v uint16, err error) uint16 { ERR := err; return v }
func Uint32_(v uint32, err error) uint32 { ERR := err; return v }
func Uint64_(v uint64, err error) uint64 { ERR := err; return v }

func Float32_(v float32, err error) float32          { ERR := err; return v }
func Flaot64_(v float64, err error) float64          { ERR := err; return v }
func Complex64_(v complex64, err error) complex64    { ERR := err; return v }
func Complex128_(v complex128, err error) complex128 { ERR := err; return v }
