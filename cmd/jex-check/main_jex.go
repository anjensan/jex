//line main.go:1
//+build !jex
//jex:off

package main

//line main.go:4
import _jex "github.com/anjensan/jex/runtime"

//line main.go:8
import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/anjensan/jex/ex"
	"github.com/anjensan/jex/transform"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: jex-check [path ...]\n")
	flag.PrintDefaults()
}

var fset *token.FileSet

func processFile_(srcfile string) {
	defer ex.Logf("file %v", srcfile)
	file, _jex_e431 := parser.ParseFile(fset, srcfile, nil, parser.ParseComments)
//line main.go:30
	_jex.Must(_jex_e431)
	errs := transform.Check(file, fset)
	if len(errs) > 0 {
		for _, e := range errs {
			p := fset.Position(e.Pos)
			fmt.Printf("%s: %s\n", p, e.Text)
		}
	}
}

func walkDir_(path string) {
	defer ex.Logf("dir %s", path)
//line main.go:41
	_jex_e724 := filepath.Walk(path, visitFile_)
	_jex.Must(_jex_e724)
}

func visitFile_(path string, f os.FileInfo, e error) (err error) {
//line main.go:45
	_jex.TryCatch(func() {
		_jex_e847 := e
		_jex.Must(_jex_e847)
		if isGoFile(f) {
			processFile_(path)
		}
//line main.go:45
		err = nil
//line main.go:51
		return
	}, func(_jex_ex _jex.Exception) {
//line main.go:52
		defer _jex.Suppress(_jex_ex)
//line main.go:45
		err = _jex_ex.
//line main.go:53
			Wrap()
//line main.go:53
		return
	})
//line main.go:54
	return
}

func isGoFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() &&
		!strings.HasPrefix(name, ".") &&
		strings.HasSuffix(name, ".go")
}

func main_() {
	fset = token.NewFileSet()
	if flag.NArg() == 0 {
		walkDir_(".")
	}
	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		dir, _jex_e1242 := os.Stat(path)
//line main.go:71
		_jex.Must(_jex_e1242)
		if dir.IsDir() {
			walkDir_(path)
		} else {
			processFile_(path)
		}
	}
}

func main() {
	flag.Usage = Usage
	flag.Parse()
//line main.go:82
	_jex.TryCatch(func() {
//line main.go:85
		main_()
	}, func(_jex_ex _jex.Exception) {
//line main.go:86
		defer _jex.Suppress(_jex_ex)
		fmt.Fprint(os.Stderr, _jex_ex)
		os.Exit(1)
	})
}

//line main.go:90
const _ = _jex.Unused
