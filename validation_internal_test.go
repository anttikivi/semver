// Copyright (c) 2025 Antti Kivi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package semver

import "testing"

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

	for b.Loop() {
		_ = IsValid(test)
	}
}

func BenchmarkIsValidShorter(b *testing.B) {
	test := "1.2.11"

	for b.Loop() {
		_ = IsValid(test)
	}
}

func BenchmarkIsValidRegex(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for b.Loop() {
		isValidRegex(test)
	}
}

// Benchmarking whether the regex for semver (from semver.org) could be faster.
// Doesn't seem like it is (at all).
func BenchmarkIsValidRegexShorter(b *testing.B) {
	test := "1.2.11"

	for b.Loop() {
		isValidRegex(test)
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

			if ok := IsValid(tt.v); ok != tt.want {
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

			if ok := IsValidLax(tt.v); ok != tt.want {
				t.Errorf("IsValidLax(%q) = %v, want %v", tt.v, ok, !ok)
			}
		})
	}
}

func isValidRegex(v string) {
	_ = versionRegex.MatchString(v)
}
