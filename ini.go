package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// statementsFromINI reads INI data from r and produces a slice of grin
// statements. The prefix is the root statement (typically just "ini").
func statementsFromINI(r io.Reader, prefix statement) (statements, error) {
	scanner := bufio.NewScanner(r)

	type kvPair struct {
		key   string
		value string
	}

	var (
		currentSection string
		sectionOrder   []string
		sectionKeys    = make(map[string][]kvPair)
		globalKeys     []kvPair
		seenSections   = make(map[string]bool)
		lineNum        int
		firstLine      = true
	)

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Strip UTF-8 BOM from the first line
		if firstLine {
			line = stripBOM(line)
			firstLine = false
		}

		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || trimmed[0] == ';' || trimmed[0] == '#' {
			continue
		}

		// Section header
		if trimmed[0] == '[' {
			end := strings.IndexByte(trimmed, ']')
			if end == -1 {
				return nil, fmt.Errorf("line %d: unclosed section header", lineNum)
			}
			sectionName := strings.TrimSpace(trimmed[1:end])
			if sectionName == "" {
				return nil, fmt.Errorf("line %d: empty section name", lineNum)
			}

			// Validate each part of the section name
			parts := strings.Split(sectionName, ".")
			for _, p := range parts {
				if !validIdentifier(p) {
					return nil, fmt.Errorf("line %d: invalid section name part %q", lineNum, p)
				}
			}

			currentSection = sectionName
			if !seenSections[sectionName] {
				seenSections[sectionName] = true
				sectionOrder = append(sectionOrder, sectionName)
			}
			continue
		}

		// Key-value pair: split on first '='
		eqIdx := strings.IndexByte(trimmed, '=')
		if eqIdx == -1 {
			return nil, fmt.Errorf("line %d: expected key = value, got %q", lineNum, trimmed)
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		value := strings.TrimSpace(trimmed[eqIdx+1:])

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}
		if !validIdentifier(key) {
			return nil, fmt.Errorf("line %d: invalid key %q", lineNum, key)
		}

		// Strip surrounding quotes from value
		value = stripINIQuotes(value)

		pair := kvPair{key: key, value: value}
		if currentSection == "" {
			globalKeys = append(globalKeys, pair)
		} else {
			sectionKeys[currentSection] = append(sectionKeys[currentSection], pair)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	// Generate statements
	ss := statements{}

	// Root: ini = {};
	ss = append(ss, prefix.withEmptyObject())

	// Global keys
	for _, kv := range globalKeys {
		ss = append(ss, prefix.withBare(kv.key).withStringValue(kv.value))
	}

	// Track emitted section objects to avoid duplicates
	emitted := make(map[string]bool)

	for _, secName := range sectionOrder {
		// Emit intermediate objects for dotted section names
		parts := strings.Split(secName, ".")
		for i := 1; i <= len(parts); i++ {
			partial := strings.Join(parts[:i], ".")
			if !emitted[partial] {
				emitted[partial] = true
				ss = append(ss, prefix.withPath(partial).withEmptyObject())
			}
		}

		// Emit key-value pairs
		for _, kv := range sectionKeys[secName] {
			ss = append(ss, prefix.withPath(secName).withBare(kv.key).withStringValue(kv.value))
		}
	}

	return ss, nil
}

// stripBOM removes a UTF-8 BOM from the beginning of a string if present.
func stripBOM(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == 0xFEFF && size > 0 {
		return s[size:]
	}
	return s
}

// stripINIQuotes removes surrounding single or double quotes from an INI value.
func stripINIQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
