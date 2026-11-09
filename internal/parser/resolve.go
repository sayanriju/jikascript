package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"jikascript/internal/ast"
)

// Loader reads the source of a .jika file given its path.
type Loader func(path string) (string, error)

// DefaultLoader reads files directly from disk.
func DefaultLoader(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ParseFile parses a .jika file from disk, resolving any `include` statements
// so the returned program contains only the flattened top-level statements.
func ParseFile(path string) (*ast.Program, error) {
	return ParseFileWithLoader(path, DefaultLoader)
}

// ParseFileWithLoader parses a .jika file using the supplied loader, resolving
// any `include` statements recursively. Includes are resolved relative to the
// directory of the file that contains them.
func ParseFileWithLoader(path string, loader Loader) (*ast.Program, error) {
	return resolveFile(path, loader, nil)
}

// resolveFile parses one file, splices all included files' top-level
// statements into the result, and reports circular includes via stack.
func resolveFile(path string, loader Loader, stack []string) (*ast.Program, error) {
	for _, active := range stack {
		if sameFile(active, path) {
			return nil, fmt.Errorf("circular include %s: %s", path, strings.Join(append(stack, path), " -> "))
		}
	}
	source, err := loader(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	program, err := ParseProgram(source)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	next := append(stack, path)
	resolved := &ast.Program{Statements: []ast.Stmt{}}
	for _, stmt := range program.Statements {
		include, ok := stmt.(*ast.IncludeStmt)
		if !ok {
			resolved.Statements = append(resolved.Statements, stmt)
			continue
		}
		if !strings.HasSuffix(include.Path, ".jika") {
			return nil, fmt.Errorf("%s: line %d:%d: included file must have a .jika extension: %q", path, include.Pos.Line, include.Pos.Column, include.Path)
		}
		includePath := filepath.Join(filepath.Dir(path), include.Path)
		sub, err := resolveFile(includePath, loader, next)
		if err != nil {
			if strings.HasPrefix(err.Error(), "circular include") {
				return nil, err
			}
			return nil, fmt.Errorf("%s: line %d:%d: %w", path, include.Pos.Line, include.Pos.Column, err)
		}
		resolved.Statements = append(resolved.Statements, sub.Statements...)
	}
	return resolved, nil
}

func sameFile(a, b string) bool {
	absA, errA := filepath.Abs(a)
	absB, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return a == b
	}
	return filepath.Clean(absA) == filepath.Clean(absB)
}