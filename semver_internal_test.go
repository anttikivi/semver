package semver

import "testing"

func BenchmarkIsValidByParse(b *testing.B) {
	test := "0.1.0-alpha.24+sha.19031c2.darwin.amd64"

	for range b.N {
		_ = isValidByParse(test)
	}
}

// isValidByParse is the old implementation of the validation function.
func isValidByParse(s string) bool {
	if _, err := parse(s); err != nil {
		return false
	}

	return true
}
