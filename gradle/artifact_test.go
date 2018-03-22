package gradle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractArtifactName(t *testing.T) {
	t.Log("simple repo")
	{
		proj := Project{
			location: "/root_dir",
			monoRepo: false,
		}
		got, err := proj.extractArtifactName(
			"/root_dir/mymodule/build/reports/myartifact.html")

		require.NoError(t, err)
		require.Equal(t, "mymodule-myartifact.html", got)
	}

	t.Log("mono repo")
	{
		proj := Project{
			location: "/root_dir",
			monoRepo: true,
		}
		got, err := proj.extractArtifactName(
			"/root_dir/mymodule/build/reports/myartifact.html")

		require.NoError(t, err)
		require.Equal(t, "root_dir-mymodule-myartifact.html", got)
	}

	t.Log("simple repo in root")
	{
		proj := Project{
			location: "/",
			monoRepo: false,
		}
		got, err := proj.extractArtifactName(
			"/mymodule/build/reports/myartifact.html")

		require.NoError(t, err)
		require.Equal(t, "mymodule-myartifact.html", got)
	}

	t.Log("mono repo in root")
	{
		proj := Project{
			location: "/",
			monoRepo: true,
		}
		got, err := proj.extractArtifactName(
			"/mymodule/build/reports/myartifact.html")

		require.NoError(t, err)
		require.Equal(t, "mymodule-myartifact.html", got)
	}
}
