package docsites

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestDocsifyContentLinksUseRootRelativePaths guards against page-relative or
// hash-prefixed content links, both of which Docsify renders as broken routes
// (for example `#/current-page?id=%2fother-page`).
func TestDocsifyContentLinksUseRootRelativePaths(t *testing.T) {
	root := filepath.Clean("../..")
	linkPattern := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)

	pages, err := filepath.Glob(filepath.Join(root, "docs", "*.md"))
	if err != nil {
		t.Fatalf("glob docs root: %v", err)
	}
	nested, err := filepath.Glob(filepath.Join(root, "docs", "*", "*.md"))
	if err != nil {
		t.Fatalf("glob docs nested: %v", err)
	}
	pages = append(pages, nested...)

	for _, page := range pages {
		if filepath.Base(page) == "_sidebar.md" {
			continue // sidebar links are resolved by Docsify from the docs root
		}
		content, err := os.ReadFile(page)
		if err != nil {
			t.Fatalf("read %s: %v", page, err)
		}
		for lineNo, line := range strings.Split(string(content), "\n") {
			for _, match := range linkPattern.FindAllStringSubmatch(line, -1) {
				target := strings.TrimSpace(match[1])
				if strings.HasPrefix(target, "#/") {
					t.Errorf("%s:%d: use root-relative .md links, not hash routes: %s", page, lineNo, target)
				}
				if strings.HasPrefix(target, "../") {
					t.Errorf("%s:%d: use root-relative .md links, not parent paths: %s", page, lineNo, target)
				}
				if strings.HasPrefix(target, "#") || strings.Contains(target, "://") {
					continue
				}
				if !strings.HasSuffix(target, ".md") {
					continue
				}
				resolved := filepath.Clean(filepath.Join(root, "docs", strings.SplitN(target, "#", 2)[0]))
				if _, err := os.Stat(resolved); err != nil {
					t.Errorf("%s:%d: link target %s does not resolve to %s: %v", page, lineNo, target, resolved, err)
				}
			}
		}
	}
}
