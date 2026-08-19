package artifactmap

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource},
			{DeployPath: demoMappingRenamed, SourcePath: paidMappingSource},
		})

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
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: strayMappingSource}})

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

// TestBuild_ReportFilesStayUnmatched: R8 writes usage.txt, seeds.txt etc.
// next to mapping.txt; when a widened filter matches them too, only the file
// literally named mapping.txt can be the variant's mapping — report files
// stay visible under unmatched.
func TestBuild_ReportFilesStayUnmatched(t *testing.T) {
	demoMappingDir := "/bitrise/src/app/build/outputs/mapping/demoRelease/"
	m, warnings := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil,
		nil,
		[]File{
			{DeployPath: deploy("usage.txt"), SourcePath: demoMappingDir + "usage.txt"},
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingDir + "mapping.txt"},
		})

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if got := m.Modules["app"]["demoRelease"].Mapping; got != "mapping.txt" {
		t.Fatalf("Mapping = %q, want the canonical mapping.txt", got)
	}
	if want := []string{"usage.txt"}; !reflect.DeepEqual(m.Unmatched.Mapping, want) {
		t.Fatalf("Unmatched.Mapping = %v, want %v", m.Unmatched.Mapping, want)
	}
}

// TestBuild_OriginalIssueLayout replays the exact customer log from SSW-3065:
// a single-variant build where three files matched a widened mapping filter —
// the Compose mapping, the R8 task's intermediates workdir copy, and the
// official outputs/ copy. Only the official file pairs with the variant,
// regardless of discovery order (the deploy-dir collision renames whichever
// arrives second); the two intermediates files are left out of the document.
func TestBuild_OriginalIssueLayout(t *testing.T) {
	const (
		composeSrc       = "/bitrise/src/app/build/intermediates/compose_mapping/productionRelease/compose-mapping.txt"
		intermediatesSrc = "/bitrise/src/app/build/intermediates/mapping/productionRelease/minifyProductionReleaseWithR8/mapping.txt"
		outputsSrc       = "/bitrise/src/app/build/outputs/mapping/productionRelease/mapping.txt"
		apkSrc           = "/bitrise/src/app/build/outputs/apk/production/release/app-production-release.apk"
	)

	cases := map[string]struct {
		mappings []File
		want     string // deploy name of the outputs-sourced file
	}{
		"customer discovery order (intermediates copied first)": {
			mappings: []File{
				{DeployPath: deploy("compose-mapping.txt"), SourcePath: composeSrc},
				{DeployPath: deploy("mapping.txt"), SourcePath: intermediatesSrc},
				{DeployPath: deploy("mapping20260703072155.txt"), SourcePath: outputsSrc},
			},
			want: "mapping20260703072155.txt",
		},
		"reversed discovery order (outputs copied first)": {
			mappings: []File{
				{DeployPath: deploy("mapping.txt"), SourcePath: outputsSrc},
				{DeployPath: deploy("mapping20260703072155.txt"), SourcePath: intermediatesSrc},
				{DeployPath: deploy("compose-mapping.txt"), SourcePath: composeSrc},
			},
			want: "mapping.txt",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			m, warnings := Build(
				[]File{{DeployPath: deploy("app-production-release.apk"), SourcePath: apkSrc}},
				nil,
				nil,
				tc.mappings)
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
			if want := []string{}; !reflect.DeepEqual(m.Unmatched.Mapping, want) {
				t.Fatalf("Unmatched.Mapping = %v, want intermediates files ignored", m.Unmatched.Mapping)
			}
		})
	}
}

func TestBuild_DuplicateMappingWarnsAndKeepsLast(t *testing.T) {
	m, warnings := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil,
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource},
			{DeployPath: otherMappingRenamed, SourcePath: demoMappingSource},
		})

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
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource},
			{DeployPath: demoMappingRenamed, SourcePath: wearMappingSource},
		})

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

// TestBuild_SameBasenameModulesStayDistinct: two modules whose directories
// share a basename (monorepo brandA/app + brandB/app) must not merge into one
// "app" key — the document keys grow parent directories until unique.
func TestBuild_SameBasenameModulesStayDistinct(t *testing.T) {
	m, warnings := Build(
		[]File{
			{DeployPath: deploy("brandA.apk"), SourcePath: "/bitrise/src/brandA/app/build/outputs/apk/release/brandA.apk"},
			{DeployPath: deploy("brandB.apk"), SourcePath: "/bitrise/src/brandB/app/build/outputs/apk/release/brandB.apk"},
		},
		nil,
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: "/bitrise/src/brandA/app/build/outputs/mapping/release/mapping.txt"},
			{DeployPath: demoMappingRenamed, SourcePath: "/bitrise/src/brandB/app/build/outputs/mapping/release/mapping.txt"},
		})

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if want := []string{"brandA/app", "brandB/app"}; !reflect.DeepEqual(m.SortedModules(), want) {
		t.Fatalf("modules = %v, want %v", m.SortedModules(), want)
	}
	if got := m.Modules["brandA/app"]["release"].Mapping; got != "mapping.txt" {
		t.Fatalf("brandA mapping = %q, want mapping.txt", got)
	}
	if got := m.Modules["brandB/app"]["release"].Mapping; got != "mapping-20260805121530.txt" {
		t.Fatalf("brandB mapping = %q, want the renamed file", got)
	}
}

// TestBuild_SourcesRecordProvenance: every referenced name — variant files and
// unmatched alike — carries its build-output path; names the document dropped
// (a replaced duplicate mapping) and ignored intermediates files do not.
func TestBuild_SourcesRecordProvenance(t *testing.T) {
	m, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil,
		nil,
		[]File{
			{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource},
			{DeployPath: otherMappingRenamed, SourcePath: demoMappingSource}, // replaces mapping.txt
			{DeployPath: deploy("stray-mapping.txt"), SourcePath: strayMappingSource},
			{DeployPath: deploy("compose-mapping.txt"), SourcePath: "/bitrise/src/app/build/intermediates/compose_mapping/demoRelease/compose-mapping.txt"},
		})

	want := map[string]string{
		"app-demo-release.apk":       demoAPKSource,
		"mapping-20260805121599.txt": demoMappingSource,
		"stray-mapping.txt":          strayMappingSource,
	}
	if !reflect.DeepEqual(m.Sources, want) {
		t.Fatalf("Sources = %v, want %v", m.Sources, want)
	}
}

// TestBuild_AARs: AGP's aar layout encodes no variant, so AARs attach to
// their module — with the same key disambiguation modules get. Files not
// under build/outputs/aar stay visible in unmatched.
func TestBuild_AARs(t *testing.T) {
	m, warnings := Build(
		nil,
		nil,
		[]File{
			{DeployPath: deploy("data-debug.aar"), SourcePath: "/bitrise/src/feature-name-1/data/build/outputs/aar/data-debug.aar"},
			{DeployPath: deploy("data-debug20260818.aar"), SourcePath: "/bitrise/src/feature-name-2/data/build/outputs/aar/data-debug.aar"},
			{DeployPath: deploy("stray.aar"), SourcePath: "/bitrise/src/custom-out/stray.aar"},
		},
		nil,
	)

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	want := map[string][]string{
		"feature-name-1/data": {"data-debug.aar"},
		"feature-name-2/data": {"data-debug20260818.aar"},
	}
	if !reflect.DeepEqual(m.ModuleAARs, want) {
		t.Fatalf("ModuleAARs = %v, want %v", m.ModuleAARs, want)
	}
	if want := []string{"stray.aar"}; !reflect.DeepEqual(m.Unmatched.AAR, want) {
		t.Fatalf("Unmatched.AAR = %v, want %v", m.Unmatched.AAR, want)
	}
}

// TestBuild_PhantomVariantNestingWarns: unexpected nesting under
// outputs/<kind>/ turns the variant into a phantom key ("demoReleaseExtra"),
// so the mapping's real variant ends up with no app artifact — the map can't
// detect the phantom itself, but the orphaned mapping is warned about.
func TestBuild_PhantomVariantNestingWarns(t *testing.T) {
	m, warnings := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: "/bitrise/src/app/build/outputs/apk/demo/release/extra/app-demo-release.apk"}},
		nil,
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)

	if len(warnings) != 1 || !strings.Contains(warnings[0], "app/demoRelease has a mapping but no app artifact") {
		t.Fatalf("expected the orphaned-mapping warning, got %v", warnings)
	}
	if got := m.Modules["app"]["demoReleaseExtra"].APK; !reflect.DeepEqual(got, []string{"app-demo-release.apk"}) {
		t.Fatalf("APK = %v, want it under the (phantom) derived variant", got)
	}
	if got := m.Modules["app"]["demoRelease"].Mapping; got != "mapping.txt" {
		t.Fatalf("Mapping = %q, want it under the real variant", got)
	}
}

func TestBuild_EmptyInput(t *testing.T) {
	m, warnings := Build(nil, nil, nil, nil)
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
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}})

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
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}})

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
		"unmatched": map[string]any{"apk": []any{}, "aab": []any{}, "aar": []any{}, "mapping": []any{}},
		"sources": map[string]any{
			"app-demo-release.apk": demoAPKSource,
			"app-demo-release.aab": demoAABSource,
			"mapping.txt":          demoMappingSource,
		},
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

// TestBuild_ModuleNamedIntermediatesIsNotDropped guards the intermediates rule's
// anchor: only Gradle's build/intermediates/ tree is task-workdir noise, not a
// module or checkout directory that happens to carry the name.
func TestBuild_ModuleNamedIntermediatesIsNotDropped(t *testing.T) {
	const (
		apkSource     = "/bitrise/src/intermediates/build/outputs/apk/demo/release/app-demo-release.apk"
		mappingSource = "/bitrise/src/intermediates/build/outputs/mapping/demoRelease/mapping.txt"
	)
	m, warnings := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: apkSource}},
		nil,
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: mappingSource}})

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	want := map[string]map[string]Entry{
		"intermediates": {
			"demoRelease": {Mapping: "mapping.txt", AAB: []string{}, APK: []string{"app-demo-release.apk"}},
		},
	}
	if !reflect.DeepEqual(m.Modules, want) {
		t.Fatalf("Modules = %+v, want %+v", m.Modules, want)
	}
}

// TestRead_NormalizesAbsentLists keeps null out of consumers' way: a foreign
// same-version document may omit the artifact lists, and jq's .apk[] errors on
// null.
func TestRead_NormalizesAbsentLists(t *testing.T) {
	path := filepath.Join(t.TempDir(), DefaultFileName)
	document := `{"version": 1, "modules": {"app": {"demoRelease": {"mapping": "mapping.txt"}}}}`
	if err := os.WriteFile(path, []byte(document), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	entry := m.Modules["app"]["demoRelease"]
	if entry.APK == nil || entry.AAB == nil {
		t.Fatalf("entry lists = %+v, want empty slices, not nil", entry)
	}
	if m.Unmatched.APK == nil || m.Unmatched.AAB == nil || m.Unmatched.AAR == nil || m.Unmatched.Mapping == nil {
		t.Fatalf("unmatched lists = %+v, want empty slices, not nil", m.Unmatched)
	}

	// a re-written document must then be null-free too
	out := filepath.Join(t.TempDir(), DefaultFileName)
	if err := Write(out, m); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "null") {
		t.Fatalf("re-written document contains null:\n%s", content)
	}
}

func TestRead_MissingModulesMapIsUsable(t *testing.T) {
	path := filepath.Join(t.TempDir(), DefaultFileName)
	if err := os.WriteFile(path, []byte(`{"version": 1}`), 0644); err != nil {
		t.Fatal(err)
	}
	m, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Modules == nil {
		t.Fatal("Modules = nil, want an empty map callers can range over and assign into")
	}
	m.Modules["app"] = map[string]Entry{"demoRelease": {}}
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
