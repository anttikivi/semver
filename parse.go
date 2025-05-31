package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidVersion is the error returned by the version parsing functions when
// they encounter invalid version string.
var ErrInvalidVersion = errors.New("invalid semantic version")

// ErrUnknown is returned when there is a problem with the parsing that is not
// directly related to the caller giving an invalid string.
var ErrUnknown = errors.New("parsing failed")

// Parse parses the given string into a Version. The version string may have
// a 'v' prefix.
func Parse(ver string) (*Version, error) {
	v, err := parse(ver, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	return v, nil
}

// MustParse parses the given string into a Version and panics if it encounters
// an error. The version string may have a 'v' prefix.
func MustParse(ver string) *Version {
	v, err := Parse(ver)
	if err != nil {
		panic(fmt.Sprintf("failed to parse the string %q into a version: %v", ver, err))
	}

	return v
}

// ParseLax parses the given string into a Version. The version number may be
// partial, i.e. it parses 'v1' into '1.0.0' and 'v1.2' into '1.2.0'.
// The version string may have a 'v' prefix.
func ParseLax(ver string) (*Version, error) {
	v, err := parse(ver, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	return v, nil
}

// MustParseLax parses the given string into a Version and panics if it
// encounters an error. The version string number may be partial, i.e. it parses
// 'v1' into '1.0.0' and 'v1.2' into '1.2.0'. The version may have a 'v' prefix.
func MustParseLax(ver string) *Version {
	v, err := ParseLax(ver)
	if err != nil {
		panic(fmt.Sprintf("failed to parse the string %q into a version: %v", ver, err))
	}

	return v
}

// NewBuildIdentifiers returns new [BuildIdentifiers] for the given strings.
func NewBuildIdentifiers(s ...string) BuildIdentifiers {
	b := make(BuildIdentifiers, 0, len(s))
	b = append(b, s...)

	return b
}

// checkPrefix parses the possible "v" prefix for the version string.
// The function returns the new position where the parsing continues.
func checkPrefix(s string) (int, error) {
	pos := 0

	c := s[0]
	if !isDigit(c) && c != 'v' {
		return pos, fmt.Errorf(
			"%w: version %q does not start with a digit or 'v'",
			ErrInvalidVersion,
			s,
		)
	}

	if c == 'v' {
		pos++
	}

	if pos == len(s) {
		return pos, fmt.Errorf("%w: %q", ErrInvalidVersion, s)
	}

	return pos, nil
}

func countDigits(i uint64) int {
	if i == 0 {
		return 1
	}

	count := 0

	for i != 0 {
		i /= 10
		count++
	}

	return count
}

func isAlphanumericIdentifier(s string) bool {
	for i := range len(s) {
		if !isIdentifierCharacter(s[i]) {
			return false
		}
	}

	return true
}

func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isIdentifierCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || c == '-'
}

func isNumericIdentifier(s string) bool {
	for i := range len(s) {
		if !isDigit(s[i]) {
			return false
		}
	}

	return true
}

func isPrereleaseSeparator(b byte) bool {
	return b == '.' || b == '+'
}

func parse(ver string, minCore int) (*Version, error) {
	if ver == "" {
		return nil, fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	if !isASCII(ver) {
		return nil, fmt.Errorf("%w: version contains non-ASCII characters", ErrInvalidVersion)
	}

	pos, err := checkPrefix(ver)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the version prefix: %w", err)
	}

	i := len(ver)

	for j := range ver[pos:] {
		c := ver[pos+j]
		if !isDigit(c) && c != '.' {
			i = pos + j

			break
		}
	}

	nums := strings.Split(ver[pos:i], ".")

	if len(nums) > 3 {
		return nil, fmt.Errorf("%w: too many core version numbers in %q", ErrInvalidVersion, ver)
	}

	if len(nums) < minCore {
		return nil, fmt.Errorf("%w: not enough core version numbers in %q", ErrInvalidVersion, ver)
	}

	major := uint64(0)
	minor := uint64(0)
	patch := uint64(0)

	for j, n := range nums {
		if n == "" {
			return nil, fmt.Errorf("%w: empty version number in %q", ErrInvalidVersion, ver)
		}

		if !isNumericIdentifier(n) {
			return nil, fmt.Errorf("%w: version number %q is not a number", ErrInvalidVersion, n)
		}

		// If the number is only a zero, we already have the correct value and
		// can move on to the next one.
		if n == "0" {
			continue
		}

		if n[0] == '0' {
			return nil, fmt.Errorf("%w: leading zero in %q", ErrInvalidVersion, n)
		}

		u, err := strconv.ParseUint(n, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert the string %q to uint64: %w", n, err)
		}

		switch j {
		case 0:
			major = u
		case 1:
			minor = u
		case 2:
			patch = u
		default:
			return nil, fmt.Errorf("%w: index when checking version number is out of bounds: %d", ErrUnknown, j)
		}
	}

	pos = i

	if pos < len(ver) && ver[pos] != '-' && ver[pos] != '+' {
		return nil, fmt.Errorf("%w: invalid char %q at %d", ErrInvalidVersion, ver[pos], pos)
	}

	var prerelease Prerelease

	if pos < len(ver) && ver[pos] == '-' {
		// The hyphen is not passed to the parser.
		pos++

		i = len(ver)

		for j := range ver[pos:] {
			c := ver[pos+j]
			if c == '+' {
				i = pos + j

				break
			}
		}

		prerelease, err = ParsePrerelease(ver[pos:i])
		if err != nil {
			return nil, fmt.Errorf("failed to parse the pre-release: %w", err)
		}

		pos = i
	}

	var build BuildIdentifiers

	if pos < len(ver) && ver[pos] == '+' {
		// Move past the '+'.
		pos++

		build, err = parseBuild(ver[pos:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse the build identifiers: %w", err)
		}
	}

	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: prerelease,
		Build:      build,
	}, nil
}

func parseBuild(s string) ([]string, error) {
	if s == "" {
		return nil, fmt.Errorf("cannot parse empty string as a build: %w", ErrInvalidVersion)
	}

	result := strings.Split(s, ".")
	for _, v := range result {
		if s == "" {
			return nil, fmt.Errorf(
				"empty string as a dot-separated build identifier: %w",
				ErrInvalidVersion,
			)
		}

		// This should be safe as all of the characters in the version must be
		// ASCII.
		if strings.ContainsFunc(
			v,
			func(r rune) bool { return !isIdentifierCharacter(byte(r)) },
		) {
			return nil, fmt.Errorf(
				"invalid rune in the build identifier %q: %w",
				v,
				ErrInvalidVersion,
			)
		}
	}

	return result, nil
}

// parseNext parses the next integer from the given string. The string should be
// a version string or the next part to parse from a version string adhering to
// the semantic versioning.
func parseNext(s string) (uint64, error) {
	if s == "" {
		return 0, fmt.Errorf("cannot parse empty string as int: %w", ErrInvalidVersion)
	}

	b := s[0]
	if b < '0' || '9' < b {
		return 0, fmt.Errorf("first character is not a digit: %w", ErrInvalidVersion)
	}

	i := 1
	for i < len(s) && '0' <= s[i] && s[i] <= '9' {
		i++
	}

	// Check that the number has no leading zeros.
	if s[0] == '0' && i != 1 {
		return 0, fmt.Errorf("the number has a leading zero: %w", ErrInvalidVersion)
	}

	u, err := strconv.ParseUint(s[:i], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert the string %s to integer: %w", s[:i], err)
	}

	return u, nil
}
