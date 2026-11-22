// Package naming provides identifier normalization shared by the checker and
// the C code generator so both agree on which source names collide.
package naming

import "strings"

// NormalizeIdentifier converts a source-level name into a C-identifier-shaped
// string. Characters that are not letters, digits, or underscores become
// underscores, and a leading digit is prefixed with an underscore. Case is
// preserved because C identifiers are case-sensitive.
func NormalizeIdentifier(name string) string {
	var b strings.Builder
	for i, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		if i == 0 && r >= '0' && r <= '9' {
			b.WriteString("_")
			b.WriteRune(r)
			continue
		}
		b.WriteRune('_')
	}
	return b.String()
}
