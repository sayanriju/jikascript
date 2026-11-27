package ch32fun

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// TestPeripheralRuntimeHandleReuse compiles the real PWM and ADC runtime files
// against a minimal CH32fun test header and checks handle allocation on the
// host. It guards the dedupe-by-pin behavior that keeps repeated direct
// constructor calls from exhausting the handle tables.
func TestDebugRuntimeRead(t *testing.T) {
	if _, err := exec.LookPath("gcc"); err != nil {
		t.Skip("gcc not available for debug runtime test")
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine source location")
	}
	runtimeDir := filepath.Join(filepath.Dir(thisFile), "runtime")
	testdataDir := filepath.Join(runtimeDir, "testdata")
	binary := filepath.Join(t.TempDir(), "debug_runtime_test")

	compile := exec.Command("gcc", "-std=c99", "-Wall", "-Wextra", "-Werror",
		"-I", runtimeDir,
		filepath.Join(testdataDir, "debug_runtime_test.c"),
		filepath.Join(runtimeDir, "debug_stub.c"),
		"-o", binary)
	if output, err := compile.CombinedOutput(); err != nil {
		t.Fatalf("C99 compile failed:\n%s", output)
	}
	if output, err := exec.Command(binary).CombinedOutput(); err != nil {
		t.Fatalf("runtime test failed:\n%s", output)
	}
}

func TestPeripheralRuntimeHandleReuse(t *testing.T) {
	if _, err := exec.LookPath("gcc"); err != nil {
		t.Skip("gcc not available for peripheral runtime test")
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine source location")
	}
	runtimeDir := filepath.Join(filepath.Dir(thisFile), "runtime")
	testdataDir := filepath.Join(runtimeDir, "testdata")
	binary := filepath.Join(t.TempDir(), "peripheral_runtime_test")

	compile := exec.Command("gcc", "-std=c99", "-Wall", "-Wextra", "-Werror",
		"-I", testdataDir, "-I", runtimeDir,
		filepath.Join(testdataDir, "peripheral_runtime_test.c"),
		filepath.Join(runtimeDir, "pwm.c"),
		filepath.Join(runtimeDir, "adc.c"),
		"-o", binary)
	if output, err := compile.CombinedOutput(); err != nil {
		t.Fatalf("C99 compile failed:\n%s", output)
	}
	if output, err := exec.Command(binary).CombinedOutput(); err != nil {
		t.Fatalf("runtime test failed:\n%s", output)
	}
}
