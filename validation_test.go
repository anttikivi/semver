package semver_test

import (
	"regexp"
	"testing"

	"github.com/anttikivi/go-semver"
)

const rawVersionRegex = `^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

var versionRegex *regexp.Regexp

var (
	isValidRegexTests []validationTestCase
	isValidTests      []validationTestCase
	laxIsValidTests   []validationTestCase
)

var baseValidationTests = []baseValidationTestCase{
	{"", false, false},

	{"0.1.0-alpha.24+sha.19031c2.darwin.amd64", true, true},
	{"0.1.0-alpha.24+sha.19031c2-darwin-amd64", true, true},

	{"1,2.3", false, false},
	{"1.2.3,pre", false, false},
	{"1.2.3-pre,hello", false, false},
	{"1.2.3-pre.hello,", false, false},
	{"1.2.3-pre.hello,wrong", false, false},
	{"bad", false, false},
	{"1-alpha.beta.gamma", false, true},
	{"1-pre", false, true},
	{"1+meta", false, true},
	{"1-pre+meta", false, true},
	{"1.2-pre", false, true},
	{"1.2+meta", false, true},
	{"1.2-pre+meta", false, true},
	{"1.0.0-alpha", true, true},
	{"1.0.0-alpha.1", true, true},
	{"1.0.0-alpha.beta", true, true},
	{"1.0.0-beta", true, true},
	{"1.0.0-beta.2", true, true},
	{"1.0.0-beta.11", true, true},
	{"1.0.0-rc.1", true, true},
	{"1", false, true},
	{"1.0", false, true},
	{"1.0.0", true, true},
	{"1.2", false, true},
	{"1.2.0", true, true},
	{"1.2.3-456", true, true},
	{"1.2.3-456.789", true, true},
	{"1.2.3-456-789", true, true},
	{"1.2.3-456a", true, true},
	{"1.2.3-pre", true, true},
	{"1.2.3-pre+meta", true, true},
	{"1.2.3-pre.1", true, true},
	{"1.2.3-zzz", true, true},
	{"1.2.3", true, true},
	{"1.2.3+meta", true, true},
	{"1.2.3+meta-pre", true, true},
	{"1.2.3+meta-pre.sha.256a", true, true},
	{"1.2.3-012a", true, true},
	{"1.2.3-0123", false, false},
	{"01.2.3", false, false},
	{"1.02.3", false, false},
	{"1.2.03", false, false},
}

type baseValidationTestCase struct {
	v          string
	wantStrict bool
	wantLax    bool
}

type validationTestCase struct {
	v    string
	want bool
}

func init() {
	versionRegex = regexp.MustCompile(rawVersionRegex)

	prefixes := map[string]bool{
		"":       true,
		"v":      true,
		"semver": false,
	}
	for prefix, allowed := range prefixes {
		for _, t := range baseValidationTests {
			input := prefix + t.v
			want := allowed && t.wantStrict

			isValidRegexTests = append(isValidRegexTests, validationTestCase{
				v:    input,
				want: want,
			})

			isValidTests = append(isValidTests, validationTestCase{
				v:    input,
				want: want,
			})

			want = allowed && t.wantLax

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

func TestIsValidRegex(t *testing.T) {
	t.Parallel()

	for _, tt := range isValidRegexTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if ok := isValidRegex(tt.v); ok != tt.want {
				t.Errorf("IsValidRegex(%q) = %v, want %v", tt.v, ok, !ok)
			}
		})
	}
}

func isValidRegex(v string) bool {
	return versionRegex.MatchString(v)
}
