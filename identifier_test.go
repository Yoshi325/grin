package main

import "testing"

func TestValidIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		// Valid identifiers
		{"key", true},
		{"section", true},
		{"my_key", true},
		{"_private", true},
		{"key1", true},
		{"my-key", true},
		{"camelCase", true},
		{"PascalCase", true},
		{"with123numbers", true},
		{"a", true},
		{"café", true},
		{"日本語", true},

		// Invalid identifiers
		{"", false},
		{"1starts-with-digit", false},
		{"-starts-with-hyphen", false},
		{"has space", false},
		{"has.dot", false},
		{"has=equals", false},
		{"has;semi", false},
		{"has[bracket", false},
		{"has]bracket", false},
		{"has\"quote", false},
	}

	for _, tt := range tests {
		if got := validIdentifier(tt.input); got != tt.want {
			t.Errorf("validIdentifier(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
