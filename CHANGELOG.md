# Changelog

All notable changes to this project will be documented in this file.

This project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `Version.FullString` function that returns the version as a string that
  includes the build metadata.

### Fixed

- Fix the version parser accepting any character after the patch version without
  returning an error.

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

[unreleased]: https://github.com/anttikivi/go-semver/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/anttikivi/go-semver/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/anttikivi/go-semver/releases/tag/v0.1.0
