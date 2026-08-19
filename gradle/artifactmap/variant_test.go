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

// TestAARVariantFromPath covers what AGP encodes in an archive's file name
// instead of its directory: the build type alone for an unflavored library, the
// flavors plus build type when the library has flavor dimensions.
func TestAARVariantFromPath(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		want   ArtifactVariant
		wantOK bool
	}{
		{
			name:   "no flavor: build type only",
			path:   "/bitrise/src/mylib/build/outputs/aar/mylib-debug.aar",
			want:   ArtifactVariant{Module: "mylib", ModulePath: "/bitrise/src/mylib", Variant: "debug"},
			wantOK: true,
		},
		{
			name:   "one flavor dimension",
			path:   "/bitrise/src/mylib/build/outputs/aar/mylib-free-release.aar",
			want:   ArtifactVariant{Module: "mylib", ModulePath: "/bitrise/src/mylib", Variant: "freeRelease"},
			wantOK: true,
		},
		{
			name:   "two flavor dimensions",
			path:   "/bitrise/src/mylib/build/outputs/aar/mylib-free-arm-release.aar",
			want:   ArtifactVariant{Module: "mylib", ModulePath: "/bitrise/src/mylib", Variant: "freeArmRelease"},
			wantOK: true,
		},
		{
			name:   "nested module keeps its full path",
			path:   "/bitrise/src/feature/feature-name-1/data/build/outputs/aar/data-paid-release.aar",
			want:   ArtifactVariant{Module: "data", ModulePath: "/bitrise/src/feature/feature-name-1/data", Variant: "paidRelease"},
			wantOK: true,
		},
		{
			name:   "customised base name with a version is not decoded",
			path:   "/bitrise/src/mylib/build/outputs/aar/mylib-1.0-release.aar",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "base name unrelated to the module dir is not decoded",
			path:   "/bitrise/src/mylib/build/outputs/aar/renamed-release.aar",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "bare module name carries no variant",
			path:   "/bitrise/src/mylib/build/outputs/aar/mylib.aar",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "nesting under outputs/aar is not an AGP layout",
			path:   "/bitrise/src/mylib/build/outputs/aar/release/mylib-release.aar",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "outside build/outputs/aar",
			path:   "/bitrise/src/custom-out/mylib-release.aar",
			want:   ArtifactVariant{},
			wantOK: false,
		},
		{
			name:   "a checkout directory named outputs cannot hijack parsing",
			path:   "/bitrise/src/outputs/aar/tools/mylib/build/outputs/aar/mylib-release.aar",
			want:   ArtifactVariant{Module: "mylib", ModulePath: "/bitrise/src/outputs/aar/tools/mylib", Variant: "release"},
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := AARVariantFromPath(tt.path)
			if ok != tt.wantOK {
				t.Fatalf("AARVariantFromPath(%q) ok = %v, want %v", tt.path, ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("AARVariantFromPath(%q) = %+v, want %+v", tt.path, got, tt.want)
			}
		})
	}
}

// TestAARVariantFromPath_PairsAARWithMapping: a self-minifying library writes
// its mapping to the standard outputs/mapping/<variant>/ tree, so the archive
// and that mapping must resolve to equal ArtifactVariants — the same guarantee
// apps get, across the two differently-shaped layouts.
func TestAARVariantFromPath_PairsAARWithMapping(t *testing.T) {
	aarVariant, ok := AARVariantFromPath("/bitrise/src/feature/data/build/outputs/aar/data-free-release.aar")
	if !ok {
		t.Fatal("expected the AAR path to resolve a variant")
	}
	mappingVariant, ok := VariantFromPath("/bitrise/src/feature/data/build/outputs/mapping/freeRelease/mapping.txt")
	if !ok {
		t.Fatal("expected the mapping path to resolve a variant")
	}
	if aarVariant != mappingVariant {
		t.Fatalf("AAR variant %+v and mapping variant %+v should match", aarVariant, mappingVariant)
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
