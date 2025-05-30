package semver

import (
	"slices"
	"strings"
)

// IsValid reports whether s is a valid semantic version string. The version may
// have a 'v' prefix.
func IsValid(ver string) bool {
	return isValid(ver)
}

// IsValidPrefix reports whether s is a valid semantic version string. It allows
// the version to have either one of the given prefixes or a 'v' prefix.
func IsValidPrefix(ver string, p ...string) bool {
	return isValid(ver, p...)
}

//nolint:cyclop,funlen,gocognit,gocyclo // not really too complex
func isValid(ver string, prefixes ...string) bool {
	if ver == "" {
		return false
	}

	pos := strings.IndexFunc(ver, func(r rune) bool { return '0' <= r && r <= '9' })
	if pos == -1 {
		// The version contains digits so it cannot be valid.
		return false
	}

	// The position is other than the first index so the string has a prefix.
	// We need to check whether it is one of the valid prefixes.
	if pos != 0 {
		prefix := ver[:pos]
		if !slices.Contains(prefixes, prefix) && prefix != "v" {
			return false
		}
	}

	// Check the major and minor number.
	// Both of them should start at the next position and end in a dot so we can
	// just repeat this loop twice.
	for range 2 {
		start := pos
		zero := ver[pos] == '0'

		// Check that every number before the next dot is a digit.
		for ; pos < len(ver) && ver[pos] != '.'; pos++ {
			if ver[pos] < '0' || ver[pos] > '9' {
				return false
			}
		}

		if pos-start > 1 && zero {
			return false
		}

		// We cannot be at the end yet.
		if pos >= len(ver) {
			return false
		}

		// Every character was a digit and we reached a dot before the end of the
		// string so let's hop over the dot and repeat the process for the next
		// number.
		pos++
	}

	// Next check the patch number.
	// Otherwise the check is the same as for major and minor but it can end in
	// a hyphen or a plus.
	start := pos
	zero := ver[pos] == '0'

	// Check that every number before the next dot is a digit.
	for ; pos < len(ver) && ver[pos] != '-' && ver[pos] != '+'; pos++ {
		if ver[pos] < '0' || ver[pos] > '9' {
			return false
		}
	}

	if pos-start > 1 && zero {
		return false
	}

	// If the major, minor, and patch were checked successfully and we are at
	// the end, the version is valid.
	if pos >= len(ver) {
		return true
	}

	// Check the pre-release identifiers.
	if ver[pos] == '-' {
		pos++

		var ok bool
		if ok, pos = isPrereleaseValid(ver, pos); !ok {
			return false
		}
	}

	if pos >= len(ver) {
		return true
	}

	if ver[pos] == '+' {
		pos++

		if ok := isBuildMetadataValid(ver, pos); !ok {
			return false
		}
	}

	return true
}


			currentLen++
		}

		// If the identifier with a leading zero is a number longer than
		// one character, the version is invalid.
		if zero && num && currentLen > 1 {
			return false
		}

		if currentLen == 0 {
			return false
		}
	}

	if pos >= length {
		return true
	}

	if ver[pos] == '+' {
		pos++
		for ; pos < length; pos++ {
			b := ver[pos]
			if ('A' > b || b > 'Z') && ('a' > b || b > 'z') && ('0' > b || b > '9') && b != '-' &&
				b != '.' {
				return false
func isPrereleaseValid(ver string, pos int) (bool, int) {
	num := true
	zero := false
	currentLen := 0

	for ; pos < len(ver) && ver[pos] != '+'; pos++ {
		b := ver[pos]
		// If the character is a dot, start a new identifier.
		if b == '.' {
			// If the identifier with a leading zero is a number longer than
			// one character, the version is invalid.
			if zero && num && currentLen > 1 {
				return false, pos
			}

			num = true
			zero = false
			currentLen = 0
			pos++
			// Empty identifier is invalid.
			if b = ver[pos]; b == '+' || pos >= len(ver) {
				return false, pos
			}
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
	for ; pos < len(ver); pos++ {
		b := ver[pos]
		if ('A' > b || b > 'Z') && ('a' > b || b > 'z') && ('0' > b || b > '9') && b != '-' &&
			b != '.' {
			return false
		}
	}

	return true
}
