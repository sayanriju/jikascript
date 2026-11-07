package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"jikascript/internal/checker"
	cgen "jikascript/internal/codegen/c"
	"jikascript/internal/parser"
	"jikascript/targets/ch32fun"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "-h", "--help", "help":
		printUsage()
	case "check":
		if err := runCheck(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "transpile":
		if err := runTranspile(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "export":
		if err := runExport(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "build":
		if err := runBuild(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  jikascript check <input.jika>")
	fmt.Println("  jikascript transpile <input.jika> [-o output.c] [--target ch32v003]")
	fmt.Println("  jikascript export <input.jika> [-o output-dir] [--target ch32v003]")
	fmt.Println("  jikascript build <input.jika|input.c> [--target ch32v003] [--optimize size|speed]")
}

func runCheck(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing input file")
	}
	path := args[0]
	if filepath.Ext(path) != ".jika" {
		return fmt.Errorf("input must be a .jika source file: %s", path)
	}
	code, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	program, err := parser.ParseProgram(string(code))
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if err := checker.Check(program); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	fmt.Printf("check: %s OK\n", path)
	return nil
}

func runTranspile(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if err := validateTarget(opts.target); err != nil {
		return err
	}
	if opts.input == "" {
		return fmt.Errorf("missing input file")
	}
	if filepath.Ext(opts.input) != ".jika" {
		return fmt.Errorf("input must be a .jika source file: %s", opts.input)
	}
	code, err := os.ReadFile(opts.input)
	if err != nil {
		return err
	}
	program, err := parser.ParseProgram(string(code))
	if err != nil {
		return fmt.Errorf("%s: %w", opts.input, err)
	}
	if err := checker.Check(program); err != nil {
		return fmt.Errorf("%s: %w", opts.input, err)
	}
	output, err := cgen.Generate(program)
	if err != nil {
		return fmt.Errorf("%s: %w", opts.input, err)
	}
	if opts.output == "" {
		fmt.Print(output)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(opts.output), 0o755); err != nil {
		return err
	}
	return os.WriteFile(opts.output, []byte(output), 0o644)
}

// runExport creates a self-contained PlatformIO project. Unlike transpile and
// build, export deliberately recreates its destination so no stale generated
// or runtime files remain in the package.
func runExport(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if err := validateTarget(opts.target); err != nil {
		return err
	}
	if opts.input == "" {
		return fmt.Errorf("missing input file")
	}
	if filepath.Ext(opts.input) != ".jika" {
		return fmt.Errorf("input must be a .jika source file: %s", opts.input)
	}

	outputDir := strings.TrimSpace(opts.output)
	if outputDir == "" {
		base := strings.TrimSuffix(filepath.Base(opts.input), filepath.Ext(opts.input))
		outputDir = filepath.Join("out", base)
	}
	if err := validateExportDir(outputDir); err != nil {
		return err
	}

	code, err := os.ReadFile(opts.input)
	if err != nil {
		return err
	}
	program, err := parser.ParseProgram(string(code))
	if err != nil {
		return fmt.Errorf("%s: %w", opts.input, err)
	}
	if err := checker.Check(program); err != nil {
		return fmt.Errorf("%s: %w", opts.input, err)
	}
	mainC, err := cgen.Generate(program)
	if err != nil {
		return fmt.Errorf("%s: %w", opts.input, err)
	}

	if err := os.RemoveAll(outputDir); err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "main.c"), []byte(mainC), 0o644); err != nil {
		return err
	}
	if err := ensureFunConfig(outputDir); err != nil {
		return err
	}
	for _, file := range exportedRuntimeFiles {
		if err := copyFile(file, filepath.Join(outputDir, filepath.Base(file))); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(outputDir, "platformio.ini"), []byte(platformIOConfig), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(outputDir, "README.md"), []byte(exportReadme), 0o644); err != nil {
		return err
	}
	fmt.Printf("export: generated PlatformIO project in %s\n", outputDir)
	return nil
}

var exportedRuntimeFiles = []string{
	"runtime/periodic/periodic.c",
	"runtime/periodic/periodic.h",
	"targets/ch32fun/runtime/gpio.c",
	"targets/ch32fun/runtime/gpio.h",
	"targets/ch32fun/runtime/periodic_ch32fun.c",
	"targets/ch32fun/runtime/periodic_ch32fun.h",
}

const platformIOConfig = `[platformio]
src_dir = .
include_dir = .

[env:ch32v003]
platform = https://github.com/Community-PIO-CH32V/platform-ch32v.git
board = genericCH32V003F4P6
framework = ch32v003fun
build_flags = -I.
`

const exportReadme = `# JikaScript export

This directory is a generated, self-contained PlatformIO source project for a
CH32V003. Open it in PlatformIO, then build or upload the ch32v003
environment. Change the board in platformio.ini if your CH32V003 package
differs from genericCH32V003F4P6.

Do not edit generated runtime files; regenerate this directory with
jikascript export after changing the JikaScript source.
`

func validateExportDir(dir string) error {
	cleaned := filepath.Clean(dir)
	if cleaned == "." || cleaned == string(filepath.Separator) {
		return fmt.Errorf("export output directory must not be %q", dir)
	}
	if volume := filepath.VolumeName(cleaned); volume != "" && cleaned == volume+string(filepath.Separator) {
		return fmt.Errorf("export output directory must not be a filesystem root")
	}
	return nil
}

func copyFile(source, destination string) error {
	contents, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return os.WriteFile(destination, contents, 0o644)
}

func runBuild(args []string) error {
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	if err := validateTarget(opts.target); err != nil {
		return err
	}
	if opts.input == "" {
		return fmt.Errorf("missing input file")
	}

	source, outputDir, base, err := prepareBuildInput(opts.input, opts.output, opts.target)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	if err := buildCH32fun(source, outputDir, base, opts.optimize); err != nil {
		return err
	}
	fmt.Printf("build: generated %s, %s, %s\n",
		filepath.Join(outputDir, base+".elf"),
		filepath.Join(outputDir, base+".hex"),
		filepath.Join(outputDir, base+".bin"),
	)
	return nil
}

func prepareBuildInput(inputPath, output, target string) (string, string, string, error) {
	ext := strings.ToLower(filepath.Ext(inputPath))
	outputDir, base := resolveBuildOutput(inputPath, output)

	switch ext {
	case ".jika":
		generatedDir := "out"
		generated := filepath.Join(generatedDir, base+".c")
		if err := runTranspile([]string{inputPath, "-o", generated, "--target", target}); err != nil {
			return "", "", "", err
		}
		if err := ensureFunConfig(generatedDir); err != nil {
			return "", "", "", err
		}
		return generated, outputDir, base, nil
	case ".c":
		if err := ensureFunConfig(filepath.Dir(inputPath)); err != nil {
			return "", "", "", err
		}
		return inputPath, outputDir, base, nil
	default:
		return "", "", "", fmt.Errorf("input must be a .jika source file or generated .c file: %s", inputPath)
	}
}

type cliOptions struct {
	input    string
	output   string
	target   string
	optimize string
}

func parseOptions(args []string) (cliOptions, error) {
	opts := cliOptions{
		target:   ch32fun.NewTarget().Name,
		optimize: ch32fun.NewTarget().DefaultOptimizeOpt,
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--target=") {
			opts.target = strings.TrimPrefix(arg, "--target=")
			continue
		}
		if strings.HasPrefix(arg, "--optimize=") {
			opts.optimize = strings.TrimPrefix(arg, "--optimize=")
			continue
		}
		switch arg {
		case "--target":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("missing value for --target")
			}
			i++
			opts.target = args[i]
		case "--optimize":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("missing value for --optimize")
			}
			i++
			opts.optimize = args[i]
		case "-o":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("missing value for -o")
			}
			i++
			opts.output = args[i]
		default:
			if strings.HasPrefix(arg, "-") {
				return opts, fmt.Errorf("unknown option: %s", arg)
			}
			if opts.input == "" {
				opts.input = arg
				continue
			}
			if opts.output == "" {
				opts.output = arg
				continue
			}
			return opts, fmt.Errorf("unexpected argument: %s", arg)
		}
	}
	if opts.optimize != "size" && opts.optimize != "speed" {
		return opts, fmt.Errorf("unsupported optimize mode %q (use size or speed)", opts.optimize)
	}
	return opts, nil
}

func validateTarget(target string) error {
	if target != ch32fun.NewTarget().Name {
		return fmt.Errorf("unsupported target %q", target)
	}
	return nil
}

func resolveBuildOutput(inputPath, output string) (string, string) {
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	if strings.TrimSpace(output) == "" {
		return filepath.Join("build"), base
	}
	trimmed := strings.TrimSpace(output)
	withoutExt := strings.TrimSuffix(trimmed, filepath.Ext(trimmed))
	return filepath.Dir(withoutExt), filepath.Base(withoutExt)
}

func buildCH32fun(sourceC, outputDir, base, optimize string) error {
	ch32funRoot := strings.TrimSpace(os.Getenv("JIKASCRIPT_CH32FUN"))
	if ch32funRoot == "" {
		return fmt.Errorf("JIKASCRIPT_CH32FUN must point to a CH32fun checkout")
	}
	frameworkDir := filepath.Join(ch32funRoot, "ch32fun")
	if err := validateCH32funCheckout(frameworkDir); err != nil {
		return err
	}
	if err := ensureFunConfig(filepath.Dir(sourceC)); err != nil {
		return err
	}

	target := ch32fun.NewTarget()
	cc := toolOrDefault("JIKASCRIPT_RISCV_GCC", target.ToolchainTriple+"-gcc")
	objcopy := toolOrDefault("JIKASCRIPT_OBJCOPY", target.ToolchainTriple+"-objcopy")

	optFlag := "-Os"
	if optimize == "speed" {
		optFlag = "-O2"
	}

	elfOut := filepath.Join(outputDir, base+".elf")
	hexOut := filepath.Join(outputDir, base+".hex")
	binOut := filepath.Join(outputDir, base+".bin")

	intermediateDir := "out"
	if err := os.MkdirAll(intermediateDir, 0o755); err != nil {
		return err
	}
	linkerScript := filepath.Join(intermediateDir, "jikascript_ch32v003.ld")
	if err := generateCH32funLinkerScript(cc, frameworkDir, linkerScript); err != nil {
		return fmt.Errorf("generate CH32fun linker script: %w", err)
	}

	args := ch32funBuildArgs(sourceC, elfOut, frameworkDir, linkerScript, optFlag)
	if err := runCmd(cc, args...); err != nil {
		return fmt.Errorf("risc-v gcc failed: %w", err)
	}

	if err := runCmd(objcopy, "-O", "ihex", elfOut, hexOut); err != nil {
		return fmt.Errorf("objcopy hex failed: %w", err)
	}
	if err := runCmd(objcopy, "-O", "binary", elfOut, binOut); err != nil {
		return fmt.Errorf("objcopy bin failed: %w", err)
	}
	return nil
}

func validateCH32funCheckout(frameworkDir string) error {
	for _, name := range []string{"ch32fun.h", "ch32fun.c", "ch32fun.ld"} {
		path := filepath.Join(frameworkDir, name)
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			return fmt.Errorf("JIKASCRIPT_CH32FUN must point to a CH32fun checkout containing ch32fun/%s", name)
		}
	}
	return nil
}

func generateCH32funLinkerScript(cc, frameworkDir, output string) error {
	return runCmd(cc,
		"-E", "-P", "-x", "c",
		"-DTARGET_MCU=CH32V003",
		"-DMCU_PACKAGE=1",
		"-DTARGET_MCU_LD=0",
		"-DTARGET_MCU_MEMORY_SPLIT=0",
		"-o", output,
		filepath.Join(frameworkDir, "ch32fun.ld"),
	)
}

// ch32funBuildArgs mirrors the CH32fun CH32V003 Makefile integration while
// keeping JikaScript's scheduler and GPIO backend as additional C sources.
func ch32funBuildArgs(sourceC, elfOut, frameworkDir, linkerScript, optFlag string) []string {
	return []string{
		optFlag,
		"-march=rv32ec",
		"-mabi=ilp32e",
		"-DCH32V003",
		"-nostdlib",
		"-static-libgcc",
		"-ffunction-sections",
		"-fdata-sections",
		"-msmall-data-limit=8",
		"-fno-tree-loop-distribute-patterns",
		"-I", filepath.Dir(sourceC),
		"-I", frameworkDir,
		"-I", filepath.Join(".", "runtime", "periodic"),
		"-I", filepath.Join(".", "targets", "ch32fun", "runtime"),
		filepath.Join(frameworkDir, "ch32fun.c"),
		sourceC,
		filepath.Join("runtime", "periodic", "periodic.c"),
		filepath.Join("targets", "ch32fun", "runtime", "gpio.c"),
		filepath.Join("targets", "ch32fun", "runtime", "periodic_ch32fun.c"),
		"-T", linkerScript,
		"-Wl,--gc-sections",
		"-Wl,--print-memory-usage",
		"-L", filepath.Join(filepath.Dir(frameworkDir), "misc"),
		"-lgcc",
		"-o", elfOut,
	}
}

func ensureFunConfig(dir string) error {
	configPath := filepath.Join(dir, "funconfig.h")
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}
	content := `#ifndef _FUNCONFIG_H
#define _FUNCONFIG_H

#define FUNCONF_USE_HSI 1
#define FUNCONF_USE_PLL 0
#define FUNCONF_USE_DEBUGPRINTF 0
#define FUNCONF_USE_UARTPRINTF 0

#endif
`
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(configPath, []byte(content), 0o644)
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func toolOrDefault(envVar, fallback string) string {
	v := strings.TrimSpace(os.Getenv(envVar))
	if v == "" {
		return fallback
	}
	return v
}
