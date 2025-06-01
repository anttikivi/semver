/*
Package semver provides utilities and a parser to work with version numbers that
adhere to [semantic versioning]. The goal of this parser is to be reliable and
performant. Reliability is ensured by using a wide range of tests and fuzzing.
Performance is achieved by implementing a custom parser instead of the common
alternative: regular expressions.

This package implements [semantic versioning 2.0.0]. Specifically, the current
capabilities of this package include:

  - Parsing version strings.
  - Checking if a string is valid version string. This check doesn’t require
    full parsing of the version.
  - Comparing versions.
  - Sorting versions.

The version strings can optionally have a "v" prefix.

# Parsing versions

The package includes two types of functions for parsing versions. There are
the [Parse] and [ParseLax] functions. [Parse] parses only full valid version
strings like "1.2.3", "1.2.3-beta.1", or "1.2.3-beta.1+darwin.amd64". [ParseLax]
works otherwise like [Parse] but it tries to coerse incomplete core version into
a full version. For example, it parses "v1" as "1.0.0" and "1.2-beta" as
"1.2.0-beta". Both functions return a pointer to the [Version] object and
an error.

They can be used as follows:

	v, err := semver.Parse("1.2.3-beta.1")

The package also offers [MustParse] and [MustParseLax] variants of these
functions. They are otherwise the same but only return the pointer to [Version].
They panic on errors.

# Validating version strings

The package includes two functions, similar to the parsing functions, for
checking if a string is a valid version string. The functions are [IsValid] and
[IsValidLax] and they return a single boolean value. The return value is
analogous to whether the matching parsing function would parse the given string.

Example usage:

	ok := semver.IsValid("1.2.3-beta.1")

# Sorting versions

The package contains the [Versions] type that supports sorting using the Go
standard library [sort] package. [Versions] is defined as []*Version.

Example usage:

	a := []string{"1.2.3", "1.0", "1.3", "2", "0.4.2"}
	slice := make(Versions, len(a))

	for i, s := range {
		slice[i] = semver.MustParseLax(s)
	}

	sort.Sort(slice)

	for _, v := range slice {
		fmt.Println(v.String())
	}

The above code would print:

	0.4.2
	1.0.0
	1.2.3
	1.3.0
	2.0.0

[semantic versioning]: https://semver.org
[semantic versioning 2.0.0]: https://semver.org/spec/v2.0.0.html
*/
package semver

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// Common errors returned by the version parser.
var (
	// ErrInvalidVersion is the error returned by the version parsing functions when
	// they encounter invalid version string.
	ErrInvalidVersion = errors.New("invalid semantic version")

	// ErrParser is returned when there is a problem with the parsing that is not
	// directly related to the caller giving an invalid string.
	ErrParser = errors.New("parsing failed")
)

// A Version is a parsed instance of a version number that adheres to the
// semantic versioning 2.0.0.
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease Prerelease
	Build      Build
}

// A Prerelease holds the pre-release identifiers of a version.
type Prerelease []PrereleaseIdentifier

// A PrereleaseIdentifier is a single pre-release identifier separated by dots.
type PrereleaseIdentifier interface {
	// String returns the string representation of the identifier.
	String() string

	// compare returns
	//
	//	-1 if this identifier is less than o,
	//	 0 if this identifier equals o,
	//	+1 if this identifier is greater than o.
	//
	// The comparison is done according to the semantic versioning specification
	// for pre-release identifiers.
	compare(o PrereleaseIdentifier) int

	// equal tells if the given PrereleaseIdentifier is equal to this one.
	equal(o PrereleaseIdentifier) bool

	// isAlphanumeric reports whether this PrereleaseIdentifier is alphanumeric.
	isAlphanumeric() bool

	// isNumeric reports whether this PrereleaseIdentifier is numeric.
	isNumeric() bool

	// len returns the length of the pre-release identifier in characters.
	len() int
}

// Build is a list of build identifiers in the Version.
type Build []string

type alphanumericIdentifier struct {
	v string
}

type numericIdentifier struct {
	v uint64
}

// MustParse parses the given string into a Version and panics if it encounters
// an error. The version string may have a 'v' prefix.
func MustParse(s string) *Version {
	v, err := Parse(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse the string %q into a version: %v", s, err))
	}

	return v
}

// MustParseLax parses the given string into a Version and panics if it
// encounters an error. The version string number may be partial, i.e. it parses
// 'v1' into '1.0.0' and 'v1.2' into '1.2.0'. The version may have a 'v' prefix.
func MustParseLax(s string) *Version {
	v, err := ParseLax(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse the string %q into a version: %v", s, err))
	}

	return v
}

// Parse parses the given string into a Version. The version string may have
// a 'v' prefix.
func Parse(s string) (*Version, error) {
	v, err := parse(s, 3) //nolint:mnd // <major>.<minor>.<patch>
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	return v, nil
}

// ParseLax parses the given string into a Version. The version number may be
// partial, i.e. it parses 'v1' into '1.0.0' and 'v1.2' into '1.2.0'.
// The version string may have a 'v' prefix.
func ParseLax(s string) (*Version, error) {
	v, err := parse(s, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	return v, nil
}

// Compare returns
//
//	-1 if v is less than w,
//	 0 if v equals w,
//	+1 if v is greater than w.
//
// The comparison is done according to the semantic versioning specification.
func (v *Version) Compare(w *Version) int {
	var d int

	if d = cmp.Compare(v.Major, w.Major); d != 0 {
		return d
	}

	if d = cmp.Compare(v.Minor, w.Minor); d != 0 {
		return d
	}

	if d = cmp.Compare(v.Patch, w.Patch); d != 0 {
		return d
	}

	if v.Prerelease == nil && w.Prerelease != nil {
		return 1
	}

	if v.Prerelease != nil && w.Prerelease == nil {
		return -1
	}

	return v.Prerelease.compare(w.Prerelease)
}

// ComparableString returns the comparable string representation of the version.
// It doesn't include the build metadata.
func (v *Version) ComparableString() string {
	var sb strings.Builder

	sb.WriteString(strconv.FormatUint(v.Major, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Minor, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Patch, 10))

	if len(v.Prerelease) > 0 {
		sb.WriteByte('-')
		sb.WriteString(v.Prerelease.String())
	}

	return sb.String()
}

// CoreString returns the core version string representation of the version. It
// doesn't include the pre-release nor the build metadata.
func (v *Version) CoreString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Equal reports whether Version w is equal to v. The two Versions are equal
// according to this function if all of their parts that are comparable in
// the semantic versioning specification are equal; this does not include
// the build metadata.
func (v *Version) Equal(w *Version) bool {
	if w == nil {
		return v == nil
	}

	return v.Major == w.Major && v.Minor == w.Minor && v.Patch == w.Patch &&
		v.Prerelease.equal(w.Prerelease)
}

// StrictEqual reports whether Version w is equal to v. The two Versions are
// equal if all of their parts are; this includes the build metadata.
func (v *Version) StrictEqual(w *Version) bool {
	if w == nil {
		return v == nil
	}

	return v.Major == w.Major && v.Minor == w.Minor && v.Patch == w.Patch &&
		v.Prerelease.equal(w.Prerelease) &&
		v.Build.equal(w.Build)
}

// String returns the string representation of v.
func (v *Version) String() string {
	var sb strings.Builder

	sb.WriteString(strconv.FormatUint(v.Major, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Minor, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Patch, 10))

	if len(v.Prerelease) > 0 {
		sb.WriteByte('-')
		sb.WriteString(v.Prerelease.String())
	}

	if len(v.Build) > 0 {
		sb.WriteByte('+')
		sb.WriteString(v.Build.String())
	}

	return sb.String()
}

// String returns the string representation of p.
func (p Prerelease) String() string {
	if len(p) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, ident := range p {
		if i > 0 {
			sb.WriteRune('.')
		}

		switch v := ident.(type) {
		case alphanumericIdentifier:
			sb.WriteString(v.v)
		case numericIdentifier:
			sb.WriteString(strconv.FormatUint(v.v, 10))
		default:
			// Internal invariant violation.
			panic(fmt.Sprintf("invalid pre-release identifier option: %[1]v (%[1]T)", v))
		}
	}

	return sb.String()
}

// String returns the string representation of b.
func (b Build) String() string {
	if len(b) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, s := range b {
		if i > 0 {
			sb.WriteRune('.')
		}

		sb.WriteString(s)
	}

	return sb.String()
}

// String returns the string representation of the identifier.
func (i alphanumericIdentifier) String() string {
	return i.v
}

// String returns the string representation of the identifier.
func (i numericIdentifier) String() string {
	return strconv.FormatUint(i.v, 10)
}

// Compare returns
//
//	-1 if v is less than w,
//	 0 if v equals w,
//	+1 if v is greater than w.
//
// The comparison is done according to the semantic versioning specification.
func Compare(v, w *Version) int {
	return v.Compare(w)
}

//nolint:cyclop,funlen,gocognit // TODO: see if worth fixing
func parse(s string, minCore int) (*Version, error) {
	if s == "" {
		return nil, fmt.Errorf("%w: empty string", ErrInvalidVersion)
	}

	if !isASCII(s) {
		return nil, fmt.Errorf("%w: version contains non-ASCII characters", ErrInvalidVersion)
	}

	pos, err := stripPrefix(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the version prefix: %w", err)
	}

	i := len(s)

	for j := range s[pos:] {
		c := s[pos+j]
		if !isDigit(c) && c != '.' {
			i = pos + j

			break
		}
	}

	nums := strings.Split(s[pos:i], ".")

	if len(nums) > 3 { //nolint:mnd // <major>.<minor>.<patch>
		return nil, fmt.Errorf("%w: too many core version numbers in %q", ErrInvalidVersion, s)
	}

	if len(nums) < minCore {
		return nil, fmt.Errorf("%w: not enough core version numbers in %q", ErrInvalidVersion, s)
	}

	major := uint64(0)
	minor := uint64(0)
	patch := uint64(0)

	for j, n := range nums {
		if n == "" {
			return nil, fmt.Errorf("%w: empty version number in %q", ErrInvalidVersion, s)
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
		case 2: //nolint:mnd // represents the index
			patch = u
		default:
			return nil, fmt.Errorf(
				"%w: index for checking version number is out of bounds: %d",
				ErrParser,
				j,
			)
		}
	}

	pos = i

	if pos < len(s) && s[pos] != '-' && s[pos] != '+' {
		return nil, fmt.Errorf("%w: invalid char %q at %d", ErrInvalidVersion, s[pos], pos)
	}

	var prerelease Prerelease

	if pos < len(s) && s[pos] == '-' {
		// The hyphen is not passed to the parser.
		pos++

		i = len(s)

		for j := range s[pos:] {
			c := s[pos+j]
			if c == '+' {
				i = pos + j

				break
			}
		}

		parts := strings.Split(s[pos:i], ".")
		prerelease = make(Prerelease, 0, len(parts))

		for _, v := range parts {
			p, err := parsePrereleaseIdentifier(v)
			if err != nil {
				return nil, fmt.Errorf("parsing prerelease %q failed: %w", s, err)
			}

			prerelease = append(prerelease, p)
		}

		pos = i
	}

	var build Build

	if pos < len(s) && s[pos] == '+' {
		// Move past the '+'.
		pos++

		build, err = parseBuild(s[pos:])
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

// newPrerelease creates new [Prerelease] from the given elements. The elements
// must be strings, ints, or uint64s.
func newPrerelease(a ...any) (Prerelease, error) {
	identifiers := make(Prerelease, 0, len(a))

	for _, v := range a {
		switch u := v.(type) {
		case int:
			if u < 0 {
				return nil, fmt.Errorf("%w: %v", ErrInvalidVersion, v)
			}

			identifiers = append(identifiers, numericIdentifier{uint64(u)})
		case uint64:
			identifiers = append(identifiers, numericIdentifier{u})
		case string:
			if !isASCII(u) {
				return nil, fmt.Errorf(
					"%w: identifier %q contains non-ASCII characters",
					ErrInvalidVersion,
					u,
				)
			}

			p, err := parsePrereleaseIdentifier(u)
			if err != nil {
				return nil, fmt.Errorf("cannot create Prerelease: %w", err)
			}

			identifiers = append(identifiers, p)
		default:
			return nil, fmt.Errorf("%w: %v", ErrInvalidVersion, v)
		}
	}

	return identifiers, nil
}

//nolint:ireturn // interface return is needed
func parsePrereleaseIdentifier(s string) (PrereleaseIdentifier, error) {
	if s == "" {
		return nil, fmt.Errorf("%w: identifier is an empty string", ErrInvalidVersion)
	}

	// Check the case for single zero early.
	if s == "0" {
		return numericIdentifier{0}, nil
	}

	switch {
	case isNumericIdentifier(s):
		// If this is a numeric identifier and the first character is zero, we
		// already know that the length is greater than 1 as the case for that
		// was checked at the start.
		if s[0] == '0' {
			return nil, fmt.Errorf(
				"%w: numeric identifier with a leading zero: %s",
				ErrInvalidVersion,
				s,
			)
		}

		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to convert pre-release identifier to integer: %w",
				err,
			)
		}

		return numericIdentifier{u}, nil
	case isAlphanumericIdentifier(s):
		return alphanumericIdentifier{s}, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidVersion, s)
	}
}

// newBuild returns new [Build] for the given strings.
func newBuild(s ...string) Build {
	b := make(Build, 0, len(s))
	b = append(b, s...)

	return b
}

// compare returns
//
//	-1 if p is less than o,
//	 0 if p equals o,
//	+1 if p is greater than o.
//
// The comparison is done according to the semantic versioning specification for
// pre-release identifiers.
func (p Prerelease) compare(o Prerelease) int {
	for i := range max(len(p), len(o)) {
		var (
			x PrereleaseIdentifier
			y PrereleaseIdentifier
		)

		if i < len(p) {
			x = p[i]
		}

		if i < len(o) {
			y = o[i]
		}

		if d := comparePrereleaseIdentifiers(x, y); d != 0 {
			return d
		}
	}

	return 0
}

// equal tells if p is equal to o.
func (p Prerelease) equal(o Prerelease) bool {
	return slices.EqualFunc(p, o, func(a, b PrereleaseIdentifier) bool {
		return a.equal(b)
	})
}

// equal tells if b is equal to a.
func (b Build) equal(a Build) bool {
	return slices.Equal(b, a)
}

// compare returns
//
//	-1 if this identifier is less than o,
//	 0 if this identifier equals o,
//	+1 if this identifier is greater than o.
//
// The comparison is done according to the semantic versioning specification for
// pre-release identifiers.
func (i alphanumericIdentifier) compare(o PrereleaseIdentifier) int {
	// Alphanumeric identifiers always have higher precedence than numeric ones.
	if o.isNumeric() {
		return 1
	}

	// Now both of the identifiers must be alphanumeric.
	j, ok := o.(alphanumericIdentifier)
	if !ok {
		panic(fmt.Sprintf("compared identifier should be alphanumeric: %v", o))
	}

	return cmp.Compare(i.v, j.v)
}

// equal tells if the given prereleaseIdentifier is equal to this one.
func (i alphanumericIdentifier) equal(o PrereleaseIdentifier) bool {
	other, ok := o.(alphanumericIdentifier)
	if !ok {
		return false
	}

	return i.v == other.v
}

// isAlphanumeric reports whether this PrereleaseIdentifier is alphanumeric.
func (i alphanumericIdentifier) isAlphanumeric() bool {
	return true
}

// isNumeric reports whether this PrereleaseIdentifier is numeric.
func (i alphanumericIdentifier) isNumeric() bool {
	return false
}

// len returns the length of the pre-release identifier in characters.
func (i alphanumericIdentifier) len() int {
	return len(i.v)
}

// compare returns
//
//	-1 if this identifier is less than o,
//	 0 if this identifier equals o,
//	+1 if this identifier is greater than o.
//
// The comparison is done according to the semantic versioning specification for
// pre-release identifiers.
func (i numericIdentifier) compare(o PrereleaseIdentifier) int {
	// Alphanumeric identifiers always have higher precedence than numeric ones.
	if o.isAlphanumeric() {
		return -1
	}

	// Now both of the identifiers must be numeric.
	j, ok := o.(numericIdentifier)
	if !ok {
		panic(fmt.Sprintf("compared identifier should be numeric: %v", o))
	}

	return cmp.Compare(i.v, j.v)
}

// equal tells if the given prereleaseIdentifier is equal to this one.
func (i numericIdentifier) equal(o PrereleaseIdentifier) bool {
	other, ok := o.(numericIdentifier)
	if !ok {
		return false
	}

	return i.v == other.v
}

// isAlphanumeric reports whether this PrereleaseIdentifier is alphanumeric.
func (i numericIdentifier) isAlphanumeric() bool {
	return false
}

// isNumeric reports whether this PrereleaseIdentifier is numeric.
func (i numericIdentifier) isNumeric() bool {
	return true
}

// len returns the length of the pre-release identifier in characters.
func (i numericIdentifier) len() int {
	return countDigits(i.v)
}

func comparePrereleaseIdentifiers(x, y PrereleaseIdentifier) int {
	if x == y {
		return 0
	}

	if x == nil {
		if y != nil {
			return -1
		}

		return 1
	}

	if y == nil {
		return 1
	}

	return x.compare(y)
}

func countDigits(u uint64) int {
	if u == 0 {
		return 1
	}

	count := 0

	for u != 0 {
		u /= 10
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

func parseBuild(s string) ([]string, error) {
	if s == "" {
		return nil, fmt.Errorf("%w: cannot parse empty string as a build", ErrInvalidVersion)
	}

	result := strings.Split(s, ".")
	for _, v := range result {
		if v == "" {
			return nil, fmt.Errorf("%w: empty string as a build identifier", ErrInvalidVersion)
		}

		// This should be safe as all of the characters in the version must be
		// ASCII.
		if !isAlphanumericIdentifier(v) {
			return nil, fmt.Errorf(
				"%w: invalid byte in the build identifier %q",
				ErrInvalidVersion,
				v,
			)
		}
	}

	return result, nil
}

// stripPrefix parses the possible "v" prefix for the version string.
// The function returns the new position where the parsing continues.
func stripPrefix(s string) (int, error) {
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
