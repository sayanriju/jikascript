package parser

import (
	"strings"
	"testing"
)

func memLoader(files map[string]string) Loader {
	return func(path string) (string, error) {
		src, ok := files[path]
		if !ok {
			return "", &notFoundError{path}
		}
		return src, nil
	}
}

type notFoundError struct{ path string }

func (e *notFoundError) Error() string { return "open " + e.path + ": no such file or directory" }

func TestResolveSingleInclude(t *testing.T) {
	loader := memLoader(map[string]string{
		"/proj/app.jika":     `include "helpers.jika";
task("blink", 100, () => { led.toggle(); });
scheduler.start();`,
		"/proj/helpers.jika": `const led = pin("PC1");
led.output();`,
	})
	program, err := ParseFileWithLoader("/proj/app.jika", loader)
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 flattened statements, got %d", len(program.Statements))
	}
}

func TestResolveTransitiveInclude(t *testing.T) {
	loader := memLoader(map[string]string{
		"/proj/a.jika": `include "b.jika";`,
		"/proj/b.jika": `include "c.jika";`,
		"/proj/c.jika": `const shared = 10;`,
	})
	program, err := ParseFileWithLoader("/proj/a.jika", loader)
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 flattened statement, got %d", len(program.Statements))
	}
}

func TestResolveRelativeSubdirectory(t *testing.T) {
	loader := memLoader(map[string]string{
		"/proj/app.jika":        `include "lib/helpers.jika";`,
		"/proj/lib/helpers.jika": `include "utils.jika";`,
		"/proj/lib/utils.jika":   `const x = 1;
const y = 2;`,
	})
	program, err := ParseFileWithLoader("/proj/app.jika", loader)
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 flattened statements, got %d", len(program.Statements))
	}
}

func TestResolveMissingFile(t *testing.T) {
	loader := memLoader(map[string]string{
		"/proj/app.jika": `include "missing.jika";`,
	})
	_, err := ParseFileWithLoader("/proj/app.jika", loader)
	if err == nil {
		t.Fatalf("expected error for missing include, got nil")
	}
	if !strings.Contains(err.Error(), "missing.jika") {
		t.Fatalf("error should mention missing file, got %v", err)
	}
}

func TestResolveNonJikaExtension(t *testing.T) {
	loader := memLoader(map[string]string{
		"/proj/app.jika": `include "helpers.txt";`,
	})
	_, err := ParseFileWithLoader("/proj/app.jika", loader)
	if err == nil {
		t.Fatalf("expected error for non-.jika include, got nil")
	}
	if !strings.Contains(err.Error(), ".jika") {
		t.Fatalf("error should mention .jika extension, got %v", err)
	}
}

func TestResolveCycleDetected(t *testing.T) {
	loader := memLoader(map[string]string{
		"/proj/a.jika": `include "b.jika";`,
		"/proj/b.jika": `include "a.jika";`,
	})
	_, err := ParseFileWithLoader("/proj/a.jika", loader)
	if err == nil {
		t.Fatalf("expected circular include error, got nil")
	}
	if !strings.Contains(err.Error(), "circular") {
		t.Fatalf("error should mention circular include, got %v", err)
	}
}

func TestResolveNestedIncludeRejectedByParser(t *testing.T) {
	loader := memLoader(map[string]string{
		"/proj/app.jika":    `function f() { include "helpers.jika"; }`,
		"/proj/helpers.jika": `const x = 1;`,
	})
	_, err := ParseFileWithLoader("/proj/app.jika", loader)
	if err == nil {
		t.Fatalf("expected error for nested include, got nil")
	}
}