// Structural reader for `.fct` library facets.
//
// This is deliberately NOT a compiler. fct's `facet build/check/inspect` already
// compile a whole application graph and enforce semantics/placement — but they
// require a runnable app: an isolated library atom like social/postcard.fct
// references host actions (like/repost) and host @client state (q/draft) that
// only exist in the host, so `facet build` on it fails ("unknown action repost").
// The facets stdlib is a catalog of exactly those non-standalone units. This
// reader parses their *declared surface* (kind, name, params, imports, sockets,
// slots, and top-level members) so the library can be listed, inspected, and
// lint-checked for authoring conventions — without pretending to be the compiler.
package main

import (
	"os"
	"regexp"
	"strings"
)

// Kind classifies a top-level `.fct` declaration.
const (
	KindApp        = "app"        // `app Name:` — a full app OR a component-wrapper namespace
	KindUI         = "ui"         // `ui Name in socket:`
	KindData       = "data"       // `data Name in socket:`
	KindWireframe  = "wireframe"  // `wireframe Name:`
	KindPlayground = "playground" // `playground Name:`
	KindComponent  = "component"  // `component Name(params):` — nested inside an app wrapper
)

// Param is a single `name: type` declaration parameter.
type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Socket is a wireframe `socket name: kind` typed slot.
type Socket struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// Facet is one reusable unit the stdlib exposes: either a top-level ui/data/
// wireframe/playground/app, or a component declared inside an app wrapper.
type Facet struct {
	Kind    string  `json:"kind"`
	Name    string  `json:"name"`
	Params  []Param `json:"params,omitempty"`
	InSock  string  `json:"in,omitempty"`  // ui/data: the `in <socket>` target
	Wrapper string  `json:"wrapper,omitempty"` // component: the enclosing `app X` name
}

// File is the parsed structural surface of one `.fct` file.
type File struct {
	Path     string   `json:"path"`
	Doc      string   `json:"doc,omitempty"`      // leading `#` comment block, joined
	Imports  []string `json:"imports,omitempty"`
	Facets   []Facet  `json:"facets"`             // the reusable units this file exposes
	Entities []string `json:"entities,omitempty"`
	Enums    []string `json:"enums,omitempty"`
	States   []string `json:"states,omitempty"`
	Actions  []string `json:"actions,omitempty"`
	Policies []string `json:"policies,omitempty"`
	Sockets  []Socket `json:"sockets,omitempty"`
	Slots    []string `json:"slots,omitempty"`
	Mounts   []string `json:"mounts,omitempty"`
	Views    []string `json:"views,omitempty"`
}

var (
	reImport    = regexp.MustCompile(`^import\s+"([^"]+)"`)
	reTop       = regexp.MustCompile(`^(app|ui|data|wireframe|playground)\s+([A-Za-z_][\w]*)\s*(?:\bin\s+([A-Za-z_][\w]*)\s*)?:`)
	reComponent = regexp.MustCompile(`^component\s+([A-Za-z_][\w]*)\s*\(([^)]*)\)\s*:`)
	reEntity    = regexp.MustCompile(`^entity\s+([A-Za-z_][\w]*)\s*:`)
	reEnum      = regexp.MustCompile(`^enum\s+([A-Za-z_][\w]*)\s*:`)
	reState     = regexp.MustCompile(`^state\s+([A-Za-z_][\w]*)\s*:`)
	reAction    = regexp.MustCompile(`^action\s+([A-Za-z_][\w]*)\s*\(`)
	rePolicy    = regexp.MustCompile(`^policy\s+([A-Za-z_][\w]*)`)
	reSocket    = regexp.MustCompile(`^socket\s+([A-Za-z_][\w]*)\s*:\s*([A-Za-z_][\w]*)`)
	reSlot      = regexp.MustCompile(`^slot\s+([A-Za-z_][\w]*)`)
	reMount     = regexp.MustCompile(`^mount\s+([A-Za-z_][\w]*)\s+at`)
	reView      = regexp.MustCompile(`^view\s+([A-Za-z_][\w]*)\s+at`)
)

// ParseFile reads and structurally parses a `.fct` file at path.
func ParseFile(path string) (*File, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := Parse(string(src))
	f.Path = path
	return f, nil
}

// Parse extracts the structural surface of `.fct` source. It is line-oriented and
// tolerant: unknown constructs are ignored rather than erroring, because the point
// is cataloguing a declared surface, not full semantic validation (that is fct's).
func Parse(src string) *File {
	f := &File{}
	lines := strings.Split(src, "\n")

	// Leading `#` comment block → doc. Stops at the first code/blank boundary once
	// any comment has been seen, so a mid-file comment never becomes the doc.
	var doc []string
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "#") {
			doc = append(doc, strings.TrimSpace(strings.TrimPrefix(t, "#")))
			continue
		}
		break
	}
	f.Doc = strings.TrimSpace(strings.Join(doc, "\n"))

	// currentWrapper tracks the enclosing `app X` so nested components can name it.
	currentWrapper := ""
	for _, ln := range lines {
		code := stripComment(ln)
		t := strings.TrimSpace(code)
		if t == "" {
			continue
		}
		indented := code != "" && (code[0] == ' ' || code[0] == '\t')

		if m := reImport.FindStringSubmatch(t); m != nil {
			f.Imports = append(f.Imports, m[1])
			continue
		}
		// Top-level facet: only at column 0 (no leading whitespace).
		if !indented {
			if m := reTop.FindStringSubmatch(t); m != nil {
				kind, name, sock := m[1], m[2], m[3]
				if kind == KindApp {
					// Could be a full app OR a component wrapper — resolved after we
					// see whether it contains components. Record the wrapper name; a
					// top-level app facet is only emitted if no component appears.
					currentWrapper = name
					f.Facets = append(f.Facets, Facet{Kind: KindApp, Name: name})
					continue
				}
				currentWrapper = ""
				f.Facets = append(f.Facets, Facet{Kind: kind, Name: name, InSock: sock})
				continue
			}
		}
		if m := reComponent.FindStringSubmatch(t); m != nil {
			f.Facets = append(f.Facets, Facet{
				Kind:    KindComponent,
				Name:    m[1],
				Params:  parseParams(m[2]),
				Wrapper: currentWrapper,
			})
			continue
		}
		switch {
		case reEntity.MatchString(t):
			f.Entities = append(f.Entities, reEntity.FindStringSubmatch(t)[1])
		case reEnum.MatchString(t):
			f.Enums = append(f.Enums, reEnum.FindStringSubmatch(t)[1])
		case reState.MatchString(t):
			f.States = append(f.States, reState.FindStringSubmatch(t)[1])
		case reAction.MatchString(t):
			f.Actions = append(f.Actions, reAction.FindStringSubmatch(t)[1])
		case rePolicy.MatchString(t):
			f.Policies = append(f.Policies, rePolicy.FindStringSubmatch(t)[1])
		case reSocket.MatchString(t):
			m := reSocket.FindStringSubmatch(t)
			f.Sockets = append(f.Sockets, Socket{Name: m[1], Kind: m[2]})
		case reSlot.MatchString(t):
			f.Slots = append(f.Slots, reSlot.FindStringSubmatch(t)[1])
		case reMount.MatchString(t):
			f.Mounts = append(f.Mounts, reMount.FindStringSubmatch(t)[1])
		case reView.MatchString(t):
			f.Views = append(f.Views, reView.FindStringSubmatch(t)[1])
		}
	}

	// If an app wrapper actually contained components, it is a namespace, not a
	// reusable unit — drop the bare `app` facet so `list` shows the components.
	f.Facets = collapseWrappers(f.Facets)
	return f
}

// Units returns the facets to expose in a catalog: components when present under a
// wrapper, otherwise the top-level facet itself.
func (f *File) Units() []Facet { return f.Facets }

// collapseWrappers removes an `app X` entry when a component names X as its wrapper
// (the app is just a namespace for the components).
func collapseWrappers(fs []Facet) []Facet {
	wrapped := map[string]bool{}
	for _, x := range fs {
		if x.Kind == KindComponent && x.Wrapper != "" {
			wrapped[x.Wrapper] = true
		}
	}
	out := fs[:0]
	for _, x := range fs {
		if x.Kind == KindApp && wrapped[x.Name] {
			continue
		}
		out = append(out, x)
	}
	return out
}

// stripComment removes a trailing/inline `#` comment while preserving leading
// indentation (used to detect column-0 vs nested declarations). A `#` inside a
// double-quoted string is not a comment.
func stripComment(ln string) string {
	inStr := false
	for i := 0; i < len(ln); i++ {
		switch ln[i] {
		case '"':
			inStr = !inStr
		case '#':
			if !inStr {
				return strings.TrimRight(ln[:i], " \t")
			}
		}
	}
	return strings.TrimRight(ln, " \t")
}

// parseParams splits a `name: type, name: type` parameter list. Empty → nil.
func parseParams(s string) []Param {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	var out []Param
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		name, typ, ok := strings.Cut(part, ":")
		p := Param{Name: strings.TrimSpace(name)}
		if ok {
			p.Type = strings.TrimSpace(typ)
		}
		out = append(out, p)
	}
	return out
}
