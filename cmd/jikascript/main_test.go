package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"jikascript/internal/version"
)

func TestParseOptionsDefaults(t *testing.T) {
	opts, err := parseOptions([]string{"examples/blink.jika"})
	if err != nil {
		t.Fatalf("parseOptions error: %v", err)
	}
	if opts.input != "examples/blink.jika" {
		t.Fatalf("unexpected input: %s", opts.input)
	}
	if opts.target != "ch32v003" {
		t.Fatalf("unexpected target: %s", opts.target)
	}
	if opts.optimize != "" {
		t.Fatalf("expected no optimize mode by default, got %q", opts.optimize)
	}
}

func TestVersionMetadata(t *testing.T) {
	if version.Compiler != "0.2.0" || version.Spec != "0.1.0" {
		t.Fatalf("unexpected version metadata: compiler=%q spec=%q", version.Compiler, version.Spec)
	}
}

func TestParseOptionsWithFlags(t *testing.T) {
	opts, err := parseOptions([]string{"examples/blink.jika", "-o", "build/blink.c", "--target", "ch32v003", "--optimize", "speed"})
	if err != nil {
		t.Fatalf("parseOptions error: %v", err)
	}
	if opts.output != "build/blink.c" {
		t.Fatalf("unexpected output: %s", opts.output)
	}
	if opts.optimize != "speed" {
		t.Fatalf("unexpected optimize mode: %s", opts.optimize)
	}
}

func TestParseOptionsRejectsUnknownOption(t *testing.T) {
	_, err := parseOptions([]string{"examples/blink.jika", "--unknown"})
	if err == nil {
		t.Fatal("expected error for unknown option")
	}
}

func TestParseOptionsAcceptsOutputFlag(t *testing.T) {
	for _, args := range [][]string{
		{"examples/blink.jika", "--output=build/blink.c"},
		{"examples/blink.jika", "--output", "build/blink.c"},
		{"examples/blink.jika", "-o", "build/blink.c"},
	} {
		opts, err := parseOptions(args)
		if err != nil {
			t.Fatalf("parseOptions(%v) error: %v", args, err)
		}
		if opts.output != "build/blink.c" {
			t.Fatalf("parseOptions(%v) output = %q, want %q", args, opts.output, "build/blink.c")
		}
	}
}

func TestParseOptionsAcceptsEqualsFormForAllFlags(t *testing.T) {
	opts, err := parseOptions([]string{
		"examples/blink.jika",
		"--output=build/blink.elf",
		"--board=genericCH32V003J4M6",
		"--env=ch32v003",
		"--optimize=speed",
		"--port=/dev/ttyUSB0",
	})
	if err != nil {
		t.Fatalf("parseOptions error: %v", err)
	}
	if opts.output != "build/blink.elf" {
		t.Fatalf("unexpected output: %s", opts.output)
	}
	if opts.board != "genericCH32V003J4M6" {
		t.Fatalf("unexpected board: %s", opts.board)
	}
	if opts.env != "ch32v003" {
		t.Fatalf("unexpected env: %s", opts.env)
	}
	if opts.optimize != "speed" {
		t.Fatalf("unexpected optimize mode: %s", opts.optimize)
	}
	if opts.port != "/dev/ttyUSB0" {
		t.Fatalf("unexpected port: %s", opts.port)
	}
}

func TestParseOptionsRejectsUnknownOptimizeMode(t *testing.T) {
	_, err := parseOptions([]string{"examples/blink.jika", "--optimize", "balanced"})
	if err == nil || !strings.Contains(err.Error(), "unsupported optimize mode") {
		t.Fatalf("expected optimize validation error, got %v", err)
	}
}

func TestRunCheckRejectsNonMicInput(t *testing.T) {
	err := runCheck([]string{"example.txt"})
	if err == nil || !strings.Contains(err.Error(), ".jika source file") {
		t.Fatalf("expected .jika extension error, got %v", err)
	}
}

func TestRunExportCreatesSelfContainedPlatformIOProject(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "blink.jika")
	outputDir := filepath.Join(tempDir, "blink-export")
	if err := os.WriteFile(inputPath, []byte("const led = pin(\"PC1\");\nled.output();\n"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("create export directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "stale.txt"), []byte("stale"), 0o644); err != nil {
		t.Fatalf("write stale file: %v", err)
	}

	if err := runExport([]string{inputPath, "-o", outputDir}); err != nil {
		t.Fatalf("runExport error: %v", err)
	}
	for _, name := range []string{
		"src/main.c",
		"src/runtime/periodic.c",
		"src/runtime/gpio.c",
		"src/runtime/periodic_ch32fun.c",
		"include/periodic.h",
		"include/gpio.h",
		"include/periodic_ch32fun.h",
		"include/funconfig.h",
		"platformio.ini",
		"README.md",
	} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err != nil {
			t.Errorf("export missing %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(outputDir, "stale.txt")); !os.IsNotExist(err) {
		t.Fatalf("export should remove stale files, stat error = %v", err)
	}
	config, err := os.ReadFile(filepath.Join(outputDir, "platformio.ini"))
	if err != nil {
		t.Fatalf("read PlatformIO configuration: %v", err)
	}
	if !strings.Contains(string(config), "framework = ch32v003fun") {
		t.Fatalf("export should select the CH32fun PlatformIO framework, got:\n%s", config)
	}
	if !strings.Contains(string(config), "src_dir = src") {
		t.Fatalf("export should use src/ as PlatformIO source dir, got:\n%s", config)
	}
	if !strings.Contains(string(config), "include_dir = include") {
		t.Fatalf("export should use include/ as PlatformIO include dir, got:\n%s", config)
	}
}

func TestValidateExportDirRejectsCurrentDirectory(t *testing.T) {
	if err := validateExportDir("."); err == nil {
		t.Fatal("expected current directory to be rejected")
	}
}

func TestValidateTargetRejectsUnknown(t *testing.T) {
	err := validateTarget("invalid")
	if err == nil || !strings.Contains(err.Error(), "unsupported target") {
		t.Fatalf("expected unsupported target error, got %v", err)
	}
}

func TestResolveBuildOutputDefault(t *testing.T) {
	dir, base := resolveBuildOutput("examples/blink.jika", "")
	if dir != "build" || base != "blink" {
		t.Fatalf("unexpected default build output: dir=%q base=%q", dir, base)
	}
}

func TestResolveBuildOutputWithCustomPath(t *testing.T) {
	dir, base := resolveBuildOutput("examples/blink.jika", "dist/fw/blink.elf")
	if dir != "dist/fw" || base != "blink" {
		t.Fatalf("unexpected custom build output: dir=%q base=%q", dir, base)
	}
}

func TestPrepareBuildInputAcceptsGeneratedC(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "blink.c")
	if err := os.WriteFile(inputPath, []byte("int main(void) { return 0; }\n"), 0o644); err != nil {
		t.Fatalf("write generated C: %v", err)
	}

	project, dir, base, err := preparePlatformIOProject(cliOptions{input: inputPath, target: "ch32v003"})
	if err != nil {
		t.Fatalf("preparePlatformIOProject error for generated C: %v", err)
	}
	if dir != "build" || base != "blink" {
		t.Fatalf("unexpected generated C output metadata: dir=%q base=%q", dir, base)
	}
	if project.Dir != filepath.Join("out", "blink-pio") {
		t.Fatalf("unexpected project dir: %q", project.Dir)
	}
	if project.Env != "ch32v003" || project.Board != "genericCH32V003F4P6" {
		t.Fatalf("unexpected project env/board: env=%q board=%q", project.Env, project.Board)
	}
	for _, name := range []string{
		"platformio.ini",
		"src/main.c",
		"src/runtime/periodic.c",
		"include/funconfig.h",
	} {
		if _, err := os.Stat(filepath.Join(project.Dir, name)); err != nil {
			t.Errorf("project missing %s: %v", name, err)
		}
	}
}

func TestPrepareBuildInputTranspilesJikaSource(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "blink.jika")
	if err := os.WriteFile(inputPath, []byte("const led = pin(\"PC1\");\nled.output();\n"), 0o644); err != nil {
		t.Fatalf("write temp .jika file: %v", err)
	}
	outputPrefix := filepath.Join(tempDir, "firmware.elf")

	project, dir, base, err := preparePlatformIOProject(cliOptions{input: inputPath, output: outputPrefix, target: "ch32v003"})
	if err != nil {
		t.Fatalf("preparePlatformIOProject error for .jika file: %v", err)
	}
	if dir != tempDir || base != "firmware" {
		t.Fatalf("unexpected .jika output metadata: dir=%q base=%q", dir, base)
	}
	if project.Dir != filepath.Join("out", "firmware-pio") {
		t.Fatalf("unexpected project dir: %q", project.Dir)
	}
	mainC, err := os.ReadFile(filepath.Join(project.Dir, "src", "main.c"))
	if err != nil {
		t.Fatalf("read generated main.c: %v", err)
	}
	if !strings.Contains(string(mainC), "int main(void)") {
		t.Fatalf("generated main.c does not look like JikaScript output:\n%s", mainC)
	}
	if _, err := os.Stat(filepath.Join(project.Dir, "include", "funconfig.h")); err != nil {
		t.Fatalf("expected generated funconfig.h for build inputs: %v", err)
	}
}

func TestPreparePlatformIOProjectRejectsUnknownBoard(t *testing.T) {
	_, _, _, err := preparePlatformIOProject(cliOptions{input: "blink.jika", board: "esp32"})
	if err == nil || !strings.Contains(err.Error(), "unsupported board") {
		t.Fatalf("expected unsupported board error, got %v", err)
	}
}

func TestPreparePlatformIOProjectHonorsEnvAndBoard(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "blink.jika")
	if err := os.WriteFile(inputPath, []byte("const led = pin(\"PC1\");\nled.output();\n"), 0o644); err != nil {
		t.Fatalf("write temp .jika file: %v", err)
	}

	project, _, _, err := preparePlatformIOProject(cliOptions{
		input: inputPath,
		env:   "custom",
		board: "genericCH32V003J4M6",
	})
	if err != nil {
		t.Fatalf("preparePlatformIOProject error: %v", err)
	}
	if project.Env != "custom" {
		t.Fatalf("unexpected env: %q", project.Env)
	}
	if project.Board != "genericCH32V003J4M6" {
		t.Fatalf("unexpected board: %q", project.Board)
	}
	config, err := os.ReadFile(filepath.Join(project.Dir, "platformio.ini"))
	if err != nil {
		t.Fatalf("read platformio.ini: %v", err)
	}
	if !strings.Contains(string(config), "[env:custom]") {
		t.Fatalf("platformio.ini missing custom env:\n%s", config)
	}
	if !strings.Contains(string(config), "board = genericCH32V003J4M6") {
		t.Fatalf("platformio.ini missing selected board:\n%s", config)
	}
}

func TestPreparePlatformIOProjectHonorsOptimize(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "blink.jika")
	if err := os.WriteFile(inputPath, []byte("const led = pin(\"PC1\");\nled.output();\n"), 0o644); err != nil {
		t.Fatalf("write temp .jika file: %v", err)
	}

	cases := map[string]string{
		"size":  "build_flags = -Os",
		"speed": "build_flags = -O2",
	}
	for mode, expected := range cases {
		project, _, _, err := preparePlatformIOProject(cliOptions{input: inputPath, optimize: mode})
		if err != nil {
			t.Fatalf("preparePlatformIOProject(%s) error: %v", mode, err)
		}
		if project.Optimize != mode {
			t.Fatalf("unexpected optimize mode: %q", project.Optimize)
		}
		config, err := os.ReadFile(filepath.Join(project.Dir, "platformio.ini"))
		if err != nil {
			t.Fatalf("read platformio.ini: %v", err)
		}
		if !strings.Contains(string(config), expected) {
			t.Fatalf("optimize=%s: platformio.ini missing %q:\n%s", mode, expected, config)
		}
	}
}

func TestEnsureFunConfigCreatesMinimalHeader(t *testing.T) {
	tempDir := t.TempDir()
	if err := ensureFunConfig(tempDir); err != nil {
		t.Fatalf("ensureFunConfig error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(tempDir, "funconfig.h"))
	if err != nil {
		t.Fatalf("read generated funconfig.h: %v", err)
	}
	if !strings.Contains(string(data), "#ifndef _FUNCONFIG_H") {
		t.Fatal("expected standard CH32fun guard in generated funconfig.h")
	}
}

func TestParseOptionsAcceptsBoardEnvAndPort(t *testing.T) {
	opts, err := parseOptions([]string{
		"examples/blink.jika",
		"--board", "genericCH32V003J4M6",
		"--env=custom",
		"--port", "/dev/ttyUSB0",
	})
	if err != nil {
		t.Fatalf("parseOptions error: %v", err)
	}
	if opts.board != "genericCH32V003J4M6" {
		t.Fatalf("unexpected board: %q", opts.board)
	}
	if opts.env != "custom" {
		t.Fatalf("unexpected env: %q", opts.env)
	}
	if opts.port != "/dev/ttyUSB0" {
		t.Fatalf("unexpected port: %q", opts.port)
	}
}

// TestPrintUsageDocumentsEqualsForm guards against the usage text drifting away
// from the flag convention the parser actually accepts.
func TestPrintUsageDocumentsEqualsForm(t *testing.T) {
	output := captureStdout(t, printUsage)

	for _, want := range []string{
		"--output=<file.c>",
		"--target=<target>",
		"--board=<board>",
		"--env=<env>",
		"--optimize=<size|speed>",
		"--port=<device>",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("usage text missing %q\n---\n%s", want, output)
		}
	}

	// The old space-separated placeholders must not come back.
	for _, unwanted := range []string{
		"--target ch32v003",
		"--board <board>",
		"--env <env>",
		"--optimize size|speed",
		"--port <device>",
		"[-o output.c]",
	} {
		if strings.Contains(output, unwanted) {
			t.Errorf("usage text still contains stale form %q\n---\n%s", unwanted, output)
		}
	}
}

// captureStdout runs fn and returns everything it wrote to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}
	os.Stdout = writer
	defer func() { os.Stdout = original }()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(data)
}

func TestRunExportWarnsAboutDeprecation(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "blink.jika")
	outputDir := filepath.Join(tempDir, "blink-export")
	if err := os.WriteFile(inputPath, []byte("const led = pin(\"PC1\");\nled.output();\n"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	stderr, err := captureStderr(t, func() error {
		return runExport([]string{inputPath, "-o", outputDir})
	})
	if err != nil {
		t.Fatalf("runExport error: %v", err)
	}
	if !strings.Contains(stderr, "deprecated") {
		t.Fatalf("expected deprecation warning on stderr, got %q", stderr)
	}
}

// captureStderr redirects os.Stderr for the duration of fn and returns what
// was written to it.
func captureStderr(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stderr pipe: %v", err)
	}
	original := os.Stderr
	os.Stderr = writer
	defer func() { os.Stderr = original }()

	runErr := fn()
	if err := writer.Close(); err != nil {
		t.Fatalf("close stderr writer: %v", err)
	}
	captured, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read captured stderr: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stderr reader: %v", err)
	}
	return string(captured), runErr
}

// fakePio installs a stub `pio` executable on PATH that records its arguments
// and creates the firmware artifacts PlatformIO would normally produce.
func fakePio(t *testing.T) string {
	t.Helper()
	binDir := t.TempDir()
	logPath := filepath.Join(binDir, "pio.log")
	script := `#!/bin/sh
echo "$@" >> "` + logPath + `"
env=""
prev=""
for arg in "$@"; do
	if [ "$prev" = "-e" ]; then env="$arg"; fi
	prev="$arg"
done
if [ -n "$env" ]; then
	mkdir -p ".pio/build/$env"
	printf 'elf' > ".pio/build/$env/firmware.elf"
	printf 'bin' > ".pio/build/$env/firmware.bin"
fi
exit 0
`
	if err := os.WriteFile(filepath.Join(binDir, "pio"), []byte(script), 0o755); err != nil {
		t.Fatalf("write fake pio: %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	return logPath
}

func TestRunFlashUploadsViaPlatformIO(t *testing.T) {
	logPath := fakePio(t)
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "blink.jika")
	if err := os.WriteFile(inputPath, []byte("const led = pin(\"PC1\");\nled.output();\n"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	outputPrefix := filepath.Join(tempDir, "blink.elf")

	if err := runFlash([]string{inputPath, "-o", outputPrefix, "--port", "/dev/ttyUSB0"}); err != nil {
		t.Fatalf("runFlash error: %v", err)
	}

	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read fake pio log: %v", err)
	}
	invocation := string(logged)
	for _, want := range []string{"run", "-e", "ch32v003", "-t", "upload", "--upload-port", "/dev/ttyUSB0"} {
		if !strings.Contains(invocation, want) {
			t.Fatalf("pio invocation missing %q:\n%s", want, invocation)
		}
	}
	for _, name := range []string{"blink.elf", "blink.bin"} {
		if _, err := os.Stat(filepath.Join(tempDir, name)); err != nil {
			t.Errorf("flash missing artifact %s: %v", name, err)
		}
	}
}

func TestRunBuildInvokesPlatformIO(t *testing.T) {
	logPath := fakePio(t)
	tempDir := t.TempDir()
	inputPath := filepath.Join(tempDir, "blink.jika")
	if err := os.WriteFile(inputPath, []byte("const led = pin(\"PC1\");\nled.output();\n"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	outputPrefix := filepath.Join(tempDir, "blink.elf")

	if err := runBuild([]string{inputPath, "-o", outputPrefix, "--env", "custom"}); err != nil {
		t.Fatalf("runBuild error: %v", err)
	}

	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read fake pio log: %v", err)
	}
	invocation := string(logged)
	if !strings.Contains(invocation, "-e custom") {
		t.Fatalf("pio invocation missing custom env:\n%s", invocation)
	}
	if strings.Contains(invocation, "upload") {
		t.Fatalf("build must not upload:\n%s", invocation)
	}
	for _, name := range []string{"blink.elf", "blink.bin"} {
		if _, err := os.Stat(filepath.Join(tempDir, name)); err != nil {
			t.Errorf("build missing artifact %s: %v", name, err)
		}
	}
}
