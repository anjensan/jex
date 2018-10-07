package transform

import (
	"go/ast"
	"runtime"
	"regexp"
)

func ignoreTypeAssertionError() {
	p := recover()
	if p != nil && !isTypeAssertionError(p) {
		panic(p)
	}
}

func isTypeAssertionError(err interface{}) bool {
	_, ok := err.(*runtime.TypeAssertionError)
	return ok
}

func isJexIdent(n ast.Node, name string) bool {
	defer ignoreTypeAssertionError()
	return n.(*ast.Ident).Name == name
}

func isThrowCall(n ast.Node) bool {
	defer ignoreTypeAssertionError()
	fn := n.(*ast.ExprStmt).X.(*ast.CallExpr).Fun
	return isJexIdent(fn, "THROW")
}

func isUnderscoredCall(n ast.Node) bool {
	f, ok := n.(*ast.CallExpr)
	return ok && (isUnderscoredIdent(f.Fun) || isThrowIdent(f.Fun))
}

func isExcCall(n ast.Node) bool {
	defer ignoreTypeAssertionError()
	fn := n.(*ast.CallExpr).Fun
	return isJexIdent(fn, "EX")
}

func isFuncDecl(n ast.Node) bool {
	_, ok := n.(*ast.FuncDecl)
	return ok
}

func isFuncLit(n ast.Node) bool {
	_, ok := n.(*ast.FuncLit)
	return ok
}

func isReturn(n ast.Node) bool {
	_, ok := n.(*ast.ReturnStmt)
	return ok
}

func isDefer(n ast.Node) bool {
	_, ok := n.(*ast.DeferStmt)
	return ok
}

func isAnyOf(ps ...func(ast.Node) bool) func(ast.Node) bool {
	return func(n ast.Node) bool {
		for _, p := range ps {
			if p(n) {
				return true
			}
		}
		return false
	}
}

func isIfTry(n ast.Node) bool {
	defer ignoreTypeAssertionError()
	ifs := n.(*ast.IfStmt)
	fn := ifs.Cond.(*ast.CallExpr).Fun
	return isJexIdent(fn, "TRY")
}

func isUnderscoredIdent(n ast.Node) bool {
	le, ok := n.(*ast.Ident)
	return ok && len(le.Name) > 1 && le.Name[len(le.Name)-1] == '_'
}

func isThrowIdent(n ast.Node) bool {
	le, ok := n.(*ast.Ident)
	return ok && le.Name == "THROW"
}

func isUnderscoredFuncDecl(n ast.Node) bool {
	defer ignoreTypeAssertionError()
	fd := n.(*ast.FuncDecl)
	nm := fd.Name.Name
	return len(nm) > 1 && nm[len(nm)-1] == '_'
}

func isErrAssign(n ast.Node) bool {
	defer ignoreTypeAssertionError()
	a := n.(*ast.AssignStmt)
	le := a.Lhs[len(a.Lhs)-1].(*ast.Ident)
	return le.Name == "ERR"
}

func isErrIdent(n ast.Node) bool {
	le, ok := n.(*ast.Ident)
	return ok && le.Name == "ERR"
}

func isFunctionMarkerCall(s ast.Node, name string) bool {
	defer ignoreTypeAssertionError()
	f := s.(*ast.ExprStmt).X.(*ast.CallExpr).Fun
	return isJexIdent(f, name)
}

var mustRE = regexp.MustCompile("[mM]ust([^a-z]|$)")

func isFuncNamedWithMust(n ast.Node) bool {
	defer ignoreTypeAssertionError()
	fd := n.(*ast.FuncDecl)
	nm := fd.Name.Name
	return mustRE.MatchString(nm)
}
