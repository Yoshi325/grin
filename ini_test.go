package main

import (
	"strings"
	"testing"
)

func TestStatementsFromINISimple(t *testing.T) {
	input := `[section]
key1 = value1
key2 = value2
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		`ini = {};`,
		`ini.section = {};`,
		`ini.section.key1 = "value1";`,
		`ini.section.key2 = "value2";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsFromINIGlobalKeys(t *testing.T) {
	input := `name = grin
version = 1.0

[section]
key = value
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		`ini = {};`,
		`ini.name = "grin";`,
		`ini.version = "1.0";`,
		`ini.section = {};`,
		`ini.section.key = "value";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsFromINIDottedSections(t *testing.T) {
	input := `[database]
host = localhost

[database.pool]
max = 10

[database.pool.overflow]
enabled = true
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		`ini = {};`,
		`ini.database = {};`,
		`ini.database.host = "localhost";`,
		`ini.database.pool = {};`,
		`ini.database.pool.max = "10";`,
		`ini.database.pool.overflow = {};`,
		`ini.database.pool.overflow.enabled = "true";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d\nGot:", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsFromINIComments(t *testing.T) {
	input := `; This is a semicolon comment
# This is a hash comment
[section]
; inline section comment
key = value
# another comment
key2 = value2
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		`ini = {};`,
		`ini.section = {};`,
		`ini.section.key = "value";`,
		`ini.section.key2 = "value2";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsFromINIQuotedValues(t *testing.T) {
	input := `[section]
double = "hello world"
single = 'hello world'
unquoted = hello world
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		`ini = {};`,
		`ini.section = {};`,
		`ini.section.double = "hello world";`,
		`ini.section.single = "hello world";`,
		`ini.section.unquoted = "hello world";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsFromINIEmpty(t *testing.T) {
	input := ""
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ss) != 1 {
		t.Fatalf("got %d statements, want 1", len(ss))
	}

	want := `ini = {};`
	got := statementToString(ss[0])
	if got != want {
		t.Errorf("statement[0] = %q, want %q", got, want)
	}
}

func TestStatementsFromINIEmptyValue(t *testing.T) {
	input := `[section]
key =
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		`ini = {};`,
		`ini.section = {};`,
		`ini.section.key = "";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsFromINIValueWithEquals(t *testing.T) {
	input := `[section]
key = val=ue
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should split on first '=' only
	want := `ini.section.key = "val=ue";`
	got := statementToString(ss[2])
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStatementsFromINIBOM(t *testing.T) {
	input := "\xEF\xBB\xBF[section]\nkey = value\n"
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `ini.section.key = "value";`
	got := statementToString(ss[2])
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStatementsFromINIDuplicateSections(t *testing.T) {
	input := `[section]
key1 = value1

[section]
key2 = value2
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Duplicate sections should merge keys
	want := []string{
		`ini = {};`,
		`ini.section = {};`,
		`ini.section.key1 = "value1";`,
		`ini.section.key2 = "value2";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsFromINIEmptySection(t *testing.T) {
	input := `[empty]

[notempty]
key = value
`
	prefix := statement{{text: "ini", typ: typBare}}
	ss, err := statementsFromINI(strings.NewReader(input), prefix)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		`ini = {};`,
		`ini.empty = {};`,
		`ini.notempty = {};`,
		`ini.notempty.key = "value";`,
	}

	if len(ss) != len(want) {
		t.Fatalf("got %d statements, want %d", len(ss), len(want))
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("statement[%d] = %q, want %q", i, got, w)
		}
	}
}

// Error cases

func TestStatementsFromINIUnclosedSection(t *testing.T) {
	input := "[section\nkey = value\n"
	prefix := statement{{text: "ini", typ: typBare}}
	_, err := statementsFromINI(strings.NewReader(input), prefix)
	if err == nil {
		t.Fatal("expected error for unclosed section, got nil")
	}
}

func TestStatementsFromINIEmptySectionName(t *testing.T) {
	input := "[]\nkey = value\n"
	prefix := statement{{text: "ini", typ: typBare}}
	_, err := statementsFromINI(strings.NewReader(input), prefix)
	if err == nil {
		t.Fatal("expected error for empty section name, got nil")
	}
}

func TestStatementsFromINIInvalidKey(t *testing.T) {
	input := "[section]\n123 = value\n"
	prefix := statement{{text: "ini", typ: typBare}}
	_, err := statementsFromINI(strings.NewReader(input), prefix)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

func TestStatementsFromININoEquals(t *testing.T) {
	input := "[section]\nnotavalidline\n"
	prefix := statement{{text: "ini", typ: typBare}}
	_, err := statementsFromINI(strings.NewReader(input), prefix)
	if err == nil {
		t.Fatal("expected error for line without =, got nil")
	}
}

func TestStatementsFromINIEmptyKey(t *testing.T) {
	input := "[section]\n = value\n"
	prefix := statement{{text: "ini", typ: typBare}}
	_, err := statementsFromINI(strings.NewReader(input), prefix)
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestStatementsFromINIInvalidSectionName(t *testing.T) {
	input := "[has space]\nkey = value\n"
	prefix := statement{{text: "ini", typ: typBare}}
	_, err := statementsFromINI(strings.NewReader(input), prefix)
	if err == nil {
		t.Fatal("expected error for invalid section name, got nil")
	}
}

func TestStripBOM(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"\xEF\xBB\xBFhello", "hello"},
		{"hello", "hello"},
		{"", ""},
		{"\xEF\xBB\xBF", ""},
	}

	for _, tt := range tests {
		if got := stripBOM(tt.input); got != tt.want {
			t.Errorf("stripBOM(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestStripINIQuotes(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`"hello"`, "hello"},
		{`'hello'`, "hello"},
		{"hello", "hello"},
		{`""`, ""},
		{`''`, ""},
		{`"mismatched'`, `"mismatched'`},
		{`'mismatched"`, `'mismatched"`},
		{`"`, `"`},
	}

	for _, tt := range tests {
		if got := stripINIQuotes(tt.input); got != tt.want {
			t.Errorf("stripINIQuotes(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
