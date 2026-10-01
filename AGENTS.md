# AGENTS.md

This file provides guidance to coding agents when working with code in this repository.

## Overview

Module `github.com/ctx42/xflag` is a small Go library that extends the standard
library `flag` package. The library code lives in the single package
`pkg/xflag`; there is no executable — it is imported by other programs.

## Commands

```bash
go test ./...                          # run all tests
go test ./pkg/xflag/                   # test the package
go test -run Test_FlagSet ./pkg/xflag/ # run a single top-level test
go test -run Test_FlagSet/flag ./pkg/xflag/  # run a subtest by name
go vet ./...                           # vet
golangci-lint run                      # lint (config is gitignored / CI-provided)
```

## Architecture

The package wraps and extends `*flag.FlagSet` rather than replacing it. Two
files split the responsibilities:

- **`pkg/xflag/flag_set.go`** — `FlagSet` embeds `*flag.FlagSet` and adds:
  - **Required flags**: `Required(name)` (panics if parsed already, flag
    unknown, or an alias), `IsRequired`, and `CheckRequired` (returns
    `ErrReqFlag` for the first missing flag by name, call after parsing).
  - **Parse errors**: `Parse` wraps failures in `*ParseError` with causes
    `ErrUndefinedFlag`, `ErrNeedsValue`, `ErrBadSyntax` or the value error.
  - **Typed accessors**: `GetBool/GetInt/GetInt64/GetUint/GetUint64/GetString/
    GetFloat64/GetDuration` (return zero value for unknown/mismatched-type
    flags) and `SetBool/SetString` (error on unknown/mismatched flag; mark the
    flag as set). `Getter`/`Valuer` expose the underlying value.
  - **`WasSet(name)`** — whether a flag was actually set.
  - Overridden `Visit`/`VisitAll` that are **alias-aware** (see below).
  - The `*SL` constructors registering long/short pairs (see below).

- **`pkg/xflag/help.go`** — the `HelpOptions`/`HelpOptionLines` methods render
  help that collapses each alias pair onto one line.

## The alias mechanism (central concept)

Long/short flag pairs (e.g. `--name` / `-n`) are implemented by registering
**two** stdlib flags that share the same value pointer. The `FlagSet` records
the link in its unexported `aliasOf` map (short name → long name); the usage
strings of both flags stay untouched.

The stdlib does not know about that link, which is why `FlagSet.Visit`/
`VisitAll` are overridden: `VisitAll` skips aliases entirely, and `Visit`
redirects an alias back to its canonical long flag so callers never see
duplicates. **Any new code that iterates flags must respect this or it will
double-count aliased flags.** The stdlib default usage output does not, so it
lists both names.

Use the `*SL` methods (`BoolSL`, `IntSL`, `Int64SL`, `Uint64SL`, `StringSL`,
`Float64SL`, `DurationSL`, `FuncSL`) on `*FlagSet` to register a long/short
pair.

## Conventions

- **Tests** use `github.com/ctx42/testing` (`pkg/assert`, `pkg/must`), not the
  standard library `testing` assertions. Follow the existing
  `// --- Given ---` / `// --- When ---` / `// --- Then ---` block structure and
  table/subtest style already in `*_test.go`.
- **Line width** is 80 columns (`.editorconfig`); multi-arg function signatures
  are wrapped accordingly.
