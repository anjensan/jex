//line must.go:1
//+build !jex
//jex:off

package must

//line must.go:4
import _jex "github.com/anjensan/jex/runtime"

//line must.go:8
func Interface_(v interface{}, err error) interface{} {
//line must.go:8
	_jex_e139 := err
//line must.go:8
	_jex.Must(_jex_e139)
//line must.go:8
	return v
//line must.go:8
}
func String_(v string, err error) string	{ _jex_e218 := err; _jex.Must(_jex_e218); return v }
func Uintptr_(v uintptr, err error) uintptr	{ _jex_e297 := err; _jex.Must(_jex_e297); return v }

func Bool_(v bool, err error) bool	{ _jex_e358 := err; _jex.Must(_jex_e358); return v }
func Byte_(v byte, err error) byte	{ _jex_e418 := err; _jex.Must(_jex_e418); return v }
func Rune_(v rune, err error) rune	{ _jex_e478 := err; _jex.Must(_jex_e478); return v }

func Int_(v int, err error) int	{ _jex_e542 := err; _jex.Must(_jex_e542); return v }
func Int8_(v int8, err error) int8	{ _jex_e605 := err; _jex.Must(_jex_e605); return v }
func Int16_(v int16, err error) int16	{ _jex_e668 := err; _jex.Must(_jex_e668); return v }
func Int32_(v int32, err error) int32	{ _jex_e731 := err; _jex.Must(_jex_e731); return v }
func Int64_(v int64, err error) int64	{ _jex_e794 := err; _jex.Must(_jex_e794); return v }

func Uint_(v uint, err error) uint	{ _jex_e861 := err; _jex.Must(_jex_e861); return v }
func Uint8_(v uint8, err error) uint8	{ _jex_e927 := err; _jex.Must(_jex_e927); return v }
func Uint16_(v uint16, err error) uint16	{ _jex_e993 := err; _jex.Must(_jex_e993); return v }
func Uint32_(v uint32, err error) uint32	{ _jex_e1059 := err; _jex.Must(_jex_e1059); return v }
func Uint64_(v uint64, err error) uint64	{ _jex_e1125 := err; _jex.Must(_jex_e1125); return v }

func Float32_(v float32, err error) float32	{ _jex_e1204 := err; _jex.Must(_jex_e1204); return v }
func Flaot64_(v float64, err error) float64	{ _jex_e1282 := err; _jex.Must(_jex_e1282); return v }
func Complex64_(v complex64, err error) complex64	{ _jex_e1360 := err; _jex.Must(_jex_e1360); return v }
func Complex128_(v complex128, err error) complex128 {
//line must.go:31
	_jex_e1438 := err
//line must.go:31
	_jex.Must(_jex_e1438)
//line must.go:31
	return v
//line must.go:31
}

//line must.go:31
const _ = _jex.Unused
