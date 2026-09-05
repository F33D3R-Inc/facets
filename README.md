# Facet standard library — `github.com/F33D3R-Inc/facets`

The reusable facets f33d3r.com (and any Facet app) is assembled from — imports, not
hand-written screens. Published via the registry, versioned independently of the language.

> **Toolchain pin.** v0.4.0 needs the compiler build that added the `route`
> builtin and the `textarea` / `checkbox` / `toggle` / `radio` control nodes, on
> top of everything v0.3.0 needed (reference parameters, component `slot`s,
> route-expression link destinations, interpolated `placeholder`/`upload`
> label/`icon` attributes). The language shipped that whole batch as `1.31.0`,
> so `facet.json` pins `>= 1.31.0` — the true lower bound, not a stand-in for it.

> **v0.1.0 — the f33d3r core batch.** Value-first, not count-first: the ~20 facets
> f33d3r's home + profile actually need, and a full exercise of the language
> (filtered `count`/`exists`, `tabs`, `match`, `richtext`/`video`, `pending`/`failed`,
> `contains`/search, `remove … where`, components, `image`/`badge`, `icon`).
>
> **v0.2.0 — the site batch.** The 17 facets a site (rather than a feed) is
> assembled from: the `layout/` category, the states a page spends most of its life
> in (loading, empty, failed, notice), and the first self-contained forms. Same
> rule as v0.1.0 — nothing that could not be written cleanly was written at all;
> what the language cannot yet express is named in the doc comment of the facet
> that wanted it.
>
> **v0.3.0 — the batch v0.2.0 could not write.** Three compiler changes landed
> together and each one unblocked a category that had been named as impossible in
> the doc comment of the facet that wanted it:
>
> * **Reference parameters** — `v: cell text`, `a: action`. A component can be
>   handed the NAME of a state cell or an action, not just a value, so `bind`,
>   `-> action(…)`, `pending(a)` and `failed(a)` are parameterisable. `forms/`
>   becomes a set of FIELDS instead of a set of whole forms.
> * **Component children** — `slot` inside a `component`, spliced from the `use`
>   line's indented block. `layout/`'s CSS vocabulary becomes real wrappers, and a
>   block handed to a component with no slot is now a compile error rather than
>   silently discarded.
> * **Parameterized link destinations** — `link "{label}" -> "{href}"`. The whole
>   `navigation/` category, which could not contain a single facet before.
>
> A fourth landed mid-batch: `placeholder`, `upload`'s label, the `form` submit
> label and the `icon` glyph interpolate now, so `TextField` carries a real
> placeholder and `SidebarItem` a real glyph.
>
> **v0.4.0 — catching up to the language, and three bugs found in real use.**
> Two compiler features landed and each one deleted a paragraph from the
> "impossible" lists below:
>
> * **`route`** — the path being rendered, readable in a view and in any component
>   a view renders. `NavLink` and `SidebarItem` MARK THEMSELVES now: no `current`
>   argument, no flag threaded through every page, and a nav written once in a
>   shared `layout` finally works, which is the only place a nav is ever written.
>   `NavLinkWhen`/`SidebarItemWhen` keep the explicit condition for the cases exact
>   match gets wrong (a section, a route parameter).
> * **`textarea`, `checkbox`, `toggle`, `radio`** — four two-way controls with
>   their cell types checked at the bind site. `TextAreaField`, `CheckboxField`,
>   `ToggleField` and `Toggle` exist; and because a `toggle` is a legal writer of a
>   `@client bool` and `overlay bind` reads one, so does **`Menu`** — the facet the
>   v0.3.0 README said could not exist at all.
>
> The batch also fixes three bugs that a real site (f33d3r.com, assembled entirely
> from this library) hit and had to work around at its own call sites. Each is
> written up where it was fixed: `SubmitButton` could not submit a form that
> carried data and is gone (`forms/form.fct`); `.x-field input` over-styled the new
> tick controls (`forms/fieldlabel.fct`); `Grid2`/`Grid3` counted columns at
> whatever width they were given, so `Grid2` was four columns across a page
> (`layout/vocabulary.fct`). Two smaller ones with it: `PageHeader` no longer
> forces `position: sticky` on every host (`StickyPageHeader` does, and stacks
> under another sticky bar via `--x-pagehead-top`), and `.x-panel-note` is a block,
> so a note that wraps keeps its inset.
>
> Promoted from that site, which wrote them locally and nominated them: `Prose` and
> `Code` (the `richtext` wrappers), `Specs`/`SpecRow` (a definition list) and
> `TextAreaField`. Its `Band`, `Rung`, `FeatureCard`, `Metric`, `SiteNav` and
> `SiteFooter` were left where they are — see "what was not promoted" below.

## Import

```
import "github.com/F33D3R-Inc/facets/social/postcard.fct"
```

Locally (this repo) the same files live under `library/` and build with `facet build`.

## What's in the library

`ui/`, `layout/`, `navigation/` and all of `forms/`'s field components build on
their own with `facet build` — parameters in, markup out. That now includes the
ones that bind and call: a component with a `cell`/`action` parameter is a
template, and the compiler checks it once against synthetic references, so
`forms/textfield.fct` type-checks its `input bind v` with no host at all. The
facets that need a host are the ones naming a host's cells and actions directly
(`social/`, `notify/`, `profile/`, `forms/searchbox.fct`); they compile inside the
app that owns them, which is what `home.fct` is for.

| Category | Facets |
|---|---|
| `ui/` | Avatar · VerifiedBadge · UserChip · Trend · Trends (ui) · Nav (icon rail, ui) · Spinner · SkeletonLine/SkeletonPost · EmptyState · ErrorState · Tag · Toast · Tooltip · **Menu/MenuItem/MenuSeparator/MenuClose** · **Disclosure** · **Prose/Code** · **Specs/SpecRow** |
| `layout/` | Layout vocabulary (page · stack · row · grid2/grid3 · split · split-3 · sticky · card · the 3-column app skeleton) · Page/PageNarrow · Stack/StackTight/StackLoose · Row/RowBetween/RowEnd · Grid2/Grid3 *(column count is a maximum now)* · Split/Split3/Sticky · Card/CardRaised · Section · SectionHeader · PageHeader *(+slot, no longer sticky)* · **StickyPageHeader** · Panel/PanelTitle/PanelNote · Divider/DividerInset · Spacer |
| `navigation/` | Link · **NavLink** *(self-marking)* · **NavLinkWhen** · **SidebarItem** *(self-marking)* · **SidebarItemWhen** · Breadcrumb/Crumb/CrumbCurrent · Pagination/PageNumbers/PageNumber · TabBar |
| `social/` | PostCard · EngagementBar · ComposeBox · FollowButton · WhoToFollow |
| `forms/` | Field *(slot)* · TextField/NumberField · **TextAreaField** · **CheckboxField/ToggleField/Toggle** · UploadField · Form · **FormActions/ActionButton** *(`SubmitButton` removed)* · SearchBox · FieldLabel · FieldHint/FieldError · NewsletterForm · ReportForm |
| `notify/` | UnreadBadge · NotificationItem |
| `profile/` | ProfileHeader |
| `data/` | Feed (a full vertical slice: entities + actions + policy + content) |
| `wireframes/` | Shell (the 3-column app skeleton) |
| — | f33d3r (playground baseplate) |

**Forms are fields again.** `bind` names a state cell that must be resolved at
compile time, and a component parameter was not one — so `component Field(label,
cell)` could not be written and the reusable unit had to be the whole form. A
`cell text` parameter passes the cell's NAME; the compiler substitutes it through
the body and lowers a copy per call site, so `input bind v` is `input bind draft`
before any check runs, and every check a hand-written input gets, the field gets:
that the cell exists, that its type matches, that it is `@client`. `TextField`,
`NumberField`, `TextAreaField`, `CheckboxField`, `ToggleField` and `UploadField`
are those fields; `Form(a: action)` wraps them in one action's `pending`/`failed`.
The long version is at the top of `forms/textfield.fct`.

**A form's submit control is the language's `form` node, not a component.** This
library shipped `SubmitButton(label, a: action)` and it did not work: it lowered
to `button "{label}" -> a`, an action reference invoked with NO arguments, and
every field in a form is a `@client` cell that a server-placed action cannot read
— so the values have to be passed as arguments. `use SubmitButton("Send",
request)` failed with *action "request" takes 6 argument(s), got 0*, not in an
unusual form but in the only kind there is. There is no varargs and no way to
forward a call site's expressions through an `action` reference, so no component
can ever be the submit control of a form that carries data. It is removed. Write

```
use Form(request):
    form "Send request" -> request(name, email, note):
        use TextField("Name", "Ada Lovelace", true, "", "", name)
        use TextAreaField("What for?", "A sentence or two.", true, "", "", note)
```

— the `form` node takes the arguments because it is at the call site, where the
arguments are, and `Form` adds the half it has no opinion about: `pending`,
`failed`, and the styling of the submit button the node renders for itself. What
was reusable survives as `FormActions()` (the strip, a slot) and `ActionButton`
(a button for an action that genuinely takes no arguments).

**What forms still cannot do.** A `select`'s and a `radio`'s options are a fixed
compile-time list — an `option`'s value is an identity the compiler checks, not a
display string — so a general `SelectField(label, options)` or
`RadioField(label, options)` is not expressible, and a `radio` bound to a plain
`text` cell with no options is refused at build time (*radio on "…" needs options
(or a `@client` enum cell to default them)*). The other half of that message is
the other half of the answer: a radio bound to an ENUM cell defaults its options
to the enum's members and needs none written — but a component parameter typed
`v: cell Layer` names a concrete enum, and a general-purpose library has no enums
to name, so the enum-bound radio group is a two-line facet in the app that
declares the enum and not a facet here. `Field` is the general answer, and always
was: it is the label/hint/error scaffold with a `slot` where the control goes, so
a `radio`, a `select`, a `tabs` or a `typeahead` gets the same field layout with
its own six lines of `option` written in place. `textarea`, `checkbox` and
`toggle` are no longer on this list — the nodes landed and the fields are real.

**One live compiler bug worth knowing before you use a `slot`.** Slot children
are evaluated in the CALLEE's scope, not the caller's: a name in a child resolves
against the wrapper component's parameters first and only falls through to the
caller's when the wrapper has no parameter of that name. So
`component Outer(b): use Wrap("x"): text "{b}"` renders the caller's `b` — until
`Wrap` grows a `b` of its own, at which point it silently renders Wrap's. It
compiles, it renders, nothing warns. It cost this batch a facet: `CheckboxField`
written as `use Field("", …): checkbox bind v label "{label}"` renders a checkbox
with no words at all, because the child's `{label}` binds to `Field`'s own first
parameter, which that call site deliberately passes as `""`. `forms/checkfield.fct`
assembles `.x-field` from `FieldLabel`/`FieldHint`/`FieldError` directly so the
case cannot arise. State cells, `actor` and `route` are unaffected — only locals
(component parameters, `for` variables) collide. Reported to `fct`.

**Layout primitives are wrappers now, and still rules.** A `component` could not
accept child nodes, so a container could only ship as a CSS class a call site
applied to a box of its own. `slot` closed that and `Page`/`Stack`/`Row`/`Grid`/
`Split`/`Card`/`Section` are real components. The vocabulary file stays, because
a `for` renders as one element of its own: a grid of query results is
`for p in Post … class "x-l-grid2"`, not a `for` inside a `use Grid2()` — which
compiles and makes the whole list one cell. One substrate, many wrappers; the
long version is at the top of `layout/vocabulary.fct`.

**Navigation destinations are checked in one of two ways.** A **path template**
starts with a literal `/` (`"/post/{p.id}"`) and is route-checked at compile
time. A **route expression** is one whole interpolation and nothing else
(`"{href}"`) and is checked by the renderers, which emit an `<a>` only if the
value names a route this app serves and inert text otherwise. Nothing in
between — `"{base}/edit"` is refused. Prefer the template wherever the call site
is writing the path; `Link`/`NavLink` earn their keep where the destination
arrives as data. Note that a component ARGUMENT is an expression, so `"{p.id}"`
there is a literal brace-p-dot-id-brace: build a destination with `+`
(`"/post/" + p.id`), which concatenates.

**A nav marks itself now.** `route` is the path being rendered, readable in a
view and in any component a view renders (and refused in a derive, a policy or an
action, which have no page in sight). `NavLink(label, href)` and
`SidebarItem(glyph, label, href)` compare `route` to the destination themselves —
no `current`/`active` argument, no flag threaded through every page, and a nav
written once inside a shared `layout` finally marks the right item, which is the
only place a nav is ever written. Exact match is the default because it is right
without being told anything; a prefix default would light "Home" on every page.
Where exact match is wrong — a "Docs" tab that should stay lit across
`/docs/getting-started`, a profile tab current when the route's handle is yours —
`NavLinkWhen`/`SidebarItemWhen` take the condition instead:
`use NavLinkWhen("Docs", "/docs", contains(route, "/docs"))`. Same markup, same
rules, so a nav can mix the two.

**A `Menu` exists, and it is a sheet.** Nothing can write a `@client` cell except
a control bound to it — that diagnosis was right and the missing piece was a
CONTROL, not a statement that assigns state. `toggle` and `overlay bind` on the
same cell are a menu: one cell, one writer, open and closed, nothing to get out
of step. What opens is a modal sheet and not an anchored dropdown, and the facet
is named for the more honest of the two words: `overlay` renders a fixed backdrop
with a centred panel, there is no anchor and no way to express one, because the
language has no measurement of where a node landed. A hover menu is still not
expressible and is not faked. `Disclosure` is the same control with the panel in
the flow instead of over it.

**What navigation still cannot do.** Page numbers are `PageNumber` atoms rather
than a loop, because `for` walks an entity or a `[T]` cell and there is no numeric
range. And nothing here can leave the site: `link "…" -> "https://…"` is refused
(a path must start with `/`) and a `Link` handed an absolute URL renders as inert
text, so a footer of repository links is a footer of monospaced coordinates until
the language grows external destinations.

**What was not promoted.** f33d3r.com wrote thirteen facets locally and nominated
them all. Four came in — `TextAreaField` (an exact mirror of `TextField`), `Prose`
and `Code` (`richtext` wrappers, where what is promoted is the measure and the
rhythm rather than the two lines), and `Specs`/`SpecRow` (a definition list every
docs page and every settings summary needs). The rest stayed where they are, and
the reason is the same each time: they carry the site's content or the site's type
scale rather than a shape. `Band` is a full-bleed strip whose value is that site's
background and vertical rhythm — `Page` already supplies the measure inside it.
`FeatureCard` and `Metric` are a `Card` and a `box` reading `x-h3`/`x-small`, a
type scale this library deliberately does not define, so promoting them would mean
promoting the scale. `Rung`/`Rungs` is the F33D3R stack diagram with a `match` arm
per layer name. `SiteNav`, `SiteFooter`, `NavItems`, `LayerTabs` and `FootCol` are
that site's destinations, arranged. The generalisable part of `SiteNav` — a bar
whose mobile drawer is a toggle over an overlay — is `Menu`, and that did come in.

## Two composition tracks, one set of atoms

The same component atoms serve both ways an app is built:

- **Plain-app track** — `library/home.fct` (`app F33D3RHome`) imports the atoms and
  assembles a home screen directly.
- **Layered (typed-brick) track** — `library/f33d3r.fct` (`playground`) → `Shell`
  (`wireframe`, typed sockets) → `Nav`/`Trends` (`ui`) + `Feed` (`data`). The `data`
  facet imports the **same** `PostCard`/`ComposeBox`/`SearchBox`/`WhoToFollow` atom
  files — a layered build can pull in component-only modules (closed in v1.17.0).

```
facet dev library/home.fct      # plain-app track
facet dev library/f33d3r.fct    # layered typed-brick track
```

## Quality bar (per facet)

- Compiles (`facet build`); runs where it makes sense (`facet dev`).
- A doc comment at the top of the file.
- Placement-sound by construction — the compiler enforces it across the import
  boundary (an imported server action can't read your `@client` state).
