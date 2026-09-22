#!/usr/bin/env python3
# Starts facetql + api/main.fct, fetches GET /api/_contract, and diffs it
# against the golden f33d3r contract: paths declared vs missing, schemas
# declared vs missing. Run from anywhere; kills its own processes on exit.
import json, os, signal, subprocess, sys, time, urllib.request

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))  # facets/
FCT_ROOT = os.path.join(os.path.dirname(ROOT), "fct")
FACET_BIN = os.environ.get("FACET_BIN",
    "/tmp/claude-1000/-home-hiiro-Documents-f33d3r-Inc/0758b9de-3cf0-4bf1-9fbe-de801e7f20ab/scratchpad/facet")
FACETQL_BIN = "/home/hiiro/Documents/f33d3r-Inc/facetql/target/release/facetql"
GOLDEN = "/home/hiiro/f33d3r/feed-engine/internal/handler/testdata/api_v2_contract.golden.json"
DATA_DIR = "/tmp/contract-diff-fqdata"
APP = os.path.join(ROOT, "api", "main.fct")
FQ_PORT = "8299"
APP_PORT = "7599"

def fetch(url, headers=None, timeout=5):
    req = urllib.request.Request(url, headers=headers or {})
    with urllib.request.urlopen(req, timeout=timeout) as r:
        return json.load(r)

def wait_for(url, tries=30):
    for _ in range(tries):
        try:
            urllib.request.urlopen(url, timeout=1)
            return True
        except Exception:
            time.sleep(0.5)
    return False

def main():
    os.makedirs(DATA_DIR, exist_ok=True)
    env = dict(os.environ)
    env.update({
        "FACETQL_DATA_DIR": DATA_DIR, "FACETQL_ENV": "development",
        "FACETQL_ALLOW_PLAINTEXT": "1", "FACETQL_TOKEN": "cdtok",
        "ENOCHIAN_TOKENS": "cdtok:root:admin",
    })
    fq = subprocess.Popen([FACETQL_BIN, "start", "--port", FQ_PORT],
                           env=env, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                           preexec_fn=os.setsid)
    app = None
    try:
        if not wait_for(f"http://127.0.0.1:{FQ_PORT}/"):
            print("facetql did not come up"); sys.exit(1)
        aenv = dict(env)
        aenv["FACET_DATABASE_URL"] = f"facetql://cdtok@127.0.0.1:{FQ_PORT}"
        app = subprocess.Popen([FACET_BIN, "run", APP, f":{APP_PORT}"],
                                env=aenv, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE,
                                preexec_fn=os.setsid)
        if not wait_for(f"http://127.0.0.1:{APP_PORT}/"):
            err = app.stderr.read().decode(errors="replace")[-2000:]
            print("app did not come up\n" + err); sys.exit(1)
        contract = fetch(f"http://127.0.0.1:{APP_PORT}/api/_contract")
        golden = json.load(open(GOLDEN))

        gp = set(golden["paths"].keys())
        dp = set(contract["paths"].keys())
        gs = set(golden["components"]["schemas"].keys())
        ds = set(contract["components"]["schemas"].keys())

        print(f"golden paths:    {len(gp)}")
        print(f"declared paths:  {len(dp & gp)} (+ {len(dp - gp)} extra/non-golden)")
        print(f"missing paths:   {len(gp - dp)}")
        print(f"golden schemas:  {len(gs)}")
        print(f"declared schemas:{len(ds & gs)}")
        print(f"missing schemas: {len(gs - ds)}")
        if "--paths" in sys.argv:
            print("\n--- missing paths ---")
            for p in sorted(gp - dp):
                print(" ", p)
        if "--schemas" in sys.argv:
            print("\n--- missing schemas ---")
            for s in sorted(gs - ds):
                print(" ", s)
    finally:
        for p in (app, fq):
            if p is None:
                continue
            try:
                os.killpg(os.getpgid(p.pid), signal.SIGTERM)
            except Exception:
                pass

if __name__ == "__main__":
    main()
