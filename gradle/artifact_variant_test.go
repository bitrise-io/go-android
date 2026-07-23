package gradle

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
			want:   ArtifactVariant{Module: "app", Variant: "demoRelease"},
			wantOK: true,
		},
		{
			name:   "APK with split flavor and build type segments",
			path:   "/bitrise/src/app/build/outputs/apk/demo/release/app-demo-release.apk",
			want:   ArtifactVariant{Module: "app", Variant: "demoRelease"},
			wantOK: true,
		},
		{
			name:   "APK without flavor",
			path:   "/bitrise/src/app/build/outputs/apk/release/app-release.apk",
			want:   ArtifactVariant{Module: "app", Variant: "release"},
			wantOK: true,
		},
		{
			name:   "mapping under outputs",
			path:   "/bitrise/src/app/build/outputs/mapping/demoRelease/mapping.txt",
			want:   ArtifactVariant{Module: "app", Variant: "demoRelease"},
			wantOK: true,
		},
		{
			name:   "mapping under intermediates with minify subdir",
			path:   "/bitrise/src/app/build/intermediates/mapping/productionRelease/minifyProductionReleaseWithR8/mapping.txt",
			want:   ArtifactVariant{Module: "app", Variant: "productionRelease"},
			wantOK: true,
		},
		{
			name:   "different module stays distinct",
			path:   "/bitrise/src/feature/build/outputs/apk/release/feature-release.apk",
			want:   ArtifactVariant{Module: "feature", Variant: "release"},
			wantOK: true,
		},
		{
			name:   "multi dimension flavor merges consistently with mapping",
			path:   "/bitrise/src/app/build/outputs/apk/demoFree/release/app.apk",
			want:   ArtifactVariant{Module: "app", Variant: "demoFreeRelease"},
			wantOK: true,
		},
		{
			name:   "flavor named like the kind marker is not mistaken for it",
			path:   "/bitrise/src/app/build/outputs/apk/apk/release/app.apk",
			want:   ArtifactVariant{Module: "app", Variant: "apkRelease"},
			wantOK: true,
		},
		{
			name:   "unrecognised path",
			path:   "/bitrise/src/some/random/file.apk",
			want:   ArtifactVariant{},
			wantOK: false,
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

// TestArtifact_Variant_PairsAPKWithMapping is the load-bearing guarantee: an APK
// and the mapping file built for the same variant, laid out by AGP in their
// differently-shaped directories, must resolve to equal ArtifactVariants so they
// can be paired.
func TestArtifact_Variant_PairsAPKWithMapping(t *testing.T) {
	apk := Artifact{Path: "/bitrise/src/app/build/outputs/apk/demo/release/app-demo-release.apk"}
	mapping := Artifact{Path: "/bitrise/src/app/build/outputs/mapping/demoRelease/mapping.txt"}

	apkVariant, ok := apk.Variant()
	if !ok {
		t.Fatal("expected APK path to resolve a variant")
	}
	mappingVariant, ok := mapping.Variant()
	if !ok {
		t.Fatal("expected mapping path to resolve a variant")
	}
	if apkVariant != mappingVariant {
		t.Fatalf("APK variant %+v and mapping variant %+v should match", apkVariant, mappingVariant)
	}
}
