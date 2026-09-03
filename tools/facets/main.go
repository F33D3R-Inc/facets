// facets — the authoring CLI for the Facet standard library (this repo).
//
// It does NOT compile, run, or serve `.fct` apps — fct's `facet build/dev/run`
// already own that, and by design an isolated stdlib atom (which references host
// actions/state) is not independently compilable. `facets` instead operates on
// the library's declared *surface*:
//
//	facets list    [dir]          enumerate the reusable facets, grouped by category
//	facets inspect <file.fct>     show one facet's kind/params/imports/members
//	facets check   [dir|file]     lint authoring conventions (doc, socket, imports)
//	facets new     <path/Name>    scaffold one new library facet from the conventions
//
// Add --json to list/inspect/check for machine-readable output.
package main

import (
	"errors"
	"fmt"
	"os"
)

var version = "0.1.0"

// errCheckFailed makes `facets check` exit non-zero without printing a duplicate
// error line (the findings were already printed).
var errCheckFailed = errors.New("check found issues")

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "list", "ls":
		err = cmdList(os.Stdout, os.Args[2:])
	case "inspect", "show":
		err = cmdInspect(os.Stdout, os.Args[2:])
	case "check", "lint":
		err = cmdCheck(os.Stdout, os.Args[2:])
	case "new":
		err = cmdNew(os.Stdout, os.Args[2:])
	case "version", "-v", "--version":
		fmt.Printf("facets %s\n", version)
		return
	case "help", "-h", "--help":
		usage()
		return
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		if errors.Is(err, errCheckFailed) {
			os.Exit(1) // findings already printed
		}
		fmt.Fprintln(os.Stderr, "facets:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "facets — authoring CLI for the Facet standard library")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  facets list    [dir]        enumerate reusable facets, grouped by category")
	fmt.Fprintln(os.Stderr, "  facets inspect <file.fct>   show one facet's kind/params/imports/members")
	fmt.Fprintln(os.Stderr, "  facets check   [dir|file]   lint authoring conventions (doc, socket, imports)")
	fmt.Fprintln(os.Stderr, "  facets new     <path/Name>  scaffold one new library facet")
	fmt.Fprintln(os.Stderr, "  facets version              print the tool version")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "flags: --json (list/inspect/check) ·")
	fmt.Fprintln(os.Stderr, "       new: --kind component|ui|data|wireframe|playground|app --socket <s> --params \"a: int, b: text\"")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "NOTE: compiling/running .fct apps is fct's job — use `facet build|dev|run`.")
}
