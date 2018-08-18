# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

For upcoming features, improvements and ideas please refer to [TODO].

## [Unreleased]

### Added

- Restore terminal state and exit when receiving a signal
- Add support for [fish]

### Fixed

- Terminal state wasn't properly restored sometimes.  After removing `stty sane` from the shell-specific input wrappers the problem disappeared.

## [v0.1.4]

### Added

- Restore terminal state after each invocation of `leader` (requires `stty`)
- Add `leader help` subcommand

### Fixed
- Fix configuration load order: project-local overrides were broken

[Unreleased]: https://github.com/dhamidi/leader/compare/v0.1.4...HEAD
[v0.1.4]: https://github.com/dhamidi/leader/compare/v0.1.3...v0.1.4
[TODO]: https://github.com/dhamidi/leader/blob/master/TODO.md
[fish]: https://fishshell.com
