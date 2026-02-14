package main

import (
	"sort"
	"testing"
)

func TestStatementToString(t *testing.T) {
	s := statement{
		{text: "ini", typ: typBare},
		{text: ".", typ: typDot},
		{text: "section", typ: typBare},
		{text: ".", typ: typDot},
		{text: "key", typ: typBare},
		{text: " = ", typ: typEquals},
		{text: `"value"`, typ: typString},
		{text: ";", typ: typSemi},
	}

	want := `ini.section.key = "value";`
	if got := statementToString(s); got != want {
		t.Errorf("statementToString() = %q, want %q", got, want)
	}
}

func TestStatementWithBare(t *testing.T) {
	base := statement{{text: "ini", typ: typBare}}
	got := base.withBare("section")
	want := "ini.section"

	result := statementToString(got)
	if result != want {
		t.Errorf("withBare() produced %q, want %q", result, want)
	}
}

func TestStatementWithPath(t *testing.T) {
	base := statement{{text: "ini", typ: typBare}}

	tests := []struct {
		path string
		want string
	}{
		{"section", "ini.section"},
		{"section.sub", "ini.section.sub"},
		{"a.b.c", "ini.a.b.c"},
	}

	for _, tt := range tests {
		got := statementToString(base.withPath(tt.path))
		if got != tt.want {
			t.Errorf("withPath(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestStatementWithEmptyObject(t *testing.T) {
	base := statement{{text: "ini", typ: typBare}}
	got := statementToString(base.withEmptyObject())
	want := "ini = {};"
	if got != want {
		t.Errorf("withEmptyObject() = %q, want %q", got, want)
	}
}

func TestStatementWithStringValue(t *testing.T) {
	base := statement{{text: "ini", typ: typBare}}.withBare("key")
	got := statementToString(base.withStringValue("hello"))
	want := `ini.key = "hello";`
	if got != want {
		t.Errorf("withStringValue() = %q, want %q", got, want)
	}
}

func TestStatementsSort(t *testing.T) {
	ss := statements{
		// ini.b = "2";
		makeAssignment("ini", "b", "2"),
		// ini.a = "1";
		makeAssignment("ini", "a", "1"),
		// ini = {};
		makeEmptyObj("ini"),
	}

	sort.Sort(ss)

	want := []string{
		`ini = {};`,
		`ini.a = "1";`,
		`ini.b = "2";`,
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("sorted[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementsSortSections(t *testing.T) {
	ss := statements{
		makeAssignment("ini", "z.key", "val"),
		makeAssignment("ini", "a.key", "val"),
		makeSectionObj("ini", "z"),
		makeSectionObj("ini", "a"),
		makeEmptyObj("ini"),
	}

	sort.Sort(ss)

	want := []string{
		`ini = {};`,
		`ini.a = {};`,
		`ini.a.key = "val";`,
		`ini.z = {};`,
		`ini.z.key = "val";`,
	}

	for i, w := range want {
		got := statementToString(ss[i])
		if got != w {
			t.Errorf("sorted[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestStatementWithBareDoesNotMutateOriginal(t *testing.T) {
	base := statement{{text: "ini", typ: typBare}}
	_ = base.withBare("section")

	if len(base) != 1 {
		t.Errorf("withBare mutated original statement: len = %d, want 1", len(base))
	}
}

// helpers

func makeAssignment(root, path, value string) statement {
	base := statement{{text: root, typ: typBare}}
	return base.withPath(path).withStringValue(value)
}

func makeEmptyObj(root string) statement {
	base := statement{{text: root, typ: typBare}}
	return base.withEmptyObject()
}

func makeSectionObj(root, section string) statement {
	base := statement{{text: root, typ: typBare}}
	return base.withPath(section).withEmptyObject()
}
