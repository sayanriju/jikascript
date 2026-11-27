package cgen

import "jikascript/internal/ast"

// arrayUsage records how a top-level const array is used. It is produced once
// per array during the code-generation pre-pass.
type arrayUsage struct {
	// mutates is set when the array (or any element/property of it) may be
	// written, or when the array is referenced in a way the analysis cannot
	// prove read-only.
	mutates bool
	// indexed records an indexed read, which needs the bounds-checked getter.
	indexed bool
	// peek and indexOf record read-only method calls so only the helpers the
	// program actually uses are emitted.
	peek    bool
	indexOf bool
}

// analyzeArrayUsage walks the whole program looking for any use of the array
// that could write to it. Only a conservative whitelist of read-only uses is
// accepted; anything unrecognized marks the array as mutated so it stays in
// RAM with the full mutable helper set.
func analyzeArrayUsage(program *ast.Program, node *ast.ConstDecl) arrayUsage {
	usage := arrayUsage{}
	consumed := map[ast.Expr]bool{}
	handled := map[ast.Expr]bool{}

	isBase := func(expr ast.Expr) bool {
		id, ok := expr.(*ast.Identifier)
		if !ok {
			return false
		}
		if node.Symbol != "" && id.Symbol != "" {
			return id.Symbol == node.Symbol
		}
		return id.Name == node.Name
	}

	visit := func(expr ast.Expr) {
		if usage.mutates {
			return
		}
		switch e := expr.(type) {
		case *ast.CallExpr:
			selector, ok := e.Function.(*ast.SelectorExpr)
			if !ok || !isBase(selector.Object) {
				return
			}
			handled[selector] = true
			consumed[selector.Object] = true
			switch selector.Field {
			case "peek":
				usage.peek = true
			case "indexOf":
				usage.indexOf = true
			default:
				usage.mutates = true
			}
		case *ast.SelectorExpr:
			if handled[e] {
				return
			}
			if isBase(e.Object) {
				consumed[e.Object] = true
				switch e.Field {
				case "length", "capacity":
					// Read-only properties.
				default:
					usage.mutates = true
				}
			}
		case *ast.IndexExpr:
			if isBase(e.Object) {
				consumed[e.Object] = true
				usage.indexed = true
			}
		case *ast.BinaryExpr:
			if isAssignmentOp(e.Op) && isArrayAssignmentTarget(e.Left, isBase) {
				usage.mutates = true
			}
		case *ast.Identifier:
			if isBase(e) && !consumed[e] {
				// A bare reference (argument, alias, whole-array assignment)
				// cannot be proven read-only.
				usage.mutates = true
			}
		}
	}

	for _, stmt := range program.Statements {
		walkStmt(stmt, visit)
		if usage.mutates {
			return usage
		}
	}
	return usage
}

// arrayInitializerIsStatic reports whether an array initializer can be written
// as a C static initializer. Arrays without an initializer are accepted and
// zero-filled.
func arrayInitializerIsStatic(value ast.Expr) bool {
	if value == nil {
		return true
	}
	lit, ok := value.(*ast.ArrayLiteralExpr)
	if !ok {
		return false
	}
	for _, element := range lit.Elements {
		if !isStaticArrayElement(element) {
			return false
		}
	}
	return true
}

// isStaticArrayElement reports whether expr is valid inside a C static
// initializer list. Object literals and identifiers are excluded: the former
// lower to compound literals, which are not constant expressions, and the
// latter would require whole-program constant propagation.
func isStaticArrayElement(expr ast.Expr) bool {
	switch node := expr.(type) {
	case *ast.NumberLiteral, *ast.BoolLiteral, *ast.StringLiteral:
		return true
	case *ast.UnaryExpr:
		if node.Op == "++" || node.Op == "--" {
			return false
		}
		return isStaticArrayElement(node.Expr)
	case *ast.BinaryExpr:
		return isStaticArrayElement(node.Left) && isStaticArrayElement(node.Right)
	case *ast.ConditionalExpr:
		return isStaticArrayElement(node.Cond) &&
			isStaticArrayElement(node.Then) &&
			isStaticArrayElement(node.Else)
	default:
		return false
	}
}

func isArrayAssignmentTarget(expr ast.Expr, isBase func(ast.Expr) bool) bool {
	switch target := expr.(type) {
	case *ast.Identifier:
		return isBase(target)
	case *ast.IndexExpr:
		return isBase(target.Object)
	case *ast.SelectorExpr:
		return isBase(target.Object)
	default:
		return false
	}
}
