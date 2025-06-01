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

## Install

    go get github.com/anttikivi/semver

## Usage

The module exports a couple of functions to use. The functions accept version
strings that adhere to the semantic versioning. The version strings may start
with a `v` prefix.

The functions also have `Prefix` counterparts that accept one or more prefix
strings that will be allowed in front of the strings. For example, if you want
your version strings to have the form `go1.2.3`, you can pass `"go"` as a prefix
to the prefix version of a function.

**`Parse`**

As you’d expect, this function takes a version string and parses it into a
`Version`. To use a custom prefix, use the `ParsePrefix` function. To panic
instead of returning an error on failure, use `MustParse` or `MustParsePrefix`
functions.

**`IsValid`**

Checks if the given string is a valid version string. To use a custom prefix,
use the `IsValidPrefix` function.

## License

Copyright (c) 2024 Antti Kivi

This project is licensed under the MIT License. For more information, see
[LICENSE](LICENSE).
