package xflag

import (
	"bytes"
	"flag"
	"fmt"
	"sort"
	"text/tabwriter"
)

// HelpOptions returns formatted help with a list of options, collapsing each
// alias onto its long-flag line. The default usage message of the standard
// library does not; see [FlagSet] for using it on -h.
func (fs *FlagSet) HelpOptions() string {
	buf := &bytes.Buffer{}
	tw := tabwriter.NewWriter(buf, 0, 8, 4, ' ', 0)
	for _, lin := range fs.HelpOptionLines() {
		_, _ = tw.Write([]byte(lin))
	}
	_ = tw.Flush()
	return buf.String()
}

// HelpOptionLines returns the help lines backing [FlagSet.HelpOptions], one per
// flag in lexicographical order with each alias collapsed onto its long-flag
// line.
func (fs *FlagSet) HelpOptionLines() []string {
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
