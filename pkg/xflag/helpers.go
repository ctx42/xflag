package xflag

import (
	"bytes"
	"flag"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

// AliasFor is a flag usage prefix indicating the flag is an alias to another
// flag.
//
// Example:
//
//	const argName = "name"
//	fp.fs.StringVar(&fp.fls.ImgName, argName, "", "docker image name")
//	fp.fs.StringVar(&fp.fls.ImgName, "n", "", xflag.AliasFor+argName)
const AliasFor = "~~alias~for~~:"

// HelpOptions returns formatted help with list of options.
func HelpOptions(fs *FlagSet) string {
	buf := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buf, 0, 8, 4, ' ', 0)
	for _, lin := range HelpOptionLines(fs) {
		_, _ = tw.Write([]byte(lin))
	}
	_ = tw.Flush()
	return buf.String()
}

// HelpOptionLines returns the help lines backing [HelpOptions], one per flag in
// lexicographical order with each alias collapsed onto its long-flag line.
func HelpOptionLines(fs *FlagSet) []string {
	var buf []string
	var names []string
	// The row array holds: 0 - name, 1 - alias, 2 - usage.
	rows := make(map[string][3]string)

	fs.FlagSet.VisitAll(func(flg *flag.Flag) {
		name := flg.Name
		if alias := IsAlias(flg.Usage); alias != "" {
			row := rows[alias]
			row[0] = alias
			row[1] = name
			rows[alias] = row
			return
		}

		names = append(names, name)
		row := rows[name]
		row[0] = name
		row[2] = flg.Usage
		rows[name] = row
	})

	sort.Strings(names)
	for _, name := range names {
		row := rows[name]
		alias := "    "
		if row[1] != "" {
			alias = "-" + row[1] + ", "
		}
		line := fmt.Sprintf("  %s--%s\t%s\n", alias, name, row[2])
		buf = append(buf, line)
	}
	return buf
}

// IsAlias returns the flag name which usage is for.
func IsAlias(usage string) string {
	if alias, ok := strings.CutPrefix(usage, AliasFor); ok && alias != "" {
		return alias
	}
	return ""
}

// BoolSL adds a bool flag with long and short names to the flag set and returns
// the pointer that stores its value. The long and short names share the pointer.
func BoolSL(fs *flag.FlagSet, long, short string, value bool, usage string) *bool {
	val := fs.Bool(long, value, usage)
	fs.BoolVar(val, short, value, AliasFor+long)
	return val
}

// IntSL adds an int flag with long and short names to the flag set and returns
// the pointer that stores its value. The long and short names share the pointer.
func IntSL(fs *flag.FlagSet, long, short string, value int, usage string) *int {
	val := fs.Int(long, value, usage)
	fs.IntVar(val, short, value, AliasFor+long)
	return val
}

// Int64SL adds a 64bit int flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share the
// pointer.
func Int64SL(fs *flag.FlagSet, long, short string, value int64, usage string) *int64 {
	val := fs.Int64(long, value, usage)
	fs.Int64Var(val, short, value, AliasFor+long)
	return val
}

// Uint64SL adds an unsigned 64-bit int flag with long and short names to the
// flag set and returns the pointer that stores its value. The long and short
// names share the pointer.
func Uint64SL(
	fs *flag.FlagSet,
	long, short string,
	value uint64,
	usage string,
) *uint64 {

	val := fs.Uint64(long, value, usage)
	fs.Uint64Var(val, short, value, AliasFor+long)
	return val
}

// StringSL adds a string flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share the
// pointer.
func StringSL(fs *flag.FlagSet, long, short, value, usage string) *string {
	val := fs.String(long, value, usage)
	fs.StringVar(val, short, value, AliasFor+long)
	return val
}

// Float64SL adds a 64bit float flag with long and short names to the flag set
// and returns the pointer that stores its value. The long and short names share
// the pointer.
func Float64SL(
	fs *flag.FlagSet,
	long, short string,
	value float64,
	usage string,
) *float64 {

	val := fs.Float64(long, value, usage)
	fs.Float64Var(val, short, value, AliasFor+long)
	return val
}

// DurationSL adds a duration flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share the
// pointer.
func DurationSL(
	fs *flag.FlagSet,
	long, short string,
	value time.Duration,
	usage string,
) *time.Duration {

	val := fs.Duration(long, value, usage)
	fs.DurationVar(val, short, value, AliasFor+long)
	return val
}

// FuncSL adds a "function" flag with long and short names to the flag set.
func FuncSL(
	fs *flag.FlagSet,
	long, short, usage string,
	fn func(string) error,
) {

	fs.Func(long, usage, fn)
	fs.Func(short, AliasFor+long, fn)
}
