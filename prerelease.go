package semver

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var errInvalidPrereleaseIndent = errors.New("invalid pre-release identifier")

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

type numericIdentifier struct {
	v uint64
}

type alphanumericIdentifier struct {
	v string
}

// NewPrerelease creates new [Prerelease] from the given elements. The elements
// must be strings or ints.
func NewPrerelease(a ...any) (Prerelease, error) {
	identifiers := make(Prerelease, 0)

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

	return Prerelease(identifiers), nil
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

func comparePrereleaseIdentifiers(x PrereleaseIdentifier, y PrereleaseIdentifier) int {
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
