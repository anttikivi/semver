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

// Parse parses the given string into a Version. The version string may have
// a 'v' prefix.
func Parse(ver string) (*Version, error) {
	if ver == "" {
		return nil, fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	pos, err := parsePrefix(ver)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the version prefix: %w", err)
	}

	major, err := parseNext(ver[pos:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse the major version: %w", err)
	}

	pos += countDigits(major)
	if pos >= len(ver) || ver[pos] != '.' {
		return nil, fmt.Errorf("no dot after the major version: %w", ErrInvalidVersion)
	}

	pos++

	minor, err := parseNext(ver[pos:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse the minor version: %w", err)
	}

	pos += countDigits(minor)
	if pos >= len(ver) || ver[pos] != '.' {
		return nil, fmt.Errorf("no dot after the minor version: %w", ErrInvalidVersion)
	}

	pos++

	patch, err := parseNext(ver[pos:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse the patch version: %w", err)
	}

	pos += countDigits(patch)

	if pos < len(ver) && ver[pos] != '-' && ver[pos] != '+' {
		return nil, fmt.Errorf("%w: invalid char %q at %d", ErrInvalidVersion, ver[pos], pos)
	}

	var prereleaseIdentifiers []prereleaseIdentifier

	if pos < len(ver) && ver[pos] == '-' {
		// The hyphen is not passed to the parser.
		pos++

		prereleaseIdentifiers, err = parsePrereleaseIdentifiers(ver[pos:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse the pre-release identifiers: %w", err)
		}

		// Move the position by the number of dots in the pre-release.
		pos += len(prereleaseIdentifiers) - 1

		for _, v := range prereleaseIdentifiers {
			pos += v.Len()
		}
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
		Prerelease: Prerelease{identifiers: prereleaseIdentifiers},
		Build:      build,
	}, nil
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
func ParseLax(ver string) (*Version, error) { //nolint:cyclop,funlen // no problem here
	if ver == "" {
		return nil, fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	pos, err := parsePrefix(ver)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the version prefix: %w", err)
	}

	major, err := parseNext(ver[pos:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse the major version: %w", err)
	}

	pos += countDigits(major)
	if pos >= len(ver) {
		return &Version{
			Major:      major,
			Minor:      0,
			Patch:      0,
			Prerelease: Prerelease{identifiers: []prereleaseIdentifier{}},
			Build:      BuildIdentifiers{},
		}, nil
	}

	minor := uint64(0)

	// Parse the minor version only if the next character is a dot.
	if ver[pos] == '.' {
		pos++

		minor, err = parseNext(ver[pos:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse the minor version: %w", err)
		}

		pos += countDigits(minor)
		if pos >= len(ver) {
			return &Version{
				Major:      major,
				Minor:      minor,
				Patch:      0,
				Prerelease: Prerelease{identifiers: []prereleaseIdentifier{}},
				Build:      BuildIdentifiers{},
			}, nil
		}
	}

	patch := uint64(0)

	// Parse the patch version only if the next character is a dot.
	if ver[pos] == '.' {
		pos++

		patch, err = parseNext(ver[pos:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse the patch version: %w", err)
		}

		pos += countDigits(minor)
		if pos >= len(ver) {
			return &Version{
				Major:      major,
				Minor:      minor,
				Patch:      patch,
				Prerelease: Prerelease{identifiers: []prereleaseIdentifier{}},
				Build:      BuildIdentifiers{},
			}, nil
		}
	}

	if pos < len(ver) && ver[pos] != '-' && ver[pos] != '+' {
		return nil, fmt.Errorf("%w: invalid char %q at %d", ErrInvalidVersion, ver[pos], pos)
	}

	var prereleaseIdentifiers []prereleaseIdentifier

	if pos < len(ver) && ver[pos] == '-' {
		// The hyphen is not passed to the parser.
		pos++

		prereleaseIdentifiers, err = parsePrereleaseIdentifiers(ver[pos:])
		if err != nil {
			return nil, fmt.Errorf("failed to parse the pre-release identifiers: %w", err)
		}

		// Move the position by the number of dots in the pre-release.
		pos += len(prereleaseIdentifiers) - 1

		for _, v := range prereleaseIdentifiers {
			pos += v.Len()
		}
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
		Prerelease: Prerelease{identifiers: prereleaseIdentifiers},
		Build:      build,
	}, nil
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

func isAlphanumericIdentifier(b byte) bool {
	return ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') || ('0' <= b && b <= '9') || b == '-'
}

func isPrereleaseSeparator(b byte) bool {
	return b == '.' || b == '+'
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
			func(r rune) bool { return !isAlphanumericIdentifier(byte(r)) },
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

// parsePrefix parses the possible "v" prefix for the version string.
// The function returns the new position where the parsing continues.
func parsePrefix(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	pos := 0

	b := s[0]
	if (b < '0' || '9' < b) && b != 'v' {
		return pos, fmt.Errorf(
			"%w: version %q does not start with a digit or 'v'",
			ErrInvalidVersion,
			s,
		)
	}

	if b == 'v' {
		pos++
	}

	if pos == len(s) {
		return pos, fmt.Errorf("%w: %q", ErrInvalidVersion, s)
	}

	return pos, nil
}

func parsePrereleaseIdentifiers(s string) ([]prereleaseIdentifier, error) {
	if s == "" {
		return nil, fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	var builder strings.Builder

	resultLen := strings.Count(s, ".") + 1
	if i := strings.IndexRune(s, '+'); i != -1 {
		resultLen -= strings.Count(s[i:], ".")
	}

	result := make([]prereleaseIdentifier, resultLen)

	i := 0

	for j := range len(s) {
		char := s[j]
		if !isPrereleaseSeparator(char) {
			builder.WriteByte(char)
		}

		if isPrereleaseSeparator(char) || j == len(s)-1 {
			current := builder.String()

			isAlphanum := strings.ContainsFunc(
				current,
				func(r rune) bool { return !unicode.IsDigit(r) },
			)

			switch {
			case isAlphanum:
				result[i] = alphanumericIdentifier{current}
			case current == "0":
				result[i] = numericIdentifier{0}
			case current[0] != '0':
				num, err := strconv.ParseUint(current, 10, 64)
				if err != nil {
					return nil, fmt.Errorf(
						"failed to convert pre-release identifier to integer: %w",
						err,
					)
				}

				result[i] = numericIdentifier{num}
			default:
				return nil, fmt.Errorf(
					"invalid pre-release identifier %q: %w",
					current,
					ErrInvalidVersion,
				)
			}

			i++

			if char == '+' {
				break
			}

			builder.Reset()
		}

		if !isAlphanumericIdentifier(char) && char != '.' {
			return nil, fmt.Errorf("invalid pre-release identifier %q: %w", char, ErrInvalidVersion)
		}
	}

	return result, nil
}
