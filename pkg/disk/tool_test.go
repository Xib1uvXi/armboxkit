package disk

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExt4PartitionsByLabel(t *testing.T) {
	ps, err := Ext4PartitionsByLabel("hificloud")
	require.NoError(t, err)
	require.True(t, len(ps) > 0)

	t.Log(ps)
}

func TestExt4Partitions(t *testing.T) {
	ps, err := Ext4Partitions()
	require.NoError(t, err)
	require.True(t, len(ps) > 0)

	t.Log(ps)
}
