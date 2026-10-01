## v0.11.0 (Thu, 01 Oct 2026 19:33:22 UTC)
- fix: wrap bad syntax and bool func failures in ParseError.
- fix: visit set flags in long flag name order.
- fix: mark flags set by SetBool and SetString as set.
- fix: report the first missing required flag deterministically.
- refactor: move *SL methods to flag_set.go and tidy style.
- style: tidy *SL test naming, Then grouping and godoc.
- feat!: make HelpOptions and HelpOptionLines methods on *FlagSet.
- fix: split parse errors at a defined flag name containing ": ".
- feat: add ErrNotParsed sentinel for CheckRequired before parsing.
- docs: document ErrBadSyntax, ErrNotParsed and setters marking set.
- test: assert aliasOf is initialized by the FlagSet constructors.
- docs: fix README example markers and add badges and prerequisites.

## v0.10.0 (Mon, 06 Jul 2026 19:17:34 UTC)
- feat: return the value pointer from *SL constructors.
- feat!: make *SL constructors methods on *FlagSet.
- feat!: track flag aliases in the FlagSet, not the usage string.
- feat: wrap parse failures in a typed *ParseError.
- test: simplify test case names.

## v0.9.0 (Mon, 06 Jul 2026 14:18:50 UTC)
- fix: reject wrong-type flags in FlagSet.SetBool.
- fix: reject wrong-type flags in FlagSet.SetString.
- fix: reject alias names in FlagSet.Required.
- test: add direct coverage for HelpOptionLines.
- fix: make FlagSet zero value safe in Required.
- test: use AliasFor constant instead of its literal value.

## v0.8.0 (Mon, 06 Jul 2026 11:19:44 UTC)
- Initial commit.
- feat: add xflag package extending the standard flag package.

