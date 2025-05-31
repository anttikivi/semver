package semver_test

import (
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/anttikivi/semver"
)

func TestVersionsSort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input []string
		want  []string
	}{
		{
			[]string{
				"1.2.3",
				"1.0",
				"1.3",
				"2",
				"0.4.2",
			},
			[]string{
				"0.4.2",
				"1.0.0",
				"1.2.3",
				"1.3.0",
				"2.0.0",
			},
		},
		{
			[]string{
				"1.2.3",
				"1.0",
				"10",
				"1.3",
				"2",
				"0.4.2",
			},
			[]string{
				"0.4.2",
				"1.0.0",
				"1.2.3",
				"1.3.0",
				"2.0.0",
				"10.0.0",
			},
		},
		{
			[]string{
				"10",
				"2",
				"12",
				"1.2",
				"1.0",
				"1",
			},
			[]string{
				"1.0.0",
				"1.0.0",
				"1.2.0",
				"2.0.0",
				"10.0.0",
				"12.0.0",
			},
		},
		{
			[]string{
				"1-beta",
				"1",
				"2-beta",
				"2.0.1",
				"1-0.beta",
				"1-alpha",
				"1-alpha.1",
				"1.3",
				"1-alpha.beta",
			},
			[]string{
				"1.0.0-0.beta",
				"1.0.0-alpha",
				"1.0.0-alpha.1",
				"1.0.0-alpha.beta",
				"1.0.0-beta",
				"1.0.0",
				"1.3.0",
				"2.0.0-beta",
				"2.0.1",
			},
		},
		{
			[]string{
				"1.0.0-beta.2",
				"1.0.0-alpha.beta",
				"1.0.0-beta.11",
				"1.0.0",
				"1.0.0-alpha",
				"1.0.0-beta",
				"1.0.0-rc.1",
				"1.0.0-alpha.1",
			},
			[]string{
				"1.0.0-alpha",
				"1.0.0-alpha.1",
				"1.0.0-alpha.beta",
				"1.0.0-beta",
				"1.0.0-beta.2",
				"1.0.0-beta.11",
				"1.0.0-rc.1",
				"1.0.0",
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			versions := make(semver.Versions, len(tt.input))

			for j, s := range tt.input {
				v, _ := semver.ParseLax(s)
				if v == nil {
					t.Fatalf("Setup error: Version is nil for input %q", s)
				}

				versions[j] = v
			}

			sort.Sort(versions)

			x := make([]string, len(versions))

			for j, v := range versions {
				x[j] = v.String()
			}

			if !reflect.DeepEqual(x, tt.want) {
				t.Errorf("sort.Sort(%#v) = %#v, want %#v", tt.input, x, tt.want)
			}
		})
	}
}
