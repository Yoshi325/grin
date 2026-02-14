package main

import (
	"testing"
)

func TestTokenIsValue(t *testing.T) {
	tests := []struct {
		tok  token
		want bool
	}{
		{token{"\"hello\"", typString}, true},
		{token{"{}", typEmptyObject}, true},
		{token{"ini", typBare}, false},
		{token{".", typDot}, false},
		{token{"=", typEquals}, false},
		{token{";", typSemi}, false},
		{token{"--", typIgnored}, false},
		{token{"err", typError}, false},
	}

	for _, tt := range tests {
		if got := tt.tok.isValue(); got != tt.want {
			t.Errorf("token{%q, %d}.isValue() = %v, want %v", tt.tok.text, tt.tok.typ, got, tt.want)
		}
	}
}

func TestTokenIsPunct(t *testing.T) {
	tests := []struct {
		tok  token
		want bool
	}{
		{token{".", typDot}, true},
		{token{"=", typEquals}, true},
		{token{";", typSemi}, true},
		{token{"ini", typBare}, false},
		{token{"\"hello\"", typString}, false},
		{token{"{}", typEmptyObject}, false},
	}

	for _, tt := range tests {
		if got := tt.tok.isPunct(); got != tt.want {
			t.Errorf("token{%q, %d}.isPunct() = %v, want %v", tt.tok.text, tt.tok.typ, got, tt.want)
		}
	}
}

func TestTokenFormat(t *testing.T) {
	tok := token{"ini", typBare}
	if got := tok.format(); got != "ini" {
		t.Errorf("format() = %q, want %q", got, "ini")
	}
}

func TestQuoteString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", `"hello"`},
		{"", `""`},
		{`say "hi"`, `"say \"hi\""`},
		{"line\nbreak", `"line\nbreak"`},
		{"tab\there", `"tab\there"`},
		{`back\slash`, `"back\\slash"`},
		{"carriage\rreturn", `"carriage\rreturn"`},
		{"control\x01char", `"control\u0001char"`},
		{"unicode café", `"unicode café"`},
	}

	for _, tt := range tests {
		if got := quoteString(tt.input); got != tt.want {
			t.Errorf("quoteString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestUnquoteString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`"hello"`, "hello"},
		{`""`, ""},
		{`"say \"hi\""`, `say "hi"`},
		{`"line\nbreak"`, "line\nbreak"},
		{`"tab\there"`, "tab\there"},
		{`"back\\slash"`, `back\slash`},
		// Not quoted — returned as-is
		{"noquotes", "noquotes"},
		// Too short
		{`"`, `"`},
	}

	for _, tt := range tests {
		if got := unquoteString(tt.input); got != tt.want {
			t.Errorf("unquoteString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestQuoteUnquoteRoundTrip(t *testing.T) {
	values := []string{
		"hello world",
		`has "quotes" inside`,
		"has\nnewlines\nand\ttabs",
		`back\slash`,
		"",
		"unicode: 日本語",
	}

	for _, v := range values {
		quoted := quoteString(v)
		unquoted := unquoteString(quoted)
		if unquoted != v {
			t.Errorf("round-trip failed: %q -> %q -> %q", v, quoted, unquoted)
		}
	}
}
