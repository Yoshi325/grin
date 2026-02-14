package main

import "strings"

// statement is an ordered list of tokens representing one grin assignment line.
// e.g. ini.section.key = "value";
type statement []token

// statements is a sortable slice of statement.
type statements []statement

func (ss statements) Len() int {
	return len(ss)
}

func (ss statements) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func (ss statements) Less(i, j int) bool {
	a := ss[i]
	b := ss[j]
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	for k := 0; k < minLen; k++ {
		at := a[k]
		bt := b[k]

		// Equals always comes first (so "ini = {}" sorts before "ini.x = ...")
		if at.typ == typEquals && bt.typ != typEquals {
			return true
		}
		if bt.typ == typEquals && at.typ != typEquals {
			return false
		}

		if at.text != bt.text {
			return at.text < bt.text
		}
	}

	return len(a) < len(b)
}

// withBare appends a dot separator and a bare identifier token.
func (s statement) withBare(key string) statement {
	new := make(statement, len(s), len(s)+2)
	copy(new, s)
	return append(new,
		token{".", typDot},
		token{key, typBare},
	)
}

// withPath appends a dotted path (e.g. "section.subsection") as a series
// of dot-separated bare tokens.
func (s statement) withPath(dotted string) statement {
	parts := strings.Split(dotted, ".")
	cur := s
	for _, p := range parts {
		cur = cur.withBare(p)
	}
	return cur
}

// withValue appends an equals sign, a value token, and a semicolon.
func (s statement) withValue(t token) statement {
	new := make(statement, len(s), len(s)+4)
	copy(new, s)
	return append(new,
		token{" = ", typEquals},
		t,
		token{";", typSemi},
	)
}

// withEmptyObject appends " = {};" to the statement.
func (s statement) withEmptyObject() statement {
	return s.withValue(token{"{}", typEmptyObject})
}

// withStringValue appends " = \"value\";" to the statement.
func (s statement) withStringValue(val string) statement {
	return s.withValue(token{quoteString(val), typString})
}

// statementToString renders a statement as a plain string.
func statementToString(s statement) string {
	var b strings.Builder
	for _, t := range s {
		b.WriteString(t.format())
	}
	return b.String()
}

// statementToColorString renders a statement with ANSI color codes.
func statementToColorString(s statement) string {
	var b strings.Builder
	for _, t := range s {
		b.WriteString(t.formatColor())
	}
	return b.String()
}

// statementconv is a function type for converting statements to strings.
type statementconv func(statement) string
