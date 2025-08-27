package sdkmanager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-android/v2/sdk"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		sdkLayout  map[string]string
		wantLegacy bool
		wantPath   string
		wantErr    bool
	}{
		{
			name: "Legacy tools",
			sdkLayout: map[string]string{
				"tools": "android",
			},
			wantLegacy: true,
			wantPath:   filepath.Join("tools", "android"),
		},
		{
			name: "bin/tools",
			sdkLayout: map[string]string{
				filepath.Join("tools", "bin"): "sdkmanager",
			},
			wantPath: filepath.Join("tools", "bin", "sdkmanager"),
		},
		{
			name: "cmdline-tools",
			sdkLayout: map[string]string{
				filepath.Join("tools"):                          "android",
				filepath.Join("cmdline-tools", "latest", "bin"): "sdkmanager",
			},
			wantPath: filepath.Join("cmdline-tools", "latest", "bin", "sdkmanager"),
		},
		{
			name: "",
			sdkLayout: map[string]string{
				filepath.Join("cmdline-tools", "latest", "bin"): "other_tool",
				filepath.Join("tools", "bin"):                   "sdkmanager",
			},
			wantPath: "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sdkRoot, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatalf("failed to eval symlink: %v", err)
			}

			for dir, file := range tt.sdkLayout {
				if err := os.MkdirAll(filepath.Join(sdkRoot, dir), 0700); err != nil {
					t.Fatalf("failed to create directory: %s", err)
				}

				if err := os.WriteFile(filepath.Join(sdkRoot, dir, file), []byte{}, 0600); err != nil {
					t.Fatalf("failed to create file: %s", err)
				}
			}

			sdk, err := sdk.New(sdkRoot)
			if err != nil {
				t.Fatalf("failed to create sdk: %v", err)
			}

			want := &Model{
				androidHome: sdkRoot,
				binPth:      filepath.Join(sdkRoot, tt.wantPath),
				legacy:      tt.wantLegacy,
				cmdFactory:  command.NewFactory(env.NewRepository()),
			}

			got, err := New(sdk, command.NewFactory(env.NewRepository()))
			if tt.wantErr {
				require.Error(t, err, "New()")
				require.Nil(t, got)
			} else {
				require.NoError(t, err, "New()")
				require.Equal(t, want, got, "New()")
			}
		})
	}
}
