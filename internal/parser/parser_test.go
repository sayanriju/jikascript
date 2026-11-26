package parser

import (
	"strings"
	"testing"

	"jikascript/internal/ast"
)

func TestParseBlinkProgram(t *testing.T) {
	program, err := ParseProgram(`const led = pin("PC1");
led.output();
task("blink", 1000, () => {
    led.toggle();
});
scheduler.start();
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(program.Statements))
	}
}

func TestParseFunctionAndIf(t *testing.T) {
	program, err := ParseProgram(`function add(a, b) { return a + b; }
if (x > 0) { y = 1; } else { y = 2; }
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

func TestParseLoopStatements(t *testing.T) {
	program, err := ParseProgram(`while (count < 3) { count = count + 1; }
for (let i = 0; i < 3; i = i + 1) { count = count + 1; }
do { count = count - 1; } while (count > 0);
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 loop statements, got %d", len(program.Statements))
	}
}

func TestParseSwitchBreakContinueAndBitwise(t *testing.T) {
	program, err := ParseProgram(`switch (mode) {
case 1:
  flags = flags | 1;
  break;
default:
  while (x < 10) { x += 1; continue; }
}
if ((a << 1) === (b ^ 2)) { c &= 3; }
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

func TestParseBoundedArraySyntax(t *testing.T) {
	program, err := ParseProgram(`let readings: int32[8] = [10, 20, 30];
readings.push(40);
const latest = readings.pop();
const count = readings.length;
readings[0] = latest;
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 5 {
		t.Fatalf("expected 5 statements, got %d", len(program.Statements))
	}
}

func TestParseHexArrayCapacity(t *testing.T) {
	program, err := ParseProgram(`let readings: int32[0x10] = [];`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	decl, ok := program.Statements[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl, got %T", program.Statements[0])
	}
	if decl.ArrayCapacity != 16 {
		t.Fatalf("ArrayCapacity = %d, want 16", decl.ArrayCapacity)
	}
}

func TestParseEmptyBoundedArrayLiteral(t *testing.T) {
	program, err := ParseProgram(`let queue: uint16[8] = [];`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}
}

func TestParseObjectLiteralAndPropertyAssignment(t *testing.T) {
	program, err := ParseProgram(`const device = { name: "sensor", pin: 3, enabled: true };
device.pin = 4;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
}

func TestParseTopLevelInclude(t *testing.T) {
	program, err := ParseProgram(`include "helpers.jika";
scheduler.start();
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
	inc, ok := program.Statements[0].(*ast.IncludeStmt)
	if !ok {
		t.Fatalf("first statement is %T, want IncludeStmt", program.Statements[0])
	}
	if inc.Path != "helpers.jika" {
		t.Fatalf("include path = %q, want %q", inc.Path, "helpers.jika")
	}
}

func TestParseNestedIncludeRejected(t *testing.T) {
	_, err := ParseProgram(`function f() { include "helpers.jika"; }
`)
	if err == nil {
		t.Fatalf("expected error for nested include, got nil")
	}
}

func TestParseIncludeMissingSemicolonRejected(t *testing.T) {
	_, err := ParseProgram(`include "helpers.jika"
`)
	if err == nil {
		t.Fatalf("expected error for missing semicolon, got nil")
	}
}

func TestParsePeripheralConstructors(t *testing.T) {
	program, err := ParseProgram(`const serial = uart("PD5", "PD6");
const dimmer = pwm("PC1");
const sensor = adc("PD4");
const led = pin("PC1");
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got %d", len(program.Statements))
	}
}

func TestParseArrowFunctionBinding(t *testing.T) {
	program, err := ParseProgram(`const square = x => x * x;
const toggle = () => { led.toggle(); };
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
	first, ok := program.Statements[0].(*ast.ConstDecl)
	if !ok {
		t.Fatalf("first statement is %T, want ConstDecl", program.Statements[0])
	}
	arrow, ok := first.Value.(*ast.ArrowFunctionExpr)
	if !ok {
		t.Fatalf("first value is %T, want ArrowFunctionExpr", first.Value)
	}
	if len(arrow.Params) != 1 || arrow.Params[0] != "x" {
		t.Fatalf("arrow params = %v, want [x]", arrow.Params)
	}
	second, ok := program.Statements[1].(*ast.ConstDecl)
	if !ok {
		t.Fatalf("second statement is %T, want ConstDecl", program.Statements[1])
	}
	blockArrow, ok := second.Value.(*ast.ArrowFunctionExpr)
	if !ok {
		t.Fatalf("second value is %T, want ArrowFunctionExpr", second.Value)
	}
	if len(blockArrow.Params) != 0 {
		t.Fatalf("block arrow params = %v, want none", blockArrow.Params)
	}
}

func TestParseNestedTaskRejected(t *testing.T) {
	_, err := ParseProgram(`if (true) {
    task("nested", 100, () => {});
}
`)
	if err == nil || !strings.Contains(err.Error(), "top level") {
		t.Fatalf("expected nested task rejection, got %v", err)
	}
}

func TestParsePostfixIncrementAndDecrement(t *testing.T) {
	program, err := ParseProgram(`i++;
i--;
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(program.Statements))
	}
	for index, op := range []string{"++", "--"} {
		stmt, ok := program.Statements[index].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("statement %d is %T, want ExprStmt", index, program.Statements[index])
		}
		unary, ok := stmt.Expr.(*ast.UnaryExpr)
		if !ok {
			t.Fatalf("expression %d is %T, want UnaryExpr", index, stmt.Expr)
		}
		if !unary.Postfix || unary.Op != op {
			t.Fatalf("expression %d = {Op:%q Postfix:%v}, want {%q true}", index, unary.Op, unary.Postfix, op)
		}
	}
}
