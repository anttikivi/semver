# Semantic Version Parser for Go

[![CI](https://github.com/anttikivi/go-semver/actions/workflows/ci.yml/badge.svg)](https://github.com/anttikivi/go-semver/actions/workflows/ci.yml)
[![Godoc](https://godoc.org/github.com/anttikivi/go-semver?status.svg)](https://godoc.org/github.com/anttikivi/go-semver)
[![Go Report Card](https://goreportcard.com/badge/github.com/anttikivi/go-semver)](https://goreportcard.com/report/github.com/anttikivi/go-semver)

The `semver` package provides utilities and a parser to work with version
numbers that adhere to [semantic versioning](https://semver.org). The goal of
this parser is to be reliable and performant. Reliability is ensured by using a
wide range of tests and fuzzing. Performance is achieved by implementing a
custom parser instead of the common alternative: regular expressions.

This package implements
[semantic versioning 2.0.0](https://semver.org/spec/v2.0.0.html). Specifically,
the current capabilities of this package include:

- Parsing version strings.
- Checking if a string is valid version string. This check doesn’t require full
  parsing of the version.
- Comparing versions.
- Sorting versions.

The version strings can optionally have a `"v"` prefix.

Future versions of this library will probably include the following planned
features:

- Version ranges and constraints.
- Wildcard versions.
- Database compatibility.
- JSON compatibility.
- TextMarshaler and TextUnmarshaler compatibility.
- See how the parser could be made faster.

## Install

    go get github.com/anttikivi/semver

## Usage

The functions accept version strings that adhere to the semantic versioning. The
version strings may start with a `v` prefix.

### Parsing versions

The package includes two types of functions for parsing versions. There are the
`Parse` and `ParseLax` functions. `Parse` parses only full valid version strings
like `"1.2.3"`, `"1.2.3-beta.1"`, or `"1.2.3-beta.1+darwin.amd64"`. `ParseLax`
works otherwise like `Parse` but it tries to coerse incomplete core version into
a full version. For example, it parses `"v1"` as `1.0.0` and `"1.2-beta"` as
`1.2.0-beta`. Both functions return a pointer to the `Version` object and an
error.

They can be used as follows:

```go
v, err := semver.Parse("1.2.3-beta.1")
```

The package also offers `MustParse` and `MustParseLax` variants of these
functions. They are otherwise the same but only return the pointer to `Version`.
They panic on errors.

### Validating version strings

The package includes two functions, similar to the parsing functions, for
checking if a string is a valid version string. The functions are `IsValid` and
`IsValidLax` and they return a single boolean value. The return value is
analogous to whether the matching parsing function would parse the given string.

Example usage:

```go
ok := semver.IsValid("1.2.3-beta.1")
```

### Sorting versions

The package contains the `Versions` type that supports sorting using the Go
standard library `sort` package. `Versions` is defined as `[]*Version`.

Example usage:

```go
a := []string{"1.2.3", "1.0", "1.3", "2", "0.4.2"}
slice := make(Versions, len(a))

for i, s := range {
  slice[i] = semver.MustParseLax(s)
}

sort.Sort(slice)

for _, v := range slice {
  fmt.Println(v.String())
}
```

The above code would print:

```
0.4.2
1.0.0
1.2.3
1.3.0
2.0.0
```

## Security

This code should be safe to use in a project and to ensure that, security is an
important consideration for the project. It includes tooling to help with
securing the code, like fuzz testing, the
[CodeQL](https://github.com/anttikivi/semver/actions/workflows/codeql.yml), and
strict
[Golangci-lint](https://github.com/anttikivi/semver/blob/main/.golangci.yml)
ruleset.

If you think you have found a security vulnerability, please disclose it
privately according to the
[security policy](https://github.com/anttikivi/semver/security/policy).

## License

Copyright (c) 2024 Antti Kivi

This project is licensed under the MIT License. For more information, see
[LICENSE](LICENSE).
