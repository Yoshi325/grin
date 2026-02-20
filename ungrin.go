package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ungrinKVPair holds a key-value pair for INI output.
type ungrinKVPair struct {
	key   string
	value string
}

// ungrinStatements reads grin assignment lines from r and returns them as
// parsed statements.
func ungrinStatements(r io.Reader) (statements, error) {
	scanner := bufio.NewScanner(r)
	var ss statements

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and grep separator lines
		if trimmed == "" || trimmed == "--" {
			continue
		}

		s, err := lexStatement(trimmed)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return ss, nil
}

// ungrinFromStatements converts parsed grin statements back into INI format.
func ungrinFromStatements(ss statements, w io.Writer) error {
	globalKeys, sections, sectionOrder, emptyOnlySections := indexStatements(ss)

	first := true

	for _, kv := range globalKeys {
		if _, err := fmt.Fprintf(w, "%s = %s\n", kv.key, kv.value); err != nil {
			return err
		}
		first = false
	}

	for _, secName := range sectionOrder {
		kvs := sections[secName]
		if len(kvs) == 0 && !emptyOnlySections[secName] {
			continue
		}
		if !first {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		first = false
		if _, err := fmt.Fprintf(w, "[%s]\n", secName); err != nil {
			return err
		}
		for _, kv := range kvs {
			if _, err := fmt.Fprintf(w, "%s = %s\n", kv.key, kv.value); err != nil {
				return err
			}
		}
	}

	return nil
}

// indexStatements processes parsed statements and returns the data structures
// needed to reconstruct INI output: global keys, per-section key lists,
// the original section order, and the set of sections that were declared
// empty (i.e., only appeared as `section = {};` with no key assignments).
func indexStatements(ss statements) (
	globalKeys []ungrinKVPair,
	sections map[string][]ungrinKVPair,
	sectionOrder []string,
	emptyOnlySections map[string]bool,
) {
	sections = make(map[string][]ungrinKVPair)
	seenSections := make(map[string]bool)
	emptyOnlySections = make(map[string]bool)

	for _, s := range ss {
		path, value, isObj := extractPathAndValue(s)
		if len(path) == 0 {
			continue
		}

		if isObj {
			secName := strings.Join(path, ".")
			if !seenSections[secName] {
				seenSections[secName] = true
				sectionOrder = append(sectionOrder, secName)
				emptyOnlySections[secName] = true
			}
			continue
		}

		if len(path) == 1 {
			globalKeys = append(globalKeys, ungrinKVPair{key: path[0], value: value})
		} else {
			secName := strings.Join(path[:len(path)-1], ".")
			key := path[len(path)-1]
			if !seenSections[secName] {
				seenSections[secName] = true
				sectionOrder = append(sectionOrder, secName)
			}
			delete(emptyOnlySections, secName)
			sections[secName] = append(sections[secName], ungrinKVPair{key: key, value: value})
		}
	}
	return
}

// extractPathAndValue extracts the path components and value from a statement.
// Returns (path, value, isObject). path is nil if the statement should be skipped.
func extractPathAndValue(s statement) (path []string, value string, isObj bool) {
	// Walk tokens to extract bare identifiers (the path) and the value
	var parts []string
	var val string
	foundEquals := false
	isObject := false

	for _, t := range s {
		switch t.typ {
		case typBare:
			if !foundEquals {
				parts = append(parts, t.text)
			}
		case typEquals:
			foundEquals = true
		case typString:
			if foundEquals {
				val = unquoteString(t.text)
			}
		case typEmptyObject:
			if foundEquals {
				isObject = true
			}
		case typDot, typSemi, typIgnored:
			// skip punctuation
		}
	}

	if len(parts) == 0 {
		return nil, "", false
	}

	// Remove the root "ini" prefix
	if parts[0] == "ini" {
		parts = parts[1:]
	}

	return parts, val, isObject
}

// Lexer for parsing grin assignment lines back into tokens.

type lexer struct {
	input  string
	pos    int
	width  int
	tokens []token
}

func newLexer(input string) *lexer {
	return &lexer{input: input}
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return -1
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += w
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) emit(typ tokenTyp, text string) {
	l.tokens = append(l.tokens, token{text: text, typ: typ})
}

func (l *lexer) skipWhitespace() {
	for {
		r := l.next()
		if r == -1 || !unicode.IsSpace(r) {
			l.backup()
			return
		}
	}
}

// lexStatement parses a single grin assignment line into tokens.
func lexStatement(line string) (statement, error) {
	l := newLexer(line)

	// Parse: Path = Value ;
	// Path: BareWord ( "." BareWord )*
	// Value: String | "{}"

	l.skipWhitespace()

	// First bare word
	if err := l.lexBareWord(); err != nil {
		return nil, fmt.Errorf("parsing statement: %w", err)
	}

	// More path components
	for {
		l.skipWhitespace()
		r := l.peek()
		if r == '.' {
			l.next()
			l.emit(typDot, ".")
			l.skipWhitespace()
			if err := l.lexBareWord(); err != nil {
				return nil, fmt.Errorf("parsing path: %w", err)
			}
		} else {
			break
		}
	}

	// Equals
	l.skipWhitespace()
	r := l.next()
	if r != '=' {
		return nil, fmt.Errorf("expected '=', got %q", string(r))
	}
	l.emit(typEquals, " = ")

	// Value
	l.skipWhitespace()
	if err := l.lexValue(); err != nil {
		return nil, fmt.Errorf("parsing value: %w", err)
	}

	// Semicolon
	l.skipWhitespace()
	r = l.next()
	if r != ';' {
		return nil, fmt.Errorf("expected ';', got %q", string(r))
	}
	l.emit(typSemi, ";")

	return statement(l.tokens), nil
}

func (l *lexer) lexBareWord() error {
	start := l.pos
	r := l.next()
	if r == -1 || (!unicode.IsLetter(r) && r != '_') {
		return fmt.Errorf("expected identifier, got %q", string(r))
	}

	for {
		r = l.next()
		if r == -1 || (!unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-') {
			l.backup()
			break
		}
	}

	l.emit(typBare, l.input[start:l.pos])
	return nil
}

func (l *lexer) lexValue() error {
	r := l.peek()

	switch r {
	case '"':
		return l.lexString()
	case '{':
		return l.lexBraces()
	default:
		return fmt.Errorf("expected value, got %q", string(r))
	}
}

func (l *lexer) lexString() error {
	start := l.pos
	r := l.next() // consume opening "
	if r != '"' {
		return fmt.Errorf("expected '\"', got %q", string(r))
	}

	for {
		r = l.next()
		if r == -1 {
			return fmt.Errorf("unterminated string")
		}
		if r == '\\' {
			// Skip escaped character
			l.next()
			continue
		}
		if r == '"' {
			break
		}
	}

	l.emit(typString, l.input[start:l.pos])
	return nil
}

func (l *lexer) lexBraces() error {
	r := l.next()
	if r != '{' {
		return fmt.Errorf("expected '{', got %q", string(r))
	}
	r = l.next()
	if r != '}' {
		return fmt.Errorf("expected '}', got %q", string(r))
	}
	l.emit(typEmptyObject, "{}")
	return nil
}
