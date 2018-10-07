//+build jex
//go:generate jex

package main

import . "github.com/anjensan/jex"

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
	file, ERR := parser.ParseFile(fset, srcfile, nil, parser.ParseComments)
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
	ERR := filepath.Walk(path, visitFile_)
}

func visitFile_(path string, f os.FileInfo, e error) (err error) {
	if TRY() {
		ERR := e
		if isGoFile(f) {
			processFile_(path)
		}
		return nil
	} else {
		return EX().Wrap()
	}
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
		dir, ERR := os.Stat(path)
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

	if TRY() {
		main_()
	} else {
		fmt.Fprint(os.Stderr, EX())
		os.Exit(1)
	}
}
