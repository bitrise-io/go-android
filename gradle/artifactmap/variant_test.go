package artifactmap

import "testing"

func TestVariantFromPath(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		want   ArtifactVariant
		wantOK bool
	}{
		{
			name:   "AAB with merged variant segment",
			path:   "/bitrise/src/app/build/outputs/bundle/demoRelease/app-demo-release.aab",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "demoRelease"},
			wantOK: true,
		},
		{
			name:   "APK with split flavor and build type segments",
			path:   "/bitrise/src/app/build/outputs/apk/demo/release/app-demo-release.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "demoRelease"},
			wantOK: true,
		},
		{
			name:   "APK without flavor",
			path:   "/bitrise/src/app/build/outputs/apk/release/app-release.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "release"},
			wantOK: true,
		},
		{
			name:   "mapping under outputs",
			path:   "/bitrise/src/app/build/outputs/mapping/demoRelease/mapping.txt",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "demoRelease"},
			wantOK: true,
		},
		{
			// only official build/outputs/ paths pair; task-workdir copies
			// stay unmatched so they can never displace the outputs copy
			name:   "mapping under intermediates is not recognised",
			path:   "/bitrise/src/app/build/intermediates/mapping/productionRelease/minifyProductionReleaseWithR8/mapping.txt",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "different module stays distinct",
			path:   "/bitrise/src/feature/build/outputs/apk/release/feature-release.apk",
			want:   ArtifactVariant{Module: "feature", ModulePath: "/bitrise/src/feature", Variant: "release"},
			wantOK: true,
		},
		{
			name:   "multi dimension flavor merges consistently with mapping",
			path:   "/bitrise/src/app/build/outputs/apk/demoFree/release/app.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "demoFreeRelease"},
			wantOK: true,
		},
		{
			name:   "flavor named like the kind marker is not mistaken for it",
			path:   "/bitrise/src/app/build/outputs/apk/apk/release/app.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "apkRelease"},
			wantOK: true,
		},
		{
			name:   "unrecognised path",
			path:   "/bitrise/src/some/random/file.apk",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			// regression: this used to panic (slice bounds out of range)
			name:   "file named like the kind marker directly under outputs",
			path:   "/x/outputs/apk",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "file named bundle directly under intermediates",
			path:   "/x/intermediates/bundle",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			// regression: the checkout's own outputs/mapping/... prefix used to
			// hijack parsing; the marker closest to the file must win
			name:   "outputs directory in the checkout prefix does not hijack",
			path:   "/bitrise/src/outputs/mapping/tools/app/build/outputs/apk/release/app.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/outputs/mapping/tools/app", Variant: "release"},
			wantOK: true,
		},
		{
			name:   "ProGuard-era split flavor mapping merges like its APK",
			path:   "/bitrise/src/app/build/outputs/mapping/demo/release/mapping.txt",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "demoRelease"},
			wantOK: true,
		},
		{
			name:   "APK under intermediates is not recognised",
			path:   "/bitrise/src/app/build/intermediates/apk/demo/release/app.apk",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "universal APK built from the bundle (AGP 7)",
			path:   "/bitrise/src/app/build/outputs/universal_apk/release/app-release-universal.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "release"},
			wantOK: true,
		},
		{
			name:   "non-ASCII flavor keeps its rune intact",
			path:   "/bitrise/src/app/build/outputs/apk/ünnep/release/app.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "ünnepRelease"},
			wantOK: true,
		},
		{
			name:   "non-ASCII build type is capitalized as a rune",
			path:   "/bitrise/src/app/build/outputs/apk/demo/ürelease/app.apk",
			want:   ArtifactVariant{Module: "app", ModulePath: "/bitrise/src/app", Variant: "demoÜrelease"},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := VariantFromPath(tt.path)
			if ok != tt.wantOK {
				t.Fatalf("VariantFromPath(%q) ok = %v, want %v", tt.path, ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("VariantFromPath(%q) = %+v, want %+v", tt.path, got, tt.want)
			}
		})
	}
}

// TestVariantFromPath_PairsAPKWithMapping is the load-bearing guarantee: an APK
// and the mapping file built for the same variant, laid out by AGP in their
// differently-shaped directories, must resolve to equal ArtifactVariants so they
// can be paired.
func TestVariantFromPath_PairsAPKWithMapping(t *testing.T) {
	apkVariant, ok := VariantFromPath("/bitrise/src/app/build/outputs/apk/demo/release/app-demo-release.apk")
	if !ok {
		t.Fatal("expected APK path to resolve a variant")
	}
	mappingVariant, ok := VariantFromPath("/bitrise/src/app/build/outputs/mapping/demoRelease/mapping.txt")
	if !ok {
		t.Fatal("expected mapping path to resolve a variant")
	}
	if apkVariant != mappingVariant {
		t.Fatalf("APK variant %+v and mapping variant %+v should match", apkVariant, mappingVariant)
	}
}
