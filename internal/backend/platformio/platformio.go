// Package platformio runs PlatformIO builds and uploads for generated
// JikaScript firmware projects.
//
// The package is deliberately small: it owns the PlatformIO project layout
// (platformio.ini, src/, include/), the mapping from JikaScript board names to
// PlatformIO board identifiers, and the invocation of the `pio` executable.
// It does not know anything about the JikaScript language or its compiler.
package platformio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DefaultBoard is the PlatformIO board identifier used when the user does not
// select one explicitly. It matches the generic CH32V003F4P6 package that the
// CH32fun ecosystem targets by default.
const DefaultBoard = "genericCH32V003F4P6"

// DefaultEnv is the PlatformIO environment name generated for a project.
const DefaultEnv = "ch32v003"

// PlatformURL is the Community PlatformIO platform that provides the CH32V
// family and the ch32v003fun framework.
const PlatformURL = "https://github.com/Community-PIO-CH32V/platform-ch32v.git"

// Framework is the PlatformIO framework name for the CH32fun runtime.
const Framework = "ch32v003fun"

// boards maps accepted board selectors to PlatformIO board identifiers.
// The CH32V003 package variants share the same MCU and differ only in pinout,
// so any of them is a valid target for the same generated firmware.
var boards = map[string]string{
	"ch32v003":            DefaultBoard,
	"genericch32v003f4p6": DefaultBoard,
	"genericch32v003f4u6": "genericCH32V003F4U6",
	"genericch32v003a4m6": "genericCH32V003A4M6",
	"genericch32v003j4m6": "genericCH32V003J4M6",
	"ch32v003f4p6_evt_r0": "ch32v003f4p6_evt_r0",
}

// ResolveBoard normalizes a user-supplied board selector into a PlatformIO
// board identifier. An empty selector resolves to DefaultBoard.
func ResolveBoard(selector string) (string, error) {
	trimmed := strings.TrimSpace(selector)
	if trimmed == "" {
		return DefaultBoard, nil
	}
	if board, ok := boards[strings.ToLower(trimmed)]; ok {
		return board, nil
	}
	return "", fmt.Errorf("unsupported board %q (supported: %s)", selector, strings.Join(SupportedBoards(), ", "))
}

// SupportedBoards returns the sorted list of accepted board selectors.
func SupportedBoards() []string {
	names := make([]string, 0, len(boards))
	for name := range boards {
		names = append(names, name)
	}
	sortStrings(names)
	return names
}

// sortStrings is a tiny insertion sort so the package avoids importing sort
// solely for a handful of board names.
func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}

// Optimize modes accepted by Project.Optimize.
const (
	// OptimizeSize favors a small firmware image.
	OptimizeSize = "size"
	// OptimizeSpeed favors execution speed.
	OptimizeSpeed = "speed"
)

// optimizeFlags maps an optimization mode to the compiler flags that
// implement it.
var optimizeFlags = map[string]string{
	OptimizeSize:  "-Os",
	OptimizeSpeed: "-O2",
}

// ResolveOptimize normalizes a user-supplied optimization mode. An empty mode
// is valid and means "leave the toolchain default alone", so it resolves to
// the empty string rather than to a concrete mode.
func ResolveOptimize(mode string) (string, error) {
	trimmed := strings.ToLower(strings.TrimSpace(mode))
	if trimmed == "" {
		return "", nil
	}
	if _, ok := optimizeFlags[trimmed]; ok {
		return trimmed, nil
	}
	return "", fmt.Errorf("unsupported optimize mode %q (use %s or %s)", mode, OptimizeSize, OptimizeSpeed)
}

// Project describes a generated PlatformIO project on disk.
type Project struct {
	// Dir is the project root that contains platformio.ini.
	Dir string
	// Env is the PlatformIO environment name to build or upload.
	Env string
	// Board is the resolved PlatformIO board identifier.
	Board string
	// Optimize selects the compiler optimization mode. Empty means the
	// PlatformIO platform default is used unchanged.
	Optimize string
}

// Config renders the platformio.ini contents for the project.
func (p Project) Config() string {
	config := fmt.Sprintf(`[platformio]
src_dir = src
include_dir = include

[env:%s]
platform = %s
board = %s
framework = %s
`, p.Env, PlatformURL, p.Board, Framework)
	if flags := p.buildFlags(); flags != "" {
		// PlatformIO appends its own default CCFLAGS after build_flags, so a
		// requested mode that differs from the platform default would be
		// overridden. Strip the default first, then apply the selected mode.
		// When the mode matches the default, leave it alone: build_unflags
		// would remove the flag we are about to add.
		if flags != defaultOptimizeFlag {
			config += fmt.Sprintf("build_unflags = %s\n", defaultOptimizeFlag)
		}
		config += fmt.Sprintf("build_flags = %s\n", flags)
	}
	return config
}

// buildFlags returns the compiler flags for the project's optimization mode,
// or an empty string when no mode was requested.
func (p Project) buildFlags() string {
	return optimizeFlags[p.Optimize]
}

// defaultOptimizeFlag is the optimization flag PlatformIO's CH32V platform
// injects by default.
const defaultOptimizeFlag = "-Os"

// WriteConfig writes platformio.ini into the project directory.
func (p Project) WriteConfig() error {
	return os.WriteFile(filepath.Join(p.Dir, "platformio.ini"), []byte(p.Config()), 0o644)
}

// ArtifactDir returns the directory where PlatformIO places build outputs.
func (p Project) ArtifactDir() string {
	return filepath.Join(p.Dir, ".pio", "build", p.Env)
}

// ArtifactPath returns the path of a PlatformIO artifact by base name, such as
// "firmware.elf".
func (p Project) ArtifactPath(name string) string {
	return filepath.Join(p.ArtifactDir(), name)
}

// Runner invokes the PlatformIO CLI. It is a struct rather than a set of free
// functions so tests can substitute a fake executable.
type Runner struct {
	// Executable is the PlatformIO command name or path. Empty means "pio".
	Executable string
	// Stdout and Stderr receive the child process output. Nil means the
	// current process streams.
	Stdout *os.File
	Stderr *os.File
}

// NewRunner returns a Runner that uses the `pio` executable on PATH.
func NewRunner() *Runner {
	return &Runner{Executable: "pio"}
}

func (r *Runner) executable() string {
	if strings.TrimSpace(r.Executable) == "" {
		return "pio"
	}
	return r.Executable
}

// EnsureAvailable verifies that the PlatformIO executable can be found. It
// returns an actionable error when PlatformIO is missing.
func (r *Runner) EnsureAvailable() error {
	exe := r.executable()
	if _, err := exec.LookPath(exe); err != nil {
		return fmt.Errorf("PlatformIO (%s) not found on PATH; install it with `pip install platformio` or see https://docs.platformio.org/en/latest/core/installation/", exe)
	}
	return nil
}

// Build runs `pio run` for the project environment.
func (r *Runner) Build(project Project) error {
	return r.run(project.Dir, "run", "-e", project.Env)
}

// Upload runs `pio run -t upload` for the project environment. When port is
// non-empty it is passed through as the upload port.
func (r *Runner) Upload(project Project, port string) error {
	args := []string{"run", "-e", project.Env, "-t", "upload"}
	if strings.TrimSpace(port) != "" {
		args = append(args, "--upload-port", strings.TrimSpace(port))
	}
	return r.run(project.Dir, args...)
}

func (r *Runner) run(dir string, args ...string) error {
	cmd := exec.Command(r.executable(), args...)
	cmd.Dir = dir
	if r.Stdout != nil {
		cmd.Stdout = r.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}
	if r.Stderr != nil {
		cmd.Stderr = r.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pio %s failed: %w", strings.Join(args, " "), err)
	}
	return nil
}

// CopyArtifacts copies the PlatformIO build outputs into destDir using the
// supplied base name. It returns the list of written paths.
//
// PlatformIO always names its outputs firmware.elf/firmware.bin, so the
// artifacts are renamed to <base>.elf and <base>.bin to match the artifact
// naming used by the rest of the JikaScript CLI.
func CopyArtifacts(project Project, destDir, base string) ([]string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}
	sources := []struct {
		from string
		to   string
	}{
		{from: "firmware.elf", to: base + ".elf"},
		{from: "firmware.bin", to: base + ".bin"},
	}
	written := make([]string, 0, len(sources))
	for _, artifact := range sources {
		sourcePath := project.ArtifactPath(artifact.from)
		contents, err := os.ReadFile(sourcePath)
		if err != nil {
			return nil, fmt.Errorf("read PlatformIO artifact %s: %w", sourcePath, err)
		}
		destination := filepath.Join(destDir, artifact.to)
		if err := os.WriteFile(destination, contents, 0o644); err != nil {
			return nil, err
		}
		written = append(written, destination)
	}
	return written, nil
}
