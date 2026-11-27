package checker

import (
	"strings"
	"testing"

	"jikascript/internal/ast"
	"jikascript/internal/parser"
)

func TestCheckAllowsLoopingCode(t *testing.T) {
	program, err := parser.ParseProgram(`let count = 0;
while (count < 2) {
  count = count + 1;
}
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsConstMutation(t *testing.T) {
	program, err := parser.ParseProgram(`const limit = 3;
limit = 5;
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "constant") {
		t.Fatalf("expected const mutation error, got %v", err)
	}
}

func TestCheckRejectsInvalidPin(t *testing.T) {
	program, err := parser.ParseProgram(`const led = pin("PX99");`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "unknown hardware pin") {
		t.Fatalf("expected invalid pin error, got %v", err)
	}
}

func TestCheckAcceptsPinNameVariants(t *testing.T) {
	for _, source := range []string{
		`const led = pin("PC1");`,
		`const led = pin("pc1");`,
		`const led = pin(" pc1 ");`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsUnsupportedPinMethod(t *testing.T) {
	program, err := parser.ParseProgram(`const led = pin("PC1"); led.pulse();`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "unsupported pin method") {
		t.Fatalf("expected unsupported pin method error, got %v", err)
	}
}

func TestCheckRejectsPinMethodArity(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const led = pin("PC1"); led.output(1);`, "output() expects no arguments"},
		{`const led = pin("PC1"); led.input(1);`, "input() expects no arguments"},
		{`const led = pin("PC1"); led.toggle(1);`, "toggle() expects no arguments"},
		{`const led = pin("PC1"); led.read(1);`, "read() expects no arguments"},
		{`const led = pin("PC1"); led.write(1, 2);`, "write() expects at most one argument"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsPinMethodArity(t *testing.T) {
	for _, source := range []string{
		`const led = pin("PC1"); led.output();`,
		`const led = pin("PC1"); led.input();`,
		`const led = pin("PC1"); led.toggle();`,
		`const led = pin("PC1"); let x = led.read();`,
		`const led = pin("PC1"); led.write(HIGH);`,
		`const led = pin("PC1"); led.write();`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckAcceptsUART(t *testing.T) {
	for _, source := range []string{
		`const serial = uart("PD5", "PD6");`,
		`const serial = uart(" pd5 ", "pd6");`,
		`const serial = uart("PD5", "PD6"); serial.write("hi");`,
		`const serial = uart("PD5", "PD6"); serial.println("ok");`,
		`const serial = uart("PD5", "PD6"); let ch = serial.read();`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidUART(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const serial = uart("PA2", "PA1");`, "unsupported uart pin pair"},
		{`const serial = uart("PD6", "PD5");`, "unsupported uart pin pair"},
		{`const serial = uart("PD5");`, "uart() expects two or three arguments"},
		{`const serial = uart("PD5", "PD6", "PC1");`, "baud argument must be a positive integer literal"},
		{`const serial = uart("PD5", 1);`, "rx argument must be a string literal"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsUnsupportedUARTMethod(t *testing.T) {
	program, err := parser.ParseProgram(`const serial = uart("PD5", "PD6"); serial.print("x");`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "unsupported uart method") {
		t.Fatalf("expected unsupported uart method error, got %v", err)
	}
}

func TestCheckRejectsUARTMethodArity(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const serial = uart("PD5", "PD6"); serial.write();`, "write() expects exactly one argument"},
		{`const serial = uart("PD5", "PD6"); serial.write("a", "b");`, "write() expects exactly one argument"},
		{`const serial = uart("PD5", "PD6"); serial.println();`, "println() expects exactly one argument"},
		{`const serial = uart("PD5", "PD6"); serial.read(1);`, "read() expects no arguments"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsPWM(t *testing.T) {
	for _, source := range []string{
		`const led = pwm("PC1");`,
		`const led = pwm(" pc1 ");`,
		`const led = pwm("PC1"); led.set(50);`,
		`const led = pwm("PC1"); led.setFrequency(1000);`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidPWM(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const led = pwm("PD4");`, "unknown pwm pin"},
		{`const led = pwm("PC1", "PC2");`, "pwm() expects exactly one argument"},
		{`const led = pwm(1);`, "pwm() argument must be a string literal"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsUnsupportedPWMMethod(t *testing.T) {
	program, err := parser.ParseProgram(`const led = pwm("PC1"); led.pulse();`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "unsupported pwm method") {
		t.Fatalf("expected unsupported pwm method error, got %v", err)
	}
}

func TestCheckRejectsPWMMethodArity(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const led = pwm("PC1"); led.set();`, "set() expects exactly one argument"},
		{`const led = pwm("PC1"); led.setFrequency();`, "setFrequency() expects exactly one argument"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsADC(t *testing.T) {
	for _, source := range []string{
		`const sensor = adc("PD4");`,
		`const sensor = adc(" pd4 ");`,
		`const sensor = adc("PD4"); let value = sensor.read();`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidADC(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const sensor = adc("PC1");`, "unknown adc pin"},
		{`const sensor = adc("PD4", "PD5");`, "adc() expects exactly one argument"},
		{`const sensor = adc(1);`, "adc() argument must be a string literal"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsUnsupportedADCMethod(t *testing.T) {
	program, err := parser.ParseProgram(`const sensor = adc("PD4"); sensor.write(1);`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "unsupported adc method") {
		t.Fatalf("expected unsupported adc method error, got %v", err)
	}
}

func TestCheckRejectsADCMethodArity(t *testing.T) {
	program, err := parser.ParseProgram(`const sensor = adc("PD4"); sensor.read(1);`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "read() expects no arguments") {
		t.Fatalf("expected adc read arity error, got %v", err)
	}
}

func TestCheckRejectsDuplicateTaskNames(t *testing.T) {
	program, err := parser.ParseProgram(`task("blink", 1000, () => { }); task("blink", 500, () => { });`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "duplicate task name") {
		t.Fatalf("expected duplicate task error, got %v", err)
	}
}

func TestCheckRejectsNonPositiveTaskInterval(t *testing.T) {
	for _, source := range []string{
		`task("blink", 0, () => { });`,
		`task("blink", -1, () => { });`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), "task interval must be greater than zero") {
			t.Fatalf("expected interval error for %q, got %v", source, err)
		}
	}
}

func TestCheckAcceptsPositiveTaskInterval(t *testing.T) {
	for _, source := range []string{
		`task("blink", 1, () => { });`,
		`task("blink", 1000, () => { });`,
		`let period = 500; task("blink", period, () => { });`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsTaskCallbackParameters(t *testing.T) {
	program, err := parser.ParseProgram(`task("blink", 1000, (elapsed) => { });`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "task callback must not declare parameters") {
		t.Fatalf("expected callback parameter error, got %v", err)
	}
}

func TestCheckRejectsBreakOutsideLoopOrSwitch(t *testing.T) {
	program, err := parser.ParseProgram(`break;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "break is only valid") {
		t.Fatalf("expected invalid break error, got %v", err)
	}
}

func TestCheckRejectsContinueOutsideLoop(t *testing.T) {
	program, err := parser.ParseProgram(`continue;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "continue is only valid") {
		t.Fatalf("expected invalid continue error, got %v", err)
	}
}

func TestCheckAllowsSwitchAndLoopFlowControl(t *testing.T) {
	program, err := parser.ParseProgram(`let mode = 1; let x = 0; switch (mode) { case 1: break; default: while (x < 2) { continue; } }`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckAllowsBoundedArrayOps(t *testing.T) {
	program, err := parser.ParseProgram(`let readings: int32[6] = [1, 2, 3];
readings.push(4);
const latest = readings.pop();
readings[0] = latest;
const count = readings.length;
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsArrayCapacityOverflow(t *testing.T) {
	program, err := parser.ParseProgram(`let readings: int32[2] = [1, 2, 3];`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "exceeds declared capacity") {
		t.Fatalf("expected capacity overflow error, got %v", err)
	}
}

func TestCheckRejectsConstArrayReassignmentButAllowsMutation(t *testing.T) {
	program, err := parser.ParseProgram(`const readings: int32[4] = [1, 2];
readings.push(3);
readings = [4, 5];
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "cannot assign to constant variable") {
		t.Fatalf("expected const reassignment error, got %v", err)
	}
}

func TestCheckRejectsArrayPushTypeMismatch(t *testing.T) {
	program, err := parser.ParseProgram(`let readings: int32[4] = [1, 2];
readings.push("oops");
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "push() type mismatch") {
		t.Fatalf("expected push type mismatch error, got %v", err)
	}
}

func TestCheckAllowsConstObjectPropertyMutation(t *testing.T) {
	program, err := parser.ParseProgram(`const device = { name: "sensor", pin: 3, enabled: true };
device.pin = 4;
device.enabled = false;
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsInvalidObjectPropertyAndType(t *testing.T) {
	for _, source := range []string{
		`let device = { pin: 3 }; device.temperature;`,
		`let device = { pin: 3 }; device.pin = "wrong";`,
		`let device = { pin: 3, pin: 4 };`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := Check(program); err == nil {
			t.Fatalf("expected checker error for %q", source)
		}
	}
}

func TestCheckObjectFunctionReturn(t *testing.T) {
	program, err := parser.ParseProgram(`function createDevice() { return { pin: 3, enabled: true }; }
const device = createDevice();
device.pin = 4;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckAllowsInferredArrayOfObjects(t *testing.T) {
	program, err := parser.ParseProgram(`let devices = [{ pin: 1 }, { pin: 2 }];`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsOutOfRangeIntegerLiteral(t *testing.T) {
	program, err := parser.ParseProgram(`let small: uint8 = 300;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "type mismatch") {
		t.Fatalf("expected declaration type mismatch error, got %v", err)
	}
}

func TestCheckAllowsInRangeIntegerLiteralForDeclaredType(t *testing.T) {
	program, err := parser.ParseProgram(`let count: uint32 = 100;
let small: uint8 = 255;
let temp: int8 = -128;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsIntegerLiteralOutsideInt32Range(t *testing.T) {
	sources := []string{
		`let x = 2147483648;`,
		`let x: int32 = 2147483648;`,
		`let x = 4294967295;`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), "out of range") {
			t.Fatalf("expected out-of-range error for %q, got %v", source, err)
		}
	}
}

func TestCheckAllowsMaxInt32Literal(t *testing.T) {
	program, err := parser.ParseProgram(`let x = 2147483647;
let y: int32 = -2147483648;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckAllowsNumericWideningDeclaration(t *testing.T) {
	program, err := parser.ParseProgram(`let count: int32 = 42;
let reading: float32 = count;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsIncompatibleScalarAssignment(t *testing.T) {
	program, err := parser.ParseProgram(`let enabled: boolean = true;
enabled = 1;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "type mismatch assigning 'enabled'") {
		t.Fatalf("expected assignment type mismatch error, got %v", err)
	}
}

func TestCheckRejectsIncompatibleEqualityComparison(t *testing.T) {
	program, err := parser.ParseProgram(`let count = 1;
let label = "x";
if (count == label) { count = 2; }`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "cannot compare") {
		t.Fatalf("expected incompatible equality error, got %v", err)
	}
}

func TestCheckRejectsLocalArrayDeclaration(t *testing.T) {
	program, err := parser.ParseProgram(`function init() {
  let values: int32[4] = [1, 2];
  return 0;
}`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "array declarations are only supported at top-level scope") {
		t.Fatalf("expected local array declaration error, got %v", err)
	}
}

func TestCheckRejectsArrayLiteralReassignment(t *testing.T) {
	program, err := parser.ParseProgram(`let values: int32[4] = [1, 2];
values = [3, 4];`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "array literals are only supported in top-level declarations") {
		t.Fatalf("expected array reassignment rejection, got %v", err)
	}
}

func TestCheckRejectsArrayUsedAsScalarValue(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const t: int32[4] = [1, 2]; const u = t;`, "cannot assign array 't' to a scalar binding"},
		{`const t: int32[4] = [1, 2]; let u = t;`, "cannot assign array 't' to a scalar binding"},
		{`const t: int32[4] = [1, 2]; const x = t + 1;`, "array 't' cannot be used as a value"},
		{`const t: int32[4] = [1, 2]; const x = t == 1;`, "array 't' cannot be used as a value"},
		{`const t: int32[4] = [1, 2]; if (t) { }`, "array 't' cannot be used as a value"},
		{`function f(a) { return a; } const t: int32[4] = [1, 2]; const x = f(t);`, "array 't' cannot be used as a value"},
		{`let t: int32[4] = [1, 2]; t = 5;`, "cannot assign to array 't'"},
		{`const t: int32[4] = [1, 2]; switch (t) { case 1: break; }`, "array 't' cannot be used as a value"},
		{`const t: int32[4] = [1, 2]; switch (1) { case t: break; }`, "array 't' cannot be used as a value"},
		{`const t: int32[4] = [1, 2]; task("x", t, () => { });`, "array 't' cannot be used as a value"},
		{`const t: int32[4] = [1, 2]; const a: int32[2] = [t, 1];`, "array 't' cannot be used as a value"},
		{`const t: int32[4] = [1, 2]; const o = { a: t };`, "array 't' cannot be used as a value"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsObjectUsedAsScalarValue(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const o = { a: 1 }; const u = o;`, "cannot assign object 'o' to a scalar binding"},
		{`const o = { a: 1 }; const x = o + 1;`, "object 'o' cannot be used as a value"},
		{`const o = { a: 1 }; if (o) { }`, "object 'o' cannot be used as a value"},
		{`let o = { a: 1 }; o = 5;`, "cannot assign to object 'o'"},
		{`const o = { a: 1 }; switch (o) { case 1: break; }`, "object 'o' cannot be used as a value"},
		{`const o = { a: 1 }; const p = { b: o };`, "object 'o' cannot be used as a value"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAllowsArrayAndObjectAccess(t *testing.T) {
	sources := []string{
		`const t: int32[4] = [1, 2]; const x = t[0];`,
		`const t: int32[4] = [1, 2]; const n = t.length;`,
		`const t: int32[4] = [1, 2]; const c = t.capacity;`,
		`const t: int32[4] = [1, 2]; const p = t.peek();`,
		`const t: int32[4] = [1, 2]; const i = t.indexOf(2);`,
		`const o = { a: 1 }; const y = o.a;`,
		`const o = { a: 1 }; o.a = 2;`,
		`function f(a) { return a + 1; } const t: int32[4] = [1, 2]; const x = f(t[0]);`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsMixedBareAndValueReturns(t *testing.T) {
	program, err := parser.ParseProgram(`function test(flag) {
  if (flag) { return 1; }
  return;
}`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "mixes bare and value returns") {
		t.Fatalf("expected mixed return contract error, got %v", err)
	}
}

func TestCheckRejectsInconsistentReturnTypes(t *testing.T) {
	program, err := parser.ParseProgram(`function test(flag) {
  if (flag) { return 1; }
  return "oops";
}`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "inconsistent return type") {
		t.Fatalf("expected inconsistent return type error, got %v", err)
	}
}

func TestCheckAllowsCompatibleNumericReturns(t *testing.T) {
	program, err := parser.ParseProgram(`function test(flag) {
  if (flag) { return 1; }
  return 2.5;
}`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckResolvesForwardFunctionCallTypes(t *testing.T) {
	program, err := parser.ParseProgram(`const result = add(1, 2);
function add(a, b) { return a + b; }`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
	fn := program.Statements[1].(*ast.FunctionDecl)
	if fn.ReturnType != "int32" || len(fn.ParamTypes) != 2 || fn.ParamTypes[0] != "int32" || fn.ParamTypes[1] != "int32" {
		t.Fatalf("unexpected resolved function types: params=%v return=%q", fn.ParamTypes, fn.ReturnType)
	}
}

func TestCheckRejectsIncompatibleFunctionArgument(t *testing.T) {
	program, err := parser.ParseProgram(`function add(a, b) { return a + b; }
const result = add(1, "wrong");`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "numeric operands") {
		t.Fatalf("expected incompatible function argument error, got %v", err)
	}
}

func TestCheckAcceptsPeripheralArgumentTypes(t *testing.T) {
	for _, source := range []string{
		`const led = pwm("PC1"); led.set(50);`,
		`const led = pwm("PC1"); led.setFrequency(1000);`,
		`const led = pwm("PC1"); led.set(1.5);`,
		`const led = pwm("PC1"); let duty = 50; led.set(duty);`,
		`const serial = uart("PD5", "PD6"); serial.write("hi");`,
		`const serial = uart("PD5", "PD6"); serial.println("ok");`,
		`const serial = uart("PD5", "PD6"); let msg = "hi"; serial.write(msg);`,
		`const led = pin("PC1"); led.write(HIGH);`,
		`const led = pin("PC1"); led.write();`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsPeripheralArgumentTypeMismatch(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const led = pwm("PC1"); led.set("50");`, "set() argument 1 expects number, got string"},
		{`const led = pwm("PC1"); led.set(true);`, "set() argument 1 expects number, got boolean"},
		{`const led = pwm("PC1"); let duty = "50"; led.set(duty);`, "set() argument 1 expects number, got string"},
		{`const led = pwm("PC1"); led.setFrequency("fast");`, "setFrequency() argument 1 expects number, got string"},
		{`const serial = uart("PD5", "PD6"); serial.write(1);`, "write() argument 1 expects string, got int32"},
		{`const serial = uart("PD5", "PD6"); serial.println(42);`, "println() argument 1 expects string, got int32"},
		{`const led = pin("PC1"); led.write("HIGH");`, "write() argument 1 expects number, got string"},
		{`const led = pin("PC1"); led.write(true);`, "write() argument 1 expects number, got boolean"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsPeripheralHandleAsValue(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const serial = uart("PD5", "PD6"); serial.write(serial);`, "write() argument 1 expects string, got uart handle"},
		{`const led = pwm("PC1"); const serial = uart("PD5", "PD6"); led.set(serial);`, "set() argument 1 expects number, got uart handle"},
		{`const serial = uart("PD5", "PD6"); let copy = serial;`, "cannot assign uart peripheral 'serial' to a scalar binding"},
		{`const led = pin("PC1"); let copy = led;`, "cannot assign pin peripheral 'led' to a scalar binding"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsPeripheralReturnValues(t *testing.T) {
	for _, source := range []string{
		`const serial = uart("PD5", "PD6"); let ch = serial.read();`,
		`const sensor = adc("PD4"); let level = sensor.read();`,
		`const button = pin("PD2"); let state = button.read();`,
		`const serial = uart("PD5", "PD6"); if (serial.read() >= 0) { serial.write("ok"); }`,
		`const sensor = adc("PD4"); if (sensor.read() > 512) { }`,
		`const serial = uart("PD5", "PD6"); let ch: int32 = serial.read();`,
		`const sensor = adc("PD4"); let level: int32 = sensor.read();`,
	} {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsVoidPeripheralResult(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const led = pwm("PC1"); let x = led.set(50);`, "expression does not produce a value"},
		{`const led = pwm("PC1"); let x = led.setFrequency(1000);`, "expression does not produce a value"},
		{`const serial = uart("PD5", "PD6"); let x = serial.write("hi");`, "expression does not produce a value"},
		{`const serial = uart("PD5", "PD6"); let x = serial.println("hi");`, "expression does not produce a value"},
		{`const led = pin("PC1"); let x = led.write(HIGH);`, "expression does not produce a value"},
		{`const led = pwm("PC1"); let x = 0; x = led.set(50);`, "expression does not produce a value"},
		{`const led = pwm("PC1"); if (led.set(50)) { }`, "expression does not produce a value"},
		{`const led = pwm("PC1"); while (led.set(50)) { }`, "expression does not produce a value"},
		{`const led = pwm("PC1"); let x = led.set(50) + 1;`, "expression does not produce a value"},
		{`const led = pwm("PC1"); let x = !led.set(50);`, "expression does not produce a value"},
		{`const serial = uart("PD5", "PD6"); serial.write(serial.write("x"));`, "write() argument 1 expects string, got void expression"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsPeripheralReturnTypeMismatch(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const sensor = adc("PD4"); let value: string = sensor.read();`, "type mismatch for declaration: expected string, got int32"},
		{`const serial = uart("PD5", "PD6"); let value: string = serial.read();`, "type mismatch for declaration: expected string, got int32"},
		{`const serial = uart("PD5", "PD6"); serial.write(serial.read());`, "write() argument 1 expects string, got int32"},
		{`const button = pin("PD2"); const serial = uart("PD5", "PD6"); serial.write(button.read());`, "write() argument 1 expects string, got int32"},
		{`let values: int32[2] = [1, 2]; const serial = uart("PD5", "PD6"); serial.write(values[0]);`, "write() argument 1 expects string, got int32"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsMembersOnScalars(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`let x = 5; x.write(1);`, "'x' has type int32 and no properties or methods"},
		{`let x = 5; let y = x.length;`, "'x' has type int32 and no properties or methods"},
		{`let s = "hi"; s.println("x");`, "unsupported string method"},
		{`let b = true; b.toggle();`, "'b' has type boolean and no properties or methods"},
		{`const c = 10; c.read();`, "'c' has type int32 and no properties or methods"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckValidatesPeripheralConstructorReceivers(t *testing.T) {
	valid := []string{
		`pin("PC1").write(1);`,
		`pin("PC1").toggle();`,
		`pwm("PC1").set(50);`,
		`pwm("PC1").setFrequency(1000);`,
		`uart("PD5", "PD6").write("hi");`,
		`adc("PD4").read();`,
	}
	for _, source := range valid {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}

	invalid := []struct {
		source string
		want   string
	}{
		{`pin("PC1").foo();`, "unsupported pin method 'foo'"},
		{`pwm("PC1").foo();`, "unsupported pwm method 'foo'"},
		{`adc("PD4").foo();`, "unsupported adc method 'foo'"},
		{`pwm("PC1").set("50");`, "set() argument 1 expects number, got string"},
		{`uart("PD5", "PD6").write(1);`, "write() argument 1 expects string, got int32"},
		{`uart("PD5", "PD6").println(1);`, "println() argument 1 expects string, got int32"},
	}
	for _, tc := range invalid {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAllowsArrowFunctionBinding(t *testing.T) {
	program, err := parser.ParseProgram(`const square = x => x * x;
let result = square(5);
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckAllowsFunctionWithoutReturn(t *testing.T) {
	program, err := parser.ParseProgram(`function noop() { let x = 1; }
noop();
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckAllowsBlockArrowWithoutReturn(t *testing.T) {
	program, err := parser.ParseProgram(`const toggle = () => { let x = 1; };
toggle();
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsShadowingPredefinedNames(t *testing.T) {
	sources := []string{
		`const HIGH = 5;`,
		`let LOW = 0;`,
		`function HIGH() { return 1; }`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), "predefined name") {
			t.Fatalf("expected predefined-name error for %q, got %v", source, err)
		}
	}
}

func TestCheckRejectsTaskNameCollisionAfterNormalization(t *testing.T) {
	program, err := parser.ParseProgram(`task("my-task", 100, () => {});
task("my_task", 200, () => {});
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "collides") {
		t.Fatalf("expected task name collision error, got %v", err)
	}
}

func TestCheckAllowsDistinctTaskNames(t *testing.T) {
	program, err := parser.ParseProgram(`task("blink", 100, () => {});
task("button", 200, () => {});
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsDuplicateFunctionDeclarations(t *testing.T) {
	program, err := parser.ParseProgram(`function foo() { return 1; }
function foo() { return 2; }
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "already declared") {
		t.Fatalf("expected duplicate function error, got %v", err)
	}
}

func TestCheckRejectsFractionalTaskInterval(t *testing.T) {
	program, err := parser.ParseProgram(`task("t", 1.5, () => {});`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "whole number") {
		t.Fatalf("expected fractional interval error, got %v", err)
	}
}

func TestCheckRejectsTaskIntervalBeyondUint32(t *testing.T) {
	program, err := parser.ParseProgram(`task("t", 4294967296, () => {});`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "uint32 range") {
		t.Fatalf("expected interval range error, got %v", err)
	}
}

func TestCheckAllowsPostfixIncrementOnVariable(t *testing.T) {
	program, err := parser.ParseProgram(`let i = 0;
i++;
i--;
`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("unexpected checker error: %v", err)
	}
}

func TestCheckRejectsPostfixOnNonTarget(t *testing.T) {
	sources := []struct {
		source string
		want   string
	}{
		{`let x = 1++;`, "requires a variable or property target"},
		{`const c = 1; c++;`, "cannot assign to constant variable"},
		{`let a = [1, 2]; a[0]++;`, "requires a variable or property target"},
	}
	for _, tc := range sources {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsDebugBuiltins(t *testing.T) {
	sources := []string{
		`printf("hello\n");`,
		`printf("value=%d\n", 42);`,
		`delay(10);`,
		`delay_us(50);`,
		`let ms = 5; delay(ms);`,
		`let ch = debugRead();`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidDebugBuiltins(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`printf();`, "printf() expects a format string"},
		{`printf(1);`, "format argument must be a string literal"},
		{`let x = printf("x");`, "expression does not produce a value"},
		{`delay();`, "delay() expects exactly one argument"},
		{`delay(1, 2);`, "delay() expects exactly one argument"},
		{`delay_us();`, "delay_us() expects exactly one argument"},
		{`delay("x");`, "delay() expects a numeric argument"},
		{`delay(true);`, "delay() expects a numeric argument"},
		{`const led = pin("PC1"); delay(led);`, "delay() expects a numeric argument"},
		{`let x = delay(1);`, "expression does not produce a value"},
		{`function printf() {}`, "cannot redeclare predefined name 'printf'"},
		{`function delay() {}`, "cannot redeclare predefined name 'delay'"},
		{`debugRead(1);`, "debugRead() expects no arguments"},
		{`function debugRead() {}`, "cannot redeclare predefined name 'debugRead'"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsHexLiterals(t *testing.T) {
	sources := []string{
		`let mask: uint32 = 0xFFFFFFFF;`,
		`let flags: uint8 = 0xFF;`,
		`let value = 0x10;`,
		`let readings: int32[0x10] = [];`,
		`let readings: int32[4] = [1, 2, 3, 4]; readings[0x2] = 5;`,
		`task("tick", 0x10, () => { });`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidHexLiterals(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`let mask: uint8 = 0x100;`, "type mismatch for declaration"},
		{`let mask = 0xFFFFFFFF;`, "integer literal out of range for int32"},
		{`let readings: int32[4] = [1, 2, 3, 4]; readings[0x10] = 1;`, "out of bounds"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsMillis(t *testing.T) {
	sources := []string{
		`let now = millis();`,
		`let now: uint32 = millis();`,
		`let last: uint32 = 0; let elapsed = millis() - last;`,
		`if (millis() - 0 > 500) { }`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsMillisArguments(t *testing.T) {
	program, err := parser.ParseProgram(`millis(1);`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "millis() expects no arguments") {
		t.Fatalf("expected millis arity error, got %v", err)
	}
}

func TestCheckInfersUint32ForMillisArithmetic(t *testing.T) {
	program, err := parser.ParseProgram(`let last: uint32 = 0; let elapsed = millis() - last;`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("check error: %v", err)
	}
	decl, ok := program.Statements[1].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl, got %T", program.Statements[1])
	}
	if decl.ResolvedType != "uint32" {
		t.Fatalf("elapsed ResolvedType = %q, want uint32", decl.ResolvedType)
	}
}

func TestCheckAcceptsPinPullMethods(t *testing.T) {
	sources := []string{
		`const button = pin("PD2"); button.pullUp();`,
		`const button = pin("PD2"); button.pullDown();`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsPinPullMethodArguments(t *testing.T) {
	program, err := parser.ParseProgram(`const button = pin("PD2"); button.pullUp(1);`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "pullUp() expects no arguments") {
		t.Fatalf("expected pullUp arity error, got %v", err)
	}
}

func TestCheckAcceptsUARTBaudAndByteMethods(t *testing.T) {
	sources := []string{
		`const serial = uart("PD5", "PD6", 9600);`,
		`const serial = uart("PD5", "PD6", 0x2580);`,
		`const serial = uart("PD5", "PD6"); serial.setBaud(19200);`,
		`const serial = uart("PD5", "PD6"); serial.writeByte(0xFF);`,
		`const serial = uart("PD5", "PD6"); let ready = serial.available();`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidUARTBaudAndByteMethods(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const serial = uart("PD5", "PD6", 0);`, "baud argument must be a positive integer literal"},
		{`const serial = uart("PD5", "PD6", -1);`, "baud argument must be a positive integer literal"},
		{`const serial = uart("PD5", "PD6", 9600, 1);`, "uart() expects two or three arguments"},
		{`const serial = uart("PD5", "PD6"); serial.writeByte();`, "writeByte() expects exactly one argument"},
		{`const serial = uart("PD5", "PD6"); serial.writeByte("x");`, "writeByte() argument 1 expects number"},
		{`const serial = uart("PD5", "PD6"); serial.available(1);`, "available() expects no arguments"},
		{`const serial = uart("PD5", "PD6"); serial.setBaud();`, "setBaud() expects exactly one argument"},
		{`const serial = uart("PD5", "PD6"); serial.setBaud("x");`, "setBaud() argument 1 expects number"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckInfersBooleanForUARTAvailable(t *testing.T) {
	program, err := parser.ParseProgram(`const serial = uart("PD5", "PD6"); let ready = serial.available();`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("check error: %v", err)
	}
	decl, ok := program.Statements[1].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected VarDecl, got %T", program.Statements[1])
	}
	if decl.ResolvedType != "boolean" {
		t.Fatalf("ready ResolvedType = %q, want boolean", decl.ResolvedType)
	}
}

func TestCheckAcceptsSPIAndI2C(t *testing.T) {
	sources := []string{
		`const bus = spi("PC5", "PC7", "PC6");`,
		`const bus = spi("PC5", "PC7", "PC6"); bus.configure(1000000, 0);`,
		`const bus = spi("PC5", "PC7", "PC6"); bus.configure(1000000, 0, 1);`,
		`const bus = spi("PC5", "PC7", "PC6"); let rx = bus.transfer(0x9F);`,
		`const bus = i2c("PC2", "PC1");`,
		`const bus = i2c("PC2", "PC1"); bus.configure(100000);`,
		`const bus = i2c("PC2", "PC1"); bus.write(0x3C, 0x00, 0xAF);`,
		`const bus = i2c("PC2", "PC1"); let v = bus.read(0x3C, 0x00);`,
		`const bus = i2c("PC2", "PC1"); let ok = bus.ping(0x3C);`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckAcceptsRandom(t *testing.T) {
	sources := []string{
		`let r = random();`,
		`let r: uint32 = random();`,
		`randomSeed(123);`,
		`randomSeed(0x12345678);`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidRandom(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`random(1);`, "random() expects no arguments"},
		{`randomSeed();`, "randomSeed() expects exactly one argument"},
		{`randomSeed("x");`, "randomSeed() expects a numeric argument"},
		{`let x = randomSeed(1);`, "expression does not produce a value"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsArrayClear(t *testing.T) {
	program, err := parser.ParseProgram(`let values: int32[4] = [1, 2, 3]; values.clear();`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err != nil {
		t.Fatalf("expected array clear to pass check, got %v", err)
	}
}

func TestCheckRejectsArrayClearArguments(t *testing.T) {
	program, err := parser.ParseProgram(`let values: int32[4] = [1, 2, 3]; values.clear(1);`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "clear() expects no arguments") {
		t.Fatalf("expected clear arity error, got %v", err)
	}
}

func TestCheckAcceptsArrayMethods(t *testing.T) {
	sources := []string{
		`let a: int32[4] = [1, 2, 3]; a.peek();`,
		`let a: int32[4] = [1, 2, 3]; a.fill(0);`,
		`let a: int32[4] = [1, 2, 3]; a.reverse();`,
		`let a: int32[4] = [1, 2, 3]; let i = a.indexOf(2);`,
		`let a: int32[4] = [1, 2, 3]; let c = a.capacity;`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidArrayMethods(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`let a: int32[4] = [1]; a.peek(1);`, "peek() expects no arguments"},
		{`let a: int32[4] = [1]; a.fill();`, "fill() expects one argument"},
		{`let a: int32[4] = [1]; a.fill("x");`, "fill() type mismatch"},
		{`let a: int32[4] = [1]; a.reverse(1);`, "reverse() expects no arguments"},
		{`let a: int32[4] = [1]; a.indexOf();`, "indexOf() expects one argument"},
		{`let a: int32[4] = [1]; a.indexOf("x");`, "indexOf() type mismatch"},
		{`let a: string[2] = ["a", "b"]; a.indexOf("a");`, "indexOf() is not supported for string arrays"},
		{`const a = [{ x: 1 }]; a.indexOf(a[0]);`, "indexOf() is not supported for object arrays"},
		{`let a: int32[4] = [1]; let x = a.clear();`, "expression does not produce a value"},
		{`let a: int32[4] = [1]; let x = a.fill(0);`, "expression does not produce a value"},
		{`let a: int32[4] = [1]; let x = a.reverse();`, "expression does not produce a value"},
		{`let a: int32[4] = [1]; a.bogus;`, "unsupported array property"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsStringLengthAndEquality(t *testing.T) {
	sources := []string{
		`let n = "hello".length;`,
		`let s = "hi"; let n = s.length;`,
		`let a = "x"; let b = "y"; if (a == b) { }`,
		`let a = "x"; let b = "y"; if (a !== b) { }`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckAcceptsMathBuiltins(t *testing.T) {
	sources := []string{
		`let a = abs(-5);`,
		`let m = min(1, 2);`,
		`let m = max(1, 2);`,
		`let c = clamp(5, 0, 10);`,
		`let y = map(512, 0, 1023, 0, 255);`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidMathBuiltins(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`abs();`, "abs() expects exactly one argument"},
		{`abs(1, 2);`, "abs() expects exactly one argument"},
		{`abs("x");`, "abs() argument 1 expects a number"},
		{`min(1);`, "min() expects exactly two arguments"},
		{`max(1, "x");`, "max() argument 2 expects a number"},
		{`clamp(1, 2);`, "clamp() expects exactly three arguments"},
		{`map(1, 2, 3, 4);`, "map() expects exactly five arguments"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsStringCharAt(t *testing.T) {
	sources := []string{
		`let c = "hello".charAt(1);`,
		`let s = "hi"; let c = s.charAt(0);`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsInvalidStringCharAt(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`let c = "x".charAt();`, "charAt() expects exactly one argument"},
		{`let c = "x".charAt(1, 2);`, "charAt() expects exactly one argument"},
		{`let c = "x".charAt("a");`, "charAt() argument 1 expects a number"},
		{`let c = "x".bogus();`, "unsupported string method"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsInvalidStringProperties(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`let n = "hello".bogus;`, "unsupported string property"},
		{`let s = "hi"; let n = s.bogus;`, "has type string and no properties or methods"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckRejectsInvalidSPIAndI2C(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{`const bus = spi("PC5", "PC6", "PC7");`, "unsupported spi pin mapping"},
		{`const bus = spi("PC5", "PC7");`, "spi() expects exactly three arguments"},
		{`const bus = spi("PC5", "PC7", "PC6"); bus.configure(1000000);`, "configure() expects two or three arguments"},
		{`const bus = spi("PC5", "PC7", "PC6"); bus.configure(1, 2, 3, 4);`, "configure() expects two or three arguments"},
		{`const bus = spi("PC5", "PC7", "PC6"); bus.transfer();`, "transfer() expects exactly one argument"},
		{`const bus = spi("PC5", "PC7", "PC6"); bus.transfer("x");`, "transfer() argument 1 expects number"},
		{`const bus = i2c("PC1", "PC2");`, "unsupported i2c pin mapping"},
		{`const bus = i2c("PC2");`, "i2c() expects exactly two arguments"},
		{`const bus = i2c("PC2", "PC1"); bus.write(0x3C, 0x00);`, "write() expects exactly three arguments"},
		{`const bus = i2c("PC2", "PC1"); bus.read(0x3C);`, "read() expects exactly two arguments"},
		{`const bus = i2c("PC2", "PC1"); bus.configure();`, "configure() expects exactly one argument"},
		{`const bus = i2c("PC2", "PC1"); bus.ping();`, "ping() expects exactly one argument"},
		{`const bus = i2c("PC2", "PC1"); bus.ping("x");`, "ping() argument 1 expects number"},
	}
	for _, tc := range cases {
		program, err := parser.ParseProgram(tc.source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", tc.source, err)
		}
		if err := Check(program); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("expected %q for %q, got %v", tc.want, tc.source, err)
		}
	}
}

func TestCheckAcceptsConditionalExpressions(t *testing.T) {
	sources := []string{
		`let x = 1 > 0 ? 1 : 2;`,
		`let s = true ? "a" : "b";`,
		`let x = 1 > 0 ? 1 : 2; let y = x > 0 ? x : 0;`,
	}
	for _, source := range sources {
		program, err := parser.ParseProgram(source)
		if err != nil {
			t.Fatalf("parse error for %q: %v", source, err)
		}
		if err := Check(program); err != nil {
			t.Fatalf("expected %q to pass check, got %v", source, err)
		}
	}
}

func TestCheckRejectsIncompatibleConditionalBranches(t *testing.T) {
	program, err := parser.ParseProgram(`let x = true ? 1 : "a";`)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := Check(program); err == nil || !strings.Contains(err.Error(), "conditional branches have incompatible types") {
		t.Fatalf("expected conditional branch type error, got %v", err)
	}
}
