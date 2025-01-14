package disk

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLsblk_GetBlockDevices(t *testing.T) {
	lsblk := NewLsblk()
	devices, err := lsblk.GetBlockDevices("")
	require.NoError(t, err)
	require.True(t, len(devices) > 0)
}

func TestLsblk_SearchBlockDevice(t *testing.T) {
	lsblk := NewLsblk()

	bds, err := lsblk.SearchBlockDevice("ext4", "hificloud", true)

	require.NoError(t, err)
	t.Logf("BlockDevices: %v", bds)
}
