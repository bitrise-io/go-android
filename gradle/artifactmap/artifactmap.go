// Package artifactmap defines the variant-keyed Android artifact map that the
// Android build steps (gradle-runner, android-build) export and the Google
// Play Deploy step consumes.
//
// The map answers the question the flat BITRISE_*_PATH outputs cannot: which
// exported APK/AAB and which mapping.txt belong to the same build variant. A
// multi-variant build produces one mapping file per variant, and Google Play
// accepts exactly one mapping per version code, so a consumer must be able to
// pair them by identity instead of guessing by export order.
//
// The map is written as a JSON file into the Bitrise deploy directory, next to
// the exported artifacts, and its path is exported as
// BITRISE_ANDROID_ARTIFACT_MAP_PATH. All file references inside the map are
// bare file names relative to the map file's own directory, so the map stays
// valid when the deploy directory is archived or moved. Use Resolve to turn an
// entry's file name into a full path.
//
// The package intentionally depends only on the standard library so a consumer
// step can vendor just this package without pulling in the rest of the gradle
// helpers.
package artifactmap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Version is the schema version this package reads and writes.
const Version = 1

// DefaultFileName is the conventional name of the map file inside the deploy
// directory. Producers may deviate (e.g. on name collision), consumers must
// not rely on it and should use the exported path instead.
const DefaultFileName = "android-artifact-map.json"

// EnvKey is the environment variable the producing steps export the map file's
// path in, and the consuming steps read it from by default.
const EnvKey = "BITRISE_ANDROID_ARTIFACT_MAP_PATH"

// Map is the top-level document: the artifacts of one build, grouped by
// variant.
type Map struct {
	// Version identifies the schema so future readers can detect
	// incompatible documents.
	Version int `json:"version"`
	// Variants is keyed by the merged Gradle variant name (for example
	// "demoRelease"). When two modules produce the same variant name, the
	// colliding keys are disambiguated as "module/variant". Consumers must
	// treat the key as informational and pair artifacts by file identity.
	Variants map[string]Entry `json:"variants"`
	// Unmatched lists exported files whose variant could not be derived
	// from their build-output path (for example artifacts written to
	// non-standard locations). They are preserved so no export is silently
	// unaccounted for.
	Unmatched Unmatched `json:"unmatched"`
}

// Entry groups the exported files of a single build variant. All file
// references are names relative to the map file's directory.
type Entry struct {
	// Module is the Gradle module the variant belongs to (for example
	// "app"). Empty when it could not be derived.
	Module string `json:"module"`
	// Mapping is the R8/ProGuard mapping file of the variant, empty when
	// the variant produced none. A variant has at most one mapping file,
	// mirroring Google Play's one-mapping-per-version-code model.
	Mapping string `json:"mapping,omitempty"`
	// AAB lists the app bundles of the variant (in practice at most one).
	AAB []string `json:"aab"`
	// APK lists the APKs of the variant; ABI/density splits make several
	// per variant legitimate.
	APK []string `json:"apk"`
}

// Unmatched lists exported files that could not be attributed to a variant.
type Unmatched struct {
	APK     []string `json:"apk"`
	AAB     []string `json:"aab"`
	Mapping []string `json:"mapping"`
}

// File pairs a file's location in the deploy directory with the build-output
// path it was copied from. DeployPath is what ends up referenced in the map
// (as its base name); SourcePath is what the variant is derived from, because
// the deploy directory is flat while the Gradle output tree encodes the
// variant.
type File struct {
	DeployPath string
	SourcePath string
}

// Build assembles a Map from the files a step exported. Variants are derived
// from each file's SourcePath; files with unrecognisable paths are collected
// under Unmatched. When several mapping files resolve to the same variant the
// last one wins and the case is reported in warnings (one message per
// overwritten file); warnings is empty on a clean build.
func Build(apks, aabs, mappings []File) (Map, []string) {
	type group struct {
		variant ArtifactVariant
		entry   Entry
	}
	groups := map[ArtifactVariant]*group{}
	ordered := []*group{} // deterministic iteration for key assignment
	var warnings []string

	grab := func(f File) (*group, bool) {
		variant, ok := VariantFromPath(f.SourcePath)
		if !ok {
			return nil, false
		}
		g, ok := groups[variant]
		if !ok {
			g = &group{variant: variant, entry: Entry{Module: variant.Module, AAB: []string{}, APK: []string{}}}
			groups[variant] = g
			ordered = append(ordered, g)
		}
		return g, true
	}

	unmatched := Unmatched{APK: []string{}, AAB: []string{}, Mapping: []string{}}
	for _, f := range apks {
		if g, ok := grab(f); ok {
			g.entry.APK = append(g.entry.APK, filepath.Base(f.DeployPath))
		} else {
			unmatched.APK = append(unmatched.APK, filepath.Base(f.DeployPath))
		}
	}
	for _, f := range aabs {
		if g, ok := grab(f); ok {
			g.entry.AAB = append(g.entry.AAB, filepath.Base(f.DeployPath))
		} else {
			unmatched.AAB = append(unmatched.AAB, filepath.Base(f.DeployPath))
		}
	}
	for _, f := range mappings {
		g, ok := grab(f)
		if !ok {
			unmatched.Mapping = append(unmatched.Mapping, filepath.Base(f.DeployPath))
			continue
		}
		name := filepath.Base(f.DeployPath)
		if g.entry.Mapping != "" && g.entry.Mapping != name {
			warnings = append(warnings, fmt.Sprintf(
				"variant %s matched several mapping files: keeping %s, dropping %s",
				g.variant.Variant, name, g.entry.Mapping))
		}
		g.entry.Mapping = name
	}

	// Key by variant name; disambiguate as module/variant only on collision.
	nameCount := map[string]int{}
	for _, g := range ordered {
		nameCount[g.variant.Variant]++
	}
	variants := map[string]Entry{}
	for _, g := range ordered {
		key := g.variant.Variant
		if nameCount[g.variant.Variant] > 1 {
			key = g.variant.Module + "/" + g.variant.Variant
		}
		// The producers discover files in filesystem-walk order, which is not a
		// meaningful contract; sort so the document is deterministic.
		sort.Strings(g.entry.APK)
		sort.Strings(g.entry.AAB)
		variants[key] = g.entry
	}
	sort.Strings(unmatched.APK)
	sort.Strings(unmatched.AAB)
	sort.Strings(unmatched.Mapping)

	return Map{Version: Version, Variants: variants, Unmatched: unmatched}, warnings
}

// IsEmpty reports whether the map carries no artifacts at all — producers can
// skip writing a file for such a build.
func (m Map) IsEmpty() bool {
	return len(m.Variants) == 0 &&
		len(m.Unmatched.APK) == 0 && len(m.Unmatched.AAB) == 0 && len(m.Unmatched.Mapping) == 0
}

// SortedVariantKeys returns the variant keys in lexical order, for
// deterministic logging and iteration.
func (m Map) SortedVariantKeys() []string {
	keys := make([]string, 0, len(m.Variants))
	for key := range m.Variants {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Resolve turns a file name referenced by the map into a full path, resolving
// it against the directory of the map file itself (file names in the map are
// relative to it). mapPath is the path the map was read from; name is an
// Entry/Unmatched value. An empty name resolves to "".
func Resolve(mapPath, name string) string {
	if name == "" {
		return ""
	}
	return filepath.Join(filepath.Dir(mapPath), name)
}

// Write marshals the map and writes it to path. The document is indented so
// the exported build artifact stays human-readable.
func Write(path string, m Map) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal artifact map: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("write artifact map: %w", err)
	}
	return nil
}

// Read loads and validates a map written by Write. It rejects documents with a
// newer schema version than this package understands, and documents that are
// not artifact maps at all (missing version).
func Read(path string) (Map, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Map{}, fmt.Errorf("read artifact map: %w", err)
	}
	var m Map
	if err := json.Unmarshal(data, &m); err != nil {
		return Map{}, fmt.Errorf("parse artifact map %s: %w", path, err)
	}
	if m.Version == 0 {
		return Map{}, fmt.Errorf("parse artifact map %s: missing schema version, not an artifact map", path)
	}
	if m.Version > Version {
		return Map{}, fmt.Errorf("artifact map %s has schema version %d, this consumer understands up to %d", path, m.Version, Version)
	}
	return m, nil
}
