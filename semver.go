// Package semver is a parser for version strings that adhere to Semantic
// Versioning 2.0.0. The primary functions to use are [Parse] and [MustParse]
// which parse the given version strings into [Version]s. To check if a string
// is a valid version, you can use the [IsValid] function.
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

	// ErrUnknown is returned when there is a problem with the parsing that is not
	// directly related to the caller giving an invalid string.
	ErrUnknown = errors.New("parsing failed")

	errInvalidPrereleaseIdent = errors.New("invalid pre-release identifier")
)

// A Version is a parsed instance of a version number that adheres to the
// semantic versioning 2.0.0.
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease Prerelease
	Build      BuildIdentifiers
}

// A Prerelease holds the pre-release identifiers of a version.
type Prerelease []PrereleaseIdentifier

// A PrereleaseIdentifier is a single pre-release identifier separated by dots.
type PrereleaseIdentifier interface {
	// Compare returns
	//
	//	-1 if this identifier is less than o,
	//	 0 if this identifier equals o,
	//	+1 if this identifier is greater than o.
	//
	// The comparison is done according to the semantic versioning specification
	// for pre-release identifiers.
	Compare(o PrereleaseIdentifier) int

	// Equal tells if the given prereleaseIdentifier is equal to this one.
	Equal(o PrereleaseIdentifier) bool

	// IsAlphanumericIdentifier reports whether this prereleaseIdentifier is
	// alphanumeric.
	IsAlphanumericIdentifier() bool

	// IsNumericIdentifier reports whether this prereleaseIdentifier is numeric.
	IsNumericIdentifier() bool

	// Len returns the length of the pre-release identifier in characters.
	Len() int

	// String returns the string representation of the identifier.
	String() string

	// Value returns the Value for the identifier.
	Value() any
}

// BuildIdentifiers is a list of build identifiers in the Version.
type BuildIdentifiers []string

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

// NewPrerelease creates new [Prerelease] from the given elements. The elements
// must be strings or ints.
func NewPrerelease(a ...any) (Prerelease, error) {
	identifiers := make(Prerelease, 0, len(a))

	for _, v := range a {
		switch u := v.(type) {
		case int:
			if u < 0 {
				return nil, fmt.Errorf("%w: %v", errInvalidPrereleaseIdent, v)
			}

			identifiers = append(identifiers, numericIdentifier{uint64(u)})
		case uint64:
			identifiers = append(identifiers, numericIdentifier{u})
		case string:
			p, err := parsePrereleaseIdentifier(u)
			if err != nil {
				return nil, fmt.Errorf("cannot create Prerelease: %w", err)
			}

			identifiers = append(identifiers, p)
		default:
			return nil, fmt.Errorf("%w: %v", errInvalidPrereleaseIdent, v)
		}
	}

	return identifiers, nil
}

// ParsePrerelease parses the given string into a Prerelease, separating
// the identifiers at dots.
func ParsePrerelease(s string) (Prerelease, error) {
	parts := strings.Split(s, ".")
	prerelease := make(Prerelease, 0, len(parts))

	for _, v := range parts {
		p, err := parsePrereleaseIdentifier(v)
		if err != nil {
			return nil, fmt.Errorf("parsing prerelease %q failed: %w", s, err)
		}

		prerelease = append(prerelease, p)
	}

	return prerelease, nil
}

// NewBuildIdentifiers returns new [BuildIdentifiers] for the given strings.
func NewBuildIdentifiers(s ...string) BuildIdentifiers {
	b := make(BuildIdentifiers, 0, len(s))
	b = append(b, s...)

	return b
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

	return v.Prerelease.Compare(w.Prerelease)
}

// Core returns the comparable string representation of the version. It doesn't
// include the build metadata.
func (v *Version) Core() string {
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
		v.Build.Equal(o.Build)
}

// String returns the string representation of the version.
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

// Compare returns
//
//	-1 if p is less than o,
//	 0 if p equals o,
//	+1 if p is greater than o.
//
// The comparison is done according to the semantic versioning specification for
// pre-release identifiers.
func (p Prerelease) Compare(o Prerelease) int {
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

// Equal tells if the given Prerelease o is equal to p.
func (p Prerelease) Equal(o Prerelease) bool {
	return slices.EqualFunc(p, o, func(a, b PrereleaseIdentifier) bool {
		return a.Equal(b)
	})
}

// String returns the string representation of the Prerelease p.
func (p Prerelease) String() string {
	var sb strings.Builder

	if len(p) > 0 {
		for _, ident := range p {
			val := ident.Value()

			switch v := val.(type) {
			case uint64:
				sb.WriteString(strconv.FormatUint(v, 10))
			case string:
				sb.WriteString(v)
			default:
				// TODO: Try not to panic, but we should never get here.
				panic(fmt.Sprintf("invalid pre-release identifier option: %[1]v (%[1]T)", val))
			}

			sb.WriteRune('.')
		}
	} else {
		return ""
	}

	s := sb.String()

	return s[:len(s)-1]
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

// Equal tells if the given BuildIdentifiers b are equal to o.
func (b BuildIdentifiers) Equal(o BuildIdentifiers) bool {
	return slices.Equal(b, o)
}

// Compare returns
//
//	-1 if this identifier is less than o,
//	 0 if this identifier equals o,
//	+1 if this identifier is greater than o.
//
// The comparison is done according to the semantic versioning specification for
// pre-release identifiers.
func (i alphanumericIdentifier) Compare(o PrereleaseIdentifier) int {
	// Alphanumeric identifiers always have higher precedence than numeric ones.
	if o.IsNumericIdentifier() {
		return 1
	}

	// Now both of the identifiers must be alphanumeric.
	j, ok := o.(alphanumericIdentifier)
	if !ok {
		panic(fmt.Sprintf("compared identifier should be alphanumeric: %v", o))
	}

	return cmp.Compare(i.v, j.v)
}

// Equal tells if the given prereleaseIdentifier is equal to this one.
func (i alphanumericIdentifier) Equal(o PrereleaseIdentifier) bool {
	other, ok := o.(alphanumericIdentifier)
	if !ok {
		return false
	}

	a, ok := i.Value().(string)
	if !ok {
		panic(fmt.Sprintf("failed to convert %[1]v (%[1]T) to string", i.Value()))
	}

	b, ok := other.Value().(string)
	if !ok {
		panic(fmt.Sprintf("failed to convert %[1]v (%[1]T) to string", other.Value()))
	}

	return a == b
}

// IsAlphanumericIdentifier reports whether this prereleaseIdentifier is
// alphanumeric.
func (i alphanumericIdentifier) IsAlphanumericIdentifier() bool {
	return true
}

// IsNumericIdentifier reports whether this prereleaseIdentifier is numeric.
func (i alphanumericIdentifier) IsNumericIdentifier() bool {
	return false
}

// Len returns the length of the pre-release identifier in characters.
func (i alphanumericIdentifier) Len() int {
	return len(i.v)
}

// String returns the string representation of the identifier.
func (i alphanumericIdentifier) String() string {
	return i.v
}

// Value returns the Value for the identifier.
func (i alphanumericIdentifier) Value() any {
	return i.v
}

// Compare returns
//
//	-1 if this identifier is less than o,
//	 0 if this identifier equals o,
//	+1 if this identifier is greater than o.
//
// The comparison is done according to the semantic versioning specification for
// pre-release identifiers.
func (i numericIdentifier) Compare(o PrereleaseIdentifier) int {
	// Alphanumeric identifiers always have higher precedence than numeric ones.
	if o.IsAlphanumericIdentifier() {
		return -1
	}

	// Now both of the identifiers must be numeric.
	j, ok := o.(numericIdentifier)
	if !ok {
		panic(fmt.Sprintf("compared identifier should be numeric: %v", o))
	}

	return cmp.Compare(i.v, j.v)
}

// Equal tells if the given prereleaseIdentifier is equal to this one.
func (i numericIdentifier) Equal(o PrereleaseIdentifier) bool {
	other, ok := o.(numericIdentifier)
	if !ok {
		return false
	}

	a, ok := i.Value().(uint64)
	if !ok {
		panic(fmt.Sprintf("failed to convert %[1]v (%[1]T) to uint64", i.Value()))
	}

	b, ok := other.Value().(uint64)
	if !ok {
		panic(fmt.Sprintf("failed to convert %[1]v (%[1]T) to uint64", other.Value()))
	}

	return a == b
}

// IsAlphanumericIdentifier reports whether this prereleaseIdentifier is
// alphanumeric.
func (i numericIdentifier) IsAlphanumericIdentifier() bool {
	return false
}

// IsNumericIdentifier reports whether this prereleaseIdentifier is numeric.
func (i numericIdentifier) IsNumericIdentifier() bool {
	return true
}

// Len returns the length of the pre-release identifier in characters.
func (i numericIdentifier) Len() int {
	return countDigits(i.v)
}

// String returns the string representation of the identifier.
func (i numericIdentifier) String() string {
	return strconv.FormatUint(i.v, 10)
}

// Value returns the Value for the identifier.
func (i numericIdentifier) Value() any {
	return i.v
}

// Compare returns
//
//	-1 if v is less than w,
//	 0 if v equals w,
//	+1 if v is greater than w.
//
// The comparison is done according to the semantic versioning specification.
func Compare(v *Version, w *Version) int {
	return v.Compare(w)
}

//nolint:cyclop,funlen,gocognit // TODO: see if worth fixing
func parse(s string, minCore int) (*Version, error) {
	if s == "" {
		return nil, fmt.Errorf("empty string: %w", ErrInvalidVersion)
	}

	if !isASCII(s) {
		return nil, fmt.Errorf("%w: version contains non-ASCII characters", ErrInvalidVersion)
	}

	pos, err := checkPrefix(s)
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
				"%w: index when checking version number is out of bounds: %d",
				ErrUnknown,
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

		prerelease, err = ParsePrerelease(s[pos:i])
		if err != nil {
			return nil, fmt.Errorf("failed to parse the pre-release: %w", err)
		}

		pos = i
	}

	var build BuildIdentifiers

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

//nolint:ireturn // interface return is needed
func parsePrereleaseIdentifier(s string) (PrereleaseIdentifier, error) {
	if s == "" {
		return nil, fmt.Errorf("%w: identifier is an empty string", errInvalidPrereleaseIdent)
	}

	if !isASCII(s) {
		return nil, fmt.Errorf(
			"%w: identifier %q contains non-ASCII characters",
			errInvalidPrereleaseIdent,
			s,
		)
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
				errInvalidPrereleaseIdent,
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
		return nil, fmt.Errorf("%w: %s", errInvalidPrereleaseIdent, s)
	}
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

	return x.Compare(y)
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
