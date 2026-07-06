package xflag

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
)

// ErrReqFlag is returned if a required flag has not been set.
var ErrReqFlag = errors.New("flag is required")

// FlagSet represents program flags.
type FlagSet struct {
	*flag.FlagSet                 // Embedded StdLib flag set.
	req           map[string]bool // Required flags.
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
	}
}

// Required marks a flag name as required. It must be called before parsing and
// flag must exist. Panics on error.
func (fs *FlagSet) Required(name string) {
	if fs.Parsed() {
		panic("flags already parsed")
	}
	if flg := fs.Lookup(name); flg != nil {
		fs.req[name] = true
		return
	}
	panic(fmt.Sprintf("flag %#q does not exist", name))
}

// IsRequired returns true when flag name (not alias) is required.
func (fs *FlagSet) IsRequired(name string) bool {
	return fs.req[name]
}

// CheckRequired checks all required flags were set. Returns [ErrReqFlag]
// wrapping one required flag that has not been set; which flag is unspecified
// when several are missing.
func (fs *FlagSet) CheckRequired() error {
	if !fs.Parsed() {
		return errors.New("flags not yet parsed")
	}
	for name, required := range fs.req {
		if required && !fs.WasSet(name) {
			return fmt.Errorf("%#q %w", name, ErrReqFlag)
		}
	}
	return nil
}

// VisitAll visits the flags in lexicographical order but skips flags which are
// aliases.
func (fs *FlagSet) VisitAll(fn func(*flag.Flag)) {
	fs.FlagSet.VisitAll(func(flg *flag.Flag) {
		if IsAlias(flg.Usage) == "" {
			fn(flg)
		}
	})
}

// Visit visits the flags in lexicographical order, calling fn for each. It
// visits only those flags that have been set and skips flags which are aliases,
// resolving each alias to its long flag and visiting that flag at most once.
func (fs *FlagSet) Visit(fn func(*flag.Flag)) {
	seen := make(map[string]bool)
	fs.FlagSet.Visit(func(flg *flag.Flag) {
		if name := IsAlias(flg.Usage); name != "" {
			if flg = fs.Lookup(name); flg == nil {
				return
			}
		}
		if seen[flg.Name] {
			return
		}
		seen[flg.Name] = true
		fn(flg)
	})
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

// SetBool sets the value of a boolean flag identified by `name`. If the flag
// is not of type `bool` or does not exist, it returns an error.
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
	if err := flg.Value.Set(strconv.FormatBool(value)); err != nil {
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

// SetString sets the value of a string flag identified by `name`. If the flag
// is not of type `string` or does not exist, it returns an error.
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
	if err := flg.Value.Set(value); err != nil {
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

// GetDuration returns parsed (or default) value of the [time.Duration] flag.
// Returns zero value for an unknown flag or when the flag is of a different
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

// Getter returns getter for a flag with a given name or nil if a flag does not
// exist or doesn't implement [flag.Getter] interface.
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
