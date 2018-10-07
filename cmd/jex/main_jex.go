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
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"

	"github.com/anjensan/jex/ex"
	"github.com/anjensan/jex/transform"
)

var (
	nocheck	= flag.Bool("nocheck", false, `disable checking of uncaught exceptions`)
	nosourcepos	= flag.Bool("nosourcepos", false, `emit //line directives to preserve original source positions`)
	printast	= flag.Bool("printast", false, `print generated ast to standard output`)
)

var fset *token.FileSet

func jexFilename(fn string) string {
	fn = strings.TrimSuffix(fn, ".go")
	if strings.HasSuffix(fn, "_test") {
		fn = strings.TrimSuffix(fn, "_test")
		return fn + "_jex_test.go"
	} else {
		return fn + "_jex.go"
	}
}

func processFile_(srcfile string) []transform.TransformErr {
	defer ex.Logf("file %s", srcfile)

	file, _jex_e877 := parser.ParseFile(fset, srcfile, nil, parser.ParseComments)
//line main.go:43
	_jex.Must(_jex_e877)

	var errs []transform.TransformErr
	if *nocheck {
		errs = transform.Check(file, fset)
	}

	dstfile := jexFilename(srcfile)
	tmpdest := "." + dstfile + "~"

	out, _jex_e1108 := os.Create(tmpdest)
//line main.go:53
	_jex.Must(_jex_e1108)
	defer out.Close()

	if len(errs) == 0 {
		errs = transform.Transform(file, fset)
	}

	if len(errs) == 0 {
		c := printer.Config{Mode: printer.RawFormat}
		if !*nosourcepos {
			c.Mode |= printer.SourcePos
		}
//line main.go:64
		var _jex_e1351 error
//line main.go:64
		_jex_e1351 = c.Fprint(out, fset, file)
		_jex.Must(_jex_e1351)
//line main.go:65
		var _jex_e1385 error
//line main.go:65
		_jex_e1385 = os.Rename(tmpdest, dstfile)
		_jex.Must(_jex_e1385)
	}
	if *printast {
		ast.Print(fset, file)
	}

	return errs
}

func main_() {
	srcfile := flag.Arg(0)
	if flag.NArg() == 0 {
		srcfile = os.Getenv("GOFILE")
	}

	ex.Assert_(srcfile != "", "unspecified source file")
	fset = token.NewFileSet()
	errs := processFile_(srcfile)

	if len(errs) > 0 {
		for _, e := range errs {
			p := fset.Position(e.Pos)
			fmt.Printf("%s: %s\n", p, e.Text)
		}
		os.Exit(1)
	}
}

func main() {
	flag.Parse()
//line main.go:95
	_jex.TryCatch(func() {
//line main.go:97
		main_()
	}, func(_jex_ex _jex.Exception) {
//line main.go:98
		defer _jex.Suppress(_jex_ex)
		fmt.Fprintln(os.Stderr, _jex_ex.Err())
		os.Exit(2)
	})
}

//line main.go:102
const _ = _jex.Unused
