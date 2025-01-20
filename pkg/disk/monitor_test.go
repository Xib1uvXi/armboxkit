package disk

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	monitor, err := NewMonitor()
	require.NoError(t, err)

	go func() {
		for i := 0; i < 10; i++ {
			infos := monitor.GetInfos()
			t.Logf("infos: %v", infos)
			time.Sleep(3 * time.Second)
		}
	}()

	time.Sleep(30 * time.Second)
}

func TestNewMonitor2(t *testing.T) {
	monitor, err := NewMonitor()
	require.NoError(t, err)

	go func() {
		for i := 0; i < 10; i++ {
			info, err := monitor.GetInfoByPath("/mnt/data/codespace")
			require.NoError(t, err)
			t.Logf("infos: %v", info)
			time.Sleep(3 * time.Second)
		}
	}()

	time.Sleep(30 * time.Second)
}
