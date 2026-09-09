# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

Entries before 0.3.0 were reconstructed after the fact from git history and the
generated GitHub release notes, so they are less detailed than entries written
at release time.

## [Unreleased]

### Changed

- Go 1.26.5 to 1.26.8, dependencies updated and `go mod tidy` run.
- `.goreleaser.yml` migrated to GoReleaser v2 configuration. The released
  artifacts and their filenames are unchanged.
- Release workflow pins GoReleaser to the v2 series and fetches full history
  through `actions/checkout` instead of a separate unshallow step.

## [0.3.0] - 2026-07-31

### Added

- `--sample`, the number of seconds to measure CPU usage over. Defaults to `1`.
- Test coverage for argument validation and for the sampling path, replacing an
  empty test stub.

### Changed

- **Breaking:** CPU percentages are now measured over a sample window instead of
  being reported as an average over the process's entire lifetime. Reported
  values reflect current usage and will differ from those of earlier versions.
  Review your `--cpu` thresholds when upgrading.
- **Breaking:** the check now sleeps for `--sample` seconds (default `1`) on
  every run, where it previously returned in milliseconds. If your check
  definition sets a `timeout`, make sure it is larger than the sample.
- **Breaking:** the plugin name and annotation keyspace changed from
  `check-cpu-usage` to `sensu-top-process`. Configuration set through
  `sensu.io/plugins/check-cpu-usage/config` annotations must be moved to
  `sensu.io/plugins/sensu-top-process/config`.
- **Breaking:** `--cpu 0` and `--memory 0` are now rejected. They previously ran
  but produced malformed output, see Fixed below.
- Go 1.24.4 to 1.26.5, dependencies updated and `go mod tidy` run.
- README substantially expanded with the metric format, the full option list and
  the measurement caveats.

### Fixed

- Zero thresholds caused processes that could not be inspected to be emitted
  with an empty name, producing a `..` in the Graphite path that most backends
  reject or store under a mangled key.
- A failure to enumerate processes now reports Unknown instead of silently
  exiting OK having evaluated nothing.
- The README documented the old `check-cpu-usage` annotation keyspace.

### Security

- `google.golang.org/grpc` 1.74.2 to 1.79.3, an indirect dependency, via
  Dependabot's `go_modules` security group.

## [0.2.2] - 2025-08-19

### Changed

- Go version and module updates.
- GitHub Actions workflow updates, including explicit `permissions` blocks.
- `golangci/golangci-lint-action` 6 to 8.

## [0.2.1] - 2025-01-10

### Changed

- Go version and module updates.
- Dependabot configuration.

## [0.2] - 2024-09-16

### Changed

- Go 1.23.1.
- goreleaser-action v6.
- `gopsutil/v3` to 3.24.5, `sensu-plugin-sdk` to 0.19.0, `sensu/core/v2` to
  2.20.0.
- Various CI action updates.

## [0.1.4] - 2023-10-16

Re-tag of 0.1.3 with no code change; both tags point at the same commit.

## [0.1.3] - 2023-08-10

### Added

- `--expand`, to report a process under its full command line instead of its
  name. Useful for interpreters such as bash or powershell, whose processes
  otherwise all share one name.

## [0.1.2] - 2023-06-30

### Fixed

- Removed the `:` that followed the metric name, which made the output invalid
  Graphite plaintext. `foo.process.cpu_percent.bar: 1.0 1688000000` became
  `foo.process.cpu_percent.bar 1.0 1688000000`.

## [0.1.1] - 2023-06-30

### Added

- `--scheme`, a required prefix prepended to every metric.
- A Unix timestamp on each metric line.

## [0.1.0] - 2023-06-29

Initial release.

[Unreleased]: https://github.com/elfranne/sensu-top-process/compare/0.3.0...HEAD
[0.3.0]: https://github.com/elfranne/sensu-top-process/compare/0.2.2...0.3.0
[0.2.2]: https://github.com/elfranne/sensu-top-process/compare/0.2.1...0.2.2
[0.2.1]: https://github.com/elfranne/sensu-top-process/compare/0.2...0.2.1
[0.2]: https://github.com/elfranne/sensu-top-process/compare/0.1.4...0.2
[0.1.4]: https://github.com/elfranne/sensu-top-process/compare/0.1.3...0.1.4
[0.1.3]: https://github.com/elfranne/sensu-top-process/compare/0.1.2...0.1.3
[0.1.2]: https://github.com/elfranne/sensu-top-process/compare/0.1.1...0.1.2
[0.1.1]: https://github.com/elfranne/sensu-top-process/compare/0.1.0...0.1.1
[0.1.0]: https://github.com/elfranne/sensu-top-process/releases/tag/0.1.0
