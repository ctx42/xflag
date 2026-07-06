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
		fs.BoolSL("mkdir", "d", false, "mkdir help")
		fs.String("fast", "fast", "fast help")
		fs.StringSL("name", "n", "project", "name help")

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
		fs.BoolSL("mkdir", "d", false, "mkdir help")
		fs.String("fast", "fast", "fast help")
		fs.StringSL("name", "n", "project", "name help")

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

func Test_SL_tabular(t *testing.T) {
	fn := func(string) error { return nil }

	tt := []struct {
		testN string

		fnSL func(*FlagSet)
	}{
		{
			"BoolSL",
			func(fs *FlagSet) { fs.BoolSL("name", "n", true, "usage") },
		},
		{
			"IntSL",
			func(fs *FlagSet) { fs.IntSL("name", "n", 1, "usage") },
		},
		{
			"Int64SL",
			func(fs *FlagSet) { fs.Int64SL("name", "n", 1, "usage") },
		},
		{
			"Uint64SL",
			func(fs *FlagSet) { fs.Uint64SL("name", "n", 1, "usage") },
		},
		{
			"StringSL",
			func(fs *FlagSet) { fs.StringSL("name", "n", "a", "usage") },
		},
		{
			"Float64SL",
			func(fs *FlagSet) { fs.Float64SL("name", "n", 1, "usage") },
		},
		{
			"DurationSL",
			func(fs *FlagSet) { fs.DurationSL("name", "n", 1, "usage") },
		},
		{
			"FuncSL",
			func(fs *FlagSet) { fs.FuncSL("name", "n", "usage", fn) },
		},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- Given ---
			fs := NewFlagSet("flg-set", flag.ContinueOnError)

			// --- When ---
			tc.fnSL(fs)

			// --- Then ---
			var names []string
			var usages []string
			fn := func(flg *flag.Flag) {
				names = append(names, flg.Name)
				usages = append(usages, flg.Usage)
			}
			fs.FlagSet.VisitAll(fn)
			assert.Equal(t, []string{"n", "name"}, names)
			// The short flag carries the real usage now, not a sentinel;
			// the alias link lives in the FlagSet, not the usage string.
			assert.Equal(t, []string{"usage", "usage"}, usages)
		})
	}
}

func Test_SL_returnsBoundPointer(t *testing.T) {
	t.Run("pointer reflects the long flag after Parse", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flg-set", flag.ContinueOnError)
		got := fs.StringSL("name", "n", "default", "usage")

		// --- When ---
		err := fs.Parse([]string{"--name", "long"})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "long", *got)
	})

	t.Run("pointer reflects the short flag after Parse", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flg-set", flag.ContinueOnError)
		got := fs.BoolSL("verbose", "v", false, "usage")

		// --- When ---
		err := fs.Parse([]string{"-v"})

		// --- Then ---
		assert.NoError(t, err)
		assert.True(t, *got)
	})

	t.Run("long and short share the pointer", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flg-set", flag.ContinueOnError)
		got := fs.IntSL("num", "n", 0, "usage")

		// --- When ---
		err := fs.Parse([]string{"--num", "1", "-n", "2"})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, 2, *got)
	})

	t.Run("records the alias on zero value construction", func(t *testing.T) {
		// --- Given ---
		ffs := flag.NewFlagSet("flg-set", flag.ContinueOnError)
		fs := &FlagSet{FlagSet: ffs} // Nil aliasOf map.

		// --- When ---
		fs.BoolSL("verbose", "v", false, "usage")

		// --- Then --- (alias recorded, so VisitAll skips the short flag)
		var names []string
		fs.VisitAll(func(flg *flag.Flag) { names = append(names, flg.Name) })
		assert.Equal(t, []string{"verbose"}, names)
	})
}
