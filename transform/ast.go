package transform

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type pred func(ast.Node) bool

type transformer struct {
	file *ast.File
	cmap ast.CommentMap
	errs []TransformErr
}

type TCursor struct {
	*astutil.Cursor
	t *transformer
}

func (c *TCursor) Replace(n ast.Node) {
	c.t.cmap.Update(c.Node(), n)
	c.Cursor.Replace(n)
}

func (c *TCursor) ForEach(p pred, s pred, f func(*TCursor)) {
	foreachNode0(c.t, c.Node(), p, s, f)
}

func (c *TCursor) ForEachSub(n ast.Node, p pred, s pred, f func(*TCursor)) {
	foreachNode0(c.t, n, p, s, f)
}

func (c *TCursor) Error(format string, a ...interface{}) {
	c.t.Error(c.Node(), fmt.Sprintf(format, a...))
}

func (t *transformer) ForEach(p pred, s pred, f func(*TCursor)) {
	foreachNode0(t, t.file, p, s, f)
}

func (t *transformer) Error(n ast.Node, Text string) {
	t.errs = append(t.errs, TransformErr{Text, n.Pos()})
}

func (t *transformer) JexComments(n ast.Node) (res []string) {
	for _, cg := range t.cmap[n] {
		for _, c := range cg.List {
			t := c.Text
			t = strings.TrimPrefix(t, "//")
			t = strings.Trim(t, " ")
			if strings.HasPrefix(t, "jex:") {
				res = append(res, strings.TrimPrefix(t, "jex:"))
			}
		}
	}
	return
}

func (t *transformer) HasJexTag(n ast.Node, tag string) bool {
	cs := t.JexComments(n)
	for _, c := range cs {
		if tag == c {
			return true
		}
	}
	return false
}

func stringsContains(c []string, t string) bool {
	return false
}

func ignoredNode(t *transformer, n ast.Node) bool {
	return t.HasJexTag(n, "off")
}

func foreachNode0(t *transformer, r ast.Node, p pred, s pred, f func(*TCursor)) {
	astutil.Apply(
		r,
		func(c *astutil.Cursor) bool {
			n := c.Node()
			if n == nil {
				return false
			}
			if ignoredNode(t, n) {
				return false
			}
			if n != r && s != nil && s(n) {
				if p(n) {
					f(&TCursor{c, t})
				}
				return false
			}
			return true
		},
		func(c *astutil.Cursor) bool {
			if p(c.Node()) {
				f(&TCursor{c, t})
			}
			return true
		})
}

func DeferStmt(f ast.Expr, args ...ast.Expr) *ast.DeferStmt {
	return &ast.DeferStmt{Call: CallExpr(f, args...)}
}

func Ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func QIdent(x, name string) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: Ident(x), Sel: Ident(name)}
}

func CallExpr(f ast.Expr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: f, Args: args}
}

func CallStmt(f ast.Expr, args ...ast.Expr) ast.Stmt {
	return &ast.ExprStmt{X: CallExpr(f, args...)}
}

func AsBlockStmt(s ast.Stmt) *ast.BlockStmt {
	if bs, ok := s.(*ast.BlockStmt); ok {
		return bs
	} else {
		return &ast.BlockStmt{List: []ast.Stmt{s}}
	}
}

func BlockStmt(s ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: s}
}

func Field(name string, t ast.Expr) *ast.Field {
	return &ast.Field{Type: t, Names: []*ast.Ident{Ident(name)}}
}

func FieldsList(fs ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: fs}
}

func FuncLit(params, results *ast.FieldList, body ast.Stmt) *ast.FuncLit {
	return &ast.FuncLit{
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
		Body: AsBlockStmt(body),
	}
}
