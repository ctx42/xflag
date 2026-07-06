package xflag

import (
	"flag"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_HelpOptions(t *testing.T) {
	t.Run("with aliases", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("test", flag.ContinueOnError)
		flgMkdir := fs.Bool("mkdir", false, "mkdir help")
		fs.BoolVar(flgMkdir, "d", false, AliasFor+"mkdir")
		fs.String("fast", "fast", "fast help")
		flgName := fs.String("name", "project", "name help")
		fs.StringVar(flgName, "n", "", AliasFor+"name")

		// --- When ---
		have := HelpOptions(fs)

		// --- Then ---
		want := "" +
			"      --fast     fast help\n" +
			"  -d, --mkdir    mkdir help\n" +
			"  -n, --name     name help\n"
		assert.Equal(t, want, have)
	})

	t.Run("no flags", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("test", flag.ContinueOnError)

		// --- When ---
		have := HelpOptions(fs)

		// --- Then ---
		assert.Empty(t, have)
	})
}

func Test_HelpOptionLines(t *testing.T) {
	t.Run("with aliases", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("test", flag.ContinueOnError)
		flgMkdir := fs.Bool("mkdir", false, "mkdir help")
		fs.BoolVar(flgMkdir, "d", false, AliasFor+"mkdir")
		fs.String("fast", "fast", "fast help")
		flgName := fs.String("name", "project", "name help")
		fs.StringVar(flgName, "n", "", AliasFor+"name")

		// --- When ---
		have := HelpOptionLines(fs)

		// --- Then ---
		want := []string{
			"      --fast\tfast help\n",
			"  -d, --mkdir\tmkdir help\n",
			"  -n, --name\tname help\n",
		}
		assert.Equal(t, want, have)
	})

	t.Run("no flags", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("test", flag.ContinueOnError)

		// --- When ---
		have := HelpOptionLines(fs)

		// --- Then ---
		assert.Empty(t, have)
	})
}

func Test_IsAlias_tabular(t *testing.T) {
	tt := []struct {
		testN string

		usage string
		want  string
	}{
		{"normal", "usage", ""},
		{"alias", AliasFor + "usage", "usage"},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := IsAlias(tc.usage)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_SL_tabular(t *testing.T) {
	fn := func(string) error { return nil }

	tt := []struct {
		testN string

		fnSL func(*flag.FlagSet)
	}{
		{
			"BoolSL",
			func(fs *flag.FlagSet) { BoolSL(fs, "name", "n", true, "usage") },
		},
		{
			"IntSL",
			func(fs *flag.FlagSet) { IntSL(fs, "name", "n", 1, "usage") },
		},
		{
			"Int64SL",
			func(fs *flag.FlagSet) { Int64SL(fs, "name", "n", 1, "usage") },
		},
		{
			"Uint64SL",
			func(fs *flag.FlagSet) { Uint64SL(fs, "name", "n", 1, "usage") },
		},
		{
			"StringSL",
			func(fs *flag.FlagSet) { StringSL(fs, "name", "n", "a", "usage") },
		},
		{
			"Float64SL",
			func(fs *flag.FlagSet) { Float64SL(fs, "name", "n", 1, "usage") },
		},
		{
			"DurationSL",
			func(fs *flag.FlagSet) { DurationSL(fs, "name", "n", 1, "usage") },
		},
		{
			"FuncSL",
			func(fs *flag.FlagSet) { FuncSL(fs, "name", "n", "usage", fn) },
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			fs := NewFlagSet("flg-set", flag.ContinueOnError)

			// --- When ---
			tc.fnSL(fs.FlagSet)

			// --- Then ---
			var names []string
			var usages []string
			fn := func(flg *flag.Flag) {
				names = append(names, flg.Name)
				usages = append(usages, flg.Usage)
			}
			fs.FlagSet.VisitAll(fn)
			assert.Equal(t, []string{"n", "name"}, names)
			assert.Equal(t, []string{AliasFor + "name", "usage"}, usages)
		})
	}
}

func Test_SL_returnsBoundPointer(t *testing.T) {
	t.Run("pointer reflects the long flag after Parse", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flg-set", flag.ContinueOnError)
		got := StringSL(fs.FlagSet, "name", "n", "default", "usage")

		// --- When ---
		err := fs.Parse([]string{"--name", "long"})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "long", *got)
	})

	t.Run("pointer reflects the short flag after Parse", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flg-set", flag.ContinueOnError)
		got := BoolSL(fs.FlagSet, "verbose", "v", false, "usage")

		// --- When ---
		err := fs.Parse([]string{"-v"})

		// --- Then ---
		assert.NoError(t, err)
		assert.True(t, *got)
	})

	t.Run("long and short share the pointer", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flg-set", flag.ContinueOnError)
		got := IntSL(fs.FlagSet, "num", "n", 0, "usage")

		// --- When ---
		err := fs.Parse([]string{"--num", "1", "-n", "2"})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, *got)
	})
}
