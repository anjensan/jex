//+build jex
//go:generate jex

package demo

import . "github.com/anjensan/jex"

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
)

var duplicateErr = errors.New("duplicateErr")

func qsortEx_(a []int) {
	if len(a) <= 1 {
		return
	}
	left, right := 0, len(a)-1
	pivot := len(a) / 2
	a[pivot], a[right] = a[right], a[pivot]
	for i := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			left++
		} else if i != right && a[i] == a[right] {
			THROW(duplicateErr)
		}
	}
	a[left], a[right] = a[right], a[left]
	qsortEx_(a[:left])
	qsortEx_(a[left+1:])
}

func qsortErr(a []int) error {
	if len(a) <= 1 {
		return nil
	}
	left, right := 0, len(a)-1
	pivot := len(a) / 2
	a[pivot], a[right] = a[right], a[pivot]
	for i := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			left++
		} else if i != right && a[i] == a[right] {
			return duplicateErr
		}
	}
	a[left], a[right] = a[right], a[left]
	if err := qsortErr(a[:left]); err != nil {
		return err
	}
	if err := qsortErr(a[left+1:]); err != nil {
		return err
	}
	return nil
}

func qsortErrFmt(a []int) error {
	if len(a) <= 1 {
		return nil
	}
	left, right := 0, len(a)-1
	pivot := len(a) / 2
	a[pivot], a[right] = a[right], a[pivot]
	for i := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			left++
		} else if i != right && a[i] == a[right] {
			return duplicateErr
		}
	}
	a[left], a[right] = a[right], a[left]
	if err := qsortErrFmt(a[:left]); err != nil {
		return fmt.Errorf("sort %d: %v", len(a), err)
	}
	if err := qsortErrFmt(a[left+1:]); err != nil {
		return fmt.Errorf("sort %d: %v", len(a), err)
	}
	return nil
}

func benchSorting(s []int) func(*testing.B) {
	return func(b *testing.B) {
		b.Run("exception", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shuffle(s)
				if TRY() {
					qsortEx_(s)
				} else {
					fmt.Sprint(EX())
				}
			}
		})
		b.Run("return err", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shuffle(s)
				err := qsortErr(s)
				fmt.Sprint(err)
			}
		})
		b.Run("fmt.errorf", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				shuffle(s)
				err := qsortErrFmt(s)
				fmt.Sprint(err)
			}
		})
	}
}

func shuffle(arr []int) {
	x := rand.Uint32()
	for i := len(arr) - 1; i > 0; i-- {
		j := int(x % uint32(i))
		x = 16807 * x % 2147483647
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func rangeSlice(n int) []int {
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = i
	}
	return r
}

func badSlice(n int) []int {
	r := rangeSlice(n)
	k := n / 3
	r[k+1] = r[k]
	return r
}

func BenchmarkNoErrors(b *testing.B) {
	for _, n := range []int{10, 100, 1000, 10000, 100000, 1000000} {
		b.Run(fmt.Sprintf("%7d", n), benchSorting(rangeSlice(n)))
	}
}

func BenchmarkOneError(b *testing.B) {
		for _, n := range []int{10, 100, 1000, 10000, 100000, 1000000} {
		b.Run(fmt.Sprintf("%7d", n), benchSorting(badSlice(n)))
	}
}
