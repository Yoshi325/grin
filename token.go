package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

type tokenTyp int

const (
	typBare        tokenTyp = iota // Bare identifier: ini, section, key
	typDot                         // .
	typEquals                      // =
	typSemi                        // ;
	typString                      // "quoted value"
	typEmptyObject                 // {}
	typIgnored                     // grep separator lines like --
	typError                       // parse error
)

type token struct {
	text string
	typ  tokenTyp
}

func (t token) isValue() bool {
	return t.typ == typString || t.typ == typEmptyObject
}

func (t token) isPunct() bool {
	return t.typ == typDot || t.typ == typEquals || t.typ == typSemi
}

var (
	bareColor  = color.New(color.FgBlue, color.Bold)
	strColor   = color.New(color.FgYellow)
	braceColor = color.New(color.FgMagenta)
	punctColor = color.New(color.FgRed)
)

func (t token) format() string {
	return t.text
}

func (t token) formatColor() string {
	switch t.typ {
	case typBare:
		return bareColor.Sprint(t.text)
	case typString:
		return strColor.Sprint(t.text)
	case typEmptyObject:
		return braceColor.Sprint(t.text)
	case typDot, typEquals, typSemi:
		return punctColor.Sprint(t.text)
	default:
		return t.text
	}
}

// quoteString returns a JSON-style quoted string for an INI value.
func quoteString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20:
			fmt.Fprintf(&b, `\u%04x`, r)
		default:
			b.WriteRune(r)
		}
		i += size
	}
	b.WriteByte('"')
	return b.String()
}

// unquoteString removes surrounding quotes and unescapes a quoted string value.
func unquoteString(s string) string {
	if len(s) < 2 {
		return s
	}
	if s[0] != '"' || s[len(s)-1] != '"' {
		return s
	}
	s = s[1 : len(s)-1]

	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '"':
				b.WriteByte('"')
				i += 2
			case '\\':
				b.WriteByte('\\')
				i += 2
			case 'n':
				b.WriteByte('\n')
				i += 2
			case 'r':
				b.WriteByte('\r')
				i += 2
			case 't':
				b.WriteByte('\t')
				i += 2
			default:
				b.WriteByte(s[i])
				i++
			}
		} else {
			r, size := utf8.DecodeRuneInString(s[i:])
			b.WriteRune(r)
			i += size
		}
	}
	return b.String()
}
