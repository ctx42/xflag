# xflag

xflag is a small wrapper around Go's standard [flag](https://pkg.go.dev/flag)
package. It adds the pieces the standard library leaves out — required flags,
typed getters and setters, and long/short aliases — while keeping everything you
already know. Your flag set behaves exactly as before, so you can adopt xflag in
an existing program and reach for the extras only when you need them.

## Features

| Feature            | What it adds                                          |
|--------------------|-------------------------------------------------------|
| Drop-in `FlagSet`  | Embeds `*flag.FlagSet`; stdlib methods still work.    |
| Required flags     | Mark `Required`; catch unset flags after parsing.     |
| Typed accessors    | Get/set without casting: `GetInt`, `SetBool`, etc.    |
| Long/short aliases | `*SL` binds `--name`/`-n`, returns the value pointer. |
| Alias-aware help   | `HelpOptions` folds each alias onto one help line.    |
| `WasSet`           | Tell an explicitly-set flag from a default value.     |

## Install

```bash
go get github.com/ctx42/xflag
```

```go
import "github.com/ctx42/xflag/pkg/xflag"
```

## Usage

### Quickstart

Define flags exactly as with the standard library, parse, and read them back
with typed accessors — no `flag.Lookup(...).Value.(...)` casting:

<!-- gmdoceg:pkg/xflag/ExampleNewFlagSet -->
```go
fs := xflag.NewFlagSet("example", flag.ContinueOnError)
fs.String("name", "world", "the name to greet")
fs.Int("count", 1, "how many times")

_ = fs.Parse([]string{"-name", "gopher", "-count", "3"})

fmt.Printf("name=%s count=%d\n", fs.GetString("name"), fs.GetInt("count"))
// Output:
// name=gopher count=3
```

### Required flags

Mark a flag required before parsing; check them all afterwards. `CheckRequired`
returns `ErrReqFlag` (wrapped, so `errors.Is` works) for the first unset flag:

<!-- gmdoceg:pkg/xflag/ExampleFlagSet_CheckRequired -->
```go
fs := xflag.NewFlagSet("example", flag.ContinueOnError)
fs.String("token", "", "auth token")
fs.Required("token")

_ = fs.Parse(nil)

fmt.Println(fs.CheckRequired())
// Output:
// `token` flag is required
```

### Long/short aliases

The `*SL` constructors register a long and a short name backed by a single
value, so `-v` and `--verbose` are interchangeable. They return the pointer
backing both names, so the value can be read directly instead of by string key:

<!-- gmdoceg:pkg/xflag/ExampleFlagSet_BoolSL -->
```go
fs := xflag.NewFlagSet("example", flag.ContinueOnError)
// The *SL constructors return the pointer backing both names, so the
// value can be read directly instead of by string key.
verbose := fs.BoolSL("verbose", "v", false, "enable verbose output")

_ = fs.Parse([]string{"-v"})

fmt.Println(*verbose)
// Output:
// true
```

`HelpOptions` renders help with each short flag collapsed onto its long-flag
line, instead of the two separate entries stdlib would print:

<!-- gmdoceg:pkg/xflag/ExampleHelpOptions -->
```go
fs := xflag.NewFlagSet("example", flag.ContinueOnError)
fs.StringSL("name", "n", "", "the name to greet")
fs.BoolSL("verbose", "v", false, "enable verbose output")

fmt.Print(xflag.HelpOptions(fs))
// Output:
// -n, --name       the name to greet
//   -v, --verbose    enable verbose output
```

### Wrapping an existing flag set

Already have a `*flag.FlagSet`? Wrap it with `NewFlagSetFrom` to gain the xflag
extensions without redefining a thing:

<!-- gmdoceg:pkg/xflag/ExampleNewFlagSetFrom -->
```go
// Wrap a *flag.FlagSet your program already defined; every stdlib
// method keeps working and you gain the xflag extensions on top.
std := flag.NewFlagSet("app", flag.ContinueOnError)
std.String("addr", ":8080", "listen address")

fs := xflag.NewFlagSetFrom(std)
fs.Required("addr")

_ = fs.Parse([]string{"-addr", ":9000"})

fmt.Println(fs.GetString("addr"))
fmt.Println(fs.CheckRequired())
// Output:
// :9000
// <nil>
```

### Default vs. explicitly set

`WasSet` reports whether a flag actually appeared on the command line — useful
when a zero value and "the user asked for zero" must be told apart:

<!-- gmdoceg:pkg/xflag/ExampleFlagSet_WasSet -->
```go
fs := xflag.NewFlagSet("example", flag.ContinueOnError)
fs.Int("port", 8080, "server port")

_ = fs.Parse(nil) // Nothing passed on the command line.

fmt.Println(fs.GetInt("port")) // Default value is returned...
fmt.Println(fs.WasSet("port")) // ...but the flag was never set.
// Output:
// 8080
// false
```

## API reference

`FlagSet` embeds `*flag.FlagSet`, so the entire stdlib API stays available. On
top of it:

- **Required flags** — `Required`, `IsRequired`, `CheckRequired` (returns
  `ErrReqFlag`), and `WasSet` to tell a set flag from a default.
- **Typed getters** (zero value for an unknown or mismatched-type flag) —
  `GetBool`, `GetInt`, `GetInt64`, `GetUint`, `GetUint64`, `GetString`,
  `GetFloat64`, `GetDuration`.
- **Typed setters** (error on an unknown or mismatched flag) — `SetBool`,
  `SetString`; low-level access via `Getter` and `Valuer`.
- **Long/short aliases** — the `*SL` constructor methods `BoolSL`, `IntSL`,
  `Int64SL`, `Uint64SL`, `StringSL`, `Float64SL`, `DurationSL`, `FuncSL`,
  rendered by `HelpOptions` / `HelpOptionLines`. Each (except `FuncSL`) returns
  the pointer backing both names, mirroring stdlib `flag.Bool`. The
  package-level free-func forms remain as deprecated shims.

Full API docs: [pkg.go.dev/github.com/ctx42/xflag/pkg/xflag](https://pkg.go.dev/github.com/ctx42/xflag/pkg/xflag).

## License

Licensed under [MIT](LICENSE.md).
