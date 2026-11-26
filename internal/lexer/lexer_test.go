package lexer

import (
	"strings"
	"testing"
)

func lexAll(t *testing.T, source string) []Token {
	t.Helper()
	l := New(source)
	var tokens []Token
	for {
		tok := l.NextToken()
		if tok.Type == EOF {
			return tokens
		}
		tokens = append(tokens, tok)
		if tok.Type == ILLEGAL {
			return tokens
		}
	}
}

func TestStringEscapeSequences(t *testing.T) {
	tokens := lexAll(t, `"a\tb\nc\\d\"e"`)
	if len(tokens) != 1 || tokens[0].Type != STRING {
		t.Fatalf("expected single STRING token, got %+v", tokens)
	}
	want := "a\tb\nc\\d\"e"
	if tokens[0].Literal != want {
		t.Fatalf("literal = %q, want %q", tokens[0].Literal, want)
	}
}

func TestStringUnknownEscapeIsIllegal(t *testing.T) {
	tokens := lexAll(t, `"bad \q escape"`)
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != ILLEGAL {
		t.Fatalf("expected ILLEGAL token for unknown escape, got %+v", tokens)
	}
}

func TestStringEscapedQuoteDoesNotTerminate(t *testing.T) {
	tokens := lexAll(t, `"say \"hi\""`)
	if len(tokens) != 1 || tokens[0].Type != STRING {
		t.Fatalf("expected single STRING token, got %+v", tokens)
	}
	if tokens[0].Literal != `say "hi"` {
		t.Fatalf("literal = %q, want %q", tokens[0].Literal, `say "hi"`)
	}
}

func TestIncludeIsKeyword(t *testing.T) {
	tokens := lexAll(t, `include`)
	if len(tokens) != 1 || tokens[0].Type != KEYWORD || tokens[0].Literal != "include" {
		t.Fatalf("expected KEYWORD token 'include', got %+v", tokens)
	}
}

func TestPeripheralBuiltinsAreKeywords(t *testing.T) {
	for _, word := range []string{"pin", "uart", "pwm", "adc"} {
		tokens := lexAll(t, word)
		if len(tokens) != 1 || tokens[0].Type != KEYWORD || tokens[0].Literal != word {
			t.Fatalf("expected KEYWORD token %q, got %+v", word, tokens)
		}
	}
}

func TestNumberLiterals(t *testing.T) {
	tokens := lexAll(t, "12 3.14 0.5")
	if len(tokens) != 3 {
		t.Fatalf("expected 3 NUMBER tokens, got %+v", tokens)
	}
	for i, want := range []string{"12", "3.14", "0.5"} {
		if tokens[i].Type != NUMBER || tokens[i].Literal != want {
			t.Fatalf("token %d = %+v, want NUMBER %q", i, tokens[i], want)
		}
	}
}

func TestHexNumberLiterals(t *testing.T) {
	tokens := lexAll(t, "0xFF 0X1a 0x0")
	if len(tokens) != 3 {
		t.Fatalf("expected 3 NUMBER tokens, got %+v", tokens)
	}
	for i, want := range []string{"0xFF", "0X1a", "0x0"} {
		if tokens[i].Type != NUMBER || tokens[i].Literal != want {
			t.Fatalf("token %d = %+v, want NUMBER %q", i, tokens[i], want)
		}
	}
}

func TestInvalidHexNumberLiterals(t *testing.T) {
	cases := []struct {
		source string
		want   string
	}{
		{"0x", "hex literal requires at least one digit"},
		{"0xG", "hex literal requires at least one digit"},
		{"0xFFg", "invalid hex literal"},
		{"0xFF_", "invalid hex literal"},
		{"0x1.5", "hex floating-point literals are not supported"},
	}
	for _, tc := range cases {
		tokens := lexAll(t, tc.source)
		if len(tokens) == 0 || tokens[len(tokens)-1].Type != ILLEGAL {
			t.Fatalf("expected ILLEGAL for %q, got %+v", tc.source, tokens)
		}
		if !strings.Contains(tokens[len(tokens)-1].Literal, tc.want) {
			t.Fatalf("expected %q for %q, got %+v", tc.want, tc.source, tokens)
		}
	}
}

func TestOperators(t *testing.T) {
	tokens := lexAll(t, "=== !== += <<= >>=")
	if len(tokens) != 5 {
		t.Fatalf("expected 5 OPERATOR tokens, got %+v", tokens)
	}
	for i, want := range []string{"===", "!==", "+=", "<<=", ">>="} {
		if tokens[i].Type != OPERATOR || tokens[i].Literal != want {
			t.Fatalf("token %d = %+v, want OPERATOR %q", i, tokens[i], want)
		}
	}
}

func TestWhitespaceAndCommentsAreSkipped(t *testing.T) {
	tokens := lexAll(t, "let x = 1; // trailing\n/* block */\nlet y = 2;")
	want := []struct {
		typ     TokenType
		literal string
	}{
		{KEYWORD, "let"}, {IDENT, "x"}, {OPERATOR, "="}, {NUMBER, "1"}, {PUNCT, ";"},
		{KEYWORD, "let"}, {IDENT, "y"}, {OPERATOR, "="}, {NUMBER, "2"}, {PUNCT, ";"},
	}
	if len(tokens) != len(want) {
		t.Fatalf("expected %d tokens, got %+v", len(want), tokens)
	}
	for i, w := range want {
		if tokens[i].Type != w.typ || tokens[i].Literal != w.literal {
			t.Fatalf("token %d = %+v, want %s %q", i, tokens[i], w.typ, w.literal)
		}
	}
}

func TestIllegalCharacter(t *testing.T) {
	tokens := lexAll(t, "@")
	if len(tokens) != 1 || tokens[0].Type != ILLEGAL || tokens[0].Literal != "@" {
		t.Fatalf("expected ILLEGAL '@', got %+v", tokens)
	}
}
