// Package semver is a parser for version strings that adhere to Semantic
// Versioning 2.0.0. The primary functions to use are [Parse] and [MustParse]
// which parse the given version strings into [Version]s. To check if a string
// is a valid version, you can use the [IsValid] function.
package semver

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

// BuildIdentifiers is a list of build identifiers in the Version.
type BuildIdentifiers []string

// A Version is a parsed instance of a version number that adheres to the
// semantic versioning 2.0.0.
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	Prerelease Prerelease
	Build      BuildIdentifiers
}

// Compare returns
//
//	-1 if v is less than u,
//	 0 if v equals u,
//	+1 if v is greater than u.
func (v *Version) Compare(u *Version) int {
	var d int

	if d = cmp.Compare(v.Major, u.Major); d != 0 {
		return d
	}

	if d = cmp.Compare(v.Minor, u.Minor); d != 0 {
		return d
	}

	if d = cmp.Compare(v.Patch, u.Patch); d != 0 {
		return d
	}

	return v.Prerelease.Compare(u.Prerelease)
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
		v.Build.equal(o.Build)
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
