#!/usr/bin/env python3
# Emits facets/api/types.fct: every golden components.schemas entry that the
# fct wire-type grammar can name, as a `type` declaration. Run from anywhere.
import json, os, re, sys

GOLDEN = "/home/hiiro/f33d3r/feed-engine/internal/handler/testdata/api_v2_contract.golden.json"
OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "..", "api", "types.fct")
IDENT = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*$")

d = json.load(open(GOLDEN))
S = d["components"]["schemas"]

# Schemas the fct runtime itself publishes (and serves) when the app declares
# `contract "/api/v2/contract"`: the contract document and its history/diff
# shapes, the stream connect frame, and the error envelope. The runtime owns
# their shape, so the app must not redeclare them.
RUNTIME_OWNED = {n for n in S if n.startswith("APIContract") or
                 (n.startswith("Contract") and n.endswith("DTO"))} | {
    "FacetCatalogBindingDTO", "HelloEventDTO", "APIErrorDTO", "apiErrorBody"}

GENERIC = re.compile(r"^(\w+)\[(?:[\w./-]+\.)?(\w+)\]$")

def fct_name(n):
    """The fct identifier for a golden schema name; the golden name, when it
    differs, is kept as the type's contract schema name (`type X as "n":`)."""
    m = GENERIC.match(n)
    if m:  # V2Page[github.com/…/handler.V2WorkDTO] -> V2PageOfV2WorkDTO
        return m.group(1) + "Of" + m.group(2)
    if IDENT.match(n) and not n[0].isupper():  # sealedKeyWire -> SealedKeyWire
        return n[0].upper() + n[1:]
    return n

def nameable(n):
    return bool(IDENT.match(fct_name(n)))

declared = {n for n, s in S.items() if n not in RUNTIME_OWNED and nameable(n) and s.get("properties")
            and all(IDENT.match(f) for f in s["properties"])}

def core(s):
    """fct wire type core for one JSON schema; None when it has no wire form."""
    if "$ref" in s:
        n = s["$ref"][len("#/components/schemas/"):]  # a generic's name has "/" in it
        return fct_name(n) if n in declared else "json"
    if s.get("oneOf"):
        parts = [x for x in s["oneOf"] if x.get("type") != "null"]
        return core(parts[0]) if len(parts) == 1 else "json"
    t = s.get("type")
    if isinstance(t, list):
        t = [x for x in t if x != "null"]
        t = t[0] if len(t) == 1 else None
    if t == "string":
        # An instant crosses as RFC 3339 text; `datetime` is fct's type for it
        # (published as format: date-time), filled with iso(<unix seconds>).
        return "datetime" if s.get("format") == "date-time" else "text"
    if t == "integer":
        return "int"
    if t == "boolean":
        return "bool"
    if t == "number":
        return "number"
    if t == "array":
        return None  # a nested list; the caller decides
    return "json"

def field(name, s, required):
    nullable = (isinstance(s.get("type"), list) and "null" in s["type"]) or \
               any(x.get("type") == "null" for x in s.get("oneOf", []))
    # Not required: `T?` (left out when empty). Required but nullable:
    # `T or null` (always present, null when empty).
    opt = "?" if not required else (" or null" if nullable else "")
    if s.get("type") == "array":
        inner = s.get("items", {})
        c = core(inner)
        if c is None:
            return f"{name}: json{opt}"
        return f"{name}: [{c}]{opt}"
    c = core(s)
    if c is None:
        c = "json"
    return f"{name}: {c}{opt}"

lines = ["app F33D3RTypes:"]
skipped = []
for n in sorted(S):
    s = S[n]
    if n not in declared:
        skipped.append(n)
        continue
    req = set(s.get("required", []))
    fn = fct_name(n)
    lines.append(f"    type {fn}:" if fn == n else f'    type {fn} as "{n}":')
    for f, fs in s["properties"].items():
        lines.append("        " + field(f, fs, f in req))
open(OUT, "w").write("\n".join(lines) + "\n")
print(f"wrote {OUT}: {len(declared)} types, {len(skipped)} schemas not declared as fct types:")
for n in skipped:
    why = "published by the fct runtime (contract declaration)" if n in RUNTIME_OWNED else \
          "name is not an identifier" if not nameable(n) else \
          ("no properties" if not S[n].get("properties") else "a field name is not an identifier")
    print(f"  {n}: {why}")
