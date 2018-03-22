package gradle

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetGradleModule(t *testing.T) {
	require.Equal(t, "", getGradleModule(""))
	require.Equal(t, ":app:", getGradleModule("app"))
}
