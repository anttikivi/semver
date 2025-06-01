package semver_test

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/anttikivi/semver"
)

const emptyName = "empty"

var testPrefixes = map[string]bool{
	"":       true,
	"v":      true,
	"semver": false,
}

var (
	parserTests          []parserTestCase
	laxParserTests       []parserTestCase
	stringerTests        []stringerTestCase
	coreStringerTests    []stringerTestCase
	laxStringerTests     []stringerTestCase
	laxCoreStringerTests []stringerTestCase
)

var baseTests = []baseTestCase{
	{"", nil, true, nil, true, "", ""},

	{
		"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
		newVersion(0, 1, 0, newPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
		false,
		newVersion(0, 1, 0, newPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
		false,
		"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
		"0.1.0-alpha.24",
	},
	{
		"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
		newVersion(0, 1, 0, newPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
		false,
		newVersion(0, 1, 0, newPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
		false,
		"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
		"0.1.0-alpha.24",
	},

	{"1,2.3", nil, true, nil, true, "", ""},
	{"1.2.3,pre", nil, true, nil, true, "", ""},
	{"1.2.3-pre,hello", nil, true, nil, true, "", ""},
	{"1.2.3-pre.hello,", nil, true, nil, true, "", ""},
	{"1.2.3-pre.hello,wrong", nil, true, nil, true, "", ""},
	{"bad", nil, true, nil, true, "", ""},
	{
		"1-alpha.beta.gamma",
		nil,
		true,
		newVersion(1, 0, 0, newPrerelease("alpha", "beta", "gamma")),
		false,
		"1.0.0-alpha.beta.gamma",
		"1.0.0-alpha.beta.gamma",
	},
	{
		"1-pre",
		nil,
		true,
		newVersion(1, 0, 0, newPrerelease("pre")),
		false,
		"1.0.0-pre",
		"1.0.0-pre",
	},
	{
		"1+meta",
		nil,
		true,
		newVersion(1, 0, 0, newPrerelease(), "meta"),
		false,
		"1.0.0+meta",
		"1.0.0",
	},
	{
		"1-pre+meta",
		nil,
		true,
		newVersion(1, 0, 0, newPrerelease("pre"), "meta"),
		false,
		"1.0.0-pre+meta",
		"1.0.0-pre",
	},
	{
		"1.2-pre",
		nil,
		true,
		newVersion(1, 2, 0, newPrerelease("pre")),
		false,
		"1.2.0-pre",
		"1.2.0-pre",
	},
	{
		"1.2+meta",
		nil,
		true,
		newVersion(1, 2, 0, newPrerelease(), "meta"),
		false,
		"1.2.0+meta",
		"1.2.0",
	},
	{
		"1.2-pre+meta",
		nil,
		true,
		newVersion(1, 2, 0, newPrerelease("pre"), "meta"),
		false,
		"1.2.0-pre+meta",
		"1.2.0-pre",
	},
	{
		"1.0.0-alpha",
		newVersion(1, 0, 0, newPrerelease("alpha")),
		false,
		newVersion(1, 0, 0, newPrerelease("alpha")),
		false,
		"1.0.0-alpha",
		"1.0.0-alpha",
	},
	{
		"1.0.0-alpha.1",
		newVersion(1, 0, 0, newPrerelease("alpha", 1)),
		false,
		newVersion(1, 0, 0, newPrerelease("alpha", 1)),
		false,
		"1.0.0-alpha.1",
		"1.0.0-alpha.1",
	},
	{
		"1.0.0-alpha.beta",
		newVersion(1, 0, 0, newPrerelease("alpha", "beta")),
		false,
		newVersion(1, 0, 0, newPrerelease("alpha", "beta")),
		false,
		"1.0.0-alpha.beta",
		"1.0.0-alpha.beta",
	},
	{
		"1.0.0-beta",
		newVersion(1, 0, 0, newPrerelease("beta")),
		false,
		newVersion(1, 0, 0, newPrerelease("beta")),
		false,
		"1.0.0-beta",
		"1.0.0-beta",
	},
	{
		"1.0.0-beta.2",
		newVersion(1, 0, 0, newPrerelease("beta", 2)),
		false,
		newVersion(1, 0, 0, newPrerelease("beta", 2)),
		false,
		"1.0.0-beta.2",
		"1.0.0-beta.2",
	},
	{
		"1.0.0-beta.11",
		newVersion(1, 0, 0, newPrerelease("beta", 11)),
		false,
		newVersion(1, 0, 0, newPrerelease("beta", 11)),
		false,
		"1.0.0-beta.11",
		"1.0.0-beta.11",
	},
	{
		"1.0.0-rc.1",
		newVersion(1, 0, 0, newPrerelease("rc", 1)),
		false,
		newVersion(1, 0, 0, newPrerelease("rc", 1)),
		false,
		"1.0.0-rc.1",
		"1.0.0-rc.1",
	},
	{"1", nil, true, newVersion(1, 0, 0, newPrerelease()), false, "1.0.0", "1.0.0"},
	{"1.0", nil, true, newVersion(1, 0, 0, newPrerelease()), false, "1.0.0", "1.0.0"},
	{
		"1.0.0",
		newVersion(1, 0, 0, newPrerelease()),
		false,
		newVersion(1, 0, 0, newPrerelease()),
		false,
		"1.0.0",
		"1.0.0",
	},
	{"1.2", nil, true, newVersion(1, 2, 0, newPrerelease()), false, "1.2.0", "1.2.0"},
	{
		"1.2.0",
		newVersion(1, 2, 0, newPrerelease()),
		false,
		newVersion(1, 2, 0, newPrerelease()),
		false,
		"1.2.0",
		"1.2.0",
	},
	{
		"1.2.3-456",
		newVersion(1, 2, 3, newPrerelease(456)),
		false,
		newVersion(1, 2, 3, newPrerelease(456)),
		false,
		"1.2.3-456",
		"1.2.3-456",
	},
	{
		"1.2.3-456.789",
		newVersion(1, 2, 3, newPrerelease(456, 789)),
		false,
		newVersion(1, 2, 3, newPrerelease(456, 789)),
		false,
		"1.2.3-456.789",
		"1.2.3-456.789",
	},
	{
		"1.2.3-456-789",
		newVersion(1, 2, 3, newPrerelease("456-789")),
		false,
		newVersion(1, 2, 3, newPrerelease("456-789")),
		false,
		"1.2.3-456-789",
		"1.2.3-456-789",
	},
	{
		"1.2.3-456a",
		newVersion(1, 2, 3, newPrerelease("456a")),
		false,
		newVersion(1, 2, 3, newPrerelease("456a")),
		false,
		"1.2.3-456a",
		"1.2.3-456a",
	},
	{
		"1.2.3-pre",
		newVersion(1, 2, 3, newPrerelease("pre")),
		false,
		newVersion(1, 2, 3, newPrerelease("pre")),
		false,
		"1.2.3-pre",
		"1.2.3-pre",
	},
	{
		"1.2.3-pre+meta",
		newVersion(1, 2, 3, newPrerelease("pre"), "meta"),
		false,
		newVersion(1, 2, 3, newPrerelease("pre"), "meta"),
		false,
		"1.2.3-pre+meta",
		"1.2.3-pre",
	},
	{
		"1.2.3-pre.1",
		newVersion(1, 2, 3, newPrerelease("pre", 1)),
		false,
		newVersion(1, 2, 3, newPrerelease("pre", 1)),
		false,
		"1.2.3-pre.1",
		"1.2.3-pre.1",
	},
	{
		"1.2.3-zzz",
		newVersion(1, 2, 3, newPrerelease("zzz")),
		false,
		newVersion(1, 2, 3, newPrerelease("zzz")),
		false,
		"1.2.3-zzz",
		"1.2.3-zzz",
	},
	{
		"1.2.3",
		newVersion(1, 2, 3, newPrerelease()),
		false,
		newVersion(1, 2, 3, newPrerelease()),
		false,
		"1.2.3",
		"1.2.3",
	},
	{
		"1.2.3+meta",
		newVersion(1, 2, 3, newPrerelease(), "meta"),
		false,
		newVersion(1, 2, 3, newPrerelease(), "meta"),
		false,
		"1.2.3+meta",
		"1.2.3",
	},
	{
		"1.2.3+meta-pre",
		newVersion(1, 2, 3, newPrerelease(), "meta-pre"),
		false,
		newVersion(1, 2, 3, newPrerelease(), "meta-pre"),
		false,
		"1.2.3+meta-pre",
		"1.2.3",
	},
	{
		"1.2.3+meta-pre.sha.256a",
		newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"),
		false,
		newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"),
		false,
		"1.2.3+meta-pre.sha.256a",
		"1.2.3",
	},
	{
		"1.2.3-012a",
		newVersion(1, 2, 3, newPrerelease("012a")),
		false,
		newVersion(1, 2, 3, newPrerelease("012a")),
		false,
		"1.2.3-012a",
		"1.2.3-012a",
	},
	{"1.2.3-0123", nil, true, nil, true, "", ""},
	{"01.2.3", nil, true, nil, true, "", ""},
	{"1.02.3", nil, true, nil, true, "", ""},
	{"1.2.03", nil, true, nil, true, "", ""},
	{"01", nil, true, nil, true, "", ""},
	{"1.02", nil, true, nil, true, "", ""},
	{"01.02", nil, true, nil, true, "", ""},
	{
		"0.0.0",
		newVersion(0, 0, 0, newPrerelease()),
		false,
		newVersion(0, 0, 0, newPrerelease()),
		false,
		"0.0.0",
		"0.0.0",
	},
	{
		"0.0.0-alpha",
		newVersion(0, 0, 0, newPrerelease("alpha")),
		false,
		newVersion(0, 0, 0, newPrerelease("alpha")),
		false,
		"0.0.0-alpha",
		"0.0.0-alpha",
	},
	{
		"0.0.0+build",
		newVersion(0, 0, 0, newPrerelease(), "build"),
		false,
		newVersion(0, 0, 0, newPrerelease(), "build"),
		false,
		"0.0.0+build",
		"0.0.0",
	},
	{
		"0.0.0-alpha+build",
		newVersion(0, 0, 0, newPrerelease("alpha"), "build"),
		false,
		newVersion(0, 0, 0, newPrerelease("alpha"), "build"),
		false,
		"0.0.0-alpha+build",
		"0.0.0-alpha",
	},
	{"1.0.0.", nil, true, nil, true, "", ""},
	{"1..0.0", nil, true, nil, true, "", ""},
	{"1..0..0", nil, true, nil, true, "", ""},
	{"1.0..0", nil, true, nil, true, "", ""},
	{"1...0.0", nil, true, nil, true, "", ""},
	{"1...0..0", nil, true, nil, true, "", ""},
	{"1.0...0", nil, true, nil, true, "", ""},
	{"1.0.0.alpha", nil, true, nil, true, "", ""},
	{"1.0.0-alpha..beta", nil, true, nil, true, "", ""},
	{"1.0.0-alpha...beta", nil, true, nil, true, "", ""},
	{"1.0.0-alpha...beta..gamma", nil, true, nil, true, "", ""},
	{"1.0.0-alpha+build..meta", nil, true, nil, true, "", ""},
	{"1.0.0-alpha+build...meta", nil, true, nil, true, "", ""},
	{"1.0.0-alpha+build...meta..data", nil, true, nil, true, "", ""},
	{"1.0.0+build..meta", nil, true, nil, true, "", ""},
	{"1.0.0+build...meta", nil, true, nil, true, "", ""},
	{"1.0.0+build...meta..data", nil, true, nil, true, "", ""},
	{"1.0.0-alpha.", nil, true, nil, true, "", ""},
	{"1.0.0-alpha.+meta", nil, true, nil, true, "", ""},
	{"1.0.0.-alpha", nil, true, nil, true, "", ""},
	{"1.0.0+meta.", nil, true, nil, true, "", ""},
	{"1.0.0.+meta.", nil, true, nil, true, "", ""},
	{"1.0.0-", nil, true, nil, true, "", ""},
	{"1.0.0+", nil, true, nil, true, "", ""},
	{"1.0.0-.+", nil, true, nil, true, "", ""},
	{"+1.0.0", nil, true, nil, true, "", ""},
	{"-1.0.0", nil, true, nil, true, "", ""},
	{
		"18446744073709551615.0.0",
		newVersion(math.MaxUint64, 0, 0, newPrerelease()),
		false,
		newVersion(math.MaxUint64, 0, 0, newPrerelease()),
		false,
		"18446744073709551615.0.0",
		"18446744073709551615.0.0",
	},
	{
		"0.18446744073709551615.0",
		newVersion(0, math.MaxUint64, 0, newPrerelease()),
		false,
		newVersion(0, math.MaxUint64, 0, newPrerelease()),
		false,
		"0.18446744073709551615.0",
		"0.18446744073709551615.0",
	},
	{
		"0.0.18446744073709551615",
		newVersion(0, 0, math.MaxUint64, newPrerelease()),
		false,
		newVersion(0, 0, math.MaxUint64, newPrerelease()),
		false,
		"0.0.18446744073709551615",
		"0.0.18446744073709551615",
	},
	{
		"18446744073709551615",
		nil,
		true,
		newVersion(math.MaxUint64, 0, 0, newPrerelease()),
		false,
		"18446744073709551615.0.0",
		"18446744073709551615.0.0",
	},
	{
		"18446744073709551615-pre.release",
		nil,
		true,
		newVersion(math.MaxUint64, 0, 0, newPrerelease("pre", "release")),
		false,
		"18446744073709551615.0.0-pre.release",
		"18446744073709551615.0.0-pre.release",
	},
	{
		"0.18446744073709551615",
		nil,
		true,
		newVersion(0, math.MaxUint64, 0, newPrerelease()),
		false,
		"0.18446744073709551615.0",
		"0.18446744073709551615.0",
	},
	{"1.0.0-a!b", nil, true, nil, true, "", ""},
	{"1.0.0+c$d", nil, true, nil, true, "", ""},
	{"1.0.0-a_b", nil, true, nil, true, "", ""},
	{
		"1.0.0-a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z",
		newVersion(
			1,
			0,
			0,
			newPrerelease(
				"a",
				"b",
				"c",
				"d",
				"e",
				"f",
				"g",
				"h",
				"i",
				"j",
				"k",
				"l",
				"m",
				"n",
				"o",
				"p",
				"q",
				"r",
				"s",
				"t",
				"u",
				"v",
				"w",
				"x",
				"y",
				"z",
			),
		),
		false,
		newVersion(
			1,
			0,
			0,
			newPrerelease(
				"a",
				"b",
				"c",
				"d",
				"e",
				"f",
				"g",
				"h",
				"i",
				"j",
				"k",
				"l",
				"m",
				"n",
				"o",
				"p",
				"q",
				"r",
				"s",
				"t",
				"u",
				"v",
				"w",
				"x",
				"y",
				"z",
			),
		),
		false,
		"1.0.0-a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z",
		"1.0.0-a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z",
	},
	{
		"1.0.0+a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z",
		newVersion(
			1,
			0,
			0,
			newPrerelease(),
			"a",
			"b",
			"c",
			"d",
			"e",
			"f",
			"g",
			"h",
			"i",
			"j",
			"k",
			"l",
			"m",
			"n",
			"o",
			"p",
			"q",
			"r",
			"s",
			"t",
			"u",
			"v",
			"w",
			"x",
			"y",
			"z",
		),
		false,
		newVersion(
			1,
			0,
			0,
			newPrerelease(),
			"a",
			"b",
			"c",
			"d",
			"e",
			"f",
			"g",
			"h",
			"i",
			"j",
			"k",
			"l",
			"m",
			"n",
			"o",
			"p",
			"q",
			"r",
			"s",
			"t",
			"u",
			"v",
			"w",
			"x",
			"y",
			"z",
		),
		false,
		"1.0.0+a.b.c.d.e.f.g.h.i.j.k.l.m.n.o.p.q.r.s.t.u.v.w.x.y.z",
		"1.0.0",
	},
	{"1.0.0-αlpha", nil, true, nil, true, "", ""}, // Non-ASCII Unicode
	{"1.0.0+bμild", nil, true, nil, true, "", ""}, // Non-ASCII Unicode
	{"1.0.0\x00", nil, true, nil, true, "", ""},   // Null byte
	{"1.0.0\xff", nil, true, nil, true, "", ""},
	{
		"1.0.0-" + strings.Repeat("a", 200) + "+" + strings.Repeat("b", 200),
		newVersion(1, 0, 0, newPrerelease(strings.Repeat("a", 200)), strings.Repeat("b", 200)),
		false,
		newVersion(1, 0, 0, newPrerelease(strings.Repeat("a", 200)), strings.Repeat("b", 200)),
		false,

		"1.0.0-" + strings.Repeat("a", 200) + "+" + strings.Repeat("b", 200),
		"1.0.0-" + strings.Repeat("a", 200),
	},
}

var cmpTests = []cmpTestCase{
	{"1.2.3", "1.3.1", -1},
	{"2.3.4", "1.2.3", 1},
	{"2.2.3", "2.2.2", 1},
	{"2.2.2", "2.2.3", -1},
	{"1", "1", 0},
	{"2.1", "2.1", 0},
	{"3.2-beta", "3.2-beta", 0},
	{"1.3", "1.1.4", 1},
	{"4.5", "4.5-beta", 1},
	{"4.5-beta", "4.5", -1},
	{"4.5-alpha", "4.5-beta", -1},
	{"4.5-alpha", "4.5-alpha", 0},
	{"4.5-beta.2", "4.5-beta.1", 1},
	{"4.5-beta2", "4.5-beta1", 1},
	{"4.5-beta", "4.5-beta.2", -1},
	{"4.5-beta", "4.5-beta.foo", -1},
	{"4.5-beta.2", "4.5-beta", 1},
	{"4.5-beta.foo", "4.5-beta", 1},
	{"1.2+bar", "1.2+baz", 0},
	{"1.0.0-beta.4", "1.0.0-beta.-2", -1},
	{"1.0.0-beta.-2", "1.0.0-beta.-3", -1},
	{"1.0.0-beta.-3", "1.0.0-beta.5", 1},
	{"4.2.3-beta+build", "4.2.3-beta+meta", 0},
}

type baseTestCase struct {
	v             string
	wantStrict    *semver.Version
	wantStrictErr bool
	wantLax       *semver.Version
	wantLaxErr    bool
	wantStr       string
	wantCoreStr   string
}

type cmpTestCase struct {
	x    string
	y    string
	want int
}

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

type stringerTestCase struct {
	v    string
	want string
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

	for prefix, allowed := range testPrefixes {
		for _, t := range baseTests {
			if !allowed || t.wantLax == nil {
				continue
			}

			input := prefix + t.v

			want := ""

			if allowed && t.wantLax != nil {
				want = t.wantStr
			}

			laxStringerTests = append(laxStringerTests, stringerTestCase{
				v:    input,
				want: want,
			})

			if allowed && t.wantLax != nil {
				want = t.wantCoreStr
			}

			laxCoreStringerTests = append(laxCoreStringerTests, stringerTestCase{
				v:    input,
				want: want,
			})

			if t.wantStrict == nil {
				continue
			}

			want = ""

			if allowed && t.wantStrict != nil {
				want = t.wantStr
			}

			stringerTests = append(stringerTests, stringerTestCase{
				v:    input,
				want: want,
			})

			if allowed && t.wantStrict != nil {
				want = t.wantCoreStr
			}

			coreStringerTests = append(coreStringerTests, stringerTestCase{
				v:    input,
				want: want,
			})
		}
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

func TestCompare(t *testing.T) {
	t.Parallel()

	for _, tt := range cmpTests {
		v, _ := semver.ParseLax(tt.x)
		w, _ := semver.ParseLax(tt.y)

		t.Run(tt.x+"/"+tt.y, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.x)
			}

			if w == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.y)
			}

			got := semver.Compare(v, w)
			if got != tt.want {
				t.Errorf("Version{%q}.Compare(%q) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
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

func TestVersionCompare(t *testing.T) {
	t.Parallel()

	for _, tt := range cmpTests {
		v, _ := semver.ParseLax(tt.x)
		w, _ := semver.ParseLax(tt.y)

		t.Run(tt.x+"/"+tt.y, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.x)
			}

			if w == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.y)
			}

			got := v.Compare(w)
			if got != tt.want {
				t.Errorf("Version{%q}.Compare(%q) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestVersionCore(t *testing.T) {
	t.Parallel()

	for _, tt := range coreStringerTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		v, _ := semver.Parse(tt.v)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.v)
			}

			got := v.Core()
			if got != tt.want {
				t.Errorf("Version{%q}.Core() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVersionCoreLax(t *testing.T) {
	t.Parallel()

	for _, tt := range laxCoreStringerTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		v, _ := semver.ParseLax(tt.v)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.v)
			}

			got := v.Core()
			if got != tt.want {
				t.Errorf("ParseLax(%q).Core() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	t.Parallel()

	for _, tt := range stringerTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		v, _ := semver.Parse(tt.v)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.v)
			}

			got := v.String()
			if got != tt.want {
				t.Errorf("Version{%q}.String() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVersionStringLax(t *testing.T) {
	t.Parallel()

	for _, tt := range laxStringerTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		v, _ := semver.ParseLax(tt.v)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.v)
			}

			got := v.String()
			if got != tt.want {
				t.Errorf("ParseLax(%q).String() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func newPrerelease(a ...any) semver.Prerelease {
	p, err := semver.NewPrerelease(a...)
	if err != nil {
		panic(err)
	}

	return p
}

func newVersion(major, minor, patch uint64, pr semver.Prerelease, b ...string) *semver.Version {
	return &semver.Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: pr,
		Build:      semver.NewBuildIdentifiers(b...),
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
