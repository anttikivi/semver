package semver_test

import (
	"regexp"
	"testing"

	"github.com/anttikivi/semver"
)

const rawVersionRegex = `^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

var versionRegex *regexp.Regexp

var (
	isValidRegexTests []validationTestCase
	isValidTests      []validationTestCase
	laxIsValidTests   []validationTestCase
)

type validationTestCase struct {
	v    string
	want bool
}

func init() {
	versionRegex = regexp.MustCompile(rawVersionRegex)

	for prefix, allowed := range testPrefixes {
		for _, t := range baseTests {
			input := prefix + t.v
			want := t.wantStrict != nil

			if !allowed {
				want = false
			}

			isValidRegexTests = append(isValidRegexTests, validationTestCase{
				v:    input,
				want: want,
			})

			isValidTests = append(isValidTests, validationTestCase{
				v:    input,
				want: want,
			})

			if allowed {
				want = t.wantLax != nil
			}

			laxIsValidTests = append(laxIsValidTests, validationTestCase{
				v:    input,
				want: want,
			})
		}
	}
}

func BenchmarkIsValid(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for range b.N {
		_ = semver.IsValid(test)
	}
}

func BenchmarkIsValidShorter(b *testing.B) {
	test := "1.2.11"

	for range b.N {
		_ = semver.IsValid(test)
	}
}

func BenchmarkIsValidRegex(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for range b.N {
		_ = isValidRegex(test)
	}
}

// Benchmarking whether the regex for semver (from semver.org) could be faster.
// Doesn't seem like it is (at all).
func BenchmarkIsValidRegexShorter(b *testing.B) {
	test := "1.2.11"

	for range b.N {
		_ = isValidRegex(test)
	}
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	for _, tt := range isValidTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if ok := semver.IsValid(tt.v); ok != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.v, ok, !ok)
			}
		})
	}
}

func TestIsValidLax(t *testing.T) {
	t.Parallel()

	for _, tt := range laxIsValidTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if ok := semver.IsValidLax(tt.v); ok != tt.want {
				t.Errorf("IsValidLax(%q) = %v, want %v", tt.v, ok, !ok)
			}
		})
	}
}

func isValidRegex(v string) bool {
	return versionRegex.MatchString(v)
}
