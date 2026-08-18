package artifactmap

import (
	"reflect"
	"testing"
)

func TestMerge_DisjointVariantsUnion(t *testing.T) {
	base, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)
	overlay, _ := Build(
		[]File{{DeployPath: deploy("app-paid-release.apk"), SourcePath: paidAPKSource}},
		nil,
		[]File{{DeployPath: demoMappingRenamed, SourcePath: paidMappingSource}},
	)

	merged, warnings := Merge(base, overlay)

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	want := map[string]map[string]Entry{
		"app": {
			"demoRelease": {Mapping: "mapping.txt", AAB: []string{}, APK: []string{"app-demo-release.apk"}},
			"paidRelease": {Mapping: "mapping-20260805121530.txt", AAB: []string{}, APK: []string{"app-paid-release.apk"}},
		},
	}
	if !reflect.DeepEqual(merged.Modules, want) {
		t.Fatalf("Modules = %+v, want %+v", merged.Modules, want)
	}
	if merged.Version != Version {
		t.Fatalf("Version = %d, want %d", merged.Version, Version)
	}
}

func TestMerge_RebuiltVariantReplacedWithWarnings(t *testing.T) {
	base, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)
	overlay, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release-20260805.apk"), SourcePath: demoAPKSource}},
		nil,
		[]File{{DeployPath: demoMappingRenamed, SourcePath: demoMappingSource}},
	)

	merged, warnings := Merge(base, overlay)

	if len(warnings) != 2 { // rebuilt APKs + rebuilt mapping
		t.Fatalf("expected 2 replacement warnings, got %v", warnings)
	}
	got := merged.Modules["app"]["demoRelease"]
	if !reflect.DeepEqual(got.APK, []string{"app-demo-release-20260805.apk"}) {
		t.Fatalf("APK = %v, want the overlay's artifact", got.APK)
	}
	if got.Mapping != "mapping-20260805121530.txt" {
		t.Fatalf("Mapping = %q, want the overlay's mapping", got.Mapping)
	}
}

// TestMerge_ApkThenAabRunsCombine: the canonical two-step workflow — one run
// builds the variant's APK, a later run its AAB. The merged entry must carry
// both; the second run must not wipe the first's artifacts.
func TestMerge_ApkThenAabRunsCombine(t *testing.T) {
	apkRun, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)
	aabRun, _ := Build(
		nil,
		[]File{{DeployPath: deploy("app-demo-release.aab"), SourcePath: demoAABSource}},
		[]File{{DeployPath: demoMappingRenamed, SourcePath: demoMappingSource}},
	)

	merged, warnings := Merge(apkRun, aabRun)

	got := merged.Modules["app"]["demoRelease"]
	if !reflect.DeepEqual(got.APK, []string{"app-demo-release.apk"}) {
		t.Fatalf("APK = %v, the first run's APK must survive the merge", got.APK)
	}
	if !reflect.DeepEqual(got.AAB, []string{"app-demo-release.aab"}) {
		t.Fatalf("AAB = %v, want the second run's AAB", got.AAB)
	}
	if got.Mapping != "mapping-20260805121530.txt" {
		t.Fatalf("Mapping = %q, want the second run's (rebuilt) mapping", got.Mapping)
	}
	if len(warnings) != 1 { // only the rebuilt mapping warns
		t.Fatalf("expected 1 warning for the rebuilt mapping, got %v", warnings)
	}
}

// TestMerge_IdenticalVariantKeepsContent: re-merging the same document must
// not lose or duplicate anything. A re-listed APK set does warn — the merge
// deliberately doesn't compare slice contents to stay simple.
func TestMerge_IdenticalVariantKeepsContent(t *testing.T) {
	m, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil, nil,
	)

	merged, warnings := Merge(m, m)

	if len(warnings) != 1 {
		t.Fatalf("expected 1 re-listing warning, got %v", warnings)
	}
	if !reflect.DeepEqual(merged.Modules, m.Modules) {
		t.Fatalf("Modules = %+v, want unchanged %+v", merged.Modules, m.Modules)
	}
}

// TestMerge_SameVariantNameAcrossModulesStaysSeparate: module app's
// demoRelease from one step and module wear's demoRelease from another must
// not be conflated — the exact case the module nesting exists for.
func TestMerge_SameVariantNameAcrossModulesStaysSeparate(t *testing.T) {
	base, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil,
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)
	overlay, _ := Build(
		[]File{{DeployPath: deploy("wear-demo-release.apk"), SourcePath: wearAPKSource}},
		nil,
		[]File{{DeployPath: demoMappingRenamed, SourcePath: wearMappingSource}},
	)

	merged, warnings := Merge(base, overlay)

	if len(warnings) != 0 {
		t.Fatalf("distinct modules must merge silently, got %v", warnings)
	}
	if got := merged.Modules["app"]["demoRelease"].Mapping; got != "mapping.txt" {
		t.Fatalf("app mapping = %q, want mapping.txt", got)
	}
	if got := merged.Modules["wear"]["demoRelease"].Mapping; got != "mapping-20260805121530.txt" {
		t.Fatalf("wear mapping = %q, want the renamed file", got)
	}
	if !reflect.DeepEqual(merged.Modules["app"]["demoRelease"].APK, []string{"app-demo-release.apk"}) {
		t.Fatalf("app APK list was disturbed by the wear merge: %+v", merged.Modules["app"]["demoRelease"].APK)
	}
}

func TestMerge_UnmatchedUnionDeduped(t *testing.T) {
	base := Map{Version: Version, Modules: map[string]map[string]Entry{}, Unmatched: Unmatched{
		APK: []string{"stray.apk"}, AAB: []string{}, Mapping: []string{"stray-mapping.txt"},
	}}
	overlay := Map{Version: Version, Modules: map[string]map[string]Entry{}, Unmatched: Unmatched{
		APK: []string{"stray.apk", "another.apk"}, AAB: []string{}, Mapping: []string{},
	}}

	merged, _ := Merge(base, overlay)

	if want := []string{"another.apk", "stray.apk"}; !reflect.DeepEqual(merged.Unmatched.APK, want) {
		t.Fatalf("Unmatched.APK = %v, want %v", merged.Unmatched.APK, want)
	}
	if want := []string{"stray-mapping.txt"}; !reflect.DeepEqual(merged.Unmatched.Mapping, want) {
		t.Fatalf("Unmatched.Mapping = %v, want %v", merged.Unmatched.Mapping, want)
	}
}

func TestReplaceFile(t *testing.T) {
	m, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		[]File{{DeployPath: deploy("app-demo-release.aab"), SourcePath: demoAABSource}},
		[]File{{DeployPath: deploy("mapping.txt"), SourcePath: demoMappingSource}},
	)

	if !m.ReplaceFile("app-demo-release.aab", "app-demo-release-bitrise-signed.aab") {
		t.Fatal("expected the AAB reference to be replaced")
	}
	if got := m.Modules["app"]["demoRelease"].AAB; !reflect.DeepEqual(got, []string{"app-demo-release-bitrise-signed.aab"}) {
		t.Fatalf("AAB = %v, want the signed name", got)
	}
	// mapping references are never touched
	if m.ReplaceFile("mapping.txt", "other.txt") {
		t.Fatal("mapping references must not be replaceable")
	}
	if m.ReplaceFile("no-such-file.apk", "x.apk") {
		t.Fatal("expected miss to report false")
	}
}

func TestReplaceFile_Unmatched(t *testing.T) {
	m := Map{Version: Version, Modules: map[string]map[string]Entry{}, Unmatched: Unmatched{
		APK: []string{"stray.apk"}, AAB: []string{}, Mapping: []string{},
	}}

	if !m.ReplaceFile("stray.apk", "stray-bitrise-signed.apk") {
		t.Fatal("expected the unmatched APK reference to be replaced")
	}
	if !reflect.DeepEqual(m.Unmatched.APK, []string{"stray-bitrise-signed.apk"}) {
		t.Fatalf("Unmatched.APK = %v", m.Unmatched.APK)
	}
}

func TestMarshal_MatchesWrite(t *testing.T) {
	m, _ := Build(
		[]File{{DeployPath: deploy("app-demo-release.apk"), SourcePath: demoAPKSource}},
		nil, nil,
	)
	data, err := Marshal(m)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatal("Marshal must render the trailing-newline document Write persists")
	}
}
