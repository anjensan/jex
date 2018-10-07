//+build jex
//go:generate jex

package main

import . "github.com/anjensan/jex"

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
	nocheck     = flag.Bool("nocheck", false, `disable checking of uncaught exceptions`)
	nosourcepos = flag.Bool("nosourcepos", false, `emit //line directives to preserve original source positions`)
	printast    = flag.Bool("printast", false, `print generated ast to standard output`)
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

	file, ERR := parser.ParseFile(fset, srcfile, nil, parser.ParseComments)

	var errs []transform.TransformErr
	if *nocheck {
		errs = transform.Check(file, fset)
	}

	dstfile := jexFilename(srcfile)
	tmpdest := "." + dstfile + "~"

	out, ERR := os.Create(tmpdest)
	defer out.Close()

	if len(errs) == 0 {
		errs = transform.Transform(file, fset)
	}

	if len(errs) == 0 {
		c := printer.Config{Mode: printer.RawFormat}
		if !*nosourcepos {
			c.Mode |= printer.SourcePos
		}
		ERR = c.Fprint(out, fset, file)
		ERR = os.Rename(tmpdest, dstfile)
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
	if TRY() {
		main_()
	} else {
		fmt.Fprintln(os.Stderr, EX().Err())
		os.Exit(2)
	}
}
