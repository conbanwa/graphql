// Package ident provides functions for parsing and converting identifier names
// between various naming convention. It has support for MixedCaps, lowerCamelCase,
// and SCREAMING_SNAKE_CASE naming conventions.
package ident

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// ParseMixedCaps parses a MixedCaps identifier name.
//
// E.g., "ClientMutationID" -> {"Client", "Mutation", "ID"}.
func ParseMixedCaps(name string) Name {
	var words Name

	// Split name at any lower -> Upper or Upper -> Upper,lower transitions.
	// Check each word for initialisms.
	runes := []rune(name)
	w, i := 0, 0 // Index of start of word, scan.
	for i+1 <= len(runes) {
		eow := false // Whether we hit the end of a word.
		if i+1 == len(runes) {
			eow = true
		} else if unicode.IsLower(runes[i]) && unicode.IsUpper(runes[i+1]) {
			// lower -> Upper.
			eow = true
		} else if i+2 < len(runes) && unicode.IsUpper(runes[i]) && unicode.IsUpper(runes[i+1]) && unicode.IsLower(runes[i+2]) {
			// Upper -> Upper,lower. End of acronym, followed by a word.
			eow = true

			if string(runes[i:i+3]) == "IDs" { // Special case, plural form of ID initialism.
				eow = false
			}
		}
		i++
		if !eow {
			continue
		}

		// [w, i) is a word.
		word := string(runes[w:i])
		if initialism, ok := isInitialism(word); ok {
			words = append(words, initialism)
		} else if i1, i2, ok := isTwoInitialisms(word); ok {
			words = append(words, i1, i2)
		} else {
			words = append(words, word)
		}
		w = i
	}
	return words
}

// ParseLowerCamelCase parses a lowerCamelCase identifier name.
//
// E.g., "clientMutationId" -> {"client", "Mutation", "Id"}.
func ParseLowerCamelCase(name string) Name {
	var words Name

	// Split name at any Upper letters.
	runes := []rune(name)
	w, i := 0, 0 // Index of start of word, scan.
	for i+1 <= len(runes) {
		eow := false // Whether we hit the end of a word.
		if i+1 == len(runes) {
			eow = true
		} else if unicode.IsUpper(runes[i+1]) {
			// Upper letter.
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w, i) is a word.
		words = append(words, string(runes[w:i]))
		w = i
	}
	return words
}

// ParseScreamingSnakeCase parses a SCREAMING_SNAKE_CASE identifier name.
//
// E.g., "CLIENT_MUTATION_ID" -> {"CLIENT", "MUTATION", "ID"}.
func ParseScreamingSnakeCase(name string) Name {
	var words Name

	// Split name at '_' characters.
	runes := []rune(name)
	w, i := 0, 0 // Index of start of word, scan.
	for i+1 <= len(runes) {
		eow := false // Whether we hit the end of a word.
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// Underscore.
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w, i) is a word.
		words = append(words, string(runes[w:i]))
		if i < len(runes) && runes[i] == '_' {
			// Skip underscore.
			i++
		}
		w = i
	}
	return words
}

// Name is an identifier name, broken up into individual words.
type Name []string

// ToMixedCaps expresses identifer name in MixedCaps naming convention.
//
// E.g., "ClientMutationID".
func (n Name) ToMixedCaps() string {
	for i, word := range n {
		if strings.EqualFold(word, "IDs") { // Special case, plural form of ID initialism.
			n[i] = "IDs"
			continue
		}
		if initialism, ok := isInitialism(word); ok {
			n[i] = initialism
			continue
		}
		if brand, ok := isBrand(word); ok {
			n[i] = brand
			continue
		}
		r, size := utf8.DecodeRuneInString(word)
		n[i] = string(unicode.ToUpper(r)) + strings.ToLower(word[size:])
	}
	return strings.Join(n, "")
}

// ToLowerCamelCase expresses identifer name in lowerCamelCase naming convention.
//
// E.g., "clientMutationId".
func (n Name) ToLowerCamelCase() string {
	for i, word := range n {
		if i == 0 {
			n[i] = strings.ToLower(word)
			continue
		}
		r, size := utf8.DecodeRuneInString(word)
		n[i] = string(unicode.ToUpper(r)) + strings.ToLower(word[size:])
	}
	return strings.Join(n, "")
}

// isInitialism reports whether word is an initialism.
func isInitialism(word string) (string, bool) {
	initialism := strings.ToUpper(word)
	_, ok := initialisms[initialism]
	return initialism, ok
}

// isTwoInitialisms reports whether word is two initialisms.
func isTwoInitialisms(word string) (string, string, bool) {
	word = strings.ToUpper(word)
	for i := 2; i <= len(word)-2; i++ { // Shortest initialism is 2 characters long.
		_, ok1 := initialisms[word[:i]]
		_, ok2 := initialisms[word[i:]]
		if ok1 && ok2 {
			return word[:i], word[i:], true
		}
	}
	return "", "", false
}

// initialisms is the set of initialisms in the MixedCaps naming convention.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var initialisms = map[string]struct{}{
	"ACL":   {},
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"HTTP":  {},
	"HTTPS": {},
	"ID":    {},
	"IP":    {},
	"JSON":  {},
	"LHS":   {},
	"QPS":   {},
	"RAM":   {},
	"RHS":   {},
	"RPC":   {},
	"RSS":   {},
	"SLA":   {},
	"SMTP":  {},
	"SQL":   {},
	"SSH":   {},
	"TCP":   {},
	"TLS":   {},
	"TTL":   {},
	"UDP":   {},
	"UI":    {},
	"UID":   {},
	"URI":   {},
	"URL":   {},
	"UTF8":  {},
	"UUID":  {},
	"VM":    {},
	"XML":   {},
	"XMPP":  {},
	"XSRF":  {},
	"XSS":   {},
}

// isBrand reports whether word is a brand.
func isBrand(word string) (string, bool) {
	brand, ok := brands[strings.ToLower(word)]
	return brand, ok
}

// brands is the map of brands in the MixedCaps naming convention;
// see https://dmitri.shuralyov.com/idiomatic-go#for-brands-or-words-with-more-than-1-capital-letter-lowercase-all-letters.
// Key is the lower case version of the brand, value is the canonical brand spelling.
// Only add entries that are highly unlikely to be non-brands.
var brands = map[string]string{
	"github": "GitHub",
	"gitlab": "GitLab",
	"devops": "DevOps", // For https://en.wikipedia.org/wiki/DevOps.
	// For https://docs.github.com/en/graphql/reference/enums#fundingplatform.
	"issuehunt": "IssueHunt",
	"lfx":       "LFX",
}

func LintName(name string) (should string) {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u, ok := isInitialism(word); ok {
			// Keep consistent case, which is lowercase only at the start.
			// if w == 0 && unicode.IsLower(runes[w]) {
			u = strings.ToLower(u)
			// }
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			u := strings.ToUpper(u[:1]) + (u[1:])
			copy(runes[w:], []rune(u))
		} else
		 if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
}
