package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestLexStatement(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`ini = {};`, `ini = {};`},
		{`ini.section = {};`, `ini.section = {};`},
		{`ini.section.key = "value";`, `ini.section.key = "value";`},
		{`ini.key = "hello world";`, `ini.key = "hello world";`},
		{`ini.key = "";`, `ini.key = "";`},
		{`ini.a.b.c = "deep";`, `ini.a.b.c = "deep";`},
		{`ini.my-key = "value";`, `ini.my-key = "value";`},
	}

	for _, tt := range tests {
		s, err := lexStatement(tt.input)
		if err != nil {
			t.Errorf("lexStatement(%q) error: %v", tt.input, err)
			continue
		}
		got := statementToString(s)
		if got != tt.want {
			t.Errorf("lexStatement(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestLexStatementErrors(t *testing.T) {
	tests := []struct {
		input string
		desc  string
	}{
		{"", "empty input"},
		{"123", "starts with digit"},
		{"ini.section", "missing equals"},
		{`ini = "no semi"`, "missing semicolon"},
		{`ini = unterminated`, "invalid value"},
		{`ini = "unterminated;`, "unterminated string"},
	}

	for _, tt := range tests {
		_, err := lexStatement(tt.input)
		if err == nil {
			t.Errorf("lexStatement(%q) [%s]: expected error, got nil", tt.input, tt.desc)
		}
	}
}

func TestExtractPathAndValue(t *testing.T) {
	tests := []struct {
		input string
		path  []string
		value string
		isObj bool
	}{
		{`ini = {};`, nil, "", true},                               // root — path becomes empty after removing "ini"
		{`ini.section = {};`, []string{"section"}, "", true},       // section declaration
		{`ini.key = "val";`, []string{"key"}, "val", false},        // global key
		{`ini.s.key = "val";`, []string{"s", "key"}, "val", false}, // section key
		{`ini.a.b.c = "d";`, []string{"a", "b", "c"}, "d", false}, // deep path
	}

	for _, tt := range tests {
		s, err := lexStatement(tt.input)
		if err != nil {
			t.Fatalf("lexStatement(%q) error: %v", tt.input, err)
		}
		path, value, isObj := extractPathAndValue(s)

		if tt.input == `ini = {};` {
			// Root: path should be empty after removing "ini"
			if len(path) != 0 {
				t.Errorf("extractPathAndValue(%q): path = %v, want empty", tt.input, path)
			}
			if !isObj {
				t.Errorf("extractPathAndValue(%q): isObj = false, want true", tt.input)
			}
			continue
		}

		if len(path) != len(tt.path) {
			t.Errorf("extractPathAndValue(%q): path len = %d, want %d", tt.input, len(path), len(tt.path))
			continue
		}
		for i, p := range tt.path {
			if path[i] != p {
				t.Errorf("extractPathAndValue(%q): path[%d] = %q, want %q", tt.input, i, path[i], p)
			}
		}
		if value != tt.value {
			t.Errorf("extractPathAndValue(%q): value = %q, want %q", tt.input, value, tt.value)
		}
		if isObj != tt.isObj {
			t.Errorf("extractPathAndValue(%q): isObj = %v, want %v", tt.input, isObj, tt.isObj)
		}
	}
}

func TestUngrinSimple(t *testing.T) {
	input := `ini = {};
ini.section = {};
ini.section.key1 = "value1";
ini.section.key2 = "value2";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "[section]\nkey1 = value1\nkey2 = value2\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}

func TestUngrinGlobalKeys(t *testing.T) {
	input := `ini = {};
ini.name = "grin";
ini.version = "1.0";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "name = grin\nversion = 1.0\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}

func TestUngrinMixed(t *testing.T) {
	input := `ini = {};
ini.app = "myapp";
ini.database = {};
ini.database.host = "localhost";
ini.database.port = "5432";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "app = myapp\n\n[database]\nhost = localhost\nport = 5432\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}

func TestUngrinEmptySection(t *testing.T) {
	input := `ini = {};
ini.empty = {};
ini.notempty = {};
ini.notempty.key = "value";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "[empty]\n\n[notempty]\nkey = value\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}

func TestUngrinDottedSections(t *testing.T) {
	input := `ini = {};
ini.database = {};
ini.database.host = "localhost";
ini.database.pool = {};
ini.database.pool.max = "10";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "[database]\nhost = localhost\n\n[database.pool]\nmax = 10\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}

func TestUngrinGrepSeparators(t *testing.T) {
	// Simulate grep output with -- separators
	input := `ini.section.key1 = "value1";
--
ini.section.key2 = "value2";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "[section]\nkey1 = value1\nkey2 = value2\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}

func TestUngrinFromFilteredOutput(t *testing.T) {
	// Simulate: grin file.ini | grep host | grin -u
	input := `ini.database.host = "localhost";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "[database]\nhost = localhost\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}

func TestUngrinEscapedStringValues(t *testing.T) {
	input := `ini = {};
ini.section = {};
ini.section.key = "value with \"quotes\"";
`
	ss, err := ungrinStatements(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ungrinStatements error: %v", err)
	}

	var buf bytes.Buffer
	if err := ungrinFromStatements(ss, &buf); err != nil {
		t.Fatalf("ungrinFromStatements error: %v", err)
	}

	want := "[section]\nkey = value with \"quotes\"\n"
	if buf.String() != want {
		t.Errorf("ungrin output = %q, want %q", buf.String(), want)
	}
}
