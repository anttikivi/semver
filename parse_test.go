package semver_test

import (
	"strconv"
	"testing"

	"github.com/anttikivi/semver"
)

var (
	parserTests    []parserTestCase
	laxParserTests []parserTestCase
)

type parserTestCase struct {
	v       string
	want    *semver.Version
	wantErr bool
}

type regexVer struct {
	major         int
	minor         int
	patch         int
	prerelease    string
	buildmetadata string
}

func init() {
	for prefix, allowed := range testPrefixes {
		for _, t := range baseTests {
			input := prefix + t.v

			want := t.wantStrict
			wantErr := t.wantStrictErr

			if !allowed {
				want = nil
				wantErr = true
			}

			parserTests = append(parserTests, parserTestCase{
				v:       input,
				want:    want,
				wantErr: wantErr,
			})

			if allowed {
				want = t.wantLax
				wantErr = t.wantLaxErr
			}

			laxParserTests = append(laxParserTests, parserTestCase{
				v:       input,
				want:    want,
				wantErr: wantErr,
			})
		}
	}
}

func TestMustParse(t *testing.T) {
	t.Parallel()

	for _, tt := range parserTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()

				if tt.wantErr {
					if r == nil {
						t.Errorf("MustParse(%q) did NOT panic as expected", tt.v)
					}
				} else {
					if r != nil {
						t.Errorf("MustParse(%q) panicked unexpectedly: %v", tt.v, r)
					}
				}
			}()

			got := semver.MustParse(tt.v)

			if tt.wantErr {
				t.Errorf("MustParse(%q) returned %v but was expected to panic", tt.v, got)

				return
			}

			if !tt.want.Equal(got) {
				t.Errorf("MustParse(%q) = %v, want %v (equal)", tt.v, got, tt.want)
			}

			if !tt.want.StrictEqual(got) {
				t.Errorf("MustParse(%q) = %v, want %v (strictly equal)", tt.v, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	for _, tt := range parserTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := semver.Parse(tt.v)

			if tt.wantErr {
				if gotErr == nil {
					t.Fatalf("Parse(%q) succeeded unexpectedly; got %v, want error", tt.v, got)
				}

				if got != nil {
					t.Errorf(
						"Parse(%q) returned a non-nil version (%+v) but also an error (%v); want (nil, error)",
						tt.v,
						got,
						gotErr,
					)
				}
			} else if gotErr != nil {
				t.Errorf("Parse(%q) failed unexpectedly: %v", tt.v, gotErr)

				return
			}

			if tt.wantErr {
				return
			}

			if !tt.want.Equal(got) {
				t.Errorf("Parse(%q) = %v, want %v (equal)", tt.v, got, tt.want)
			}

			if !tt.want.StrictEqual(got) {
				t.Errorf("Parse(%q) = %v, want %v (strictly equal)", tt.v, got, tt.want)
			}
		})
	}
}

func TestMustParseLax(t *testing.T) {
	t.Parallel()

	for _, tt := range laxParserTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()

				if tt.wantErr {
					if r == nil {
						t.Errorf("MustParseLax(%q) did NOT panic as expected", tt.v)
					}
				} else {
					if r != nil {
						t.Errorf("MustParseLax(%q) panicked unexpectedly: %v", tt.v, r)
					}
				}
			}()

			got := semver.MustParseLax(tt.v)

			if tt.wantErr {
				t.Errorf("MustParseLax(%q) returned %v but was expected to panic", tt.v, got)

				return
			}

			if !tt.want.Equal(got) {
				t.Errorf("MustParseLax(%q) = %v, want %v (equal)", tt.v, got, tt.want)
			}

			if !tt.want.StrictEqual(got) {
				t.Errorf("MustParseLax(%q) = %v, want %v (strictly equal)", tt.v, got, tt.want)
			}
		})
	}
}

func TestParseLax(t *testing.T) {
	t.Parallel()

	for _, tt := range laxParserTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := semver.ParseLax(tt.v)

			if tt.wantErr {
				if gotErr == nil {
					t.Fatalf("ParseLax(%q) succeeded unexpectedly; got %v, want error", tt.v, got)
				}

				if got != nil {
					t.Errorf(
						"ParseLax(%q) returned a non-nil version (%+v) but also an error (%v); want (nil, error)",
						tt.v,
						got,
						gotErr,
					)
				}
			} else if gotErr != nil {
				t.Errorf("ParseLax(%q) failed unexpectedly: %v", tt.v, gotErr)

				return
			}

			if tt.wantErr {
				return
			}

			if !tt.want.Equal(got) {
				t.Errorf("ParseLax(%q) = %v, want %v (equal)", tt.v, got, tt.want)
			}

			if !tt.want.StrictEqual(got) {
				t.Errorf("ParseLax(%q) = %v, want %v (strictly equal)", tt.v, got, tt.want)
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for range b.N {
		_, _ = semver.Parse(test)
	}
}

// To test whether using regexes is faster, looks like its not.
func BenchmarkParseRegex(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for range b.N {
		_ = parseRegex(test)
	}
}

func parseRegex(v string) *regexVer {
	match := versionRegex.FindStringSubmatch(v)
	if match == nil {
		// fmt.Println("No match found!")
		return nil
	}

	names := versionRegex.SubexpNames()

	result := make(map[string]string)

	for i, name := range names {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	major, _ := strconv.Atoi(result["major"])
	minor, _ := strconv.Atoi(result["minor"])
	patch, _ := strconv.Atoi(result["patch"])

	prerelease := ""
	if p, ok := result["prerelease"]; ok {
		prerelease = p
	}

	buildmetadata := ""
	if b, ok := result["buildmetadata"]; ok {
		buildmetadata = b
	}

	return &regexVer{major, minor, patch, prerelease, buildmetadata}
}
