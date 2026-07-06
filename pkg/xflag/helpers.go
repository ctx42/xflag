package xflag

import (
	"bytes"
	"flag"
	"fmt"
	"sort"
	"text/tabwriter"
	"time"
)

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
		if long := fs.aliasOf[name]; long != "" {
			row := rows[long]
			row[0] = long
			row[1] = name
			rows[long] = row
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

// BoolSL adds a bool flag with long and short names to the flag set and returns
// the pointer that stores its value. The long and short names share the pointer.
func (fs *FlagSet) BoolSL(long, short string, value bool, usage string) *bool {
	val := fs.Bool(long, value, usage)
	fs.BoolVar(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// IntSL adds an int flag with long and short names to the flag set and returns
// the pointer that stores its value. The long and short names share the pointer.
func (fs *FlagSet) IntSL(long, short string, value int, usage string) *int {
	val := fs.Int(long, value, usage)
	fs.IntVar(val, short, value, usage)
	fs.recordAlias(short, long)
	return val
}

// Int64SL adds a 64bit int flag with long and short names to the flag set and
// returns the pointer that stores its value. The long and short names share the
// pointer.
func (fs *FlagSet) Int64SL(long, short string, value int64, usage string) *int64 {
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
