package semver_test

import (
	"strconv"
	"testing"

	"github.com/anttikivi/go-semver"
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
				if r := recover(); tt.wantErr == (r == nil) {
					t.Errorf("MustParse(%q) did not panic", tt.v)
				}
			}()

			if got := semver.MustParse(tt.v); !tt.want.StrictEqual(got) {
				t.Errorf("MustParse(%q) = %v, want %v", tt.v, got, tt.want)
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
			if gotErr == nil && tt.wantErr {
				t.Fatalf("Parse(%q) succeeded unexpectedly", tt.v)
			}

			if gotErr != nil && !tt.wantErr {
				t.Errorf("Parse(%q) failed: %v", tt.v, gotErr)
			}

			if !tt.want.StrictEqual(got) {
				t.Errorf("Parse(%q) = %v, want %v", tt.v, got, tt.want)
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
				if r := recover(); tt.wantErr == (r == nil) {
					t.Errorf("MustParseLax(%q) did not panic", tt.v)
				}
			}()

			if got := semver.MustParseLax(tt.v); !tt.want.StrictEqual(got) {
				t.Errorf("MustParseLax(%q) = %v, want %v", tt.v, got, tt.want)
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
			if gotErr == nil && tt.wantErr {
				t.Fatalf("ParseLax(%q) succeeded unexpectedly", tt.v)
			}

			if gotErr != nil && !tt.wantErr {
				t.Errorf("ParseLax(%q) failed: %v", tt.v, gotErr)
			}

			if !tt.want.StrictEqual(got) {
				t.Errorf("ParseLax(%q) = %v, want %v", tt.v, got, tt.want)
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
