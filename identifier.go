package main

import "unicode"

// validIdentifier reports whether s is a valid bare-word identifier
// for use in grin output. Valid identifiers start with a letter or
// underscore, and contain only letters, digits, underscores, or hyphens.
func validIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
				return false
			}
		}
	}
	return true
}
