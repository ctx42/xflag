package xflag_test

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/ctx42/xflag/pkg/xflag"
)

func ExampleNewFlagSet() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	fs.String("name", "world", "the name to greet")
	fs.Int("count", 1, "how many times")

	_ = fs.Parse([]string{"-name", "gopher", "-count", "3"})

	fmt.Printf("name=%s count=%d\n", fs.GetString("name"), fs.GetInt("count"))
	// Output:
	// name=gopher count=3
}

func ExampleNewFlagSetFrom() {
	// Wrap a *flag.FlagSet your program already defined; every stdlib
	// method keeps working and you gain the xflag extensions on top.
	std := flag.NewFlagSet("app", flag.ContinueOnError)
	std.String("addr", ":8080", "listen address")

	fs := xflag.NewFlagSetFrom(std)
	fs.Required("addr")

	_ = fs.Parse([]string{"-addr", ":9000"})

	fmt.Println(fs.GetString("addr"))
	fmt.Println(fs.CheckRequired())
	// Output:
	// :9000
	// <nil>
}

func ExampleFlagSet_WasSet() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	fs.Int("port", 8080, "server port")

	_ = fs.Parse(nil) // Nothing passed on the command line.

	fmt.Println(fs.GetInt("port")) // Default value is returned...
	fmt.Println(fs.WasSet("port")) // ...but the flag was never set.
	// Output:
	// 8080
	// false
}

func ExampleFlagSet_CheckRequired() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	fs.String("token", "", "auth token")
	fs.Required("token")

	_ = fs.Parse(nil)

	fmt.Println(fs.CheckRequired())
	// Output:
	// `token` flag is required
}

func ExampleFlagSet_BoolSL() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	// The *SL constructors return the pointer backing both names, so the
	// value can be read directly instead of by string key.
	verbose := fs.BoolSL("verbose", "v", false, "enable verbose output")

	_ = fs.Parse([]string{"-v"})

	fmt.Println(*verbose)
	// Output:
	// true
}

func ExampleFlagSet_Parse() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // Silence the standard usage output.
	fs.Int("timeout", 0, "seconds to wait")

	err := fs.Parse([]string{"-timeout", "soon"})

	// errors.As recovers the flag and value without matching error strings.
	if pe, ok := errors.AsType[*xflag.ParseError](err); ok {
		fmt.Printf("flag %q rejected value %q\n", pe.Flag, pe.Value)
	}
	// Output:
	// flag "timeout" rejected value "soon"
}

func ExampleHelpOptions() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	fs.StringSL("name", "n", "", "the name to greet")
	fs.BoolSL("verbose", "v", false, "enable verbose output")

	fmt.Print(xflag.HelpOptions(fs))
	// Output:
	// -n, --name       the name to greet
	//   -v, --verbose    enable verbose output
}
