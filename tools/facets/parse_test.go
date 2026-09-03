package main

import (
	"strings"
	"testing"
)

const componentSrc = `# PostCard — a whole post: identity, kind-aware body, engagement bar.
# Second doc line.
import "../ui/userchip.fct"
import "engagementbar.fct"
app FacetSocial:
    component PostCard(id: int, author: text, verified: bool):
        box:
            use UserChip(author)  # inline comment ignored
`

func TestParseComponent(t *testing.T) {
	f := Parse(componentSrc)
	if !strings.HasPrefix(f.Doc, "PostCard — a whole post") {
		t.Fatalf("doc = %q", f.Doc)
	}
	if !strings.Contains(f.Doc, "Second doc line.") {
		t.Fatalf("second doc line dropped: %q", f.Doc)
	}
	if len(f.Imports) != 2 || f.Imports[0] != "../ui/userchip.fct" {
		t.Fatalf("imports = %v", f.Imports)
	}
	units := f.Units()
	if len(units) != 1 {
		t.Fatalf("want 1 unit (wrapper collapsed), got %d: %+v", len(units), units)
	}
	u := units[0]
	if u.Kind != KindComponent || u.Name != "PostCard" || u.Wrapper != "FacetSocial" {
		t.Fatalf("unit = %+v", u)
	}
	if len(u.Params) != 3 || u.Params[0] != (Param{"id", "int"}) || u.Params[2] != (Param{"verified", "bool"}) {
		t.Fatalf("params = %+v", u.Params)
	}
}

func TestParseUIWithSocket(t *testing.T) {
	f := Parse("# Nav\nui Nav in nav:\n    content:\n        text \"x\"\n")
	units := f.Units()
	if len(units) != 1 || units[0].Kind != KindUI || units[0].InSock != "nav" {
		t.Fatalf("units = %+v", units)
	}
}

func TestParseWireframe(t *testing.T) {
	src := "# Shell\nwireframe Shell:\n    socket nav: ui\n    socket feed: data\n    frame:\n        slot nav\n        slot feed\n"
	f := Parse(src)
	if len(f.Units()) != 1 || f.Units()[0].Kind != KindWireframe {
		t.Fatalf("units = %+v", f.Units())
	}
	if len(f.Sockets) != 2 || f.Sockets[0].Name != "nav" || f.Sockets[0].Kind != "ui" {
		t.Fatalf("sockets = %+v", f.Sockets)
	}
	if len(f.Slots) != 2 {
		t.Fatalf("slots = %v", f.Slots)
	}
}

func TestParseDataFacetMembers(t *testing.T) {
	src := `# Feed
data Feed in feed:
    enum Kind: text, video
    entity Tweet:
        id: int
    entity Like:
        id: int
    state draft: text = "" @client
    policy member:
        actor != "guest"
    action post(body: text):
        add Tweet { }
`
	f := Parse(src)
	if len(f.Units()) != 1 || f.Units()[0].InSock != "feed" {
		t.Fatalf("unit = %+v", f.Units())
	}
	if strings.Join(f.Entities, ",") != "Tweet,Like" {
		t.Fatalf("entities = %v", f.Entities)
	}
	if strings.Join(f.Enums, ",") != "Kind" || strings.Join(f.States, ",") != "draft" ||
		strings.Join(f.Actions, ",") != "post" || strings.Join(f.Policies, ",") != "member" {
		t.Fatalf("members: enums=%v states=%v actions=%v policies=%v", f.Enums, f.States, f.Actions, f.Policies)
	}
}

func TestParsePlaygroundAndApp(t *testing.T) {
	pg := Parse("# pg\nplayground F33D3R:\n    mount Shell at \"/\"\n")
	if len(pg.Units()) != 1 || pg.Units()[0].Kind != KindPlayground || len(pg.Mounts) != 1 {
		t.Fatalf("playground = %+v mounts=%v", pg.Units(), pg.Mounts)
	}
	app := Parse("# app\napp Home:\n    view Main at \"/\":\n        text \"x\"\n")
	if len(app.Units()) != 1 || app.Units()[0].Kind != KindApp {
		t.Fatalf("app should remain a unit (no components): %+v", app.Units())
	}
	if len(app.Views) != 1 || app.Views[0] != "Main" {
		t.Fatalf("views = %v", app.Views)
	}
}

func TestStripComment(t *testing.T) {
	cases := map[string]string{
		`    text "x"  # trailing`:  `    text "x"`,
		`    text "a # b"`:          `    text "a # b"`, // # inside string kept
		`# whole line`:              ``,
		`    use Foo()`:             `    use Foo()`,
	}
	for in, want := range cases {
		if got := stripComment(in); got != want {
			t.Errorf("stripComment(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseParams(t *testing.T) {
	if p := parseParams(""); p != nil {
		t.Fatalf("empty should be nil, got %+v", p)
	}
	p := parseParams("a: int, b:text ,  c : bool")
	if len(p) != 3 || p[1] != (Param{"b", "text"}) || p[2] != (Param{"c", "bool"}) {
		t.Fatalf("params = %+v", p)
	}
	// Untyped param.
	u := parseParams("bare")
	if len(u) != 1 || u[0].Type != "" {
		t.Fatalf("untyped = %+v", u)
	}
}
