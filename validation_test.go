package semver_test

import (
	"regexp"
	"testing"

	"github.com/anttikivi/go-semver"
)

const rawVersionRegex = `^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

var versionRegex *regexp.Regexp

func init() { //nolint:gochecknoinits // needed for these tests
	versionRegex = regexp.MustCompile(rawVersionRegex)
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

	tests := []struct {
		v    string
		want bool
	}{
		{"", false},

		{"0.1.0-alpha.24+sha.19031c2.darwin.amd64", true},
		{"0.1.0-alpha.24+sha.19031c2-darwin-amd64", true},

		{"1,2.3", false},
		{"1.2.3,pre", false},
		{"1.2.3-pre,hello", false},
		{"1.2.3-pre.hello,", false},
		{"1.2.3-pre.hello,wrong", false},
		{"bad", false},
		{"1-alpha.beta.gamma", false},
		{"1-pre", false},
		{"1+meta", false},
		{"1-pre+meta", false},
		{"1.2-pre", false},
		{"1.2+meta", false},
		{"1.2-pre+meta", false},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-alpha.beta", true},
		{"1.0.0-beta", true},
		{"1.0.0-beta.2", true},
		{"1.0.0-beta.11", true},
		{"1.0.0-rc.1", true},
		{"1", false},
		{"1.0", false},
		{"1.0.0", true},
		{"1.2", false},
		{"1.2.0", true},
		{"1.2.3-456", true},
		{"1.2.3-456.789", true},
		{"1.2.3-456-789", true},
		{"1.2.3-456a", true},
		{"1.2.3-pre", true},
		{"1.2.3-pre+meta", true},
		{"1.2.3-pre.1", true},
		{"1.2.3-zzz", true},
		{"1.2.3", true},
		{"1.2.3+meta", true},
		{"1.2.3+meta-pre", true},
		{"1.2.3+meta-pre.sha.256a", true},
		{"1.2.3-012a", true},
		{"1.2.3-0123", false},
		{"01.2.3", false},
		{"1.02.3", false},
		{"1.2.03", false},

		{"v", false},
		{"vbad", false},
		{"v1-alpha.beta.gamma", false},
		{"v1-pre", false},
		{"v1+meta", false},
		{"v1-pre+meta", false},
		{"v1.2-pre", false},
		{"v1.2+meta", false},
		{"v1.2-pre+meta", false},
		{"v1.0.0-alpha", true},
		{"v1.0.0-alpha.1", true},
		{"v1.0.0-alpha.beta", true},
		{"v1.0.0-beta", true},
		{"v1.0.0-beta.2", true},
		{"v1.0.0-beta.11", true},
		{"v1.0.0-rc.1", true},
		{"v1", false},
		{"v1.0", false},
		{"v1.0.0", true},
		{"v1.2", false},
		{"v1.2.0", true},
		{"v1.2.3-456", true},
		{"v1.2.3-456.789", true},
		{"v1.2.3-456-789", true},
		{"v1.2.3-456a", true},
		{"v1.2.3-pre", true},
		{"v1.2.3-pre+meta", true},
		{"v1.2.3-pre.1", true},
		{"v1.2.3-zzz", true},
		{"v1.2.3", true},
		{"v1.2.3+meta", true},
		{"v1.2.3+meta-pre", true},
		{"v1.2.3+meta-pre.sha.256a", true},
		{"v1.2.3-012a", true},
		{"v1.2.3-0123", false},
		{"v01.2.3", false},
		{"v1.02.3", false},
		{"v1.2.03", false},

		{"semver", false},
		{"semverbad", false},
		{"semver1-alpha.beta.gamma", false},
		{"semver1-pre", false},
		{"semver1+meta", false},
		{"semver1-pre+meta", false},
		{"semver1.2-pre", false},
		{"semver1.2+meta", false},
		{"semver1.2-pre+meta", false},
		{"semver1.0.0-alpha", false},
		{"semver1.0.0-alpha.1", false},
		{"semver1.0.0-alpha.beta", false},
		{"semver1.0.0-beta", false},
		{"semver1.0.0-beta.2", false},
		{"semver1.0.0-beta.11", false},
		{"semver1.0.0-rc.1", false},
		{"semver1", false},
		{"semver1.0", false},
		{"semver1.0.0", false},
		{"semver1.2", false},
		{"semver1.2.0", false},
		{"semver1.2.3-456", false},
		{"semver1.2.3-456.789", false},
		{"semver1.2.3-456-789", false},
		{"semver1.2.3-456a", false},
		{"semver1.2.3-pre", false},
		{"semver1.2.3-pre+meta", false},
		{"semver1.2.3-pre.1", false},
		{"semver1.2.3-zzz", false},
		{"semver1.2.3", false},
		{"semver1.2.3+meta", false},
		{"semver1.2.3+meta-pre", false},
		{"semver1.2.3+meta-pre.sha.256a", false},
		{"semver1.2.3-012a", false},
		{"semver1.2.3-0123", false},
		{"semver01.2.3", false},
		{"semver1.02.3", false},
		{"semver1.2.03", false},
	}
	for _, tt := range tests {
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

	tests := []struct {
		v    string
		want bool
	}{
		{"", false},

		{"0.1.0-alpha.24+sha.19031c2.darwin.amd64", true},
		{"0.1.0-alpha.24+sha.19031c2-darwin-amd64", true},

		{"1,2.3", false},
		{"1.2.3,pre", false},
		{"1.2.3-pre,hello", false},
		{"1.2.3-pre.hello,", false},
		{"1.2.3-pre.hello,wrong", false},
		{"bad", false},
		{"1-alpha.beta.gamma", true},
		{"1-pre", true},
		{"1+meta", true},
		{"1-pre+meta", true},
		{"1.2-pre", true},
		{"1.2+meta", true},
		{"1.2-pre+meta", true},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-alpha.beta", true},
		{"1.0.0-beta", true},
		{"1.0.0-beta.2", true},
		{"1.0.0-beta.11", true},
		{"1.0.0-rc.1", true},
		{"1", true},
		{"1.0", true},
		{"1.0.0", true},
		{"1.2", true},
		{"1.2.0", true},
		{"1.2.3-456", true},
		{"1.2.3-456.789", true},
		{"1.2.3-456-789", true},
		{"1.2.3-456a", true},
		{"1.2.3-pre", true},
		{"1.2.3-pre+meta", true},
		{"1.2.3-pre.1", true},
		{"1.2.3-zzz", true},
		{"1.2.3", true},
		{"1.2.3+meta", true},
		{"1.2.3+meta-pre", true},
		{"1.2.3+meta-pre.sha.256a", true},
		{"1.2.3-012a", true},
		{"1.2.3-0123", false},
		{"01.2.3", false},
		{"1.02.3", false},
		{"1.2.03", false},
		{"01", false},
		{"1.02", false},
		{"01.02", false},

		{"v", false},
		{"vbad", false},
		{"v1-alpha.beta.gamma", true},
		{"v1-pre", true},
		{"v1+meta", true},
		{"v1-pre+meta", true},
		{"v1.2-pre", true},
		{"v1.2+meta", true},
		{"v1.2-pre+meta", true},
		{"v1.0.0-alpha", true},
		{"v1.0.0-alpha.1", true},
		{"v1.0.0-alpha.beta", true},
		{"v1.0.0-beta", true},
		{"v1.0.0-beta.2", true},
		{"v1.0.0-beta.11", true},
		{"v1.0.0-rc.1", true},
		{"v1", true},
		{"v1.0", true},
		{"v1.0.0", true},
		{"v1.2", true},
		{"v1.2.0", true},
		{"v1.2.3-456", true},
		{"v1.2.3-456.789", true},
		{"v1.2.3-456-789", true},
		{"v1.2.3-456a", true},
		{"v1.2.3-pre", true},
		{"v1.2.3-pre+meta", true},
		{"v1.2.3-pre.1", true},
		{"v1.2.3-zzz", true},
		{"v1.2.3", true},
		{"v1.2.3+meta", true},
		{"v1.2.3+meta-pre", true},
		{"v1.2.3+meta-pre.sha.256a", true},
		{"v1.2.3-012a", true},
		{"v1.2.3-0123", false},
		{"v01.2.3", false},
		{"v1.02.3", false},
		{"v1.2.03", false},
		{"v01", false},
		{"v1.02", false},
		{"v01.02", false},

		{"semver", false},
		{"semverbad", false},
		{"semver1-alpha.beta.gamma", false},
		{"semver1-pre", false},
		{"semver1+meta", false},
		{"semver1-pre+meta", false},
		{"semver1.2-pre", false},
		{"semver1.2+meta", false},
		{"semver1.2-pre+meta", false},
		{"semver1.0.0-alpha", false},
		{"semver1.0.0-alpha.1", false},
		{"semver1.0.0-alpha.beta", false},
		{"semver1.0.0-beta", false},
		{"semver1.0.0-beta.2", false},
		{"semver1.0.0-beta.11", false},
		{"semver1.0.0-rc.1", false},
		{"semver1", false},
		{"semver1.0", false},
		{"semver1.0.0", false},
		{"semver1.2", false},
		{"semver1.2.0", false},
		{"semver1.2.3-456", false},
		{"semver1.2.3-456.789", false},
		{"semver1.2.3-456-789", false},
		{"semver1.2.3-456a", false},
		{"semver1.2.3-pre", false},
		{"semver1.2.3-pre+meta", false},
		{"semver1.2.3-pre.1", false},
		{"semver1.2.3-zzz", false},
		{"semver1.2.3", false},
		{"semver1.2.3+meta", false},
		{"semver1.2.3+meta-pre", false},
		{"semver1.2.3+meta-pre.sha.256a", false},
		{"semver1.2.3-012a", false},
		{"semver1.2.3-0123", false},
		{"semver01.2.3", false},
		{"semver1.02.3", false},
		{"semver1.2.03", false},
		{"semver01", false},
		{"semver1.02", false},
		{"semver01.02", false},
	}
	for _, tt := range tests {
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

	tests := []struct {
		v    string
		want bool
	}{
		{"", false},

		{"0.1.0-alpha.24+sha.19031c2.darwin.amd64", true},
		{"0.1.0-alpha.24+sha.19031c2-darwin-amd64", true},

		{"bad", false},
		{"1-alpha.beta.gamma", false},
		{"1-pre", false},
		{"1+meta", false},
		{"1-pre+meta", false},
		{"1.2-pre", false},
		{"1.2+meta", false},
		{"1.2-pre+meta", false},
		{"1.0.0-alpha", true},
		{"1.0.0-alpha.1", true},
		{"1.0.0-alpha.beta", true},
		{"1.0.0-beta", true},
		{"1.0.0-beta.2", true},
		{"1.0.0-beta.11", true},
		{"1.0.0-rc.1", true},
		{"1", false},
		{"1.0", false},
		{"1.0.0", true},
		{"1.2", false},
		{"1.2.0", true},
		{"1.2.3-456", true},
		{"1.2.3-456.789", true},
		{"1.2.3-456-789", true},
		{"1.2.3-456a", true},
		{"1.2.3-pre", true},
		{"1.2.3-pre+meta", true},
		{"1.2.3-pre.1", true},
		{"1.2.3-zzz", true},
		{"1.2.3", true},
		{"1.2.3+meta", true},
		{"1.2.3+meta-pre", true},
		{"1.2.3+meta-pre.sha.256a", true},
		{"1.2.3-012a", true},
		{"1.2.3-0123", false},
		{"01.2.3", false},
		{"1.02.3", false},
		{"1.2.03", false},
		{"01", false},
		{"1.02", false},
		{"01.02", false},

		{"v", false},
		{"vbad", false},
		{"v1-alpha.beta.gamma", false},
		{"v1-pre", false},
		{"v1+meta", false},
		{"v1-pre+meta", false},
		{"v1.2-pre", false},
		{"v1.2+meta", false},
		{"v1.2-pre+meta", false},
		{"v1.0.0-alpha", true},
		{"v1.0.0-alpha.1", true},
		{"v1.0.0-alpha.beta", true},
		{"v1.0.0-beta", true},
		{"v1.0.0-beta.2", true},
		{"v1.0.0-beta.11", true},
		{"v1.0.0-rc.1", true},
		{"v1", false},
		{"v1.0", false},
		{"v1.0.0", true},
		{"v1.2", false},
		{"v1.2.0", true},
		{"v1.2.3-456", true},
		{"v1.2.3-456.789", true},
		{"v1.2.3-456-789", true},
		{"v1.2.3-456a", true},
		{"v1.2.3-pre", true},
		{"v1.2.3-pre+meta", true},
		{"v1.2.3-pre.1", true},
		{"v1.2.3-zzz", true},
		{"v1.2.3", true},
		{"v1.2.3+meta", true},
		{"v1.2.3+meta-pre", true},
		{"v1.2.3+meta-pre.sha.256a", true},
		{"v1.2.3-012a", true},
		{"v1.2.3-0123", false},
		{"v01.2.3", false},
		{"v1.02.3", false},
		{"v1.2.03", false},
		{"v01", false},
		{"v1.02", false},
		{"v01.02", false},

		{"semver", false},
		{"semverbad", false},
		{"semver1-alpha.beta.gamma", false},
		{"semver1-pre", false},
		{"semver1+meta", false},
		{"semver1-pre+meta", false},
		{"semver1.2-pre", false},
		{"semver1.2+meta", false},
		{"semver1.2-pre+meta", false},
		{"semver1.0.0-alpha", false},
		{"semver1.0.0-alpha.1", false},
		{"semver1.0.0-alpha.beta", false},
		{"semver1.0.0-beta", false},
		{"semver1.0.0-beta.2", false},
		{"semver1.0.0-beta.11", false},
		{"semver1.0.0-rc.1", false},
		{"semver1", false},
		{"semver1.0", false},
		{"semver1.0.0", false},
		{"semver1.2", false},
		{"semver1.2.0", false},
		{"semver1.2.3-456", false},
		{"semver1.2.3-456.789", false},
		{"semver1.2.3-456-789", false},
		{"semver1.2.3-456a", false},
		{"semver1.2.3-pre", false},
		{"semver1.2.3-pre+meta", false},
		{"semver1.2.3-pre.1", false},
		{"semver1.2.3-zzz", false},
		{"semver1.2.3", false},
		{"semver1.2.3+meta", false},
		{"semver1.2.3+meta-pre", false},
		{"semver1.2.3+meta-pre.sha.256a", false},
		{"semver1.2.3-012a", false},
		{"semver1.2.3-0123", false},
		{"semver01.2.3", false},
		{"semver1.02.3", false},
		{"semver1.2.03", false},
		{"semver01", false},
		{"semver1.02", false},
		{"semver01.02", false},
	}
	for _, tt := range tests {
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
