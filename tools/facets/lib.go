package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Manifest is the subset of facet.json this tool reads.
type Manifest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// Library is a discovered facets stdlib repo rooted at Root.
type Library struct {
	Root     string
	Manifest Manifest // zero value if facet.json is absent/unreadable
}

// FindRoot locates the library root by walking up from start looking for
// facet.json. If none is found, start (made absolute) is returned — the tool then
// operates on whatever `.fct` files live there.
func FindRoot(start string) (string, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	dir := abs
	if fi, err := os.Stat(dir); err == nil && !fi.IsDir() {
		dir = filepath.Dir(dir)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "facet.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return abs, nil // no manifest anywhere; fall back to the given dir
		}
		dir = parent
	}
}

// OpenLibrary discovers the library root from start and loads its manifest.
func OpenLibrary(start string) (*Library, error) {
	root, err := FindRoot(start)
	if err != nil {
		return nil, err
	}
	lib := &Library{Root: root}
	if b, err := os.ReadFile(filepath.Join(root, "facet.json")); err == nil {
		_ = json.Unmarshal(b, &lib.Manifest) // best-effort; malformed manifest is non-fatal
	}
	return lib, nil
}

// FctFiles returns every `.fct` file under the library root, sorted, relative to
// the root. Hidden dirs and the tools/ dir are skipped.
func (l *Library) FctFiles() ([]string, error) {
	var out []string
	err := filepath.WalkDir(l.Root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if p != l.Root && (strings.HasPrefix(name, ".") || name == "tools" || name == "facet_modules") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".fct") {
			rel, err := filepath.Rel(l.Root, p)
			if err != nil {
				return err
			}
			out = append(out, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}

// Category is the top-level directory of a repo-relative path ("" for root files).
func Category(rel string) string {
	if i := strings.IndexByte(rel, '/'); i >= 0 {
		return rel[:i]
	}
	return ""
}

// resolveImport joins a relative import target against the importing file's
// directory and returns a cleaned, root-contained absolute path. It rejects
// targets that escape the library root (path-traversal guard) and remote refs
// (github.com/...), which this structural tool does not fetch.
func resolveImport(root, fromFile, target string) (string, error) {
	if strings.Contains(target, "://") || strings.HasPrefix(target, "github.com/") {
		return "", errors.New("remote import (not resolved by this tool)")
	}
	base := filepath.Dir(fromFile)
	joined := filepath.Clean(filepath.Join(base, target))
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(absRoot, absJoined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("import escapes library root")
	}
	return absJoined, nil
}
