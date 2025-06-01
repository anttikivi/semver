# Changelog

All notable changes to this project will be documented in this file.

This project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `Version.CoreString` function that returns the version as a string without the
  pre-release or the build metadata.
- `Version.ComparableString` function that returns the version as a string
  without the build metadata.
- `Version.StrictEqual` function that compares the whole version data to
  determine the equality of two versions, including build metadata.
- `IsValidLax` for checking if partial version strings are valid.
- `ParseLax` and `MustParseLax` for parsing partial version strings.
- `Compare` and `Version.Compare` for comparing versions.
- `Versions` slice type the is a slice of `*Version` and implements
  `sort.Interface` for sorting version numbers.

### Changed

- **BREAKING:** Change all of the number values in the versions to `uint64`s.
- **BREAKING:** `Version.Equal` function to only compare the version parts up to
  the build metadata as the build metadata is not comparable in the semantic
  versioning specification.
- **BREAKING:** `Version.String` to include the build metadata in the string.
- **BREAKING:** Change the `Prerelease` type to be a simple slice of pre-release
  identifiers.
- Make `PrereleaseIdentifier` public.

### Removed

- **BREAKING:** The `Prefix` variants of the functions: `IsValidPrefix`,
  `ParsePrefix`, and `MustParsePrefix` as the Go standard library offers an easy
  way to remove prefixes from strings.
- Private `rawStr` field from `Version` struct.

### Fixed

- Fix the version parser accepting any character after the patch version without
  returning an error.
- Fix `IsValid` accepting version strings that had leading zero.
- Fix many edge cases in `IsValid` and `IsValidLax`.
- Add the missing date to the changelog entry for v0.3.0.

## [0.3.0] - 2025-05-31

### Changed

- **BREAKING:** Change the module name from `github.com/anttikivi/go-semver` to
  `github.com/anttikivi/semver`.

## [0.2.0] - 2025-01-01

### Added

- New implementation for functions `IsValid` and `IsValidPrefix` that doesn’t
  parse the version but just checks for the validity. This drastically speeds up
  checking whether a string is a valid version.

## [0.1.0] - 2024-12-31

- Initial release of the project.
- Functions `Parse` and `MustParse` for parsing version strings.
- Functions `ParsePrefix` and `MustParsePrefix` for parsing version strings with
  optional prefixes.

[unreleased]: https://github.com/anttikivi/semver/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/anttikivi/semver/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/anttikivi/go-semver/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/anttikivi/go-semver/releases/tag/v0.1.0
