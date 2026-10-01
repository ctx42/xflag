package xflag

import (
	"errors"
	"flag"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ErrReqFlag is returned if a required flag has not been set.
var ErrReqFlag = errors.New("flag is required")

// Compile-time check that ParseError implements the error interface.
var _ error = (*ParseError)(nil)

// Parse error causes reported through [ParseError.Err] and matchable with
// [errors.Is].
var (
	// ErrUndefinedFlag is the cause when an argument names an unknown flag.
	ErrUndefinedFlag = errors.New("flag provided but not defined")

	// ErrNeedsValue is the cause when a flag is given without its argument.
	ErrNeedsValue = errors.New("flag needs an argument")

	// ErrBadSyntax is the cause when an argument is not a well-formed flag,
	// such as "---x" or "-=x".
	ErrBadSyntax = errors.New("bad flag syntax")
)

// FlagSet represents program flags.
//
// The default usage message of the embedded flag set, printed on -h and on a
// parse failure, lists a long flag and its short alias as two flags. Set its
// Usage field to a function printing [HelpOptions] for alias-collapsed help.
type FlagSet struct {
	*flag.FlagSet                   // Embedded StdLib flag set.
	req           map[string]bool   // Required flags.
	aliasOf       map[string]string // Maps a short flag to its long flag.
}

// NewFlagSet returns a new instance of FlagSet. It has the same arguments as
// [flag.NewFlagSet].
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	return NewFlagSetFrom(flag.NewFlagSet(name, errorHandling))
}

// NewFlagSetFrom returns a new instance of FlagSet created from an existing
// *flag.FlagSet.
func NewFlagSetFrom(fs *flag.FlagSet) *FlagSet {
	return &FlagSet{
		FlagSet: fs,
		req:     make(map[string]bool),
		aliasOf: make(map[string]string),
	}
}

// recordAlias records that the short flag is an alias of the long flag. It is
// zero-value safe: the backing map is created on first use.
func (fs *FlagSet) recordAlias(short, long string) {
	if fs.aliasOf == nil {
		fs.aliasOf = make(map[string]string)
	}
	fs.aliasOf[short] = long
}

// ParseError is returned by [FlagSet.Parse] when parsing an argument fails. It
// exposes the offending flag and value so callers can react without matching
// error strings; retrieve it with [errors.As].
type ParseError struct {
	// Flag is the flag name that failed, without leading dashes; empty for
	// [ErrBadSyntax].
	Flag string

	// Value is the offending value, empty when the failure carries none (an
	// undefined flag, a flag missing its argument, or a bool flag given
	// without one). For [ErrBadSyntax] it is the whole malformed argument.
	Value string

	// Err is the underlying cause. For an unparsable value it is the flag
	// package's message; for the other cases it is [ErrUndefinedFlag],
	// [ErrNeedsValue] or [ErrBadSyntax].
	Err error

	// msg is the original flag package message, preserved so Error reads
	// exactly as the standard library's did.
	msg string
}

// Error returns the original flag package message; a manually built value
// without one falls back to a reconstruction.
func (pe *ParseError) Error() string {
	if pe.msg != "" {
		return pe.msg
	}
	if errors.Is(pe.Err, ErrBadSyntax) {
		return fmt.Sprintf("%v: %s", pe.Err, pe.Value)
	}
	if errors.Is(pe.Err, ErrUndefinedFlag) || errors.Is(pe.Err, ErrNeedsValue) {
		return fmt.Sprintf("%v: %q", pe.Err, pe.Flag)
	}
	const format = "invalid value %q for flag %q: %v"
	return fmt.Sprintf(format, pe.Value, pe.Flag, pe.Err)
}

// Unwrap returns the underlying cause, enabling [errors.Is] and [errors.As].
func (pe *ParseError) Unwrap() error { return pe.Err }

// Parse parses flag definitions from the argument list like
// [flag.FlagSet.Parse], but wraps a parse failure in a [ParseError] that names
// the offending flag and value. [flag.ErrHelp] is returned unwrapped, so
// errors.Is(err, flag.ErrHelp) keeps working. Wrapping happens only under
// [flag.ContinueOnError]; the ExitOnError and PanicOnError modes exit or panic
// inside the embedded parser before this method can wrap.
func (fs *FlagSet) Parse(arguments []string) error {
	err := fs.FlagSet.Parse(arguments)
	if err == nil || errors.Is(err, flag.ErrHelp) {
		return err
	}
	return wrapParseError(err)
}

// wrapParseError converts a standard library flag parse error into a
// [ParseError]. It returns err unchanged when the message is not a known form.
func wrapParseError(err error) error {
	msg := err.Error()

	// The bad-value message differs for bool flags ("invalid boolean value
	// ... for -name") from every other type ("invalid value ... for flag
	// -name").
	forms := [][2]string{
		{"invalid value ", " for flag -"},
		{"invalid boolean value ", " for -"},
	}
	for _, form := range forms {
		if pe := invalidValueError(msg, form[0], form[1]); pe != nil {
			return pe
		}
	}

	if name, ok := strings.CutPrefix(
		msg,
		"flag provided but not defined: -",
	); ok {
		return &ParseError{Flag: name, Err: ErrUndefinedFlag, msg: msg}
	}
	if name, ok := strings.CutPrefix(msg, "flag needs an argument: -"); ok {
		return &ParseError{Flag: name, Err: ErrNeedsValue, msg: msg}
	}
	if arg, ok := strings.CutPrefix(msg, "bad flag syntax: "); ok {
		return &ParseError{Value: arg, Err: ErrBadSyntax, msg: msg}
	}

	// A bool flag given without a value whose Set("true") fails, e.g. a
	// [flag.FlagSet.BoolFunc] returning an error.
	if rest, ok := strings.CutPrefix(msg, "invalid boolean flag "); ok {
		if name, cause, ok := strings.Cut(rest, ": "); ok {
			return &ParseError{Flag: name, Err: errors.New(cause), msg: msg}
		}
	}

	return err
}

// invalidValueError parses a flag bad-value message shaped as
// prefix + quoted-value + sep + name + ": " + cause. It returns nil when msg
// does not match that shape.
func invalidValueError(msg, prefix, sep string) *ParseError {
	rest, ok := strings.CutPrefix(msg, prefix)
	if !ok {
		return nil
	}
	quoted, err := strconv.QuotedPrefix(rest)
	if err != nil {
		return nil
	}
	value, err := strconv.Unquote(quoted)
	if err != nil {
		return nil
	}
	if rest, ok = strings.CutPrefix(rest[len(quoted):], sep); !ok {
		return nil
	}
	name, cause, ok := strings.Cut(rest, ": ")
	if !ok {
		return nil
	}
	return &ParseError{
		Flag:  name,
		Value: value,
		Err:   errors.New(cause),
		msg:   msg,
	}
}

// Required marks a flag name as required. It must be called before parsing and
// the flag must exist and not be an alias; mark the long flag instead. Panics
// on error.
func (fs *FlagSet) Required(name string) {
	if fs.Parsed() {
		panic("flags already parsed")
	}
	if fs.Lookup(name) == nil {
		panic(fmt.Sprintf("flag %#q does not exist", name))
	}
	if long := fs.aliasOf[name]; long != "" {
		panic(fmt.Sprintf("flag %#q is an alias for flag %#q", name, long))
	}
	if fs.req == nil {
		fs.req = make(map[string]bool)
	}
	fs.req[name] = true
}

// IsRequired returns true when flag name (not alias) is required.
func (fs *FlagSet) IsRequired(name string) bool {
	return fs.req[name]
}

// CheckRequired checks all required flags were set. Returns [ErrReqFlag]
// wrapping the lexicographically first required flag that has not been set.
func (fs *FlagSet) CheckRequired() error {
	if !fs.Parsed() {
		return errors.New("flags not yet parsed")
	}
	for _, name := range slices.Sorted(maps.Keys(fs.req)) {
		if fs.req[name] && !fs.WasSet(name) {
			return fmt.Errorf("%#q %w", name, ErrReqFlag)
		}
	}
	return nil
}

// VisitAll visits the flags in lexicographical order but skips flags which are
// aliases.
func (fs *FlagSet) VisitAll(fn func(*flag.Flag)) {
	fs.FlagSet.VisitAll(func(flg *flag.Flag) {
		if fs.aliasOf[flg.Name] == "" {
			fn(flg)
		}
	})
}

// Visit visits the flags that have been set, calling fn for each. An alias is
// resolved to its long flag, each flag is visited at most once, and the order
// is lexicographical by the long flag name.
func (fs *FlagSet) Visit(fn func(*flag.Flag)) {
	set := make(map[string]*flag.Flag)
	fs.FlagSet.Visit(func(flg *flag.Flag) {
		if long := fs.aliasOf[flg.Name]; long != "" {
			if flg = fs.Lookup(long); flg == nil {
				return
			}
		}
		set[flg.Name] = flg
	})
	for _, name := range slices.Sorted(maps.Keys(set)) {
		fn(set[name])
	}
}

// WasSet returns true if the flag name (not alias) was set.
func (fs *FlagSet) WasSet(name string) bool {
	var set bool
	fs.Visit(func(flg *flag.Flag) {
		if name == flg.Name {
			set = true
		}
	})
	return set
}

// GetBool returns the parsed (or default) value of the bool flag. Returns
// false for an unknown flag or when the flag is of a different type.
func (fs *FlagSet) GetBool(name string) bool {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(bool); ok {
				return val
			}
		}
	}
	return false
}

// SetBool sets the value of a boolean flag identified by `name` and marks it
// as set, like [flag.FlagSet.Set]. If the flag is not of type `bool` or does
// not exist, it returns an error.
func (fs *FlagSet) SetBool(name string, value bool) error {
	flg := fs.Lookup(name)
	if flg == nil {
		return fmt.Errorf("cannot set not existing flag %q", name)
	}
	get, ok := flg.Value.(flag.Getter)
	if !ok {
		return fmt.Errorf("flag %#q is not a bool", name)
	}
	if _, ok := get.Get().(bool); !ok {
		return fmt.Errorf("flag %#q is not a bool", name)
	}
	if err := fs.FlagSet.Set(name, strconv.FormatBool(value)); err != nil {
		return fmt.Errorf("flag %#q %w", name, err)
	}
	return nil
}

// GetInt returns the parsed (or default) value of the int flag. Returns
// zero for an unknown flag or when the flag is of a different type.
func (fs *FlagSet) GetInt(name string) int {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(int); ok {
				return val
			}
		}
	}
	return 0
}

// GetInt64 returns the parsed (or default) value of the int64 flag. Returns
// zero for an unknown flag or when the flag is of a different type.
func (fs *FlagSet) GetInt64(name string) int64 {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(int64); ok {
				return val
			}
		}
	}
	return 0
}

// GetUint returns the parsed (or default) value of the uint flag. Returns
// zero for an unknown flag or when the flag is of a different type.
func (fs *FlagSet) GetUint(name string) uint {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(uint); ok {
				return val
			}
		}
	}
	return 0
}

// GetUint64 returns the parsed (or default) value of the uint64 flag. Returns
// zero for an unknown flag or when the flag is of a different type.
func (fs *FlagSet) GetUint64(name string) uint64 {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(uint64); ok {
				return val
			}
		}
	}
	return 0
}

// GetString returns the parsed (or default) value of the string flag. Returns
// empty string for an unknown flag or when the flag is of a different type.
func (fs *FlagSet) GetString(name string) string {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(string); ok {
				return val
			}
		}
	}
	return ""
}

// SetString sets the value of a string flag identified by `name` and marks it
// as set, like [flag.FlagSet.Set]. If the flag is not of type `string` or does
// not exist, it returns an error.
func (fs *FlagSet) SetString(name, value string) error {
	flg := fs.Lookup(name)
	if flg == nil {
		return fmt.Errorf("cannot set not existing flag %q", name)
	}
	get, ok := flg.Value.(flag.Getter)
	if !ok {
		return fmt.Errorf("flag %#q is not a string", name)
	}
	if _, ok := get.Get().(string); !ok {
		return fmt.Errorf("flag %#q is not a string", name)
	}
	if err := fs.FlagSet.Set(name, value); err != nil {
		return fmt.Errorf("flag %#q %w", name, err)
	}
	return nil
}

// GetFloat64 returns the parsed (or default) value of the float64 flag.
// Returns zero for an unknown flag or when the flag is of a different type.
func (fs *FlagSet) GetFloat64(name string) float64 {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(float64); ok {
				return val
			}
		}
	}
	return 0
}

// GetDuration returns the parsed (or default) value of the [time.Duration]
// flag. Returns zero for an unknown flag or when the flag is of a different
// type.
func (fs *FlagSet) GetDuration(name string) time.Duration {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			if val, ok := get.Get().(time.Duration); ok {
				return val
			}
		}
	}
	return 0
}

// Getter returns the getter for a flag with a given name or nil if the flag
// does not exist or doesn't implement the [flag.Getter] interface.
func (fs *FlagSet) Getter(name string) flag.Getter {
	if flg := fs.Lookup(name); flg != nil {
		if get, ok := flg.Value.(flag.Getter); ok {
			return get
		}
	}
	return nil
}

// Valuer returns [flag.Value] for a flag with a given name or nil if the flag
// does not exist.
func (fs *FlagSet) Valuer(name string) flag.Value {
	if flg := fs.Lookup(name); flg != nil {
		return flg.Value
	}
	return nil
}

// BoolSL adds a bool flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share
// the pointer.
func (fs *FlagSet) BoolSL(long, short string, value bool, usage string) *bool {
	val := fs.Bool(long, value, usage)
	fs.BoolVar(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// IntSL adds an int flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share
// the pointer.
func (fs *FlagSet) IntSL(long, short string, value int, usage string) *int {
	val := fs.Int(long, value, usage)
	fs.IntVar(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// Int64SL adds a 64bit int flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share the
// pointer.
func (fs *FlagSet) Int64SL(
	long, short string,
	value int64,
	usage string,
) *int64 {

	val := fs.Int64(long, value, usage)
	fs.Int64Var(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// Uint64SL adds an unsigned 64-bit int flag with long and short names to the
// flag set and returns the pointer that stores its value. The long and short
// names share the pointer.
func (fs *FlagSet) Uint64SL(
	long, short string,
	value uint64,
	usage string,
) *uint64 {

	val := fs.Uint64(long, value, usage)
	fs.Uint64Var(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// StringSL adds a string flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share the
// pointer.
func (fs *FlagSet) StringSL(long, short, value, usage string) *string {
	val := fs.String(long, value, usage)
	fs.StringVar(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// Float64SL adds a 64bit float flag with long and short names to the flag set
// and returns the pointer that stores its value. The long and short names share
// the pointer.
func (fs *FlagSet) Float64SL(
	long, short string,
	value float64,
	usage string,
) *float64 {

	val := fs.Float64(long, value, usage)
	fs.Float64Var(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// DurationSL adds a duration flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share the
// pointer.
func (fs *FlagSet) DurationSL(
	long, short string,
	value time.Duration,
	usage string,
) *time.Duration {

	val := fs.Duration(long, value, usage)
	fs.DurationVar(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// FuncSL adds a "function" flag with long and short names to the flag set.
func (fs *FlagSet) FuncSL(long, short, usage string, fn func(string) error) {
	fs.Func(long, usage, fn)
	fs.Func(short, usage, fn)
	fs.recordAlias(short, long)
}
