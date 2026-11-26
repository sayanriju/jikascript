package cgen

import "jikascript/internal/ast"

// UsedPeripherals reports the hardware runtime families referenced by a
// program. The result always contains "gpio" and "periodic" because generated
// main() initializes the GPIO backend and the scheduler. "pins" is present
// whenever the program uses a pin, UART, PWM, or ADC constructor, because those
// constructors emit JIKA_PIN_* constants. The remaining keys are "uart", "pwm",
// and "adc". "ch32fun" is present when the program calls printf, delay, or
// delay_us, which lower directly to CH32fun functions and macros; "stdio" is
// present when printf is used.
//
// The CLI uses this set to stage only the runtime sources a project needs.
// Generate uses it to emit only the matching headers.
func UsedPeripherals(program *ast.Program) map[string]bool {
	used := map[string]bool{
		"gpio":     true,
		"periodic": true,
	}

	var visitExpr func(ast.Expr)
	var visitStmt func(ast.Stmt)

	visitExpr = func(expr ast.Expr) {
		if expr == nil {
			return
		}
		switch node := expr.(type) {
		case *ast.CallExpr:
			if ident, ok := node.Function.(*ast.Identifier); ok {
				switch ident.Name {
				case "pin", "uart", "pwm", "adc", "spi", "i2c":
					used["pins"] = true
					if ident.Name != "pin" {
						used[ident.Name] = true
					}
				case "printf":
					used["ch32fun"] = true
					used["stdio"] = true
				case "delay", "delay_us":
					used["ch32fun"] = true
				case "random", "randomSeed":
					used["random"] = true
				case "map":
					used["math"] = true
				}
			}
			visitExpr(node.Function)
			for _, arg := range node.Args {
				visitExpr(arg)
			}
		case *ast.SelectorExpr:
			if node.StringLength {
				used["string"] = true
			}
			visitExpr(node.Object)
		case *ast.IndexExpr:
			visitExpr(node.Object)
			visitExpr(node.Index)
		case *ast.BinaryExpr:
			if node.StringCompare {
				used["string"] = true
			}
			visitExpr(node.Left)
			visitExpr(node.Right)
		case *ast.UnaryExpr:
			visitExpr(node.Expr)
		case *ast.ArrayLiteralExpr:
			for _, element := range node.Elements {
				visitExpr(element)
			}
		case *ast.ObjectLiteralExpr:
			for _, property := range node.Properties {
				visitExpr(property.Value)
			}
		case *ast.ArrowFunctionExpr:
			for _, stmt := range node.Body {
				visitStmt(stmt)
			}
		}
	}

	visitStmt = func(stmt ast.Stmt) {
		if stmt == nil {
			return
		}
		switch node := stmt.(type) {
		case *ast.VarDecl:
			visitExpr(node.Value)
		case *ast.ConstDecl:
			visitExpr(node.Value)
		case *ast.FunctionDecl:
			for _, inner := range node.Body {
				visitStmt(inner)
			}
		case *ast.ReturnStmt:
			visitExpr(node.Value)
		case *ast.ExprStmt:
			visitExpr(node.Expr)
		case *ast.AssignStmt:
			visitExpr(node.Value)
		case *ast.IfStmt:
			visitExpr(node.Cond)
			for _, inner := range node.Then {
				visitStmt(inner)
			}
			for _, inner := range node.Else {
				visitStmt(inner)
			}
		case *ast.WhileStmt:
			visitExpr(node.Cond)
			for _, inner := range node.Body {
				visitStmt(inner)
			}
		case *ast.DoWhileStmt:
			for _, inner := range node.Body {
				visitStmt(inner)
			}
			visitExpr(node.Cond)
		case *ast.ForStmt:
			visitStmt(node.Init)
			visitExpr(node.Cond)
			visitExpr(node.Post)
			for _, inner := range node.Body {
				visitStmt(inner)
			}
		case *ast.SwitchStmt:
			visitExpr(node.Expr)
			for _, clause := range node.Cases {
				visitExpr(clause.Value)
				for _, inner := range clause.Body {
					visitStmt(inner)
				}
			}
		case *ast.TaskDecl:
			visitExpr(node.Interval)
			visitExpr(node.Callback)
		case *ast.BlockStmt:
			for _, inner := range node.Statements {
				visitStmt(inner)
			}
		}
	}

	for _, stmt := range program.Statements {
		visitStmt(stmt)
	}
	return used
}
