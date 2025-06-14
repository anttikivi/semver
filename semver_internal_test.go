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

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

const emptyName = "empty"

const rawVersionRegex = `^v?(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

var testPrefixes = map[string]bool{
	"":       true,
	"v":      true,
	"semver": false,
}

var (
	versionRegex         *regexp.Regexp
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
		newVersion(0, 1, 0, newTestPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
		false,
		newVersion(0, 1, 0, newTestPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
		false,
		"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
		"0.1.0-alpha.24",
	},
	{
		"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
		newVersion(0, 1, 0, newTestPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
		false,
		newVersion(0, 1, 0, newTestPrerelease("alpha", 24), "sha", "19031c2", "darwin", "amd64"),
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
		newVersion(1, 0, 0, newTestPrerelease("alpha", "beta", "gamma")),
		false,
		"1.0.0-alpha.beta.gamma",
		"1.0.0-alpha.beta.gamma",
	},
	{
		"1-pre",
		nil,
		true,
		newVersion(1, 0, 0, newTestPrerelease("pre")),
		false,
		"1.0.0-pre",
		"1.0.0-pre",
	},
	{
		"1+meta",
		nil,
		true,
		newVersion(1, 0, 0, newTestPrerelease(), "meta"),
		false,
		"1.0.0+meta",
		"1.0.0",
	},
	{
		"1-pre+meta",
		nil,
		true,
		newVersion(1, 0, 0, newTestPrerelease("pre"), "meta"),
		false,
		"1.0.0-pre+meta",
		"1.0.0-pre",
	},
	{
		"1.2-pre",
		nil,
		true,
		newVersion(1, 2, 0, newTestPrerelease("pre")),
		false,
		"1.2.0-pre",
		"1.2.0-pre",
	},
	{
		"1.2+meta",
		nil,
		true,
		newVersion(1, 2, 0, newTestPrerelease(), "meta"),
		false,
		"1.2.0+meta",
		"1.2.0",
	},
	{
		"1.2-pre+meta",
		nil,
		true,
		newVersion(1, 2, 0, newTestPrerelease("pre"), "meta"),
		false,
		"1.2.0-pre+meta",
		"1.2.0-pre",
	},
	{
		"1.0.0-alpha",
		newVersion(1, 0, 0, newTestPrerelease("alpha")),
		false,
		newVersion(1, 0, 0, newTestPrerelease("alpha")),
		false,
		"1.0.0-alpha",
		"1.0.0-alpha",
	},
	{
		"1.0.0-alpha.1",
		newVersion(1, 0, 0, newTestPrerelease("alpha", 1)),
		false,
		newVersion(1, 0, 0, newTestPrerelease("alpha", 1)),
		false,
		"1.0.0-alpha.1",
		"1.0.0-alpha.1",
	},
	{
		"1.0.0-alpha.beta",
		newVersion(1, 0, 0, newTestPrerelease("alpha", "beta")),
		false,
		newVersion(1, 0, 0, newTestPrerelease("alpha", "beta")),
		false,
		"1.0.0-alpha.beta",
		"1.0.0-alpha.beta",
	},
	{
		"1.0.0-beta",
		newVersion(1, 0, 0, newTestPrerelease("beta")),
		false,
		newVersion(1, 0, 0, newTestPrerelease("beta")),
		false,
		"1.0.0-beta",
		"1.0.0-beta",
	},
	{
		"1.0.0-beta.2",
		newVersion(1, 0, 0, newTestPrerelease("beta", 2)),
		false,
		newVersion(1, 0, 0, newTestPrerelease("beta", 2)),
		false,
		"1.0.0-beta.2",
		"1.0.0-beta.2",
	},
	{
		"1.0.0-beta.11",
		newVersion(1, 0, 0, newTestPrerelease("beta", 11)),
		false,
		newVersion(1, 0, 0, newTestPrerelease("beta", 11)),
		false,
		"1.0.0-beta.11",
		"1.0.0-beta.11",
	},
	{
		"1.0.0-rc.1",
		newVersion(1, 0, 0, newTestPrerelease("rc", 1)),
		false,
		newVersion(1, 0, 0, newTestPrerelease("rc", 1)),
		false,
		"1.0.0-rc.1",
		"1.0.0-rc.1",
	},
	{"1", nil, true, newVersion(1, 0, 0, newTestPrerelease()), false, "1.0.0", "1.0.0"},
	{"1.0", nil, true, newVersion(1, 0, 0, newTestPrerelease()), false, "1.0.0", "1.0.0"},
	{
		"1.0.0",
		newVersion(1, 0, 0, newTestPrerelease()),
		false,
		newVersion(1, 0, 0, newTestPrerelease()),
		false,
		"1.0.0",
		"1.0.0",
	},
	{"1.2", nil, true, newVersion(1, 2, 0, newTestPrerelease()), false, "1.2.0", "1.2.0"},
	{
		"1.2.0",
		newVersion(1, 2, 0, newTestPrerelease()),
		false,
		newVersion(1, 2, 0, newTestPrerelease()),
		false,
		"1.2.0",
		"1.2.0",
	},
	{
		"1.2.3-456",
		newVersion(1, 2, 3, newTestPrerelease(456)),
		false,
		newVersion(1, 2, 3, newTestPrerelease(456)),
		false,
		"1.2.3-456",
		"1.2.3-456",
	},
	{
		"1.2.3-456.789",
		newVersion(1, 2, 3, newTestPrerelease(456, 789)),
		false,
		newVersion(1, 2, 3, newTestPrerelease(456, 789)),
		false,
		"1.2.3-456.789",
		"1.2.3-456.789",
	},
	{
		"1.2.3-456-789",
		newVersion(1, 2, 3, newTestPrerelease("456-789")),
		false,
		newVersion(1, 2, 3, newTestPrerelease("456-789")),
		false,
		"1.2.3-456-789",
		"1.2.3-456-789",
	},
	{
		"1.2.3-456a",
		newVersion(1, 2, 3, newTestPrerelease("456a")),
		false,
		newVersion(1, 2, 3, newTestPrerelease("456a")),
		false,
		"1.2.3-456a",
		"1.2.3-456a",
	},
	{
		"1.2.3-pre",
		newVersion(1, 2, 3, newTestPrerelease("pre")),
		false,
		newVersion(1, 2, 3, newTestPrerelease("pre")),
		false,
		"1.2.3-pre",
		"1.2.3-pre",
	},
	{
		"1.2.3-pre+meta",
		newVersion(1, 2, 3, newTestPrerelease("pre"), "meta"),
		false,
		newVersion(1, 2, 3, newTestPrerelease("pre"), "meta"),
		false,
		"1.2.3-pre+meta",
		"1.2.3-pre",
	},
	{
		"1.2.3-pre.1",
		newVersion(1, 2, 3, newTestPrerelease("pre", 1)),
		false,
		newVersion(1, 2, 3, newTestPrerelease("pre", 1)),
		false,
		"1.2.3-pre.1",
		"1.2.3-pre.1",
	},
	{
		"1.2.3-zzz",
		newVersion(1, 2, 3, newTestPrerelease("zzz")),
		false,
		newVersion(1, 2, 3, newTestPrerelease("zzz")),
		false,
		"1.2.3-zzz",
		"1.2.3-zzz",
	},
	{
		"1.2.3",
		newVersion(1, 2, 3, newTestPrerelease()),
		false,
		newVersion(1, 2, 3, newTestPrerelease()),
		false,
		"1.2.3",
		"1.2.3",
	},
	{
		"1.2.3+meta",
		newVersion(1, 2, 3, newTestPrerelease(), "meta"),
		false,
		newVersion(1, 2, 3, newTestPrerelease(), "meta"),
		false,
		"1.2.3+meta",
		"1.2.3",
	},
	{
		"1.2.3+meta-pre",
		newVersion(1, 2, 3, newTestPrerelease(), "meta-pre"),
		false,
		newVersion(1, 2, 3, newTestPrerelease(), "meta-pre"),
		false,
		"1.2.3+meta-pre",
		"1.2.3",
	},
	{
		"1.2.3+meta-pre.sha.256a",
		newVersion(1, 2, 3, newTestPrerelease(), "meta-pre", "sha", "256a"),
		false,
		newVersion(1, 2, 3, newTestPrerelease(), "meta-pre", "sha", "256a"),
		false,
		"1.2.3+meta-pre.sha.256a",
		"1.2.3",
	},
	{
		"1.2.3-012a",
		newVersion(1, 2, 3, newTestPrerelease("012a")),
		false,
		newVersion(1, 2, 3, newTestPrerelease("012a")),
		false,
		"1.2.3-012a",
		"1.2.3-012a",
	},
	{"-1.2.3", nil, true, nil, true, "", ""},
	{"1.-2.3", nil, true, nil, true, "", ""},
	{"1.2.-3", nil, true, nil, true, "", ""},
	{"1.2.3-0123", nil, true, nil, true, "", ""},
	{"01.2.3", nil, true, nil, true, "", ""},
	{"1.02.3", nil, true, nil, true, "", ""},
	{"1.2.03", nil, true, nil, true, "", ""},
	{"01", nil, true, nil, true, "", ""},
	{"1.02", nil, true, nil, true, "", ""},
	{"01.02", nil, true, nil, true, "", ""},
	{
		"0.0.0",
		newVersion(0, 0, 0, newTestPrerelease()),
		false,
		newVersion(0, 0, 0, newTestPrerelease()),
		false,
		"0.0.0",
		"0.0.0",
	},
	{
		"0.0.0-alpha",
		newVersion(0, 0, 0, newTestPrerelease("alpha")),
		false,
		newVersion(0, 0, 0, newTestPrerelease("alpha")),
		false,
		"0.0.0-alpha",
		"0.0.0-alpha",
	},
	{
		"0.0.0+build",
		newVersion(0, 0, 0, newTestPrerelease(), "build"),
		false,
		newVersion(0, 0, 0, newTestPrerelease(), "build"),
		false,
		"0.0.0+build",
		"0.0.0",
	},
	{
		"0.0.0-alpha+build",
		newVersion(0, 0, 0, newTestPrerelease("alpha"), "build"),
		false,
		newVersion(0, 0, 0, newTestPrerelease("alpha"), "build"),
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
		newVersion(math.MaxUint64, 0, 0, newTestPrerelease()),
		false,
		newVersion(math.MaxUint64, 0, 0, newTestPrerelease()),
		false,
		"18446744073709551615.0.0",
		"18446744073709551615.0.0",
	},
	{
		"0.18446744073709551615.0",
		newVersion(0, math.MaxUint64, 0, newTestPrerelease()),
		false,
		newVersion(0, math.MaxUint64, 0, newTestPrerelease()),
		false,
		"0.18446744073709551615.0",
		"0.18446744073709551615.0",
	},
	{
		"0.0.18446744073709551615",
		newVersion(0, 0, math.MaxUint64, newTestPrerelease()),
		false,
		newVersion(0, 0, math.MaxUint64, newTestPrerelease()),
		false,
		"0.0.18446744073709551615",
		"0.0.18446744073709551615",
	},
	{
		"18446744073709551615",
		nil,
		true,
		newVersion(math.MaxUint64, 0, 0, newTestPrerelease()),
		false,
		"18446744073709551615.0.0",
		"18446744073709551615.0.0",
	},
	{
		"18446744073709551615-pre.release",
		nil,
		true,
		newVersion(math.MaxUint64, 0, 0, newTestPrerelease("pre", "release")),
		false,
		"18446744073709551615.0.0-pre.release",
		"18446744073709551615.0.0-pre.release",
	},
	{
		"0.18446744073709551615",
		nil,
		true,
		newVersion(0, math.MaxUint64, 0, newTestPrerelease()),
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
			newTestPrerelease(
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
			newTestPrerelease(
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
			newTestPrerelease(),
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
			newTestPrerelease(),
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
		newVersion(1, 0, 0, newTestPrerelease(strings.Repeat("a", 200)), strings.Repeat("b", 200)),
		false,
		newVersion(1, 0, 0, newTestPrerelease(strings.Repeat("a", 200)), strings.Repeat("b", 200)),
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
	wantStrict    *Version
	wantStrictErr bool
	wantLax       *Version
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
	want    *Version
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
	versionRegex = regexp.MustCompile(rawVersionRegex)

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

func BenchmarkIsValidByParse(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for b.Loop() {
		_ = isValidByParse(test)
	}
}

func BenchmarkParse(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for b.Loop() {
		_, _ = Parse(test)
	}
}

func BenchmarkParseLax(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for b.Loop() {
		_, _ = ParseLax(test)
	}
}

// To test whether using regexes is faster, looks like its not.
func BenchmarkParseRegex(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for b.Loop() {
		_ = parseRegex(test)
	}
}

func FuzzParse(f *testing.F) {
	f.Add("1.0.0")
	f.Add("v2.3.4")
	f.Add("0.0.1")
	f.Add("10.20.30")
	f.Add("1.2.3-alpha")
	f.Add("1.2.3-alpha.1")
	f.Add("1.2.3-0.3.7")
	f.Add("1.2.3-x.7.z.92")
	f.Add("1.2.3+build")
	f.Add("1.2.3+build.123")
	f.Add("1.2.3-beta+exp.sha.5114f85")
	f.Add("1.0.0-alpha+001")
	f.Add("1.0.0+001")
	f.Add("1.0.0-0.3.7+build")

	f.Add("")
	f.Add("v")
	f.Add("1")
	f.Add("1.2")
	f.Add("1.2.3.4")
	f.Add("a.b.c")
	f.Add("1.0.0-alpha..1")
	f.Add("1.0.0-alpha_beta")
	f.Add("1.0.0+build..meta")
	f.Add("1.0.0+build_meta")
	f.Add("1.0.0-01")
	f.Add("01.0.0")
	f.Add("1.01.0")
	f.Add("1.0.01")
	f.Add("1.2.3--")
	f.Add("1.2.3-+")
	f.Add("1.2.3++")
	f.Add("1.2.3+ ")
	f.Add("1.2.3- ")
	f.Add(
		strings.Repeat("1", 100) + "." + strings.Repeat("2", 100) + "." + strings.Repeat("3", 100),
	)
	f.Add("1.2.3-" + strings.Repeat("a", 200))
	f.Add("1.2.3+" + strings.Repeat("b", 200))

	for _, tt := range baseTests {
		f.Add(tt.v)
	}

	f.Fuzz(func(t *testing.T, a string) {
		v, err := Parse(a)
		if err == nil { //nolint:nestif // must be complex
			if v == nil {
				t.Errorf("Parse(%q) returned nil error and a nil version", a)

				return
			}

			s := v.String()

			v2, err2 := Parse(s)
			if err2 != nil {
				t.Errorf(
					"Parse(v.String()) failed for original %q (v.String() = %q): %v",
					a,
					s,
					err2,
				)

				return
			}

			if !v.Equal(v2) {
				t.Errorf(
					"Parse(v.String()) resulted in non-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v,
					s,
					v2,
				)
			}

			if !v.StrictEqual(v2) {
				t.Errorf(
					"Parse(v.String()) resulted in non-strictly-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v,
					s,
					v2,
				)
			}

			if v.Compare(v2) != 0 {
				t.Errorf("v.Compare(%+v) != 0 for %q, parsed: %+v", v2, a, v)
			}

			b := a
			if strings.HasPrefix(a, "v") {
				b = a[1:]
			}

			if b != v.String() {
				t.Errorf(
					"Original input %q (canonical form %q) does not match v.String() %q. This may indicate an unintended canonicalization or parser behavior",
					a,
					b,
					v.String(),
				)
			}

			return
		}

		if err.Error() == "" {
			t.Errorf("Parse(%q) returned an error but the error string was empty", a)
		}

		if !errors.Is(err, ErrInvalidVersion) && !errors.Is(err, strconv.ErrRange) {
			t.Errorf(
				"Parse(%q) returned an error other than ErrInvalidVersion, which indicates an error within the parser: %v",
				a,
				err,
			)
		}
	})
}

func FuzzParseLax(f *testing.F) {
	f.Add("1.0.0")
	f.Add("v2.3.4")
	f.Add("0.0.1")
	f.Add("10.20.30")
	f.Add("1.2.3-alpha")
	f.Add("1.2.3-alpha.1")
	f.Add("1.2.3-0.3.7")
	f.Add("1.2.3-x.7.z.92")
	f.Add("1.2.3+build")
	f.Add("1.2.3+build.123")
	f.Add("1.2.3-beta+exp.sha.5114f85")
	f.Add("1.0.0-alpha+001")
	f.Add("1.0.0+001")
	f.Add("1.0.0-0.3.7+build")

	f.Add("")
	f.Add("v")
	f.Add("1")
	f.Add("1.2")
	f.Add("1.2.3.4")
	f.Add("a.b.c")
	f.Add("1.0.0-alpha..1")
	f.Add("1.0.0-alpha_beta")
	f.Add("1.0.0+build..meta")
	f.Add("1.0.0+build_meta")
	f.Add("1.0.0-01")
	f.Add("01.0.0")
	f.Add("1.01.0")
	f.Add("1.0.01")
	f.Add("1.2.3--")
	f.Add("1.2.3-+")
	f.Add("1.2.3++")
	f.Add("1.2.3+ ")
	f.Add("1.2.3- ")
	f.Add(
		strings.Repeat("1", 100) + "." + strings.Repeat("2", 100) + "." + strings.Repeat("3", 100),
	)
	f.Add("1.2.3-" + strings.Repeat("a", 200))
	f.Add("1.2.3+" + strings.Repeat("b", 200))

	for _, tt := range baseTests {
		f.Add(tt.v)
	}

	f.Fuzz(func(t *testing.T, a string) {
		v, err := ParseLax(a)
		if err == nil { //nolint:nestif // must be complex
			if v == nil {
				t.Errorf("ParseLax(%q) returned nil error and a nil version", a)

				return
			}

			s := v.String()

			v2, err2 := Parse(s)
			if err2 != nil {
				t.Errorf(
					"Parse(v.String()) failed for original %q (v.String() = %q): %v",
					a,
					s,
					err2,
				)

				return
			}

			if !v.Equal(v2) {
				t.Errorf(
					"Parse(v.String()) resulted in non-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v,
					s,
					v2,
				)
			}

			if !v.StrictEqual(v2) {
				t.Errorf(
					"Parse(v.String()) resulted in non-strictly-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v,
					s,
					v2,
				)
			}

			if v.Compare(v2) != 0 {
				t.Errorf("v.Compare(%+v) != 0 for %q, parsed: %+v", v2, a, v)
			}

			v3, err3 := ParseLax(s)
			if err3 != nil {
				t.Errorf(
					"Parse(v.String()) failed for original %q (v.String() = %q): %v",
					a,
					s,
					err3,
				)

				return
			}

			if !v.Equal(v3) {
				t.Errorf(
					"ParseLax(v.String()) resulted in non-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v,
					s,
					v3,
				)
			}

			if !v2.Equal(v3) {
				t.Errorf(
					"ParseLax(v.String()) resulted in non-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v2,
					s,
					v3,
				)
			}

			if !v.StrictEqual(v3) {
				t.Errorf(
					"ParseLax(v.String()) resulted in non-strictly-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v,
					s,
					v3,
				)
			}

			if !v2.StrictEqual(v3) {
				t.Errorf(
					"ParseLax(v.String()) resulted in non-strictly-equal version for %q\nOriginal parsed: %+v\nv.String() = %q\nParse(v.String()) = %+v",
					a,
					v2,
					s,
					v3,
				)
			}

			if v.Compare(v3) != 0 {
				t.Errorf("v.Compare(%+v) != 0 for %q, parsed: %+v", v3, a, v)
			}

			if v2.Compare(v3) != 0 {
				t.Errorf("v2.Compare(%+v) != 0 for %q, parsed: %+v", v3, a, v2)
			}

			return
		}

		if err.Error() == "" {
			t.Errorf("Parse(%q) returned an error but the error string was empty", a)
		}

		if !errors.Is(err, ErrInvalidVersion) && !errors.Is(err, strconv.ErrRange) {
			t.Errorf(
				"Parse(%q) returned an error other than ErrInvalidVersion, which indicates an error within the parser: %v",
				a,
				err,
			)
		}
	})
}

func TestCompare(t *testing.T) {
	t.Parallel()

	for _, tt := range cmpTests {
		v, _ := ParseLax(tt.x)
		w, _ := ParseLax(tt.y)

		t.Run(tt.x+"/"+tt.y, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.x)
			}

			if w == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.y)
			}

			got := Compare(v, w)
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

			got := MustParse(tt.v)

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

			s := got.CoreString()

			rv, rerr := Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()), v.CoreString() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv := &Version{Major: got.Major, Minor: got.Minor, Patch: got.Patch}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.ComparableString()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv = &Version{
				Major:      got.Major,
				Minor:      got.Minor,
				Patch:      got.Patch,
				Prerelease: got.Prerelease,
			}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.String()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			if !got.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (equal)",
					got,
					rv,
					s,
					got,
				)
			}

			if !got.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (strictly equal)",
					got,
					rv,
					s,
					got,
				)
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

			got := MustParseLax(tt.v)

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

			s := got.CoreString()

			rv, rerr := Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()), v.CoreString() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv := &Version{Major: got.Major, Minor: got.Minor, Patch: got.Patch}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.ComparableString()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv = &Version{
				Major:      got.Major,
				Minor:      got.Minor,
				Patch:      got.Patch,
				Prerelease: got.Prerelease,
			}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.String()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			if !got.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (equal)",
					got,
					rv,
					s,
					got,
				)
			}

			if !got.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (strictly equal)",
					got,
					rv,
					s,
					got,
				)
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

			got, gotErr := Parse(tt.v)

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

			s := got.CoreString()

			rv, rerr := Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()), v.CoreString() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv := &Version{Major: got.Major, Minor: got.Minor, Patch: got.Patch}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.ComparableString()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv = &Version{
				Major:      got.Major,
				Minor:      got.Minor,
				Patch:      got.Patch,
				Prerelease: got.Prerelease,
			}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.String()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			if !got.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (equal)",
					got,
					rv,
					s,
					got,
				)
			}

			if !got.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (strictly equal)",
					got,
					rv,
					s,
					got,
				)
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

			got, gotErr := ParseLax(tt.v)

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

			s := got.CoreString()

			rv, rerr := Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()), v.CoreString() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv := &Version{Major: got.Major, Minor: got.Minor, Patch: got.Patch}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.CoreString()) = %v, v.CoreString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.ComparableString()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			tv = &Version{
				Major:      got.Major,
				Minor:      got.Minor,
				Patch:      got.Patch,
				Prerelease: got.Prerelease,
			}
			if !tv.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			if !tv.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.ComparableString()) = %v, v.ComparableString() = %s, want %v (strictly equal)",
					tv,
					rv,
					s,
					got,
				)
			}

			s = got.String()

			rv, rerr = Parse(s)
			if rerr != nil {
				t.Errorf(
					"(round-trip) Parse(%#v.String()), v.String() = %s, failed unexpectedly: %v",
					got,
					s,
					rerr,
				)

				return
			}

			if !got.Equal(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (equal)",
					got,
					rv,
					s,
					got,
				)
			}

			if !got.StrictEqual(rv) {
				t.Errorf(
					"(round-trip) Parse(%#v.String()) = %v, v.String() = %s, want %v (strictly equal)",
					got,
					rv,
					s,
					got,
				)
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	t.Parallel()

	for _, tt := range cmpTests {
		v, _ := ParseLax(tt.x)
		w, _ := ParseLax(tt.y)

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

func TestVersionComparableString(t *testing.T) {
	t.Parallel()

	for _, tt := range coreStringerTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		v, _ := Parse(tt.v)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.v)
			}

			got := v.ComparableString()
			if got != tt.want {
				t.Errorf("Version{%q}.Core() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVersionComparableStringLax(t *testing.T) {
	t.Parallel()

	for _, tt := range laxCoreStringerTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		v, _ := ParseLax(tt.v)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if v == nil {
				t.Fatalf("Setup error: Version is nil for input %q", tt.v)
			}

			got := v.ComparableString()
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

		v, _ := Parse(tt.v)

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

		v, _ := ParseLax(tt.v)

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

// isValidByParse is the old implementation of the validation function.
func isValidByParse(s string) bool {
	if _, err := Parse(s); err != nil {
		return false
	}

	return true
}

func newTestPrerelease(a ...any) Prerelease {
	p, err := newPrerelease(a...)
	if err != nil {
		panic(err)
	}

	return p
}

func newVersion(major, minor, patch uint64, pr Prerelease, b ...string) *Version {
	return &Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: pr,
		Build:      newBuild(b...),
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
