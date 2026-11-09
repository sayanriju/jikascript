# Changelog

All notable changes to JikaScript are documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

The compiler version and the language specification version are tracked
separately. A compiler release may advance without changing the language.

## [0.3.0] - 2026-11-09

Language specification: `0.2.0`

This release adds multi-file program support and removes the deprecated
`jikascript export` subcommand.

### Added

- Multi-file programs via the `include` statement. A `.jika` file can splice
  another file's top-level statements into the program before checking and code
  generation, resolving paths relative to the including file.
- Include files are limited to top-level use, must end in `.jika`, and reject
  circular references to avoid ambiguous program structure.

### Removed

- The deprecated `jikascript export` subcommand and its remaining code and
  documentation. `jikascript build` and `jikascript flash` are the supported
  ways to generate and build PlatformIO projects directly.

## [0.2.0] - 2026-11-09

Language specification: `0.1.0`

This release replaces the hand-rolled CH32fun toolchain invocation with a
PlatformIO-driven build, and adds a `flash` command.

### Added

- `jikascript build` compiles `.jika` sources or generated `.c` files into
  firmware through PlatformIO, producing `.elf` and `.bin` artifacts.
- `jikascript flash` builds and uploads firmware with `pio run -t upload`,
  accepting `--port` to select the upload device.
- `--board` selects the CH32V003 package variant, defaulting to `ch32v003`
  (`genericCH32V003F4P6`).
- `--env` overrides the PlatformIO environment name.
- `--optimize` selects `size` (`-Os`) or `speed` (`-O2`). It is optional; when
  omitted, the PlatformIO platform default is used unchanged.
- `JIKASCRIPT_PIO` overrides the PlatformIO executable when `pio` is not on
  `PATH`.
- `--output` is accepted as a long form of `-o`.
- `internal/backend/platformio` isolates all PlatformIO knowledge from the
  compiler and CLI.

### Changed

- `jikascript export` is deprecated. `build` and `flash` drive PlatformIO
  directly, so `export` is retained only for users who want a project folder to
  open in the PlatformIO IDE. It prints a deprecation warning to stderr and will
  be removed in a future release.
- `export` now writes the standard PlatformIO layout (`platformio.ini`, `src/`,
  `include/`) instead of a flat directory.
- Build output is `.elf` and `.bin` only; the `.hex` artifact was dropped.
- CLI flags are consistently documented and accepted as `--flag=value`.
- `just` recipes expose their parameters as long options, so they are invoked
  as `just build --source=... --board=...` rather than positionally.

### Removed

- The `JIKASCRIPT_CH32FUN` environment variable. PlatformIO supplies the
  CH32fun framework, so no separate checkout is required.
- The `just get-ch32libs` recipe.

### Notes

- Firmware upload has not been verified against physical hardware. `flash` is
  covered by tests using a stub `pio` executable only.

## [0.1.0]

Initial implementation: lexer, parser, semantic checker, C99 code generator,
standalone cooperative scheduler, and CH32fun GPIO and timer adapters.

[0.3.0]: https://github.com/sayanriju/jikascript/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/sayanriju/jikascript/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/sayanriju/jikascript/releases/tag/v0.1.0
