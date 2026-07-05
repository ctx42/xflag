package xflag_test

import (
	"flag"
	"fmt"

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

func ExampleBoolSL() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	xflag.BoolSL(fs.FlagSet, "verbose", "v", false, "enable verbose output")

	_ = fs.Parse([]string{"-v"})

	fmt.Println(fs.GetBool("verbose"))
	// Output:
	// true
}

func ExampleHelpOptions() {
	fs := xflag.NewFlagSet("example", flag.ContinueOnError)
	xflag.StringSL(fs.FlagSet, "name", "n", "", "the name to greet")
	xflag.BoolSL(fs.FlagSet, "verbose", "v", false, "enable verbose output")

	fmt.Print(xflag.HelpOptions(fs))
	// Output:
	// -n, --name       the name to greet
	//   -v, --verbose    enable verbose output
}
