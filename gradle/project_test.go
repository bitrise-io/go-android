package gradle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractArtifactName(t *testing.T) {
	tests := []struct {
		name                string
		projectLocation     string
		monoRepo            bool
		path                string
		includeModuleInName bool
		want                string
	}{
		{
			name:                "top-level module without includeModule",
			projectLocation:     "/repo",
			path:                "/repo/app/build/outputs/apk/debug/app-debug.apk",
			includeModuleInName: false,
			want:                "app-debug.apk",
		},
		{
			name:                "top-level module with includeModule",
			projectLocation:     "/repo",
			path:                "/repo/app/build/outputs/apk/debug/app-debug.apk",
			includeModuleInName: true,
			want:                "app-app-debug.apk",
		},
		{
			name:                "nested module with includeModule",
			projectLocation:     "/repo",
			path:                "/repo/feature/feature-name-1/data/build/outputs/aar/data-debug.aar",
			includeModuleInName: true,
			want:                "feature-feature-name-1-data-data-debug.aar",
		},
		{
			name:                "nested module without includeModule",
			projectLocation:     "/repo",
			path:                "/repo/feature/feature-name-1/data/build/outputs/aar/data-debug.aar",
			includeModuleInName: false,
			want:                "data-debug.aar",
		},
		{
			name:                "monorepo top-level module",
			projectLocation:     "/repo/android",
			monoRepo:            true,
			path:                "/repo/android/app/build/outputs/apk/debug/app-debug.apk",
			includeModuleInName: true,
			want:                "android-app-app-debug.apk",
		},
		{
			name:                "monorepo nested module",
			projectLocation:     "/repo/android",
			monoRepo:            true,
			path:                "/repo/android/feature/auth/data/build/outputs/aar/data-release.aar",
			includeModuleInName: true,
			want:                "android-feature-auth-data-data-release.aar",
		},
		{
			name:                "artifact not under /build/ falls back to top-level segment",
			projectLocation:     "/repo",
			path:                "/repo/app/extra/foo.txt",
			includeModuleInName: true,
			want:                "app-foo.txt",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proj := Project{location: tc.projectLocation, monoRepo: tc.monoRepo}
			got, err := proj.extractArtifactName(tc.path, tc.includeModuleInName)
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
