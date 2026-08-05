package artifactmap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	demoAPKSource       = "/bitrise/src/app/build/outputs/apk/demo/release/app-demo-release.apk"
	demoSplitAPKSource  = "/bitrise/src/app/build/outputs/apk/demo/release/app-demo-arm64-v8a-release.apk"
	demoAABSource       = "/bitrise/src/app/build/outputs/bundle/demoRelease/app-demo-release.aab"
	demoMappingSource   = "/bitrise/src/app/build/outputs/mapping/demoRelease/mapping.txt"
	paidAPKSource       = "/bitrise/src/app/build/outputs/apk/paid/release/app-paid-release.apk"
	paidMappingSource   = "/bitrise/src/app/build/outputs/mapping/paidRelease/mapping.txt"
	strayMappingSource  = "/bitrise/src/app/custom-out/mapping.txt"
	wearAPKSource       = "/bitrise/src/wear/build/outputs/apk/demo/release/wear-demo-release.apk"
	wearMappingSource   = "/bitrise/src/wear/build/outputs/mapping/demoRelease/mapping.txt"
	deployDir           = "/bitrise/deploy"
	demoMappingRenamed  = "/bitrise/deploy/mapping-20260805121530.txt"
	otherMappingRenamed = "/bitrise/deploy/mapping-20260805121599.txt"
)

func deploy(name string) string { return filepath.Join(deployDir, name) }

func TestBuild_GroupsByVariant(t *testing.T) {
	m, warnings := Build(
		[]File{
			{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource},
			{DeployPath: deploy("app-demo-arm64-v8a-release.apk"), SourcePath: demoSplitAPKSource},
			{DeployPath: deploy("app-paid-release.apk"), SourcePath: paidAPKSource},
		},
		[]File{
			{DeployPath: deploy("app-demo-release.aab"), SourcePath: demoAABSource},
		},
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource},
			{DeployPath: demoMappingRenamed, SourcePath: paidMappingSource},
		},
	)

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	want := map[string]Entry{
		"demoRelease": {
			Module:  "app",
			Mapping: "mapping.txt",
			AAB:     []string{"app-demo-release.aab"},
			// names are sorted, not discovery-ordered
			APK: []string{"app-demo-arm64-v8a-release.apk", "app-demo-release.apk"},
		},
		"paidRelease": {
			Module:  "app",
			Mapping: "mapping-20260805121530.txt",
			AAB:     []string{},
			APK:     []string{"app-paid-release.apk"},
		},
	}
	if !reflect.DeepEqual(m.Variants, want) {
		t.Fatalf("Variants = %+v, want %+v", m.Variants, want)
	}
	if m.Version != Version {
		t.Fatalf("Version = %d, want %d", m.Version, Version)
	}
	if len(m.Unmatched.APK)+len(m.Unmatched.AAB)+len(m.Unmatched.Mapping) != 0 {
		t.Fatalf("expected no unmatched files, got %+v", m.Unmatched)
	}
}

func TestBuild_UnderivableVariantGoesToUnmatched(t *testing.T) {
	m, warnings := Build(
		nil,
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: strayMappingSource}},
	)

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(m.Variants) != 0 {
		t.Fatalf("expected no variants, got %+v", m.Variants)
	}
	if want := []string{"mapping.txt"}; !reflect.DeepEqual(m.Unmatched.Mapping, want) {
		t.Fatalf("Unmatched.Mapping = %v, want %v", m.Unmatched.Mapping, want)
	}
}

func TestBuild_DuplicateMappingWarnsAndKeepsLast(t *testing.T) {
	m, warnings := Build(
		nil,
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource},
			{DeployPath: otherMappingRenamed, SourcePath: demoMappingSource},
		},
	)

	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %v", warnings)
	}
	if got := m.Variants["demoRelease"].Mapping; got != "mapping-20260805121599.txt" {
		t.Fatalf("Mapping = %q, want the last one", got)
	}
}

func TestBuild_SameVariantNameAcrossModulesGetsModulePrefixedKeys(t *testing.T) {
	m, _ := Build(
		[]File{
			{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource},
			{DeployPath: deploy("wear-demo-release.apk"), SourcePath: wearAPKSource},
		},
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource},
			{DeployPath: demoMappingRenamed, SourcePath: wearMappingSource},
		},
	)

	wantKeys := []string{"app/demoRelease", "wear/demoRelease"}
	if got := m.SortedVariantKeys(); !reflect.DeepEqual(got, wantKeys) {
		t.Fatalf("keys = %v, want %v", got, wantKeys)
	}
	if got := m.Variants["app/demoRelease"].Mapping; got != "mapping.txt" {
		t.Fatalf("app mapping = %q, want mapping.txt", got)
	}
	if got := m.Variants["wear/demoRelease"].Mapping; got != "mapping-20260805121530.txt" {
		t.Fatalf("wear mapping = %q, want the renamed file", got)
	}
}

func TestBuild_EmptyInput(t *testing.T) {
	m, warnings := Build(nil, nil, nil)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if !m.IsEmpty() {
		t.Fatalf("expected empty map, got %+v", m)
	}
}

func TestWriteRead_RoundTrip(t *testing.T) {
	m, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		[]File{{DeployPath: deploy("app-demo-release.aab"), SourcePath: demoAABSource}},
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)

	path := filepath.Join(t.TempDir(), DefaultFileName)
	if err := Write(path, m); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !reflect.DeepEqual(got, m) {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", got, m)
	}
}

// TestWrite_DocumentShape locks the on-disk contract: this is what consumers
// outside this module parse, so a change here is a schema version bump.
func TestWrite_DocumentShape(t *testing.T) {
	m, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		[]File{{DeployPath: deploy("app-demo-release.aab"), SourcePath: demoAABSource}},
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)

	path := filepath.Join(t.TempDir(), DefaultFileName)
	if err := Write(path, m); err != nil {
		t.Fatalf("Write: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("written document is not valid JSON: %v", err)
	}
	want := map[string]any{
		"version": float64(1),
		"variants": map[string]any{
			"demoRelease": map[string]any{
				"module":  "app",
				"mapping": "mapping.txt",
				"aab":     []any{"app-demo-release.aab"},
				"apk":     []any{"app-demo-release.apk"},
			},
		},
		"unmatched": map[string]any{"apk": []any{}, "aab": []any{}, "mapping": []any{}},
	}
	if !reflect.DeepEqual(doc, want) {
		t.Fatalf("document shape changed:\n got %#v\nwant %#v", doc, want)
	}
}

func TestRead_RejectsNewerVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), DefaultFileName)
	if err := os.WriteFile(path, []byte(`{"version": 99, "variants": {}}`), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(path); err == nil {
		t.Fatal("expected error for newer schema version")
	}
}

func TestRead_RejectsForeignJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "something.json")
	if err := os.WriteFile(path, []byte(`{"name": "not an artifact map"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(path); err == nil {
		t.Fatal("expected error for JSON without schema version")
	}
}

func TestResolve(t *testing.T) {
	mapPath := "/bitrise/deploy/android-artifact-map.json"
	if got, want := Resolve(mapPath, "mapping.txt"), "/bitrise/deploy/mapping.txt"; got != want {
		t.Fatalf("Resolve = %q, want %q", got, want)
	}
	if got := Resolve(mapPath, ""); got != "" {
		t.Fatalf("Resolve of empty name = %q, want empty", got)
	}
}
