package xflag

import (
	"flag"
	"testing"
	"time"

	"github.com/ctx42/testing/pkg/assert"
	"github.com/ctx42/testing/pkg/must"
)

func Test_NewFlagSet(t *testing.T) {
	// --- When ---
	have := NewFlagSet("flag-set", flag.ContinueOnError)

	// --- Then ---
	assert.Equal(t, "flag-set", have.Name())
	assert.NotNil(t, have.req)
}

func Test_NewFlagSetFrom(t *testing.T) {
	// --- Given ---
	ffs := flag.NewFlagSet("flag-set", flag.ExitOnError)

	// --- When ---
	have := NewFlagSetFrom(ffs)

	// --- Then ---
	assert.Same(t, ffs, have.FlagSet)
	assert.Equal(t, "flag-set", have.Name())
	assert.NotNil(t, have.req)
}

func Test_FlagSet(t *testing.T) {
	t.Run("alias then long overrides", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		flgName := fs.String("name", "default", "usage")
		fs.StringVar(flgName, "n", "default", "usage")

		// --- When ---
		err := fs.Parse([]string{"-n", "short", "--name", "long"})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "long", fs.GetString("name"))
	})

	t.Run("long then alias overrides", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		flgName := fs.String("name", "default", "usage")
		fs.StringVar(flgName, "n", "default", "usage")

		// --- When ---
		err := fs.Parse([]string{"--name", "long", "-n", "short"})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "short", fs.GetString("name"))
	})
}

func Test_FlagSet_Required_IsRequired(t *testing.T) {
	t.Run("set and get required", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")

		// --- When ---
		fs.Required("name")

		// --- Then ---
		assert.True(t, fs.IsRequired("name"))
	})

	t.Run("error - when called with a not existing flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")

		// --- When ---
		msg := assert.PanicMsg(t, func() { fs.Required("unknown") })

		// --- Then ---
		assert.Equal(t, "flag `unknown` does not exist", *msg)
	})

	t.Run("error - when called after parsing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")
		must.Nil(fs.Parse(nil))

		// --- When ---
		msg := assert.PanicMsg(t, func() { fs.Required("unknown") })

		// --- Then ---
		assert.Equal(t, "flags already parsed", *msg)
	})
}

func Test_FlagSet_IsRequired(t *testing.T) {
	t.Run("existing not required flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")

		// --- When ---
		have := fs.IsRequired("name")

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("existing required flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")
		fs.Required("name")

		// --- When ---
		have := fs.IsRequired("name")

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("not existing flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")

		// --- When ---
		have := fs.IsRequired("unknown")

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_FlagSet_CheckRequired(t *testing.T) {
	t.Run("required flag set", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name0", "abc", "usage0")
		fs.Required("name0")
		fs.String("name1", "def", "usage1")
		must.Nil(fs.Parse([]string{"--name0", "xyz"}))

		// --- When ---
		err := fs.CheckRequired()

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "xyz", fs.GetString("name0"))
	})

	t.Run("error - when a required flag is not set", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name0", "abc", "usage0")
		fs.Required("name0")
		fs.String("name1", "def", "usage1")
		must.Nil(fs.Parse(nil))

		// --- When ---
		err := fs.CheckRequired()

		// --- Then ---
		assert.ErrorIs(t, ErrReqFlag, err)
		assert.ErrorEqual(t, "`name0` flag is required", err)
		assert.Equal(t, "abc", fs.GetString("name0"))
	})

	t.Run("error - flags not yet parsed", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		err := fs.CheckRequired()

		// --- Then ---
		assert.ErrorEqual(t, "flags not yet parsed", err)
	})
}

func Test_FlagSet_Parse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name0", "abc", "usage0")
		fs.Required("name0")
		fs.String("name1", "def", "usage1")

		// --- When ---
		err := fs.Parse([]string{"--name0", "xyz"})

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "xyz", fs.GetString("name0"))
	})

	t.Run("error - parsing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Int("name0", 123, "usage0")

		// --- When ---
		err := fs.Parse([]string{"--name0", "abc"})

		// --- Then ---
		wMsg := `invalid value "abc" for flag -name0: parse error`
		assert.ErrorEqual(t, wMsg, err)
		assert.Equal(t, 0, fs.GetInt("name0"))
	})
}

func Test_FlagSet_VisitAll(t *testing.T) {
	t.Run("visiting set values", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		flgA := fs.String("flg-a", "flg-a-def", "flg-a help")
		fs.StringVar(flgA, "a", "flg-a-def", AliasFor+"flg-a")
		fs.String("flg-b", "flg-b-def", "flg-b help")
		fs.String("flg-c", "flg-c-def", "flg-c help")
		fs.String("rouge", "abc", "rouge usage")
		must.Nil(fs.Parse(nil))

		var have []string
		fn := func(flg *flag.Flag) { have = append(have, flg.Name) }

		// --- When ---
		fs.VisitAll(fn)

		// --- Then ---
		assert.Equal(t, []string{"flg-a", "flg-b", "flg-c", "rouge"}, have)
	})
}

func Test_FlagSet_Visit(t *testing.T) {
	t.Run("visiting set values", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		flgA := fs.String("flg-a", "flg-a-def", "flg-a help")
		fs.StringVar(flgA, "a", "flg-a-def", AliasFor+"flg-a")
		fs.String("flg-b", "flg-b-def", "flg-b help")
		fs.String("flg-c", "flg-c-def", "flg-c help")
		fs.String("rouge", "abc", "rouge usage")
		must.Nil(fs.Parse([]string{"--flg-b", "abc", "--rouge", "xyz"}))

		var have []string
		fn := func(flg *flag.Flag) { have = append(have, flg.Name) }

		// --- When ---
		fs.Visit(fn)

		// --- Then ---
		assert.Equal(t, []string{"flg-b", "rouge"}, have)
	})

	t.Run("alias is resolved", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		flgA := fs.String("flg-a", "flg-a-def", "flg-a help")
		fs.StringVar(flgA, "a", "flg-a-def", AliasFor+"flg-a")
		fs.String("flg-b", "flg-b-def", "flg-b help")
		fs.String("flg-c", "flg-c-def", "flg-c help")
		fs.String("rouge", "abc", "rouge usage")
		must.Nil(fs.Parse([]string{"-a", "abc", "--rouge", "xyz"}))

		var have []string
		fn := func(flg *flag.Flag) { have = append(have, flg.Name) }

		// --- When ---
		fs.Visit(fn)

		// --- Then ---
		assert.Equal(t, []string{"flg-a", "rouge"}, have)
	})

	t.Run("alias and long both set visits canonical once", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		flgA := fs.String("flg-a", "flg-a-def", "flg-a help")
		fs.StringVar(flgA, "a", "flg-a-def", AliasFor+"flg-a")
		must.Nil(fs.Parse([]string{"-a", "abc", "--flg-a", "xyz"}))

		var have []string
		fn := func(flg *flag.Flag) { have = append(have, flg.Name) }

		// --- When ---
		fs.Visit(fn)

		// --- Then ---
		assert.Equal(t, []string{"flg-a"}, have)
	})

	t.Run("alias to missing long flag is skipped", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("a", "a-def", AliasFor+"missing")
		fs.String("flg-b", "flg-b-def", "flg-b help")
		must.Nil(fs.Parse([]string{"-a", "abc", "--flg-b", "xyz"}))

		var have []string
		fn := func(flg *flag.Flag) { have = append(have, flg.Name) }

		// --- When ---
		fs.Visit(fn)

		// --- Then ---
		assert.Equal(t, []string{"flg-b"}, have)
	})
}

func Test_FlagSet_WasSet_tabular(t *testing.T) {
	// --- Given ---
	fs := NewFlagSet("flag-set", flag.ContinueOnError)
	flgA := fs.String("flg-a", "flg-a-def", "flg-a help")
	fs.StringVar(flgA, "a", "flg-a-def", AliasFor+"flg-a")
	flgB := fs.String("flg-b", "flg-b-def", "flg-b help")
	fs.StringVar(flgB, "b", "flg-b-def", AliasFor+"flg-b")
	flgC := fs.String("flg-c", "flg-c-def", "flg-c help")
	fs.StringVar(flgC, "c", "flg-c-def", AliasFor+"flg-c")

	fs.String("rouge", "abc", "rouge usage")
	assert.NoError(t, fs.Parse([]string{"-a", "abc", "--flg-b", "def"}))

	tt := []struct {
		testN string

		want bool
	}{
		{"flg-a", true},
		{"flg-b", true},
		{"flg-c", false},
		{"a", false},
	}

	for _, tc := range tt {
		t.Run(tc.testN, func(t *testing.T) {
			// --- When ---
			have := fs.WasSet(tc.testN)

			// --- Then ---
			assert.Equal(t, tc.want, have)
		})
	}
}

func Test_FlagSet_GetBool(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Bool("name", true, "usage")

		// --- When ---
		have := fs.GetBool("name")

		// --- Then ---
		assert.True(t, have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetBool("name")

		// --- Then ---
		assert.False(t, have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetBool("name")

		// --- Then ---
		assert.False(t, have)
	})
}

func Test_FlagSet_SetBool(t *testing.T) {
	t.Run("set existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Bool("name", false, "usage")

		// --- When ---
		err := fs.SetBool("name", true)

		// --- Then ---
		assert.NoError(t, err)
		assert.True(t, fs.GetBool("name"))
		assert.Equal(t, "true", fs.Lookup("name").Value.String())
	})

	t.Run("error - set not existing flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		err := fs.SetBool("name", true)

		// --- Then ---
		assert.ErrorEqual(t, `cannot set not existing flag "name"`, err)
	})

	t.Run("error - set existing of a different type", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Int("name", 42, "usage")

		// --- When ---
		err := fs.SetBool("name", true)

		// --- Then ---
		assert.ErrorEqual(t, "flag `name` is not a bool", err)
		assert.Equal(t, 42, fs.GetInt("name"))
		assert.Equal(t, "42", fs.Lookup("name").Value.String())
	})

	t.Run("error - set string flag as bool", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")

		// --- When ---
		err := fs.SetBool("name", true)

		// --- Then ---
		assert.ErrorEqual(t, "flag `name` is not a bool", err)
		assert.Equal(t, "abc", fs.GetString("name"))
	})
}

func Test_FlagSet_GetInt(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Int("name", 123, "usage")

		// --- When ---
		have := fs.GetInt("name")

		// --- Then ---
		assert.Equal(t, 123, have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetInt("name")

		// --- Then ---
		assert.Equal(t, 0, have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetInt("name")

		// --- Then ---
		assert.Equal(t, 0, have)
	})
}

func Test_FlagSet_GetInt64(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Int64("name", 123, "usage")

		// --- When ---
		have := fs.GetInt64("name")

		// --- Then ---
		assert.Equal(t, int64(123), have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetInt64("name")

		// --- Then ---
		assert.Equal(t, int64(0), have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetInt64("name")

		// --- Then ---
		assert.Equal(t, int64(0), have)
	})
}

func Test_FlagSet_GetUint(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Uint("name", 123, "usage")

		// --- When ---
		have := fs.GetUint("name")

		// --- Then ---
		assert.Equal(t, uint(123), have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetUint("name")

		// --- Then ---
		assert.Equal(t, uint(0), have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetUint("name")

		// --- Then ---
		assert.Equal(t, uint(0), have)
	})
}

func Test_FlagSet_GetUint64(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Uint64("name", 123, "usage")

		// --- When ---
		have := fs.GetUint64("name")

		// --- Then ---
		assert.Equal(t, uint64(123), have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetUint64("name")

		// --- Then ---
		assert.Equal(t, uint64(0), have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetUint64("name")

		// --- Then ---
		assert.Equal(t, uint64(0), have)
	})
}

func Test_FlagSet_GetString(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetString("name")

		// --- Then ---
		assert.Equal(t, "default", have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetString("name")

		// --- Then ---
		assert.Equal(t, "", have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Int("name", 123, "usage")

		// --- When ---
		have := fs.GetString("name")

		// --- Then ---
		assert.Equal(t, "", have)
	})
}

func Test_FlagSet_SetString(t *testing.T) {
	t.Run("set existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "abc", "usage")

		// --- When ---
		err := fs.SetString("name", "xyz")

		// --- Then ---
		assert.NoError(t, err)
		assert.Equal(t, "xyz", fs.GetString("name"))
		assert.Equal(t, "xyz", fs.Lookup("name").Value.String())
	})

	t.Run("set not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		err := fs.SetString("name", "abc")

		// --- Then ---
		assert.ErrorEqual(t, `cannot set not existing flag "name"`, err)
	})

	t.Run("error - set existing of a different type", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Int("name", 42, "usage")

		// --- When ---
		err := fs.SetString("name", "abc")

		// --- Then ---
		assert.ErrorEqual(t, "flag `name` is not a string", err)
		assert.Equal(t, 42, fs.GetInt("name"))
		assert.Equal(t, "42", fs.Lookup("name").Value.String())
	})

	t.Run("error - int flag with parseable value", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Int("name", 42, "usage")

		// --- When ---
		err := fs.SetString("name", "123")

		// --- Then ---
		assert.ErrorEqual(t, "flag `name` is not a string", err)
		assert.Equal(t, 42, fs.GetInt("name"))
	})
}

func Test_FlagSet_GetFloat64(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Float64("name", 1.23, "usage")

		// --- When ---
		have := fs.GetFloat64("name")

		// --- Then ---
		assert.Equal(t, 1.23, have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetFloat64("name")

		// --- Then ---
		assert.Equal(t, 0.0, have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetFloat64("name")

		// --- Then ---
		assert.Equal(t, 0.0, have)
	})
}

func Test_FlagSet_GetDuration(t *testing.T) {
	t.Run("get existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.Duration("name", 100*time.Millisecond, "usage")

		// --- When ---
		have := fs.GetDuration("name")

		// --- Then ---
		assert.Equal(t, 100*time.Millisecond, have)
	})

	t.Run("get not existing", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)

		// --- When ---
		have := fs.GetDuration("name")

		// --- Then ---
		assert.Equal(t, time.Duration(0), have)
	})

	t.Run("get a different type than the flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("flag-set", flag.ContinueOnError)
		fs.String("name", "default", "usage")

		// --- When ---
		have := fs.GetDuration("name")

		// --- Then ---
		assert.Equal(t, time.Duration(0), have)
	})
}

func Test_FlagSet_Getter(t *testing.T) {
	t.Run("flag with getter", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("name", flag.ContinueOnError)
		fs.String("name0", "default0", "usage0")
		fs.String("name1", "default1", "usage1")

		// --- When ---
		have := fs.Getter("name0")

		// --- Then ---
		assert.NotNil(t, have)
		assert.Equal(t, "default0", have.Get())
	})

	t.Run("get a not existing flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("name", flag.ContinueOnError)
		fs.String("name0", "default0", "usage0")
		fs.String("name1", "default1", "usage1")

		// --- When ---
		have := fs.Getter("not-existing")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("get flag not implementing Getter", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("name", flag.ContinueOnError)
		fs.String("name0", "default0", "usage0")
		fs.Func("name1", "usage1", func(s string) error { return nil })

		// --- When ---
		have := fs.Getter("name1")

		// --- Then ---
		assert.Nil(t, have)
	})
}

func Test_FlagSet_Valuer(t *testing.T) {
	t.Run("flag with getter", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("name", flag.ContinueOnError)
		fs.String("name0", "default0", "usage0")
		fs.String("name1", "default1", "usage1")

		// --- When ---
		have := fs.Valuer("name0")

		// --- Then ---
		assert.NotNil(t, have)
		assert.Equal(t, "default0", have.String())
	})

	t.Run("get a not existing flag", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("name", flag.ContinueOnError)
		fs.String("name0", "default0", "usage0")
		fs.String("name1", "default1", "usage1")

		// --- When ---
		have := fs.Valuer("not-existing")

		// --- Then ---
		assert.Nil(t, have)
	})

	t.Run("function value field", func(t *testing.T) {
		// --- Given ---
		fs := NewFlagSet("name", flag.ContinueOnError)
		fs.String("name0", "default0", "usage0")
		var value string

		fn := func(s string) error { value = s; return nil }
		fs.Func("name1", "usage1", fn)

		// --- When ---
		have := fs.Valuer("name1")

		// --- Then ---
		assert.NotNil(t, have)
		assert.NoError(t, have.Set("abc"))
		assert.Equal(t, "abc", value)
	})
}
