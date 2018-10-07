package transform

// This code has been adapted from go/types/return.go.

import (
	"go/ast"
	"go/token"
)

// unparen returns e with any enclosing parentheses stripped.
func unparen(e ast.Expr) ast.Expr {
	for {
		p, ok := e.(*ast.ParenExpr)
		if !ok {
			return e
		}
		e = p.X
	}
}

func isTerminating(s ast.Stmt, label string) bool {
	switch s := s.(type) {

	case *ast.BadStmt, *ast.DeclStmt, *ast.EmptyStmt, *ast.SendStmt,
		*ast.IncDecStmt, *ast.AssignStmt, *ast.GoStmt, *ast.DeferStmt,
		*ast.RangeStmt:
		// no chance

	case *ast.LabeledStmt:
		return isTerminating(s.Stmt, s.Label.Name)

	case *ast.ExprStmt:
		// the predeclared (possibly parenthesized) panic() function is terminating
		if call, _ := unparen(s.X).(*ast.CallExpr); call != nil {
			if id, _ := call.Fun.(*ast.Ident); id != nil {
				// TODO: Respect scope & lookups.
				// if _, obj := scope.LookupParent(id.Name, token.NoPos); obj != nil {
				// 	if b, _ := obj.(*Builtin); b != nil && b.id == _Panic {
				// 		return true
				// 	}
				// }
				return id.Name == "panic"
			}
		}

	case *ast.ReturnStmt:
		return true

	case *ast.BranchStmt:
		if s.Tok == token.GOTO || s.Tok == token.FALLTHROUGH {
			return true
		}

	case *ast.BlockStmt:
		return isTerminatingList(s.List, "")

	case *ast.IfStmt:
		if s.Else != nil &&
			isTerminating(s.Body, "") &&
			isTerminating(s.Else, "") {
			return true
		}

	case *ast.SwitchStmt:
		return isTerminatingSwitch(s.Body, label)

	case *ast.TypeSwitchStmt:
		return isTerminatingSwitch(s.Body, label)

	case *ast.SelectStmt:
		for _, s := range s.Body.List {
			cc := s.(*ast.CommClause)
			if !isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
				return false
			}

		}
		return true

	case *ast.ForStmt:
		if s.Cond == nil && !hasBreak(s.Body, label, true) {
			return true
		}
	}
	return false
}

func isTerminatingList(list []ast.Stmt, label string) bool {
	// trailing empty statements are permitted - skip them
	for i := len(list) - 1; i >= 0; i-- {
		if _, ok := list[i].(*ast.EmptyStmt); !ok {
			return isTerminating(list[i], label)
		}
	}
	return false // all statements are empty
}

func isTerminatingSwitch(body *ast.BlockStmt, label string) bool {
	hasDefault := false
	for _, s := range body.List {
		cc := s.(*ast.CaseClause)
		if cc.List == nil {
			hasDefault = true
		}
		if !isTerminatingList(cc.Body, "") || hasBreakList(cc.Body, label, true) {
			return false
		}
	}
	return hasDefault
}

// hasBreak reports if s is or contains a break statement
// referring to the label-ed statement or implicit-ly the
// closest outer breakable statement.
func hasBreak(s ast.Stmt, label string, implicit bool) bool {
	switch s := s.(type) {

	case *ast.BadStmt, *ast.DeclStmt, *ast.EmptyStmt, *ast.ExprStmt,
		*ast.SendStmt, *ast.IncDecStmt, *ast.AssignStmt, *ast.GoStmt,
		*ast.DeferStmt, *ast.ReturnStmt:
		// no chance

	case *ast.LabeledStmt:
		return hasBreak(s.Stmt, label, implicit)

	case *ast.BranchStmt:
		if s.Tok == token.BREAK {
			if s.Label == nil {
				return implicit
			}
			if s.Label.Name == label {
				return true
			}
		}
	case *ast.BlockStmt:
		return hasBreakList(s.List, label, implicit)

	case *ast.IfStmt:
		if hasBreak(s.Body, label, implicit) ||
			s.Else != nil && hasBreak(s.Else, label, implicit) {
			return true
		}
	case *ast.CaseClause:
		return hasBreakList(s.Body, label, implicit)

	case *ast.SwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.TypeSwitchStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.CommClause:
		return hasBreakList(s.Body, label, implicit)

	case *ast.SelectStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.ForStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}

	case *ast.RangeStmt:
		if label != "" && hasBreak(s.Body, label, false) {
			return true
		}
	}

	return false
}

func hasBreakList(list []ast.Stmt, label string, implicit bool) bool {
	for _, s := range list {
		if hasBreak(s, label, implicit) {
			return true
		}
	}
	return false
}
