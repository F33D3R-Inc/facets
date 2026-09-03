package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// newOpts are the parsed flags for `facets new`.
type newOpts struct {
	kind   string
	socket string
	params string
}

// cmdNew scaffolds ONE new library facet file following repo conventions (a doc
// comment + a single declaration). This is deliberately narrower than fct's
// `facet new`, which scaffolds a whole runnable APP PROJECT (app.fct + Dockerfile
// + compose + seed/test). This writes a single stdlib atom into the library.
func cmdNew(w io.Writer, args []string) error {
	opts := newOpts{kind: KindComponent}
	var positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--kind" && i+1 < len(args):
			opts.kind = args[i+1]
			i++
		case strings.HasPrefix(a, "--kind="):
			opts.kind = strings.TrimPrefix(a, "--kind=")
		case a == "--socket" && i+1 < len(args):
			opts.socket = args[i+1]
			i++
		case strings.HasPrefix(a, "--socket="):
			opts.socket = strings.TrimPrefix(a, "--socket=")
		case a == "--params" && i+1 < len(args):
			opts.params = args[i+1]
			i++
		case strings.HasPrefix(a, "--params="):
			opts.params = strings.TrimPrefix(a, "--params=")
		default:
			positional = append(positional, a)
		}
	}
	if len(positional) == 0 {
		return fmt.Errorf("usage: facets new <path/Name> [--kind component|ui|data|wireframe|playground|app] [--socket s] [--params \"a: int\"]")
	}
	spec := positional[0]

	lib, err := OpenLibrary(".")
	if err != nil {
		return err
	}

	rel, name, err := splitSpec(spec)
	if err != nil {
		return err
	}

	// Destination path, guarded against traversal outside the library root.
	dest, err := safeJoin(lib.Root, rel)
	if err != nil {
		return err
	}
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("%s already exists — refusing to overwrite", rel)
	}

	body, err := render(opts, name, Category(rel))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(dest, []byte(body), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(w, "created %s\n", rel)
	fmt.Fprintf(w, "  facets inspect %s\n", rel)
	fmt.Fprintf(w, "  facets check %s\n", rel)
	return nil
}

// splitSpec turns "social/ReplyButton" (or "social/replybutton.fct") into a
// repo-relative file path and a PascalCase facet name.
func splitSpec(spec string) (rel, name string, err error) {
	spec = filepath.ToSlash(spec)
	dir, base := path_split(spec)
	base = strings.TrimSuffix(base, ".fct")
	if base == "" {
		return "", "", fmt.Errorf("no facet name in %q", spec)
	}
	name = pascal(base)
	// Filename follows repo convention: lowercase, no separators (e.g. FollowButton
	// → followbutton.fct). Deriving from the PascalCase name strips hyphens/underscores.
	file := strings.ToLower(name) + ".fct"
	rel = filepath.ToSlash(filepath.Join(dir, file))
	return rel, name, nil
}

func path_split(p string) (dir, base string) {
	if i := strings.LastIndexByte(p, '/'); i >= 0 {
		return p[:i], p[i+1:]
	}
	return "", p
}

// render produces the file body for the requested kind.
func render(o newOpts, name, category string) (string, error) {
	params := parseParams(o.params)
	switch o.kind {
	case KindComponent:
		wrapper := "Facet" + pascal(category)
		if category == "" {
			wrapper = "FacetLib"
		}
		return fmt.Sprintf(
			"# %s — TODO: one-line description of what this component renders.\n"+
				"app %s:\n"+
				"    component %s(%s):\n"+
				"        box:\n"+
				"            text \"%s\"\n",
			name, wrapper, name, paramList(params), name), nil
	case KindUI, KindData:
		sock := o.socket
		if sock == "" {
			return "", fmt.Errorf("--kind %s requires --socket <name>", o.kind)
		}
		return fmt.Sprintf(
			"# %s — TODO: one-line description. A `%s` facet for the `%s` socket.\n"+
				"%s %s in %s:\n"+
				"    content:\n"+
				"        text \"%s\"\n",
			name, o.kind, sock, o.kind, name, sock, name), nil
	case KindWireframe:
		return fmt.Sprintf(
			"# %s — TODO: describe this wireframe skeleton and its sockets.\n"+
				"wireframe %s:\n"+
				"    socket main: ui\n"+
				"    frame:\n"+
				"        box:\n"+
				"            slot main\n",
			name, name), nil
	case KindPlayground:
		return fmt.Sprintf(
			"# %s — TODO: the baseplate. Global concerns + a wireframe mount.\n"+
				"playground %s:\n"+
				"    auth\n"+
				"    mount Main at \"/\"\n",
			name, name), nil
	case KindApp:
		return fmt.Sprintf(
			"# %s — TODO: describe this app.\n"+
				"app %s:\n"+
				"    view Home at \"/\":\n"+
				"        box:\n"+
				"            text \"%s\"\n",
			name, name, name), nil
	default:
		return "", fmt.Errorf("unknown kind %q (want component|ui|data|wireframe|playground|app)", o.kind)
	}
}

// pascal converts "reply-button" / "reply_button" / "replybutton" to "Replybutton"
// / "ReplyButton" — capitalising each hyphen/underscore-separated segment and the
// first letter.
func pascal(s string) string {
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == '-' || r == '_' || r == ' ' })
	if len(fields) == 0 {
		return ""
	}
	var b strings.Builder
	for _, f := range fields {
		if f == "" {
			continue
		}
		b.WriteString(strings.ToUpper(f[:1]))
		b.WriteString(f[1:])
	}
	return b.String()
}

// safeJoin joins rel onto root and rejects any result outside root (traversal).
func safeJoin(root, rel string) (string, error) {
	joined := filepath.Clean(filepath.Join(root, rel))
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	r, err := filepath.Rel(absRoot, absJoined)
	if err != nil || r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("destination %q escapes the library root", rel)
	}
	return absJoined, nil
}
