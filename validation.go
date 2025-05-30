package semver

import (
	"fmt"
	"slices"
	"strings"
)

// Values for number mode.
const (
	dot numberMode = iota
	ending
	loose
)

// A numberMode is a helper parameter type for parsing the numbers in
// the version strings.
type numberMode int

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

// IsValidPartial reports whether s is a valid semantic version string even if
// it is only a partial version. In other words, this function reads `v1` and
// `v1.2` as valid versions. The version may have a 'v' prefix.
func IsValidPartial(ver string) bool {
	return isValidPartial(ver)
}

// IsValidPartialPrefix reports whether s is a valid semantic version string
// even if it is only a partial version. In other words, this function reads
// `v1` and `v1.2` as valid versions. It allows the version to have either one
// of the given prefixes or a 'v' prefix.
func IsValidPartialPrefix(ver string, p ...string) bool {
	return isValidPartial(ver, p...)
}

func isValid(ver string, prefixes ...string) bool {
	ok, pos := isStartValid(ver, prefixes...)
	if !ok {
		return false
	}

	// Check the major and minor number.
	// Both of them should start at the next position and end in a dot so we can
	// just repeat this loop twice.
	for range 2 {
		if ok, pos = isVersionNumberValid(ver, pos, dot); !ok {
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

	// Next check the patch number. Otherwise the check is the same as for major
	// and minor but it can end in a hyphen or a plus.
	if ok, pos = isVersionNumberValid(ver, pos, ending); !ok {
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

func isValidPartial(ver string, prefixes ...string) bool {
	ok, pos := isStartValid(ver, prefixes...)
	if !ok {
		return false
	}

	if ok, pos = isVersionNumberValid(ver, pos, loose); !ok {
		return false
	}

	// Only major version number is provided.
	if pos >= len(ver) {
		return true
	}

	// The next character is a dot so the minor version number should be
	// provided. We can branch off to checking it.
	if ver[pos] == '.' {
		pos++
		if ok, pos = isVersionNumberValid(ver, pos, loose); !ok {
			return false
		}

		// Only major and minor version numbers are provided.
		if pos >= len(ver) {
			return true
		}
	}

	// The next character is a dot so the patch version number should be
	// provided. We can branch off to checking it.
	if ver[pos] == '.' {
		pos++
		if ok, pos = isVersionNumberValid(ver, pos, loose); !ok {
			return false
		}

		// Only major and minor version numbers are provided.
		if pos >= len(ver) {
			return true
		}
	}

	// Check the pre-release identifiers.
	if ver[pos] == '-' {
		// Skip the hyphen.
		pos++

		if ok, pos = isPrereleaseValid(ver, pos); !ok {
			return false
		}
	}

	if pos >= len(ver) {
		return true
	}

	if ver[pos] == '+' {
		pos++

		if ok = isBuildMetadataValid(ver, pos); !ok {
			return false
		}
	}

	return true
}

// isStartValid checks if the start of the version string is valid. The function
// checks the prefix and that the string actually contains numbers. It returns
// the status of the check and the current index.
func isStartValid(ver string, prefixes ...string) (bool, int) {
	if ver == "" {
		return false, 0
	}

	pos := strings.IndexFunc(ver, func(r rune) bool { return '0' <= r && r <= '9' })
	if pos == -1 {
		// The version does not contain digits so it cannot be valid.
		return false, pos
	}

	// The first number was found at position other than the first so
	// the version string has a prefix. We need to check if it is one of
	// the valid prefixes.
	if pos != 0 {
		prefix := ver[:pos]
		if !slices.Contains(prefixes, prefix) && prefix != "v" {
			return false, pos
		}
	}

	return true, pos
}

// isVersionNumberValid reports whether the next number in the version string is
// valid. It returns the status of the check and the new position. If final is
// set to true, the number may end in a "-" or a "+" instead of a ".".
//
//nolint:cyclop // no problem
func isVersionNumberValid(ver string, pos int, mode numberMode) (bool, int) {
	start := pos

	// Check that every number before the next character that ends it is
	// a digit.
	switch mode {
	case dot:
		for ; pos < len(ver) && ver[pos] != '.'; pos++ {
			if ver[pos] < '0' || ver[pos] > '9' {
				return false, pos
			}
		}
	case ending:
		for ; pos < len(ver) && ver[pos] != '-' && ver[pos] != '+'; pos++ {
			if ver[pos] < '0' || ver[pos] > '9' {
				return false, pos
			}
		}
	case loose:
		for ; pos < len(ver) && ver[pos] != '.' && ver[pos] != '-' && ver[pos] != '+'; pos++ {
			if ver[pos] < '0' || ver[pos] > '9' {
				return false, pos
			}
		}
	default:
		panic(fmt.Sprintf("invalid number mode: %d", mode))
	}

	if pos-start > 1 && ver[start] == '0' {
		return false, pos
	}

	return true, pos
}

func isPrereleaseValid(ver string, pos int) (bool, int) { //nolint:cyclop // no problem
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
