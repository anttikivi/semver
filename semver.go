// Package semver is a parser for version strings that adhere to Semantic
// Versioning 2.0.0. The primary functions to use are [Parse] and [MustParse]
// which parse the given version strings into [Version]s. To check if a string
// is a valid version, you can use the [IsValid] function.
package semver

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidVersion is the error returned by the version parsing functions when
// they encounter invalid version string.
var ErrInvalidVersion = errors.New("invalid semantic version")

// BuildIdentifiers is a list of build identifiers in the Version.
type BuildIdentifiers []string

// A Version is a parsed instance of a version number that adheres to the
// semantic versioning 2.0.0.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease Prerelease
	Build      BuildIdentifiers
	rawStr     string
}

// Parse parses the given string into a Version. The version may have a 'v'
// prefix.
func Parse(ver string) (*Version, error) {
	v, err := parse(ver)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return v, nil
}

// ParsePrefix parses the given string into a Version. It allows the version to
// have either one of the given prefixes or a 'v' prefix.
func ParsePrefix(ver string, prefixes ...string) (*Version, error) {
	v, err := parse(ver, prefixes...)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return v, nil
}

// MustParse parses the given string into a Version and panics if it encounters
// an error. The version may have a 'v' prefix.
func MustParse(ver string) *Version {
	v, err := parse(ver)
	if err != nil {
		panic(fmt.Sprintf("failed to parse the string %q into a version: %v", ver, err))
	}

	return v
}

// MustParsePrefix parses the given string into a Version and panics if it
// encounters an error. It allows the version to have either one of the given
// prefixes or a 'v' prefix.
func MustParsePrefix(ver string, prefixes ...string) *Version {
	v, err := parse(ver, prefixes...)
	if err != nil {
		panic(fmt.Sprintf("failed to parse the string %q into a version: %v", ver, err))
	}

	return v
}

// Equal reports whether Version o is equal to v. The two Versions are equal
// according to this function if all of their parts that are comparable in
// the semantic versioning specification are equal; this does not include
// the build metadata.
func (v *Version) Equal(o *Version) bool {
	if o == nil {
		return v == nil
	}

	return v.Major == o.Major && v.Minor == o.Minor && v.Patch == o.Patch &&
		v.Prerelease.Equal(o.Prerelease)
}

// StrictEqual reports whether Version o is equal to v. The two Versions are
// equal if all of their parts are; this includes the build metadata.
func (v *Version) StrictEqual(o *Version) bool {
	if o == nil {
		return v == nil
	}

	return v.Major == o.Major && v.Minor == o.Minor && v.Patch == o.Patch &&
		v.Prerelease.Equal(o.Prerelease) &&
		v.Build.equal(o.Build)
}

// String returns the normal string representation of the version. It doesn't
// include the build metadata.
func (v *Version) String() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(v.Major))
	sb.WriteByte('.')
	sb.WriteString(strconv.Itoa(v.Minor))
	sb.WriteByte('.')
	sb.WriteString(strconv.Itoa(v.Patch))

	if len(v.Prerelease.identifiers) > 0 {
		sb.WriteByte('-')
		sb.WriteString(v.Prerelease.String())
	}

	return sb.String()
}

// FullString returns the full string representation of the version including
// the build metadata.
func (v *Version) FullString() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(v.Major))
	sb.WriteByte('.')
	sb.WriteString(strconv.Itoa(v.Minor))
	sb.WriteByte('.')
	sb.WriteString(strconv.Itoa(v.Patch))

	if len(v.Prerelease.identifiers) > 0 {
		sb.WriteByte('-')
		sb.WriteString(v.Prerelease.String())
	}

	if len(v.Build) > 0 {
		sb.WriteByte('+')
		sb.WriteString(v.Build.String())
	}

	return sb.String()
}

// NewBuildIdentifiers returns new [BuildIdentifiers] for the given strings.
func NewBuildIdentifiers(s ...string) BuildIdentifiers {
	b := make(BuildIdentifiers, 0, len(s))
	b = append(b, s...)

	return b
}

// String returns the string representation of the BuildIdentifiers b.
func (b BuildIdentifiers) String() string {
	var sb strings.Builder

	if len(b) > 0 {
		for _, s := range b {
			sb.WriteString(s)
			sb.WriteRune('.')
		}
	} else {
		return ""
	}

	s := sb.String()

	return s[:len(s)-1]
}

func (b BuildIdentifiers) equal(o BuildIdentifiers) bool {
	return slices.Equal(b, o)
}

func parse(ver string, prefixes ...string) (*Version, error) {
	if ver == "" {
		return nil, fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	pos := 0

	prefix, err := parsePrefix(ver, prefixes...)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the version prefix: %w", err)
	}

	pos += len(prefix)

	major, err := parseNextInt(ver[pos:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse the major version: %w", err)
	}

	pos += countDigits(major)
	if pos >= len(ver) || ver[pos] != '.' {
		return nil, fmt.Errorf("no dot after the major version: %w", ErrInvalidVersion)
	}

	pos++

	minor, err := parseNextInt(ver[pos:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse the minor version: %w", err)
	}

	pos += countDigits(minor)
	if pos >= len(ver) || ver[pos] != '.' {
		return nil, fmt.Errorf("no dot after the minor version: %w", ErrInvalidVersion)
	}

	pos++

	patch, err := parseNextInt(ver[pos:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse the patch version: %w", err)
	}

	var prereleaseIdentifiers []prereleaseIdentifier

	pos += countDigits(patch)

	if pos < len(ver) && ver[pos] != '-' && ver[pos] != '+' {
		return nil, fmt.Errorf("%w: invalid char %q at %d", ErrInvalidVersion, ver[pos], pos)
	}

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
		rawStr:     ver,
	}, nil
}

func countDigits(i int) int {
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

func isAlphanumericIdentifier(c rune) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || unicode.IsDigit(c) || c == '-'
}

func isPrereleaseSeparator(c rune) bool {
	return c == '.' || c == '+'
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

		if strings.ContainsFunc(v, func(r rune) bool { return !isAlphanumericIdentifier(r) }) {
			return nil, fmt.Errorf(
				"invalid rune in the build identifier %q: %w",
				v,
				ErrInvalidVersion,
			)
		}
	}

	return result, nil
}

// parseNextInt parses the next integer from the given string. The string should
// be a version string or the next part to parse from a version string adhering
// to the semantic versioning. The first return value is the parsed interger, or
// -1 if the parsing fails. The second return value is an error or nil.
func parseNextInt(s string) (int, error) {
	if s == "" {
		return -1, fmt.Errorf("cannot parse empty string as int: %w", ErrInvalidVersion)
	}

	if !unicode.IsDigit(rune(s[0])) {
		return -1, fmt.Errorf("first character is not a digit: %w", ErrInvalidVersion)
	}

	i := 1
	for i < len(s) && unicode.IsDigit(rune(s[i])) {
		i++
	}

	// Check that the number has no leading zeros.
	if s[0] == '0' && i != 1 {
		return -1, fmt.Errorf("the number has a leading zero: %w", ErrInvalidVersion)
	}

	n, err := strconv.Atoi(s[:i])
	if err != nil {
		return -1, fmt.Errorf("failed to convert the string %s to integer: %w", s[:i], err)
	}

	return n, nil
}

// parsePrefix parses the possible prefixes for the version string. The program
// allows using either a custom prefix or 'v' as a prefix in the version string.
func parsePrefix(s string, p ...string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	i := strings.IndexFunc(s, unicode.IsDigit)
	if i == -1 {
		return "", fmt.Errorf(
			"version string %q has no digits after the possible prefix: %w",
			s,
			ErrInvalidVersion,
		)
	}

	if i == 0 {
		return "", nil
	}

	prefix := s[:i]
	if !slices.Contains(p, prefix) && prefix != "v" {
		return "", fmt.Errorf("invalid prefix %q: %w", prefix, ErrInvalidVersion)
	}

	return prefix, nil
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
		if !isPrereleaseSeparator(rune(char)) {
			builder.WriteByte(char)
		}

		if isPrereleaseSeparator(rune(char)) || j == len(s)-1 {
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
				num, err := strconv.Atoi(current)
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

		if !isAlphanumericIdentifier(rune(char)) && char != '.' {
			return nil, fmt.Errorf("invalid pre-release identifier %q: %w", char, ErrInvalidVersion)
		}
	}

	return result, nil
}
