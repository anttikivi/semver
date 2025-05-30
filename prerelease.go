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

	// Value returns the Value for the identifier.
	Value() any
}

type numericIdentifier struct {
	v uint64
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
			if u < 0 {
				return Prerelease{}, fmt.Errorf("%w: %v", errInvalidPrereleaseIndent, v)
			}

			identifiers = append(identifiers, numericIdentifier{uint64(u)})
		case uint64:
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

// Equal tells if the given prereleaseIdentifier is equal to this one.
func (i numericIdentifier) Equal(o prereleaseIdentifier) bool {
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

// Equal tells if the given prereleaseIdentifier is equal to this one.
func (i alphanumericIdentifier) Equal(o prereleaseIdentifier) bool {
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
