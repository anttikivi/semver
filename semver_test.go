package semver_test

import (
	"testing"

	"github.com/anttikivi/go-semver"
)

const emptyName = "empty"

var testPrefixes = map[string]bool{
	"":       true,
	"v":      true,
	"semver": false,
}

var (
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

type stringerTestCase struct {
	v    string
	want string
}

func init() {
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
