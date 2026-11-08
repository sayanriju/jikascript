package platformio

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveBoardDefaults(t *testing.T) {
	board, err := ResolveBoard("")
	if err != nil {
		t.Fatalf("ResolveBoard error: %v", err)
	}
	if board != DefaultBoard {
		t.Fatalf("unexpected default board: %q", board)
	}
}

func TestResolveBoardAcceptsAliases(t *testing.T) {
	cases := map[string]string{
		"ch32v003":            DefaultBoard,
		"genericCH32V003F4P6": DefaultBoard,
		"genericCH32V003F4U6": "genericCH32V003F4U6",
		"genericCH32V003J4M6": "genericCH32V003J4M6",
	}
	for selector, want := range cases {
		got, err := ResolveBoard(selector)
		if err != nil {
			t.Fatalf("ResolveBoard(%q) error: %v", selector, err)
		}
		if got != want {
			t.Errorf("ResolveBoard(%q) = %q, want %q", selector, got, want)
		}
	}
}

func TestResolveBoardRejectsUnknown(t *testing.T) {
	_, err := ResolveBoard("esp32")
	if err == nil || !strings.Contains(err.Error(), "unsupported board") {
		t.Fatalf("expected unsupported board error, got %v", err)
	}
}

func TestProjectConfigContainsEnvironment(t *testing.T) {
	project := Project{Dir: t.TempDir(), Env: DefaultEnv, Board: DefaultBoard}
	config := project.Config()
	for _, expected := range []string{
		"src_dir = src",
		"include_dir = include",
		"[env:ch32v003]",
		"board = genericCH32V003F4P6",
		"framework = ch32v003fun",
		PlatformURL,
	} {
		if !strings.Contains(config, expected) {
			t.Errorf("platformio.ini missing %q:\n%s", expected, config)
		}
	}
}

func TestResolveOptimizeDefaultsToEmpty(t *testing.T) {
	mode, err := ResolveOptimize("")
	if err != nil {
		t.Fatalf("ResolveOptimize error: %v", err)
	}
	if mode != "" {
		t.Fatalf("ResolveOptimize(\"\") = %q, want empty", mode)
	}
}

func TestResolveOptimizeAcceptsKnownModes(t *testing.T) {
	for _, mode := range []string{"size", "speed", "SIZE", " Speed "} {
		got, err := ResolveOptimize(mode)
		if err != nil {
			t.Fatalf("ResolveOptimize(%q) error: %v", mode, err)
		}
		if got != strings.ToLower(strings.TrimSpace(mode)) {
			t.Fatalf("ResolveOptimize(%q) = %q", mode, got)
		}
	}
}

func TestResolveOptimizeRejectsUnknown(t *testing.T) {
	_, err := ResolveOptimize("balanced")
	if err == nil || !strings.Contains(err.Error(), "unsupported optimize mode") {
		t.Fatalf("expected unsupported optimize mode error, got %v", err)
	}
}

func TestProjectConfigEmitsOptimizeFlags(t *testing.T) {
	cases := map[string]string{
		OptimizeSize:  "-Os",
		OptimizeSpeed: "-O2",
	}
	for mode, flag := range cases {
		project := Project{Dir: t.TempDir(), Env: DefaultEnv, Board: DefaultBoard, Optimize: mode}
		config := project.Config()
		if !strings.Contains(config, "build_flags = "+flag) {
			t.Errorf("optimize=%s: platformio.ini missing %q:\n%s", mode, "build_flags = "+flag, config)
		}
	}
}

func TestProjectConfigUnflagsDefaultOnlyWhenOverriding(t *testing.T) {
	// The platform injects -Os by default. Overriding it requires stripping
	// the default first, otherwise PlatformIO's flag wins.
	speed := Project{Dir: t.TempDir(), Env: DefaultEnv, Board: DefaultBoard, Optimize: OptimizeSpeed}
	if config := speed.Config(); !strings.Contains(config, "build_unflags = -Os") {
		t.Fatalf("speed mode must unflag the platform default:\n%s", config)
	}

	// Requesting the default must not unflag it, or the flag would be
	// stripped and never re-added, leaving the build unoptimized.
	size := Project{Dir: t.TempDir(), Env: DefaultEnv, Board: DefaultBoard, Optimize: OptimizeSize}
	if config := size.Config(); strings.Contains(config, "build_unflags") {
		t.Fatalf("size mode must not unflag its own flag:\n%s", config)
	}
}

func TestProjectConfigOmitsFlagsWhenOptimizeUnset(t *testing.T) {
	project := Project{Dir: t.TempDir(), Env: DefaultEnv, Board: DefaultBoard}
	config := project.Config()
	if strings.Contains(config, "build_flags") {
		t.Fatalf("unset optimize must not emit build_flags:\n%s", config)
	}
	if strings.Contains(config, "build_unflags") {
		t.Fatalf("unset optimize must not emit build_unflags:\n%s", config)
	}
}

func TestProjectWriteConfig(t *testing.T) {
	dir := t.TempDir()
	project := Project{Dir: dir, Env: DefaultEnv, Board: DefaultBoard}
	if err := project.WriteConfig(); err != nil {
		t.Fatalf("WriteConfig error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "platformio.ini"))
	if err != nil {
		t.Fatalf("read platformio.ini: %v", err)
	}
	if !strings.Contains(string(data), "[env:ch32v003]") {
		t.Fatalf("unexpected platformio.ini contents:\n%s", data)
	}
}

func TestProjectArtifactPaths(t *testing.T) {
	project := Project{Dir: "/tmp/project", Env: "ch32v003"}
	if got, want := project.ArtifactDir(), filepath.Join("/tmp/project", ".pio", "build", "ch32v003"); got != want {
		t.Fatalf("ArtifactDir() = %q, want %q", got, want)
	}
	if got, want := project.ArtifactPath("firmware.elf"), filepath.Join("/tmp/project", ".pio", "build", "ch32v003", "firmware.elf"); got != want {
		t.Fatalf("ArtifactPath() = %q, want %q", got, want)
	}
}

func TestCopyArtifactsRenamesFirmware(t *testing.T) {
	projectDir := t.TempDir()
	project := Project{Dir: projectDir, Env: DefaultEnv}
	if err := os.MkdirAll(project.ArtifactDir(), 0o755); err != nil {
		t.Fatalf("create artifact dir: %v", err)
	}
	for name, contents := range map[string]string{
		"firmware.elf": "elf-bytes",
		"firmware.bin": "bin-bytes",
	} {
		if err := os.WriteFile(project.ArtifactPath(name), []byte(contents), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	destDir := filepath.Join(t.TempDir(), "build")
	written, err := CopyArtifacts(project, destDir, "blink")
	if err != nil {
		t.Fatalf("CopyArtifacts error: %v", err)
	}
	if len(written) != 2 {
		t.Fatalf("expected 2 artifacts, got %d", len(written))
	}
	for name, want := range map[string]string{
		"blink.elf": "elf-bytes",
		"blink.bin": "bin-bytes",
	} {
		data, err := os.ReadFile(filepath.Join(destDir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if string(data) != want {
			t.Errorf("%s contents = %q, want %q", name, data, want)
		}
	}
}

func TestCopyArtifactsReportsMissingArtifact(t *testing.T) {
	project := Project{Dir: t.TempDir(), Env: DefaultEnv}
	_, err := CopyArtifacts(project, t.TempDir(), "blink")
	if err == nil || !strings.Contains(err.Error(), "firmware.elf") {
		t.Fatalf("expected missing artifact error, got %v", err)
	}
}

func TestEnsureAvailableReportsMissingExecutable(t *testing.T) {
	runner := &Runner{Executable: "definitely-not-a-real-pio-binary"}
	err := runner.EnsureAvailable()
	if err == nil || !strings.Contains(err.Error(), "not found on PATH") {
		t.Fatalf("expected missing PlatformIO error, got %v", err)
	}
}

func TestRunnerBuildInvokesPioWithEnv(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "pio.log")
	script := filepath.Join(dir, "fake-pio")
	contents := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + logPath + "\n"
	if err := os.WriteFile(script, []byte(contents), 0o755); err != nil {
		t.Fatalf("write fake pio: %v", err)
	}

	project := Project{Dir: dir, Env: DefaultEnv, Board: DefaultBoard}
	runner := &Runner{Executable: script}
	if err := runner.Build(project); err != nil {
		t.Fatalf("Build error: %v", err)
	}
	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read pio log: %v", err)
	}
	if got, want := strings.TrimSpace(string(logged)), "run\n-e\nch32v003"; got != want {
		t.Fatalf("pio build args = %q, want %q", got, want)
	}
}

func TestRunnerUploadInvokesPioWithUploadTarget(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "pio.log")
	script := filepath.Join(dir, "fake-pio")
	contents := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + logPath + "\n"
	if err := os.WriteFile(script, []byte(contents), 0o755); err != nil {
		t.Fatalf("write fake pio: %v", err)
	}

	project := Project{Dir: dir, Env: DefaultEnv, Board: DefaultBoard}
	runner := &Runner{Executable: script}
	if err := runner.Upload(project, "/dev/ttyUSB0"); err != nil {
		t.Fatalf("Upload error: %v", err)
	}
	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read pio log: %v", err)
	}
	got := strings.TrimSpace(string(logged))
	for _, expected := range []string{"run", "-e", "ch32v003", "-t", "upload", "--upload-port", "/dev/ttyUSB0"} {
		if !strings.Contains(got, expected) {
			t.Errorf("pio upload args missing %q: %q", expected, got)
		}
	}
}

func TestRunnerUploadOmitsPortWhenEmpty(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "pio.log")
	script := filepath.Join(dir, "fake-pio")
	contents := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + logPath + "\n"
	if err := os.WriteFile(script, []byte(contents), 0o755); err != nil {
		t.Fatalf("write fake pio: %v", err)
	}

	project := Project{Dir: dir, Env: DefaultEnv, Board: DefaultBoard}
	runner := &Runner{Executable: script}
	if err := runner.Upload(project, ""); err != nil {
		t.Fatalf("Upload error: %v", err)
	}
	logged, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read pio log: %v", err)
	}
	if strings.Contains(string(logged), "--upload-port") {
		t.Fatalf("upload args should omit --upload-port when no port is set: %q", logged)
	}
}

func TestRunnerBuildReportsFailure(t *testing.T) {
	dir := t.TempDir()
	script := filepath.Join(dir, "fake-pio")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nexit 3\n"), 0o755); err != nil {
		t.Fatalf("write fake pio: %v", err)
	}

	project := Project{Dir: dir, Env: DefaultEnv, Board: DefaultBoard}
	runner := &Runner{Executable: script}
	err := runner.Build(project)
	if err == nil || !strings.Contains(err.Error(), "pio run") {
		t.Fatalf("expected pio run failure, got %v", err)
	}
}
