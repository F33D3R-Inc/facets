package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// cmdList enumerates every reusable facet in the library, grouped by category.
func cmdList(w io.Writer, args []string) error {
	jsonOut, rest := popFlag(args, "--json")
	dir := "."
	if len(rest) > 0 {
		dir = rest[0]
	}
	lib, err := OpenLibrary(dir)
	if err != nil {
		return err
	}
	files, err := lib.FctFiles()
	if err != nil {
		return err
	}

	type row struct {
		Category string  `json:"category"`
		File     string  `json:"file"`
		Facet    Facet   `json:"facet"`
	}
	var rows []row
	for _, rel := range files {
		f, err := ParseFile(filepath.Join(lib.Root, rel))
		if err != nil {
			return err
		}
		for _, u := range f.Units() {
			rows = append(rows, row{Category: Category(rel), File: rel, Facet: u})
		}
	}

	if jsonOut {
		out := map[string]any{
			"name":    lib.Manifest.Name,
			"version": lib.Manifest.Version,
			"facets":  rows,
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(out)
	}

	name := lib.Manifest.Name
	if name == "" {
		name = lib.Root
	}
	fmt.Fprintf(w, "%s", name)
	if lib.Manifest.Version != "" {
		fmt.Fprintf(w, " v%s", lib.Manifest.Version)
	}
	fmt.Fprintf(w, "\n%d file(s) · %d facet(s)\n\n", len(files), len(rows))

	// Group rows by category, preserving directory order.
	var cats []string
	byCat := map[string][]row{}
	for _, r := range rows {
		if _, ok := byCat[r.Category]; !ok {
			cats = append(cats, r.Category)
		}
		byCat[r.Category] = append(byCat[r.Category], r)
	}
	sort.Strings(cats)
	for _, c := range cats {
		label := c
		if label == "" {
			label = "(root)"
		}
		fmt.Fprintf(w, "%s/\n", label)
		for _, r := range byCat[c] {
			fmt.Fprintf(w, "  %-16s %-10s %-40s %s\n",
				r.Facet.Name, r.Facet.Kind, facetSurface(r.Facet), r.File)
		}
	}
	return nil
}

// facetSurface renders a one-line summary of a facet's declared interface.
func facetSurface(f Facet) string {
	switch f.Kind {
	case KindComponent:
		return "(" + paramList(f.Params) + ")"
	case KindUI, KindData:
		if f.InSock != "" {
			return "→ " + f.InSock
		}
	}
	return ""
}

func paramList(ps []Param) string {
	var parts []string
	for _, p := range ps {
		if p.Type != "" {
			parts = append(parts, p.Name+": "+p.Type)
		} else {
			parts = append(parts, p.Name)
		}
	}
	return strings.Join(parts, ", ")
}

// cmdInspect prints the full structural surface of one `.fct` file.
func cmdInspect(w io.Writer, args []string) error {
	jsonOut, rest := popFlag(args, "--json")
	if len(rest) == 0 {
		return fmt.Errorf("usage: facets inspect <file.fct> [--json]")
	}
	f, err := ParseFile(rest[0])
	if err != nil {
		return err
	}
	if jsonOut {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(f)
	}

	fmt.Fprintf(w, "%s\n", f.Path)
	if f.Doc != "" {
		fmt.Fprintf(w, "  doc       %s\n", firstLine(f.Doc))
	}
	if len(f.Imports) > 0 {
		fmt.Fprintf(w, "  imports   %s\n", strings.Join(f.Imports, ", "))
	}
	fmt.Fprintf(w, "  facets:\n")
	for _, u := range f.Units() {
		switch u.Kind {
		case KindComponent:
			fmt.Fprintf(w, "    component %s(%s)\n", u.Name, paramList(u.Params))
		case KindUI, KindData:
			in := ""
			if u.InSock != "" {
				in = " in " + u.InSock
			}
			fmt.Fprintf(w, "    %s %s%s\n", u.Kind, u.Name, in)
		default:
			fmt.Fprintf(w, "    %s %s\n", u.Kind, u.Name)
		}
	}
	printLabelled(w, "enums", f.Enums)
	printLabelled(w, "entities", f.Entities)
	printLabelled(w, "state", f.States)
	printLabelled(w, "actions", f.Actions)
	printLabelled(w, "policies", f.Policies)
	if len(f.Sockets) > 0 {
		var ss []string
		for _, s := range f.Sockets {
			ss = append(ss, s.Name+":"+s.Kind)
		}
		printLabelled(w, "sockets", ss)
	}
	printLabelled(w, "slots", f.Slots)
	printLabelled(w, "mounts", f.Mounts)
	printLabelled(w, "views", f.Views)
	return nil
}

func printLabelled(w io.Writer, label string, items []string) {
	if len(items) > 0 {
		fmt.Fprintf(w, "  %-9s %s\n", label, strings.Join(items, ", "))
	}
}

// Finding is one lint problem found by check.
type Finding struct {
	File string `json:"file"`
	Rule string `json:"rule"`
	Msg  string `json:"msg"`
}

// cmdCheck lints library facet files against authoring conventions that fct's
// compiler does NOT (and cannot, in isolation) enforce: a doc comment (the
// README quality bar), exactly one declared unit surface, ui/data facets naming a
// socket, and relative imports resolving to real files inside the library. It is
// complementary to `facet build`, never a replacement — deep semantic/placement
// checking remains fct's job and requires a whole compilable app.
func cmdCheck(w io.Writer, args []string) error {
	jsonOut, rest := popFlag(args, "--json")
	target := "."
	if len(rest) > 0 {
		target = rest[0]
	}

	lib, err := OpenLibrary(target)
	if err != nil {
		return err
	}

	// Determine the set of files to check.
	var files []string
	if fi, err := os.Stat(target); err == nil && !fi.IsDir() {
		abs, _ := filepath.Abs(target)
		files = []string{abs}
	} else {
		rels, err := lib.FctFiles()
		if err != nil {
			return err
		}
		for _, r := range rels {
			files = append(files, filepath.Join(lib.Root, r))
		}
	}

	var findings []Finding
	for _, path := range files {
		rel, _ := filepath.Rel(lib.Root, path)
		rel = filepath.ToSlash(rel)
		f, err := ParseFile(path)
		if err != nil {
			findings = append(findings, Finding{rel, "read", err.Error()})
			continue
		}
		findings = append(findings, checkFile(lib.Root, path, rel, f)...)
	}

	if jsonOut {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(map[string]any{
			"checked":  len(files),
			"findings": findings,
		}); err != nil {
			return err
		}
	} else {
		for _, fd := range findings {
			fmt.Fprintf(w, "%s: [%s] %s\n", fd.File, fd.Rule, fd.Msg)
		}
		if len(findings) == 0 {
			fmt.Fprintf(w, "ok — %d file(s), no convention issues\n", len(files))
		} else {
			fmt.Fprintf(w, "\n%d issue(s) across %d file(s)\n", len(findings), len(files))
		}
	}
	if len(findings) > 0 {
		return errCheckFailed
	}
	return nil
}

func checkFile(root, path, rel string, f *File) []Finding {
	var out []Finding
	if f.Doc == "" {
		out = append(out, Finding{rel, "doc", "missing leading doc comment (quality bar)"})
	}
	if len(f.Units()) == 0 {
		out = append(out, Finding{rel, "empty", "declares no facet unit"})
	}
	for _, u := range f.Units() {
		if (u.Kind == KindUI || u.Kind == KindData) && u.InSock == "" {
			out = append(out, Finding{rel, "socket",
				fmt.Sprintf("%s %s must declare `in <socket>`", u.Kind, u.Name)})
		}
		for _, p := range u.Params {
			if p.Type == "" {
				out = append(out, Finding{rel, "param",
					fmt.Sprintf("%s param %q has no type", u.Name, p.Name)})
			}
		}
	}
	for _, imp := range f.Imports {
		resolved, err := resolveImport(root, path, imp)
		if err != nil {
			// Remote/escaping imports: report escapes; skip remote refs quietly.
			if strings.Contains(err.Error(), "escape") {
				out = append(out, Finding{rel, "import", fmt.Sprintf("%q escapes library root", imp)})
			}
			continue
		}
		if _, err := os.Stat(resolved); err != nil {
			out = append(out, Finding{rel, "import", fmt.Sprintf("%q does not resolve to a file", imp)})
		}
	}
	return out
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// popFlag removes a boolean flag from args, returning whether it was present.
func popFlag(args []string, flag string) (bool, []string) {
	var rest []string
	found := false
	for _, a := range args {
		if a == flag {
			found = true
			continue
		}
		rest = append(rest, a)
	}
	return found, rest
}
