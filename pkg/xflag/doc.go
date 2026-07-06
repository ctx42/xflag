// Package xflag extends the standard library [flag] package with required
// flags, typed accessors, and long/short flag aliases.
//
// [FlagSet] embeds a *[flag.FlagSet] and adds required-flag tracking
// ([FlagSet.Required], [FlagSet.CheckRequired]), typed getters and setters
// ([FlagSet.GetBool], [FlagSet.GetInt], [FlagSet.SetString], ...), and
// alias-aware iteration.
//
// Long/short pairs (e.g. --name / -n) are registered as two standard flags
// sharing one value; the alias link is tracked inside the [FlagSet]. Use the
// *SL constructor methods ([FlagSet.BoolSL], [FlagSet.StringSL], ...) to declare
// a pair and [HelpOptions] to render alias-collapsed help.
package xflag
