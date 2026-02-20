package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// iniKVPair holds a parsed key-value pair from an INI line.
type iniKVPair struct {
	key   string
	value string
}

// statementsFromINI reads INI data from r and produces a slice of grin
// statements. The prefix is the root statement (typically just "ini").
func statementsFromINI(r io.Reader, prefix statement) (statements, error) {
	scanner := bufio.NewScanner(r)

	var (
		currentSection string
		sectionOrder   []string
		sectionKeys    = make(map[string][]iniKVPair)
		globalKeys     []iniKVPair
		seenSections   = make(map[string]bool)
		lineNum        int
		firstLine      = true
	)

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if firstLine {
			line = stripBOM(line)
			firstLine = false
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed[0] == ';' || trimmed[0] == '#' {
			continue
		}

		if trimmed[0] == '[' {
			sectionName, err := parseSectionHeader(trimmed, lineNum)
			if err != nil {
				return nil, err
			}
			currentSection = sectionName
			if !seenSections[sectionName] {
				seenSections[sectionName] = true
				sectionOrder = append(sectionOrder, sectionName)
			}
			continue
		}

		key, value, err := parseINIKeyValue(trimmed, lineNum)
		if err != nil {
			return nil, err
		}

		pair := iniKVPair{key: key, value: value}
		if currentSection == "" {
			globalKeys = append(globalKeys, pair)
		} else {
			sectionKeys[currentSection] = append(sectionKeys[currentSection], pair)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	return buildINIStatements(prefix, globalKeys, sectionOrder, sectionKeys), nil
}

// buildINIStatements assembles the final statement slice from the parsed INI
// data: a root object, global keys, and ordered sections with their keys.
func buildINIStatements(prefix statement, globalKeys []iniKVPair, sectionOrder []string, sectionKeys map[string][]iniKVPair) statements {
	ss := statements{}
	ss = append(ss, prefix.withEmptyObject())

	for _, kv := range globalKeys {
		ss = append(ss, prefix.withBare(kv.key).withStringValue(kv.value))
	}

	emitted := make(map[string]bool)
	for _, secName := range sectionOrder {
		parts := strings.Split(secName, ".")
		for i := 1; i <= len(parts); i++ {
			partial := strings.Join(parts[:i], ".")
			if !emitted[partial] {
				emitted[partial] = true
				ss = append(ss, prefix.withPath(partial).withEmptyObject())
			}
		}
		for _, kv := range sectionKeys[secName] {
			ss = append(ss, prefix.withPath(secName).withBare(kv.key).withStringValue(kv.value))
		}
	}

	return ss
}

// parseSectionHeader parses a trimmed section-header line like "[a.b.c]".
// Returns the section name or an error.
func parseSectionHeader(trimmed string, lineNum int) (string, error) {
	end := strings.IndexByte(trimmed, ']')
	if end == -1 {
		return "", fmt.Errorf("line %d: unclosed section header", lineNum)
	}
	sectionName := strings.TrimSpace(trimmed[1:end])
	if sectionName == "" {
		return "", fmt.Errorf("line %d: empty section name", lineNum)
	}
	for _, p := range strings.Split(sectionName, ".") {
		if !validIdentifier(p) {
			return "", fmt.Errorf("line %d: invalid section name part %q", lineNum, p)
		}
	}
	return sectionName, nil
}

// parseINIKeyValue parses a trimmed key=value line.
// Returns (key, stripped-value) or an error.
func parseINIKeyValue(trimmed string, lineNum int) (string, string, error) {
	eqIdx := strings.IndexByte(trimmed, '=')
	if eqIdx == -1 {
		return "", "", fmt.Errorf("line %d: expected key = value, got %q", lineNum, trimmed)
	}
	key := strings.TrimSpace(trimmed[:eqIdx])
	value := strings.TrimSpace(trimmed[eqIdx+1:])
	if key == "" {
		return "", "", fmt.Errorf("line %d: empty key", lineNum)
	}
	if !validIdentifier(key) {
		return "", "", fmt.Errorf("line %d: invalid key %q", lineNum, key)
	}
	return key, stripINIQuotes(value), nil
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
