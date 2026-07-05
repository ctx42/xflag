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
  - **Required flags**: `Required(name)` (panics if parsed already or flag
    unknown), `IsRequired`, and `CheckRequired` (returns `ErrReqFlag`, call
    after parsing).
  - **Typed accessors**: `GetBool/GetInt/GetInt64/GetUint/GetUint64/GetString/
    GetFloat64/GetDuration` (return zero value for unknown/mismatched-type
    flags) and `SetBool/SetString` (error on unknown/mismatched flag).
    `Getter`/`Valuer` expose the underlying value.
  - **`WasSet(name)`** — whether a flag was actually set on the command line.
  - Overridden `Visit`/`VisitAll` that are **alias-aware** (see below).

- **`pkg/xflag/helpers.go`** — the **alias mechanism** and help formatting.

## The alias mechanism (central concept)

Long/short flag pairs (e.g. `--name` / `-n`) are implemented by registering
**two** stdlib flags that share the same value pointer. The short flag is
marked as an alias by prefixing its *usage string* with the sentinel
`AliasFor` (`"~~alias~for~~:"`) followed by the long name. `IsAlias(usage)`
decodes this (returns the long name, or `""` if not an alias).

This sentinel-in-usage convention is why `FlagSet.Visit`/`VisitAll` are
overridden: `VisitAll` skips aliases entirely, and `Visit` redirects an alias
back to its canonical long flag so callers never see duplicates. **Any new code
that iterates flags must respect this or it will double-count aliased flags.**

Use the `*SL` constructors (`BoolSL`, `IntSL`, `Int64SL`, `Uint64SL`,
`StringSL`, `Float64SL`, `DurationSL`, `FuncSL`) to register a long/short pair;
they operate on a plain `*flag.FlagSet`. `HelpOptions`/`HelpOptionLines` render
tab-aligned help that collapses each alias pair onto one line.

## Conventions

- **Tests** use `github.com/ctx42/testing` (`pkg/assert`, `pkg/must`), not the
  standard library `testing` assertions. Follow the existing
  `// --- Given ---` / `// --- When ---` / `// --- Then ---` block structure and
  table/subtest style already in `*_test.go`.
- **Line width** is 80 columns (`.editorconfig`); multi-arg function signatures
  are wrapped accordingly.
