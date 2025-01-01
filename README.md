# Version Parser for Go

This is my parser for version strings that adhere to [semantic versioning](https://semver.org) in Go. It should implement the semantic versioning spec and be simple to use.

There are quite a few other parsers in Go; I implemented my own because I wanted to. Some of the existing implementations also use regular expressions to parse the version and, according to my benchmarks in this repository, this implementation without regexps is a lot faster.

## Install

    go get github.com/anttikivi/go-semver

## Usage

The module exports a couple of functions to use. The functions accept version strings that adhere to the semantic versioning. The version strings may start with a `v` prefix.

The functions also have `Prefix` counterparts that accept one or more prefix strings that will be allowed in front of the strings. For example, if you want your version strings to have the form `go1.2.3`, you can pass `"go"` as a prefix to the prefix version of a function.

**`Parse`**

As you’d expect, this function takes a version string and parses it into a `Version`. To use a custom prefix, use the `ParsePrefix` function. To panic instead of returning an error on failure, use `MustParse` or `MustParsePrefix` functions.

**`IsValid`**

Checks if the given string is a valid version string. To use a custom prefix, use the `IsValidPrefix` function.

## License

Copyright (c) 2024 Antti Kivi

This project is licensed under the MIT License. For more information, see [LICENSE](LICENSE).
