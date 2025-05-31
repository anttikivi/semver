package semver_test

import (
	"testing"

	"github.com/anttikivi/semver"
)

const emptyName = "empty"

var (
	parserTests = []struct { //nolint:gochecknoglobals // test cases are shared between multiple tests
		v       string
		want    *semver.Version
		wantErr bool
	}{
		{"", nil, true},

		{
			"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},
		{
			"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},

		{"bad", nil, true},
		{"1-alpha.beta.gamma", nil, true},
		{"1-pre", nil, true},
		{"1+meta", nil, true},
		{"1-pre+meta", nil, true},
		{"1.2-pre", nil, true},
		{"1.2+meta", nil, true},
		{"1.2-pre+meta", nil, true},
		{"1.0.0-alpha", newVersion(1, 0, 0, newPrerelease("alpha")), false},
		{"1.0.0-alpha.1", newVersion(1, 0, 0, newPrerelease("alpha", 1)), false},
		{"1.0.0-alpha.beta", newVersion(1, 0, 0, newPrerelease("alpha", "beta")), false},
		{"1.0.0-beta", newVersion(1, 0, 0, newPrerelease("beta")), false},
		{"1.0.0-beta.2", newVersion(1, 0, 0, newPrerelease("beta", 2)), false},
		{"1.0.0-beta.11", newVersion(1, 0, 0, newPrerelease("beta", 11)), false},
		{"1.0.0-rc.1", newVersion(1, 0, 0, newPrerelease("rc", 1)), false},
		{"1", nil, true},
		{"1.0", nil, true},
		{"1.0.0", newVersion(1, 0, 0, newPrerelease()), false},
		{"1.2", nil, true},
		{"1.2.0", newVersion(1, 2, 0, newPrerelease()), false},
		{"1.2.3-456", newVersion(1, 2, 3, newPrerelease(456)), false},
		{"1.2.3-456.789", newVersion(1, 2, 3, newPrerelease(456, 789)), false},
		{"1.2.3-456-789", newVersion(1, 2, 3, newPrerelease("456-789")), false},
		{"1.2.3-456a", newVersion(1, 2, 3, newPrerelease("456a")), false},
		{"1.2.3-pre", newVersion(1, 2, 3, newPrerelease("pre")), false},
		{"1.2.3-pre+meta", newVersion(1, 2, 3, newPrerelease("pre"), "meta"), false},
		{"1.2.3-pre.1", newVersion(1, 2, 3, newPrerelease("pre", 1)), false},
		{"1.2.3-zzz", newVersion(1, 2, 3, newPrerelease("zzz")), false},
		{"1.2.3", newVersion(1, 2, 3, newPrerelease()), false},
		{"1.2.3+meta", newVersion(1, 2, 3, newPrerelease(), "meta"), false},
		{"1.2.3+meta-pre", newVersion(1, 2, 3, newPrerelease(), "meta-pre"), false},
		{"1.2.3+meta-pre.sha.256a", newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"), false},
		{"1.2.3-012a", newVersion(1, 2, 3, newPrerelease("012a")), false},
		{"1.2.3-0123", nil, true},

		{
			"v0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},
		{
			"v0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},

		{"v", nil, true},
		{"vbad", nil, true},
		{"v1-alpha.beta.gamma", nil, true},
		{"v1-pre", nil, true},
		{"v1+meta", nil, true},
		{"v1-pre+meta", nil, true},
		{"v1.2-pre", nil, true},
		{"v1.2+meta", nil, true},
		{"v1.2-pre+meta", nil, true},
		{"v1.0.0-alpha", newVersion(1, 0, 0, newPrerelease("alpha")), false},
		{"v1.0.0-alpha.1", newVersion(1, 0, 0, newPrerelease("alpha", 1)), false},
		{"v1.0.0-alpha.beta", newVersion(1, 0, 0, newPrerelease("alpha", "beta")), false},
		{"v1.0.0-beta", newVersion(1, 0, 0, newPrerelease("beta")), false},
		{"v1.0.0-beta.2", newVersion(1, 0, 0, newPrerelease("beta", 2)), false},
		{"v1.0.0-beta.11", newVersion(1, 0, 0, newPrerelease("beta", 11)), false},
		{"v1.0.0-rc.1", newVersion(1, 0, 0, newPrerelease("rc", 1)), false},
		{"v1", nil, true},
		{"v1.0", nil, true},
		{"v1.0.0", newVersion(1, 0, 0, newPrerelease()), false},
		{"v1.2", nil, true},
		{"v1.2.0", newVersion(1, 2, 0, newPrerelease()), false},
		{"v1.2.3-456", newVersion(1, 2, 3, newPrerelease(456)), false},
		{"v1.2.3-456.789", newVersion(1, 2, 3, newPrerelease(456, 789)), false},
		{"v1.2.3-456-789", newVersion(1, 2, 3, newPrerelease("456-789")), false},
		{"v1.2.3-456a", newVersion(1, 2, 3, newPrerelease("456a")), false},
		{"v1.2.3-pre", newVersion(1, 2, 3, newPrerelease("pre")), false},
		{"v1.2.3-pre+meta", newVersion(1, 2, 3, newPrerelease("pre"), "meta"), false},
		{"v1.2.3-pre.1", newVersion(1, 2, 3, newPrerelease("pre", 1)), false},
		{"v1.2.3-zzz", newVersion(1, 2, 3, newPrerelease("zzz")), false},
		{"v1.2.3", newVersion(1, 2, 3, newPrerelease()), false},
		{"v1.2.3+meta", newVersion(1, 2, 3, newPrerelease(), "meta"), false},
		{"v1.2.3+meta-pre", newVersion(1, 2, 3, newPrerelease(), "meta-pre"), false},
		{"v1.2.3+meta-pre.sha.256a", newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"), false},
		{"v1.2.3-012a", newVersion(1, 2, 3, newPrerelease("012a")), false},
		{"v1.2.3-0123", nil, true},

		{"semver0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},
		{"semver0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},

		{"semver", nil, true},
		{"semverbad", nil, true},
		{"semver1-alpha.beta.gamma", nil, true},
		{"semver1-pre", nil, true},
		{"semver1+meta", nil, true},
		{"semver1-pre+meta", nil, true},
		{"semver1.2-pre", nil, true},
		{"semver1.2+meta", nil, true},
		{"semver1.2-pre+meta", nil, true},
		{"semver1.0.0-alpha", nil, true},
		{"semver1.0.0-alpha.1", nil, true},
		{"semver1.0.0-alpha.beta", nil, true},
		{"semver1.0.0-beta", nil, true},
		{"semver1.0.0-beta.2", nil, true},
		{"semver1.0.0-beta.11", nil, true},
		{"semver1.0.0-rc.1", nil, true},
		{"semver1", nil, true},
		{"semver1.0", nil, true},
		{"semver1.0.0", nil, true},
		{"semver1.2", nil, true},
		{"semver1.2.0", nil, true},
		{"semver1.2.3-456", nil, true},
		{"semver1.2.3-456.789", nil, true},
		{"semver1.2.3-456-789", nil, true},
		{"semver1.2.3-456a", nil, true},
		{"semver1.2.3-pre", nil, true},
		{"semver1.2.3-pre+meta", nil, true},
		{"semver1.2.3-pre.1", nil, true},
		{"semver1.2.3-zzz", nil, true},
		{"semver1.2.3", nil, true},
		{"semver1.2.3+meta", nil, true},
		{"semver1.2.3+meta-pre", nil, true},
		{"semver1.2.3+meta-pre.sha.256a", nil, true},
		{"semver1.2.3-012a", nil, true},
		{"semver1.2.3-0123", nil, true},

		{"sem0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},
		{"sem0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},

		{"sem", nil, true},
		{"sembad", nil, true},
		{"sem1-alpha.beta.gamma", nil, true},
		{"sem1-pre", nil, true},
		{"sem1+meta", nil, true},
		{"sem1-pre+meta", nil, true},
		{"sem1.2-pre", nil, true},
		{"sem1.2+meta", nil, true},
		{"sem1.2-pre+meta", nil, true},
		{"sem1.0.0-alpha", nil, true},
		{"sem1.0.0-alpha.1", nil, true},
		{"sem1.0.0-alpha.beta", nil, true},
		{"sem1.0.0-beta", nil, true},
		{"sem1.0.0-beta.2", nil, true},
		{"sem1.0.0-beta.11", nil, true},
		{"sem1.0.0-rc.1", nil, true},
		{"sem1", nil, true},
		{"sem1.0", nil, true},
		{"sem1.0.0", nil, true},
		{"sem1.2", nil, true},
		{"sem1.2.0", nil, true},
		{"sem1.2.3-456", nil, true},
		{"sem1.2.3-456.789", nil, true},
		{"sem1.2.3-456-789", nil, true},
		{"sem1.2.3-456a", nil, true},
		{"sem1.2.3-pre", nil, true},
		{"sem1.2.3-pre+meta", nil, true},
		{"sem1.2.3-pre.1", nil, true},
		{"sem1.2.3-zzz", nil, true},
		{"sem1.2.3", nil, true},
		{"sem1.2.3+meta", nil, true},
		{"sem1.2.3+meta-pre", nil, true},
		{"sem1.2.3+meta-pre.sha.256a", nil, true},
		{"sem1.2.3-012a", nil, true},
		{"sem1.2.3-0123", nil, true},

		{"se0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},
		{"se0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},

		{"se", nil, true},
		{"sebad", nil, true},
		{"se1-alpha.beta.gamma", nil, true},
		{"se1-pre", nil, true},
		{"se1+meta", nil, true},
		{"se1-pre+meta", nil, true},
		{"se1.2-pre", nil, true},
		{"se1.2+meta", nil, true},
		{"se1.2-pre+meta", nil, true},
		{"se1.0.0-alpha", nil, true},
		{"se1.0.0-alpha.1", nil, true},
		{"se1.0.0-alpha.beta", nil, true},
		{"se1.0.0-beta", nil, true},
		{"se1.0.0-beta.2", nil, true},
		{"se1.0.0-beta.11", nil, true},
		{"se1.0.0-rc.1", nil, true},
		{"se1", nil, true},
		{"se1.0", nil, true},
		{"se1.0.0", nil, true},
		{"se1.2", nil, true},
		{"se1.2.0", nil, true},
		{"se1.2.3-456", nil, true},
		{"se1.2.3-456.789", nil, true},
		{"se1.2.3-456-789", nil, true},
		{"se1.2.3-456a", nil, true},
		{"se1.2.3-pre", nil, true},
		{"se1.2.3-pre+meta", nil, true},
		{"se1.2.3-pre.1", nil, true},
		{"se1.2.3-zzz", nil, true},
		{"se1.2.3", nil, true},
		{"se1.2.3+meta", nil, true},
		{"se1.2.3+meta-pre", nil, true},
		{"se1.2.3+meta-pre.sha.256a", nil, true},
		{"se1.2.3-012a", nil, true},
		{"se1.2.3-0123", nil, true},
	}
	prefixTests = []struct { //nolint:gochecknoglobals // test cases are shared between multiple tests
		v       string
		want    *semver.Version
		wantErr bool
	}{
		{"", nil, true},

		{
			"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},
		{
			"0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},

		{"bad", nil, true},
		{"1-alpha.beta.gamma", nil, true},
		{"1-pre", nil, true},
		{"1+meta", nil, true},
		{"1-pre+meta", nil, true},
		{"1.2-pre", nil, true},
		{"1.2+meta", nil, true},
		{"1.2-pre+meta", nil, true},
		{"1.0.0-alpha", newVersion(1, 0, 0, newPrerelease("alpha")), false},
		{"1.0.0-alpha.1", newVersion(1, 0, 0, newPrerelease("alpha", 1)), false},
		{"1.0.0-alpha.beta", newVersion(1, 0, 0, newPrerelease("alpha", "beta")), false},
		{"1.0.0-beta", newVersion(1, 0, 0, newPrerelease("beta")), false},
		{"1.0.0-beta.2", newVersion(1, 0, 0, newPrerelease("beta", 2)), false},
		{"1.0.0-beta.11", newVersion(1, 0, 0, newPrerelease("beta", 11)), false},
		{"1.0.0-rc.1", newVersion(1, 0, 0, newPrerelease("rc", 1)), false},
		{"1", nil, true},
		{"1.0", nil, true},
		{"1.0.0", newVersion(1, 0, 0, newPrerelease()), false},
		{"1.2", nil, true},
		{"1.2.0", newVersion(1, 2, 0, newPrerelease()), false},
		{"1.2.3-456", newVersion(1, 2, 3, newPrerelease(456)), false},
		{"1.2.3-456.789", newVersion(1, 2, 3, newPrerelease(456, 789)), false},
		{"1.2.3-456-789", newVersion(1, 2, 3, newPrerelease("456-789")), false},
		{"1.2.3-456a", newVersion(1, 2, 3, newPrerelease("456a")), false},
		{"1.2.3-pre", newVersion(1, 2, 3, newPrerelease("pre")), false},
		{"1.2.3-pre+meta", newVersion(1, 2, 3, newPrerelease("pre"), "meta"), false},
		{"1.2.3-pre.1", newVersion(1, 2, 3, newPrerelease("pre", 1)), false},
		{"1.2.3-zzz", newVersion(1, 2, 3, newPrerelease("zzz")), false},
		{"1.2.3", newVersion(1, 2, 3, newPrerelease()), false},
		{"1.2.3+meta", newVersion(1, 2, 3, newPrerelease(), "meta"), false},
		{"1.2.3+meta-pre", newVersion(1, 2, 3, newPrerelease(), "meta-pre"), false},
		{"1.2.3+meta-pre.sha.256a", newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"), false},
		{"1.2.3-012a", newVersion(1, 2, 3, newPrerelease("012a")), false},
		{"1.2.3-0123", nil, true},

		{
			"v0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},
		{
			"v0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},

		{"v", nil, true},
		{"vbad", nil, true},
		{"v1-alpha.beta.gamma", nil, true},
		{"v1-pre", nil, true},
		{"v1+meta", nil, true},
		{"v1-pre+meta", nil, true},
		{"v1.2-pre", nil, true},
		{"v1.2+meta", nil, true},
		{"v1.2-pre+meta", nil, true},
		{"v1.0.0-alpha", newVersion(1, 0, 0, newPrerelease("alpha")), false},
		{"v1.0.0-alpha.1", newVersion(1, 0, 0, newPrerelease("alpha", 1)), false},
		{"v1.0.0-alpha.beta", newVersion(1, 0, 0, newPrerelease("alpha", "beta")), false},
		{"v1.0.0-beta", newVersion(1, 0, 0, newPrerelease("beta")), false},
		{"v1.0.0-beta.2", newVersion(1, 0, 0, newPrerelease("beta", 2)), false},
		{"v1.0.0-beta.11", newVersion(1, 0, 0, newPrerelease("beta", 11)), false},
		{"v1.0.0-rc.1", newVersion(1, 0, 0, newPrerelease("rc", 1)), false},
		{"v1", nil, true},
		{"v1.0", nil, true},
		{"v1.0.0", newVersion(1, 0, 0, newPrerelease()), false},
		{"v1.2", nil, true},
		{"v1.2.0", newVersion(1, 2, 0, newPrerelease()), false},
		{"v1.2.3-456", newVersion(1, 2, 3, newPrerelease(456)), false},
		{"v1.2.3-456.789", newVersion(1, 2, 3, newPrerelease(456, 789)), false},
		{"v1.2.3-456-789", newVersion(1, 2, 3, newPrerelease("456-789")), false},
		{"v1.2.3-456a", newVersion(1, 2, 3, newPrerelease("456a")), false},
		{"v1.2.3-pre", newVersion(1, 2, 3, newPrerelease("pre")), false},
		{"v1.2.3-pre+meta", newVersion(1, 2, 3, newPrerelease("pre"), "meta"), false},
		{"v1.2.3-pre.1", newVersion(1, 2, 3, newPrerelease("pre", 1)), false},
		{"v1.2.3-zzz", newVersion(1, 2, 3, newPrerelease("zzz")), false},
		{"v1.2.3", newVersion(1, 2, 3, newPrerelease()), false},
		{"v1.2.3+meta", newVersion(1, 2, 3, newPrerelease(), "meta"), false},
		{"v1.2.3+meta-pre", newVersion(1, 2, 3, newPrerelease(), "meta-pre"), false},
		{"v1.2.3+meta-pre.sha.256a", newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"), false},
		{"v1.2.3-012a", newVersion(1, 2, 3, newPrerelease("012a")), false},
		{"v1.2.3-0123", nil, true},

		{
			"semver0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},
		{
			"semver0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},

		{"semver", nil, true},
		{"semverbad", nil, true},
		{"semver1-alpha.beta.gamma", nil, true},
		{"semver1-pre", nil, true},
		{"semver1+meta", nil, true},
		{"semver1-pre+meta", nil, true},
		{"semver1.2-pre", nil, true},
		{"semver1.2+meta", nil, true},
		{"semver1.2-pre+meta", nil, true},
		{"semver1.0.0-alpha", newVersion(1, 0, 0, newPrerelease("alpha")), false},
		{"semver1.0.0-alpha.1", newVersion(1, 0, 0, newPrerelease("alpha", 1)), false},
		{"semver1.0.0-alpha.beta", newVersion(1, 0, 0, newPrerelease("alpha", "beta")), false},
		{"semver1.0.0-beta", newVersion(1, 0, 0, newPrerelease("beta")), false},
		{"semver1.0.0-beta.2", newVersion(1, 0, 0, newPrerelease("beta", 2)), false},
		{"semver1.0.0-beta.11", newVersion(1, 0, 0, newPrerelease("beta", 11)), false},
		{"semver1.0.0-rc.1", newVersion(1, 0, 0, newPrerelease("rc", 1)), false},
		{"semver1", nil, true},
		{"semver1.0", nil, true},
		{"semver1.0.0", newVersion(1, 0, 0, newPrerelease()), false},
		{"semver1.2", nil, true},
		{"semver1.2.0", newVersion(1, 2, 0, newPrerelease()), false},
		{"semver1.2.3-456", newVersion(1, 2, 3, newPrerelease(456)), false},
		{"semver1.2.3-456.789", newVersion(1, 2, 3, newPrerelease(456, 789)), false},
		{"semver1.2.3-456-789", newVersion(1, 2, 3, newPrerelease("456-789")), false},
		{"semver1.2.3-456a", newVersion(1, 2, 3, newPrerelease("456a")), false},
		{"semver1.2.3-pre", newVersion(1, 2, 3, newPrerelease("pre")), false},
		{"semver1.2.3-pre+meta", newVersion(1, 2, 3, newPrerelease("pre"), "meta"), false},
		{"semver1.2.3-pre.1", newVersion(1, 2, 3, newPrerelease("pre", 1)), false},
		{"semver1.2.3-zzz", newVersion(1, 2, 3, newPrerelease("zzz")), false},
		{"semver1.2.3", newVersion(1, 2, 3, newPrerelease()), false},
		{"semver1.2.3+meta", newVersion(1, 2, 3, newPrerelease(), "meta"), false},
		{"semver1.2.3+meta-pre", newVersion(1, 2, 3, newPrerelease(), "meta-pre"), false},
		{"semver1.2.3+meta-pre.sha.256a", newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"), false},
		{"semver1.2.3-012a", newVersion(1, 2, 3, newPrerelease("012a")), false},
		{"semver1.2.3-0123", nil, true},

		{"sem0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},
		{"sem0.1.0-alpha.24+sha.19031c2.darwin.amd64", nil, true},

		{"sem", nil, true},
		{"sembad", nil, true},
		{"sem1-alpha.beta.gamma", nil, true},
		{"sem1-pre", nil, true},
		{"sem1+meta", nil, true},
		{"sem1-pre+meta", nil, true},
		{"sem1.2-pre", nil, true},
		{"sem1.2+meta", nil, true},
		{"sem1.2-pre+meta", nil, true},
		{"sem1.0.0-alpha", nil, true},
		{"sem1.0.0-alpha.1", nil, true},
		{"sem1.0.0-alpha.beta", nil, true},
		{"sem1.0.0-beta", nil, true},
		{"sem1.0.0-beta.2", nil, true},
		{"sem1.0.0-beta.11", nil, true},
		{"sem1.0.0-rc.1", nil, true},
		{"sem1", nil, true},
		{"sem1.0", nil, true},
		{"sem1.0.0", nil, true},
		{"sem1.2", nil, true},
		{"sem1.2.0", nil, true},
		{"sem1.2.3-456", nil, true},
		{"sem1.2.3-456.789", nil, true},
		{"sem1.2.3-456-789", nil, true},
		{"sem1.2.3-456a", nil, true},
		{"sem1.2.3-pre", nil, true},
		{"sem1.2.3-pre+meta", nil, true},
		{"sem1.2.3-pre.1", nil, true},
		{"sem1.2.3-zzz", nil, true},
		{"sem1.2.3", nil, true},
		{"sem1.2.3+meta", nil, true},
		{"sem1.2.3+meta-pre", nil, true},
		{"sem1.2.3+meta-pre.sha.256a", nil, true},
		{"sem1.2.3-012a", nil, true},
		{"sem1.2.3-0123", nil, true},

		{
			"se0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},
		{
			"se0.1.0-alpha.24+sha.19031c2.darwin.amd64",
			newVersion(
				0,
				1,
				0,
				newPrerelease("alpha", 24),
				"sha", "19031c2", "darwin", "amd64",
			),
			false,
		},

		{"se", nil, true},
		{"sebad", nil, true},
		{"se1-alpha.beta.gamma", nil, true},
		{"se1-pre", nil, true},
		{"se1+meta", nil, true},
		{"se1-pre+meta", nil, true},
		{"se1.2-pre", nil, true},
		{"se1.2+meta", nil, true},
		{"se1.2-pre+meta", nil, true},
		{"se1.0.0-alpha", newVersion(1, 0, 0, newPrerelease("alpha")), false},
		{"se1.0.0-alpha.1", newVersion(1, 0, 0, newPrerelease("alpha", 1)), false},
		{"se1.0.0-alpha.beta", newVersion(1, 0, 0, newPrerelease("alpha", "beta")), false},
		{"se1.0.0-beta", newVersion(1, 0, 0, newPrerelease("beta")), false},
		{"se1.0.0-beta.2", newVersion(1, 0, 0, newPrerelease("beta", 2)), false},
		{"se1.0.0-beta.11", newVersion(1, 0, 0, newPrerelease("beta", 11)), false},
		{"se1.0.0-rc.1", newVersion(1, 0, 0, newPrerelease("rc", 1)), false},
		{"se1", nil, true},
		{"se1.0", nil, true},
		{"se1.0.0", newVersion(1, 0, 0, newPrerelease()), false},
		{"se1.2", nil, true},
		{"se1.2.0", newVersion(1, 2, 0, newPrerelease()), false},
		{"se1.2.3-456", newVersion(1, 2, 3, newPrerelease(456)), false},
		{"se1.2.3-456.789", newVersion(1, 2, 3, newPrerelease(456, 789)), false},
		{"se1.2.3-456-789", newVersion(1, 2, 3, newPrerelease("456-789")), false},
		{"se1.2.3-456a", newVersion(1, 2, 3, newPrerelease("456a")), false},
		{"se1.2.3-pre", newVersion(1, 2, 3, newPrerelease("pre")), false},
		{"se1.2.3-pre+meta", newVersion(1, 2, 3, newPrerelease("pre"), "meta"), false},
		{"se1.2.3-pre.1", newVersion(1, 2, 3, newPrerelease("pre", 1)), false},
		{"se1.2.3-zzz", newVersion(1, 2, 3, newPrerelease("zzz")), false},
		{"se1.2.3", newVersion(1, 2, 3, newPrerelease()), false},
		{"se1.2.3+meta", newVersion(1, 2, 3, newPrerelease(), "meta"), false},
		{"se1.2.3+meta-pre", newVersion(1, 2, 3, newPrerelease(), "meta-pre"), false},
		{"se1.2.3+meta-pre.sha.256a", newVersion(1, 2, 3, newPrerelease(), "meta-pre", "sha", "256a"), false},
		{"se1.2.3-012a", newVersion(1, 2, 3, newPrerelease("012a")), false},
		{"se1.2.3-0123", nil, true},
	}
)

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

			if got := semver.MustParse(tt.v); !tt.want.Equal(got) {
				t.Errorf("MustParse(%q) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestMustParsePrefix(t *testing.T) {
	t.Parallel()

	for _, tt := range prefixTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				if r := recover(); tt.wantErr == (r == nil) {
					t.Errorf("MustParsePrefix(%q, %q, %q) did not panic", tt.v, "semver", "se")
				}
			}()

			if got := semver.MustParsePrefix(tt.v, "semver", "se"); !tt.want.Equal(got) {
				t.Errorf("MustParse(%q, %q, %q) = %v, want %v", tt.v, "semver", "se", got, tt.want)
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

			if !tt.want.Equal(got) {
				t.Errorf("Parse(%q) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestParsePrefix(t *testing.T) {
	t.Parallel()

	for _, tt := range prefixTests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := semver.ParsePrefix(tt.v, "semver", "se")
			if gotErr == nil && tt.wantErr {
				t.Fatalf("ParsePrefix(%q, %q, %q) succeeded unexpectedly", tt.v, "semver", "se")
			}

			if gotErr != nil && !tt.wantErr {
				t.Errorf("Parse(%q, %q, %q) failed: %v", tt.v, "semver", "se", gotErr)
			}

			if !tt.want.Equal(got) {
				t.Errorf("Parse(%q, %q, %q) = %v, want %v", tt.v, "semver", "se", got, tt.want)
			}
		})
	}
}

func TestVersionString(t *testing.T) { //nolint:funlen // lot's of test cases
	t.Parallel()

	tests := []struct { //nolint:dupl // tests differ between tests
		v    string
		want string
	}{
		{"", ""},

		{"0.1.0-alpha.24+sha.19031c2.darwin.amd64", "0.1.0-alpha.24"},
		{"0.1.0-alpha.24+sha.19031c2-darwin-amd64", "0.1.0-alpha.24"},

		{"bad", ""},
		{"1-alpha.beta.gamma", ""},
		{"1-pre", ""},
		{"1+meta", ""},
		{"1-pre+meta", ""},
		{"1.2-pre", ""},
		{"1.2+meta", ""},
		{"1.2-pre+meta", ""},
		{"1.0.0-alpha", "1.0.0-alpha"},
		{"1.0.0-alpha.1", "1.0.0-alpha.1"},
		{"1.0.0-alpha.beta", "1.0.0-alpha.beta"},
		{"1.0.0-beta", "1.0.0-beta"},
		{"1.0.0-beta.2", "1.0.0-beta.2"},
		{"1.0.0-beta.11", "1.0.0-beta.11"},
		{"1.0.0-rc.1", "1.0.0-rc.1"},
		{"1", ""},
		{"1.0", ""},
		{"1.0.0", "1.0.0"},
		{"1.2", ""},
		{"1.2.0", "1.2.0"},
		{"1.2.3-456", "1.2.3-456"},
		{"1.2.3-456.789", "1.2.3-456.789"},
		{"1.2.3-456-789", "1.2.3-456-789"},
		{"1.2.3-456a", "1.2.3-456a"},
		{"1.2.3-pre", "1.2.3-pre"},
		{"1.2.3-pre+meta", "1.2.3-pre"},
		{"1.2.3-pre.1", "1.2.3-pre.1"},
		{"1.2.3-zzz", "1.2.3-zzz"},
		{"1.2.3", "1.2.3"},
		{"1.2.3+meta", "1.2.3"},
		{"1.2.3+meta-pre", "1.2.3"},
		{"1.2.3+meta-pre.sha.256a", "1.2.3"},
		{"1.2.3-012a", "1.2.3-012a"},
		{"1.2.3-0123", ""},

		{"v", ""},
		{"vbad", ""},
		{"v1-alpha.beta.gamma", ""},
		{"v1-pre", ""},
		{"v1+meta", ""},
		{"v1-pre+meta", ""},
		{"v1.2-pre", ""},
		{"v1.2+meta", ""},
		{"v1.2-pre+meta", ""},
		{"v1.0.0-alpha", "1.0.0-alpha"},
		{"v1.0.0-alpha.1", "1.0.0-alpha.1"},
		{"v1.0.0-alpha.beta", "1.0.0-alpha.beta"},
		{"v1.0.0-beta", "1.0.0-beta"},
		{"v1.0.0-beta.2", "1.0.0-beta.2"},
		{"v1.0.0-beta.11", "1.0.0-beta.11"},
		{"v1.0.0-rc.1", "1.0.0-rc.1"},
		{"v1", ""},
		{"v1.0", ""},
		{"v1.0.0", "1.0.0"},
		{"v1.2", ""},
		{"v1.2.0", "1.2.0"},
		{"v1.2.3-456", "1.2.3-456"},
		{"v1.2.3-456.789", "1.2.3-456.789"},
		{"v1.2.3-456-789", "1.2.3-456-789"},
		{"v1.2.3-456a", "1.2.3-456a"},
		{"v1.2.3-pre", "1.2.3-pre"},
		{"v1.2.3-pre+meta", "1.2.3-pre"},
		{"v1.2.3-pre.1", "1.2.3-pre.1"},
		{"v1.2.3-zzz", "1.2.3-zzz"},
		{"v1.2.3", "1.2.3"},
		{"v1.2.3+meta", "1.2.3"},
		{"v1.2.3+meta-pre", "1.2.3"},
		{"v1.2.3+meta-pre.sha.256a", "1.2.3"},
		{"v1.2.3-012a", "1.2.3-012a"},
		{"v1.2.3-0123", ""},

		{"semverbad", ""},
		{"semver1-alpha.beta.gamma", ""},
		{"semver1-pre", ""},
		{"semver1+meta", ""},
		{"semver1-pre+meta", ""},
		{"semver1.2-pre", ""},
		{"semver1.2+meta", ""},
		{"semver1.2-pre+meta", ""},
		{"semver1.0.0-alpha", ""},
		{"semver1.0.0-alpha.1", ""},
		{"semver1.0.0-alpha.beta", ""},
		{"semver1.0.0-beta", ""},
		{"semver1.0.0-beta.2", ""},
		{"semver1.0.0-beta.11", ""},
		{"semver1.0.0-rc.1", ""},
		{"semver1", ""},
		{"semver1.0", ""},
		{"semver1.0.0", ""},
		{"semver1.2", ""},
		{"semver1.2.0", ""},
		{"semver1.2.3-456", ""},
		{"semver1.2.3-456.789", ""},
		{"semver1.2.3-456-789", ""},
		{"semver1.2.3-456a", ""},
		{"semver1.2.3-pre", ""},
		{"semver1.2.3-pre+meta", ""},
		{"semver1.2.3-pre.1", ""},
		{"semver1.2.3-zzz", ""},
		{"semver1.2.3", ""},
		{"semver1.2.3+meta", ""},
		{"semver1.2.3+meta-pre", ""},
		{"semver1.2.3+meta-pre.sha.256a", ""},
		{"semver1.2.3-012a", ""},
		{"semver1.2.3-0123", ""},
	}
	for _, tt := range tests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := semver.Parse(tt.v)
			if tt.want == "" && got != nil {
				t.Fatalf("Parse(%q) succeeded unexpectedly in the string test", tt.v)
			}

			if got != nil && got.String() != tt.want {
				t.Errorf("Version{%q}.String() = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestVersionStringWithPrefix(t *testing.T) { //nolint:funlen // lot's of test cases
	t.Parallel()

	tests := []struct { //nolint:dupl // tests differ between tests
		v    string
		want string
	}{
		{"", ""},

		{"0.1.0-alpha.24+sha.19031c2.darwin.amd64", "0.1.0-alpha.24"},
		{"0.1.0-alpha.24+sha.19031c2-darwin-amd64", "0.1.0-alpha.24"},

		{"bad", ""},
		{"1-alpha.beta.gamma", ""},
		{"1-pre", ""},
		{"1+meta", ""},
		{"1-pre+meta", ""},
		{"1.2-pre", ""},
		{"1.2+meta", ""},
		{"1.2-pre+meta", ""},
		{"1.0.0-alpha", "1.0.0-alpha"},
		{"1.0.0-alpha.1", "1.0.0-alpha.1"},
		{"1.0.0-alpha.beta", "1.0.0-alpha.beta"},
		{"1.0.0-beta", "1.0.0-beta"},
		{"1.0.0-beta.2", "1.0.0-beta.2"},
		{"1.0.0-beta.11", "1.0.0-beta.11"},
		{"1.0.0-rc.1", "1.0.0-rc.1"},
		{"1", ""},
		{"1.0", ""},
		{"1.0.0", "1.0.0"},
		{"1.2", ""},
		{"1.2.0", "1.2.0"},
		{"1.2.3-456", "1.2.3-456"},
		{"1.2.3-456.789", "1.2.3-456.789"},
		{"1.2.3-456-789", "1.2.3-456-789"},
		{"1.2.3-456a", "1.2.3-456a"},
		{"1.2.3-pre", "1.2.3-pre"},
		{"1.2.3-pre+meta", "1.2.3-pre"},
		{"1.2.3-pre.1", "1.2.3-pre.1"},
		{"1.2.3-zzz", "1.2.3-zzz"},
		{"1.2.3", "1.2.3"},
		{"1.2.3+meta", "1.2.3"},
		{"1.2.3+meta-pre", "1.2.3"},
		{"1.2.3+meta-pre.sha.256a", "1.2.3"},
		{"1.2.3-012a", "1.2.3-012a"},
		{"1.2.3-0123", ""},

		{"v", ""},
		{"vbad", ""},
		{"v1-alpha.beta.gamma", ""},
		{"v1-pre", ""},
		{"v1+meta", ""},
		{"v1-pre+meta", ""},
		{"v1.2-pre", ""},
		{"v1.2+meta", ""},
		{"v1.2-pre+meta", ""},
		{"v1.0.0-alpha", "1.0.0-alpha"},
		{"v1.0.0-alpha.1", "1.0.0-alpha.1"},
		{"v1.0.0-alpha.beta", "1.0.0-alpha.beta"},
		{"v1.0.0-beta", "1.0.0-beta"},
		{"v1.0.0-beta.2", "1.0.0-beta.2"},
		{"v1.0.0-beta.11", "1.0.0-beta.11"},
		{"v1.0.0-rc.1", "1.0.0-rc.1"},
		{"v1", ""},
		{"v1.0", ""},
		{"v1.0.0", "1.0.0"},
		{"v1.2", ""},
		{"v1.2.0", "1.2.0"},
		{"v1.2.3-456", "1.2.3-456"},
		{"v1.2.3-456.789", "1.2.3-456.789"},
		{"v1.2.3-456-789", "1.2.3-456-789"},
		{"v1.2.3-456a", "1.2.3-456a"},
		{"v1.2.3-pre", "1.2.3-pre"},
		{"v1.2.3-pre+meta", "1.2.3-pre"},
		{"v1.2.3-pre.1", "1.2.3-pre.1"},
		{"v1.2.3-zzz", "1.2.3-zzz"},
		{"v1.2.3", "1.2.3"},
		{"v1.2.3+meta", "1.2.3"},
		{"v1.2.3+meta-pre", "1.2.3"},
		{"v1.2.3+meta-pre.sha.256a", "1.2.3"},
		{"v1.2.3-012a", "1.2.3-012a"},
		{"v1.2.3-0123", ""},

		{"semverbad", ""},
		{"semver1-alpha.beta.gamma", ""},
		{"semver1-pre", ""},
		{"semver1+meta", ""},
		{"semver1-pre+meta", ""},
		{"semver1.2-pre", ""},
		{"semver1.2+meta", ""},
		{"semver1.2-pre+meta", ""},
		{"semver1.0.0-alpha", "1.0.0-alpha"},
		{"semver1.0.0-alpha.1", "1.0.0-alpha.1"},
		{"semver1.0.0-alpha.beta", "1.0.0-alpha.beta"},
		{"semver1.0.0-beta", "1.0.0-beta"},
		{"semver1.0.0-beta.2", "1.0.0-beta.2"},
		{"semver1.0.0-beta.11", "1.0.0-beta.11"},
		{"semver1.0.0-rc.1", "1.0.0-rc.1"},
		{"semver1", ""},
		{"semver1.0", ""},
		{"semver1.0.0", "1.0.0"},
		{"semver1.2", ""},
		{"semver1.2.0", "1.2.0"},
		{"semver1.2.3-456", "1.2.3-456"},
		{"semver1.2.3-456.789", "1.2.3-456.789"},
		{"semver1.2.3-456-789", "1.2.3-456-789"},
		{"semver1.2.3-456a", "1.2.3-456a"},
		{"semver1.2.3-pre", "1.2.3-pre"},
		{"semver1.2.3-pre+meta", "1.2.3-pre"},
		{"semver1.2.3-pre.1", "1.2.3-pre.1"},
		{"semver1.2.3-zzz", "1.2.3-zzz"},
		{"semver1.2.3", "1.2.3"},
		{"semver1.2.3+meta", "1.2.3"},
		{"semver1.2.3+meta-pre", "1.2.3"},
		{"semver1.2.3+meta-pre.sha.256a", "1.2.3"},
		{"semver1.2.3-012a", "1.2.3-012a"},
		{"semver1.2.3-0123", ""},
	}
	for _, tt := range tests {
		name := tt.v
		if name == "" {
			name = emptyName
		}

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, _ := semver.ParsePrefix(tt.v, "semver")
			if tt.want == "" && got != nil {
				t.Fatalf("ParsePrefix(%q, %q) succeeded unexpectedly in the string test", tt.v, "semver")
			}

			if got != nil && got.String() != tt.want {
				t.Errorf("ParsePrefix(%q, %q).String() = %v, want %v", tt.v, "semver", got, tt.want)
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

func newVersion(major, minor, patch int, pr semver.Prerelease, b ...string) *semver.Version {
	return &semver.Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: pr,
		Build:      semver.NewBuildIdentifiers(b...),
	}
}
