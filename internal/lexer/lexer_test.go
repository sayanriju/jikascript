package lexer

import (
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
