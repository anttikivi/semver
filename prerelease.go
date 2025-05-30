package semver

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var errInvalidPrereleaseIndent = errors.New("invalid pre-release identifier")

// A Prerelease holds the pre-release identifiers of a version.
type Prerelease struct {
	identifiers []prereleaseIdentifier
}

// A prereleaseIdentifier is a single pre-release identifier separated by dots.
type prereleaseIdentifier interface {
	// Equal tells if the given prereleaseIdentifier is equal to this one.
	Equal(o prereleaseIdentifier) bool

	// Len returns the length of the pre-release identifier in characters.
	Len() int

	// String returns the string representation of the identifier.
	String() string

	// Value returns the Value for the identifier. If the identifier is a
	// numeric one, the Value is returned using the first return Value and the
	// second return Value is an empty string. If the identifier is an
	// alphanumeric identifier, the Value is returned using the second return
	// Value and the first return Value is -1.
	Value() (n int, s string)
}

type numericIdentifier struct {
	v int
}

type alphanumericIdentifier struct {
	v string
}

// NewPrerelease creates new [Prerelease] from the given elements. The elements
// must be strings or ints.
func NewPrerelease(a ...any) (Prerelease, error) {
	identifiers := make([]prereleaseIdentifier, 0)

	for _, v := range a {
		switch u := v.(type) {
		case int:
			identifiers = append(identifiers, numericIdentifier{u})
		case string:
			identifiers = append(identifiers, alphanumericIdentifier{u})
		default:
			return Prerelease{}, fmt.Errorf("%w: %v", errInvalidPrereleaseIndent, v)
		}
	}

	return Prerelease{identifiers}, nil
}

// Equal tells if the given Prerelease o is equal to p.
func (p Prerelease) Equal(o Prerelease) bool {
	return slices.EqualFunc(p.identifiers, o.identifiers, func(a, b prereleaseIdentifier) bool {
		return a.Equal(b)
	})
}

// String returns the string representation of the Prerelease p.
func (p Prerelease) String() string {
	var sb strings.Builder

	if len(p.identifiers) > 0 {
		for _, ident := range p.identifiers {
			i, s := ident.Value()

			switch {
			case i >= 0 && s == "":
				sb.WriteString(strconv.Itoa(i))
			case i == -1 && s != "":
				sb.WriteString(s)
			default:
				// TODO: Try not to panic, but we should never get here.
				panic(fmt.Sprintf("invalid pre-release identifier options: %d and %s", i, s))
			}

			sb.WriteRune('.')
		}
	} else {
		return ""
	}

	s := sb.String()

	return s[:len(s)-1]
}

// Equal tells if the given prereleaseIdentifier is equal to this one.
func (i numericIdentifier) Equal(o prereleaseIdentifier) bool {
	other, ok := o.(numericIdentifier)
	if !ok {
		return false
	}

	v1, _ := i.Value()
	v2, _ := other.Value()

	return v1 == v2
}

// Len returns the length of the pre-release identifier in characters.
func (i numericIdentifier) Len() int {
	return countDigits(i.v)
}

// String returns the string representation of the identifier.
func (i numericIdentifier) String() string {
	return strconv.Itoa(i.v)
}

// Value returns the Value for the identifier. If the identifier is a numeric
// one, the Value is returned using the first return Value and the second return
// Value is an empty string. If the identifier is an alphanumeric identifier,
// the Value is returned using the second return Value and the first return
// Value is -1.
func (i numericIdentifier) Value() (int, string) {
	return i.v, ""
}

// Equal tells if the given prereleaseIdentifier is equal to this one.
func (i alphanumericIdentifier) Equal(o prereleaseIdentifier) bool {
	other, ok := o.(alphanumericIdentifier)
	if !ok {
		return false
	}

	_, v1 := i.Value()
	_, v2 := other.Value()

	return v1 == v2
}

// Len returns the length of the pre-release identifier in characters.
func (i alphanumericIdentifier) Len() int {
	return len(i.v)
}

// String returns the string representation of the identifier.
func (i alphanumericIdentifier) String() string {
	return i.v
}

// Value returns the Value for the identifier. If the identifier is a numeric
// one, the Value is returned using the first return Value and the second return
// Value is an empty string. If the identifier is an alphanumeric identifier,
// the Value is returned using the second return Value and the first return
// Value is -1.
func (i alphanumericIdentifier) Value() (int, string) {
	return -1, i.v
}
