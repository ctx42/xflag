package xflag

import (
	"flag"
	"testing"

	"github.com/ctx42/testing/pkg/assert"
)

func Test_FlagSet_HelpOptions(t *testing.T) {
	t.Run("with aliases", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("test", flag.ContinueOnError)
		fs.BoolSL("mkdir", "d", false, "mkdir help")
		fs.String("fast", "fast", "fast help")
		fs.StringSL("name", "n", "project", "name help")

		// --- When ---
		have := fs.HelpOptions()

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
		have := fs.HelpOptions()

		// --- Then ---
		assert.Empty(t, have)
	})
}

func Test_FlagSet_HelpOptionLines(t *testing.T) {
	t.Run("with aliases", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("test", flag.ContinueOnError)
		fs.BoolSL("mkdir", "d", false, "mkdir help")
		fs.String("fast", "fast", "fast help")
		fs.StringSL("name", "n", "project", "name help")

		// --- When ---
		have := fs.HelpOptionLines()

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
		have := fs.HelpOptionLines()

		// --- Then ---
		assert.Empty(t, have)
	})
}
