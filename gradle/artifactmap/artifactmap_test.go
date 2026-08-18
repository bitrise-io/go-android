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

func TestBuild_GroupsByModuleAndVariant(t *testing.T) {
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
	want := map[string]map[string]Entry{
		"app": {
			"demoRelease": {
				Mapping: "mapping.txt",
				AAB:     []string{"app-demo-release.aab"},
				// names are sorted, not discovery-ordered
				APK: []string{"app-demo-arm64-v8a-release.apk", "app-demo-release.apk"},
			},
			"paidRelease": {
				Mapping: "mapping-20260805121530.txt",
				AAB:     []string{},
				APK:     []string{"app-paid-release.apk"},
			},
		},
	}
	if !reflect.DeepEqual(m.Modules, want) {
		t.Fatalf("Modules = %+v, want %+v", m.Modules, want)
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
	if len(m.Modules) != 0 {
		t.Fatalf("expected no modules, got %+v", m.Modules)
	}
	if want := []string{"mapping.txt"}; !reflect.DeepEqual(m.Unmatched.Mapping, want) {
		t.Fatalf("Unmatched.Mapping = %v, want %v", m.Unmatched.Mapping, want)
	}
}

// TestBuild_CanonicalMappingTxtWinsOverReportFiles: R8 writes usage.txt,
// seeds.txt etc. next to mapping.txt; when a widened filter matches them too,
// the file literally named mapping.txt must stay the variant's mapping no
// matter the discovery order.
func TestBuild_CanonicalMappingTxtWinsOverReportFiles(t *testing.T) {
	demoMappingDir := "/bitrise/src/app/build/outputs/mapping/demoRelease/"
	m, warnings := Build(
		nil,
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingDir + "mapping.txt"},
			{DeployPath: deploy("usage.txt"), SourcePath: demoMappingDir + "usage.txt"},
		},
	)

	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning about the dropped report file, got %v", warnings)
	}
	if got := m.Modules["app"]["demoRelease"].Mapping; got != "mapping.txt" {
		t.Fatalf("Mapping = %q, want the canonical mapping.txt", got)
	}

	// same result when the report file is discovered first
	m, _ = Build(
		nil,
		nil,
		[]File{
			{DeployPath: deploy("usage.txt"), SourcePath: demoMappingDir + "usage.txt"},
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingDir + "mapping.txt"},
		},
	)
	if got := m.Modules["app"]["demoRelease"].Mapping; got != "mapping.txt" {
		t.Fatalf("Mapping = %q, want the canonical mapping.txt regardless of order", got)
	}
}

// TestBuild_OriginalIssueLayout replays the exact customer log from SSW-3065:
// a single-variant build where three files matched a widened mapping filter —
// the Compose mapping, the R8 task's intermediates workdir copy, and the
// official outputs/ copy. Only the official file pairs with the variant,
// regardless of discovery order (the deploy-dir collision renames whichever
// arrives second); the other two stay visible under unmatched.
func TestBuild_OriginalIssueLayout(t *testing.T) {
	const (
		composeSrc       = "/bitrise/src/app/build/intermediates/compose_mapping/productionRelease/compose-mapping.txt"
		intermediatesSrc = "/bitrise/src/app/build/intermediates/mapping/productionRelease/minifyProductionReleaseWithR8/mapping.txt"
		outputsSrc       = "/bitrise/src/app/build/outputs/mapping/productionRelease/mapping.txt"
		apkSrc           = "/bitrise/src/app/build/outputs/apk/production/release/app-production-release.apk"
	)

	cases := map[string]struct {
		mappings      []File
		want          string   // deploy name of the outputs-sourced file
		wantUnmatched []string // sorted
	}{
		"customer discovery order (intermediates copied first)": {
			mappings: []File{
				{DeployPath: deploy("compose-mapping.txt"), SourcePath: composeSrc},
				{DeployPath: deploy("mapping.txt"), SourcePath: intermediatesSrc},
				{DeployPath: deploy("mapping20260703072155.txt"), SourcePath: outputsSrc},
			},
			want:          "mapping20260703072155.txt",
			wantUnmatched: []string{"compose-mapping.txt", "mapping.txt"},
		},
		"reversed discovery order (outputs copied first)": {
			mappings: []File{
				{DeployPath: deploy("mapping.txt"), SourcePath: outputsSrc},
				{DeployPath: deploy("mapping20260703072155.txt"), SourcePath: intermediatesSrc},
				{DeployPath: deploy("compose-mapping.txt"), SourcePath: composeSrc},
			},
			want:          "mapping.txt",
			wantUnmatched: []string{"compose-mapping.txt", "mapping20260703072155.txt"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			m, warnings := Build(
				[]File{{DeployPath: deploy("app-production-release.apk"), SourcePath: apkSrc}},
				nil,
				tc.mappings,
			)
			if len(warnings) != 0 {
				t.Fatalf("unexpected warnings: %v", warnings)
			}
			entry := m.Modules["app"]["productionRelease"]
			if entry.Mapping != tc.want {
				t.Fatalf("Mapping = %q, want the outputs-sourced %q", entry.Mapping, tc.want)
			}
			if !reflect.DeepEqual(entry.APK, []string{"app-production-release.apk"}) {
				t.Fatalf("APK = %v, want the variant's APK paired", entry.APK)
			}
			if !reflect.DeepEqual(m.Unmatched.Mapping, tc.wantUnmatched) {
				t.Fatalf("Unmatched.Mapping = %v, want %v", m.Unmatched.Mapping, tc.wantUnmatched)
			}
		})
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
	if got := m.Modules["app"]["demoRelease"].Mapping; got != "mapping-20260805121599.txt" {
		t.Fatalf("Mapping = %q, want the last one", got)
	}
}

// TestBuild_SameVariantNameAcrossModules: two modules declaring the same
// variant name must land under their own module keys.
func TestBuild_SameVariantNameAcrossModules(t *testing.T) {
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

	if want := []string{"app", "wear"}; !reflect.DeepEqual(m.SortedModules(), want) {
		t.Fatalf("modules = %v, want %v", m.SortedModules(), want)
	}
	if got := m.Modules["app"]["demoRelease"].Mapping; got != "mapping.txt" {
		t.Fatalf("app mapping = %q, want mapping.txt", got)
	}
	if got := m.Modules["wear"]["demoRelease"].Mapping; got != "mapping-20260805121530.txt" {
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
		"modules": map[string]any{
			"app": map[string]any{
				"demoRelease": map[string]any{
					"mapping": "mapping.txt",
					"aab":     []any{"app-demo-release.aab"},
					"apk":     []any{"app-demo-release.apk"},
				},
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
	if err := os.WriteFile(path, []byte(`{"version": 99, "modules": {}}`), 0644); err != nil {
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

func TestLabel(t *testing.T) {
	if got := Label("app", "demoRelease"); got != "app/demoRelease" {
		t.Fatalf("Label = %q", got)
	}
	if got := Label("", "demoRelease"); got != "demoRelease" {
		t.Fatalf("Label with empty module = %q", got)
	}
}
