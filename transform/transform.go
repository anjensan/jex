package transform

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type TransformErr struct {
	Text string
	Pos  token.Pos
}

func updateJexImport(root *transformer) {
	// Remove '.' import
	for i, d := range root.file.Decls {
		dd, ok := d.(*ast.GenDecl)
		if !ok || dd.Tok != token.IMPORT {
			continue
		}
		specs := make([]ast.Spec, 0, len(dd.Specs))
		for _, s := range dd.Specs {
			ss, ok := s.(*ast.ImportSpec)
			if !ok {
				continue
			}
			if ss.Name == nil || ss.Name.Name != "." || ss.Path.Value != `"github.com/anjensan/jex"` {
				specs = append(specs, ss)
			}
		}

		dd.Specs = specs
		if len(dd.Specs) == 0 {
			root.file.Decls = append(root.file.Decls[:i], root.file.Decls[i+1:]...)
		}
	}

	// Add '_jex' import
	d := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Name: &ast.Ident{Name: "_jex"},
				Path: &ast.BasicLit{Value: `"github.com/anjensan/jex/runtime"`},
			}}}
	root.file.Decls = append([]ast.Decl{d}, root.file.Decls...)
}

func addUnusedConst(root *transformer) {
	// Add '_jex' import
	c := &ast.GenDecl{
		Tok: token.CONST,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{Ident("_")},
				Values: []ast.Expr{QIdent("_jex", "Unused")},
			},
		},
	}
	root.file.Decls = append(root.file.Decls, c)
}

func transformBuildComments(root *transformer) {
	for _, cg := range root.file.Comments {
		for _, c := range cg.List {
			t := c.Text
			if strings.HasPrefix(t, "//+build jex") {
				c.Text = "//+build !jex"
				continue
			}
			if strings.HasPrefix(t, "//go:generate jex") {
				c.Text = "//jex:off"
				continue
			}
		}
	}
}

func expandMacro_THROW(root *transformer) {
	root.ForEach(
		isThrowCall,
		nil,
		func(c *TCursor) {
			c.Replace(
				CallStmt(
					Ident("panic"),
					CallExpr(
						QIdent("_jex", "NewException"),
						c.Node().(*ast.ExprStmt).X.(*ast.CallExpr).Args...)))
		})
}

func getFuncBody(n ast.Node) *ast.BlockStmt {
	switch f := n.(type) {
	case *ast.FuncDecl:
		return f.Body
	case *ast.FuncLit:
		return f.Body
	}
	return nil
}

func getFuncType(n ast.Node) *ast.FuncType {
	switch f := n.(type) {
	case *ast.FuncDecl:
		return f.Type
	case *ast.FuncLit:
		return f.Type
	}
	return nil
}

func checkNoUnderscoredCalls(root *transformer) {
	ignored := func(n ast.Node) bool {
		return false ||
			root.HasJexTag(n, "nocheck") ||
			isUnderscoredFuncDecl(n) ||
			isFuncNamedWithMust(n)
	}
	root.ForEach(
		isAnyOf(isFuncDecl, isFuncLit),
		ignored,
		func(c *TCursor) {
			var check func(n ast.Node)
			check = func(n ast.Node) {
				if ignored(n) {
					return
				}
				c.ForEachSub(
					n,
					isAnyOf(isUnderscoredCall, isIfTry, isErrAssign),
					isAnyOf(isFuncLit, isIfTry, ignored),
					func(n *TCursor) {
						if isIfTry(n.Node()) {
							e := n.Node().(*ast.IfStmt).Else
							if e != nil {
								check(e)
							}
							return
						}
						if isUnderscoredCall(n.Node()) || isErrAssign(n.Node()) {
							n.Error("uncaught exception")
							return
						}
					})
			}
			check(c.Node())
		})
}

func eliminateDefers(root *transformer) {
	root.ForEach(
		isAnyOf(isFuncDecl, isFuncLit),
		nil,
		func(f *TCursor) {
			f.ForEach(
				isIfTry,
				isFuncLit,
				func(c *TCursor) {
					hasdefer := false
					c.ForEach(
						isDefer,
						isFuncLit,
						func(n *TCursor) { hasdefer = true })
					if !hasdefer {
						return
					}

					md := fmt.Sprintf("_jex_md%d", c.Node().Pos())

					s := &ast.ValueSpec{Names: []*ast.Ident{Ident(md)}, Type: QIdent("_jex", "MultiDefer")}
					dcl := &ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{s}}}
					c.InsertBefore(dcl)
					c.InsertBefore(DeferStmt(QIdent(md, "Run")))

					c.ForEach(
						isDefer,
						isFuncLit,
						func(c *TCursor) {
							if !isDefer(c.Node()) {
								return
							}
							// TODO: Check for 'recover' calls.

							dd := c.Node().(*ast.DeferStmt)
							fa := &ast.AssignStmt{Tok: token.DEFINE}
							fa.Lhs = append(fa.Lhs, Ident("_f"))
							fa.Rhs = append(fa.Rhs, dd.Call.Fun)

							ps := []ast.Expr{}
							for i, a := range dd.Call.Args {
								p := Ident(fmt.Sprintf("_p%d", i))
								fa.Lhs = append(fa.Lhs, p)
								fa.Rhs = append(fa.Rhs, a)
								ps = append(ps, p)
							}
							cs := CallStmt(
								QIdent(md, "Defer"),
								&ast.FuncLit{
									Type: &ast.FuncType{Params: &ast.FieldList{}},
									Body: BlockStmt(CallStmt(Ident("_f"), ps...))})
							c.Replace(BlockStmt(fa, cs))
						})
				})
		})
}

func eliminateReturns(root *transformer) {
	root.ForEach(
		isAnyOf(isFuncDecl, isFuncLit),
		nil,
		func(f *TCursor) {
			hasreturn := false
			f.ForEach(
				isIfTry,
				isFuncLit,
				func(it *TCursor) {
					it.ForEach(
						isReturn,
						isFuncLit,
						func(n *TCursor) { hasreturn = true })
				})
			if !hasreturn {
				return
			}

			// Generate vars  _jex_r0, _jex_r1...
			res := getFuncType(f.Node()).Results
			body := getFuncBody(f.Node())

			acount := 0
			needjexret := false
			if res != nil {
				acount = len(res.List)
				for i, r := range res.List {
					if len(r.Names) == 0 {
						r.Names = []*ast.Ident{Ident(fmt.Sprintf("_jex_r%d", i))}
					}
				}
			}

			// Add 'if _jex_ret { return }' after all if-try's
			f.ForEach(
				isIfTry,
				isFuncLit,
				func(c *TCursor) {
					ff := c.Node().(*ast.IfStmt)
					ret := &ast.ReturnStmt{}
					if isTerminating(ff, "") {
						c.InsertAfter(ret)
					} else {
						needjexret = true
						c.InsertAfter(
							&ast.IfStmt{
								Cond: Ident("_jex_ret"),
								Body: BlockStmt(ret),
							})
					}
				})

			// Maybe add _jex_ret
			if needjexret {
				s := &ast.ValueSpec{Names: []*ast.Ident{Ident("_jex_ret")}, Type: Ident("bool")}
				dcl := &ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{s}}}
				body.List = append([]ast.Stmt{dcl}, body.List...)
			}

			// Replace 'return'
			f.ForEach(
				isIfTry,
				isFuncLit,
				func(t *TCursor) {
					t.ForEach(
						isReturn,
						isFuncLit,
						func(c *TCursor) {
							rs := c.Node().(*ast.ReturnStmt)
							ra := &ast.AssignStmt{Tok: token.ASSIGN}
							if needjexret {
								ra.Lhs = append(ra.Lhs, Ident("_jex_ret"))
								ra.Rhs = append(ra.Rhs, Ident("true"))
							}
							if len(rs.Results) > 0 {
								if len(rs.Results) != acount {
									c.Error("wrong number of returned args, got %d, expected %d", len(rs.Results), acount)
									return
								}
								for i := 0; i < acount; i++ {
									ra.Lhs = append(ra.Lhs, res.List[i].Names[0])
									ra.Rhs = append(ra.Rhs, rs.Results[i])
								}
							}
							c.Replace(ra)
							c.InsertAfter(&ast.ReturnStmt{})
						})
				})
		})
}

func expandMacro_CATCH_EX(root *transformer) {
	root.ForEach(
		isIfTry,
		nil,
		func(c *TCursor) {
			it := c.Node().(*ast.IfStmt)
			if it.Else == nil {
				return
			}
			c.ForEachSub(
				it.Else,
				isExcCall,
				isFuncLit,
				func(c *TCursor) {
					args := c.Node().(*ast.CallExpr).Args
					if len(args) != 0 {
						c.Error("invalid number of args passed to EX()")
						return
					}
					c.Replace(Ident("_jex_ex"))
				})
		})
}

func expandMacro_CATCH_THROW(root *transformer) {
	root.ForEach(
		isIfTry,
		nil,
		func(c *TCursor) {
			it := c.Node().(*ast.IfStmt)
			if it.Else == nil {
				return
			}
			c.ForEachSub(
				it.Else,
				isThrowCall,
				isFuncLit,
				func(c *TCursor) {
					args := c.Node().(*ast.ExprStmt).X.(*ast.CallExpr).Args
					if len(args) == 0 {
						c.Replace(
							CallStmt(
								Ident("panic"),
								Ident("_jex_ex")))
					}
				})
		})
}

func expandMacro_TRY(root *transformer) {
	root.ForEach(
		isIfTry,
		nil,
		func(c *TCursor) {
			ifs := c.Node().(*ast.IfStmt)
			if ifs.Init != nil {
				c.Error("initialization blocks are not supported in `if TRY()`")
				return
			}
			if ifs.Else == nil {
				c.Error("missing `else` branch in `if TRY()`")
				return
			}

			els := AsBlockStmt(ifs.Else)
			els.List = append([]ast.Stmt{
				DeferStmt(QIdent("_jex", "Suppress"), Ident("_jex_ex")),
			}, els.List...)

			c.Replace(
				CallStmt(
					QIdent("_jex", "TryCatch"),
					FuncLit(
						FieldsList(),
						FieldsList(),
						ifs.Body),
					FuncLit(
						FieldsList(
							Field("_jex_ex", QIdent("_jex", "Exception"))),
						FieldsList(),
						els)))
		})
}

func expandMacro_ERR(root *transformer) {
	root.ForEach(
		isErrAssign,
		nil,
		func(c *TCursor) {
			en := fmt.Sprintf("_jex_e%d", c.Node().Pos())
			// replace `ERR` with `err%d`
			c.ForEach(
				isErrIdent,
				nil,
				func(c *TCursor) {
					c.Replace(Ident(en))
				})

			// prepend `var err%d error`
			a := c.Node().(*ast.AssignStmt)
			if a.Tok == token.ASSIGN {
				c.InsertBefore(&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{Ident(en)},
								Type:  Ident("error"),
							},
						},
					},
				})
			}

			// add `_jex.Must(err%d)`
			c.InsertAfter(
				CallStmt(QIdent("_jex", "Must"), Ident(en)))
		})
}

func check(r *transformer) {
	checkNoUnderscoredCalls(r)
}

func transform(r *transformer) {
	updateJexImport(r)
	expandMacro_CATCH_THROW(r)
	expandMacro_CATCH_EX(r)
	expandMacro_THROW(r)
	expandMacro_ERR(r)
	eliminateDefers(r)
	eliminateReturns(r)
	expandMacro_TRY(r)
	transformBuildComments(r)
	addUnusedConst(r)
}

func Transform(file *ast.File, fset *token.FileSet) []TransformErr {
	cmap := ast.NewCommentMap(fset, file, file.Comments)
	t := &transformer{file: file, cmap: cmap}
	transform(t)
	file.Comments = cmap.Filter(file).Comments()
	return t.errs
}

func Check(file *ast.File, fset *token.FileSet) []TransformErr {
	cmap := ast.NewCommentMap(fset, file, file.Comments)
	t := &transformer{file: file, cmap: cmap}
	check(t)
	return t.errs
}
