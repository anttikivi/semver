// Copyright (c) 2025 Antti Kivi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package semver

import "fmt"

// Values for number mode.
const (
	dot numberMode = iota
	end
	loose
)

// Values for validation mode.
const (
	strict validationMode = iota
	lax
)

// A numberMode is a helper parameter type for parsing the numbers in
// the version strings.
type numberMode int

// A validationMode is a helper parameter type for telling if the validation
// parser should be lax about the version number.
type validationMode int

// IsValid reports whether s is a valid semantic version string. The version may
// have a 'v' prefix.
func IsValid(s string) bool {
	return isValid(s, strict)
}

// IsValidLax reports whether s is a valid semantic version string even if it is
// only a partial version. In other words, this function reads `v1` and `v1.2`
// as valid versions. The version may have a 'v' prefix.
func IsValidLax(s string) bool {
	return isValid(s, lax)
}

func isValid(s string, mode validationMode) bool {
	ok, pos := isStartValid(s)
	if !ok {
		return false
	}

	if ok, pos = isCoreValid(s, pos, mode); !ok {
		return false
	}

	if pos >= len(s) {
		return true
	}

	if len(s)-pos < 2 { //nolint:mnd // not enough characters left
		return false
	}

	if s[pos] != '-' && s[pos] != '+' {
		return false
	}

	// Check the pre-release identifiers.
	if s[pos] == '-' {
		pos++

		var ok bool
		if ok, pos = isPrereleaseValid(s, pos); !ok {
			return false
		}
	}

	if pos >= len(s) {
		return true
	}

	if s[pos] == '+' {
		pos++

		if ok := isBuildMetadataValid(s, pos); !ok {
			return false
		}
	}

	return true
}

// isStartValid checks if the start of the version string is valid. The function
// checks the prefix and that the string actually contains numbers. It returns
// the status of the check and the current index.
func isStartValid(ver string) (bool, int) {
	if ver == "" {
		return false, 0
	}

	pos := 0

	c := ver[0]
	if !isDigit(c) && c != 'v' {
		return false, pos
	}

	if c == 'v' {
		pos++
	}

	if pos == len(ver) {
		return false, pos
	}

	return true, pos
}

func isCoreValid(s string, pos int, mode validationMode) (bool, int) {
	var ok bool

	numMode := dot
	if mode == lax {
		numMode = loose
	}

	if ok, pos = isVersionNumberValid(s, pos, numMode); !ok {
		return false, pos
	}

	if pos >= len(s) {
		return mode == lax, pos
	}

	if mode == lax && (s[pos] == '-' || s[pos] == '+') {
		return true, pos
	}

	if s[pos] != '.' {
		return false, pos
	}

	pos++
	if ok, pos = isVersionNumberValid(s, pos, numMode); !ok {
		return false, pos
	}

	if pos >= len(s) {
		return mode == lax, pos
	}

	if mode == lax && (s[pos] == '-' || s[pos] == '+') {
		return true, pos
	}

	if s[pos] != '.' {
		return false, pos
	}

	if mode == strict {
		numMode = end
	}

	pos++
	if ok, pos = isVersionNumberValid(s, pos, numMode); !ok {
		return false, pos
	}

	// Only major and minor version numbers are provided.
	if pos >= len(s) {
		return true, pos
	}

	return true, pos
}

// isVersionNumberValid reports whether the next number in the version string is
// valid. It returns the status of the check and the new position. If final is
// set to true, the number may end in a "-" or a "+" instead of a ".".
func isVersionNumberValid(s string, pos int, mode numberMode) (bool, int) {
	start := pos

	// Check that every number before the next character that ends it is
	// a digit.
	switch mode {
	case dot:
		for ; pos < len(s) && s[pos] != '.'; pos++ {
			if !isDigit(s[pos]) {
				return false, pos
			}
		}
	case end:
		for ; pos < len(s) && s[pos] != '-' && s[pos] != '+'; pos++ {
			if !isDigit(s[pos]) {
				return false, pos
			}
		}
	case loose:
		for ; pos < len(s) && s[pos] != '.' && s[pos] != '-' && s[pos] != '+'; pos++ {
			if !isDigit(s[pos]) {
				return false, pos
			}
		}
	default:
		panic(fmt.Sprintf("invalid number mode: %d", mode))
	}

	if pos == start {
		return false, pos
	}

	if pos-start > 1 && s[start] == '0' {
		return false, pos
	}

	return true, pos
}

func isPrereleaseValid(ver string, pos int) (bool, int) { //nolint:cyclop,gocognit // no problem
	num := true
	zero := false
	currentLen := 0

	for ; pos < len(ver) && ver[pos] != '+'; pos++ {
		b := ver[pos]
		// If the character is a dot, start a new identifier.
		if b == '.' {
			if currentLen == 0 {
				return false, pos
			}

			// If the identifier with a leading zero is a number longer than
			// one character, the version is invalid.
			if zero && num && currentLen > 1 {
				return false, pos
			}

			num = true
			zero = false
			currentLen = 0

			if pos+1 >= len(ver) || ver[pos+1] == '+' || ver[pos+1] == '.' {
				return false, pos + 1
			}

			continue
		}

		if b == '0' && currentLen == 0 {
			zero = true
		}

		// If the identifier is still a number but we encounter a non-digit
		// character, the identifier is no longer a number.
		if num && ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') || b == '-' {
			num = false
		}

		// Otherwise just check that the character is valid.
		if ('A' > b || b > 'Z') && ('a' > b || b > 'z') && ('0' > b || b > '9') && b != '-' {
			return false, pos
		}

		currentLen++
	}

	// If the identifier with a leading zero is a number longer than
	// one character, the version is invalid.
	if zero && num && currentLen > 1 {
		return false, pos
	}

	if currentLen == 0 {
		return false, pos
	}

	return true, pos
}

func isBuildMetadataValid(ver string, pos int) bool {
	if pos >= len(ver) {
		return false
	}

	l := 0

	for ; pos < len(ver); pos++ {
		c := ver[pos]
		if c == '.' {
			if l == 0 {
				return false
			}

			l = 0

			if pos == len(ver)-1 {
				return false
			}

			continue
		}

		if !isIdentifierCharacter(c) {
			return false
		}

		l++
	}

	return l > 0
}
