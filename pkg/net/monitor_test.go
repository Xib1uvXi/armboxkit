package net

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewMonitor(t *testing.T) {
	monitor, err := NewMonitor()
	require.NoError(t, err)
	defer monitor.Close()

	go func() {
		for i := 0; i < 10; i++ {
			infos := monitor.GetInfos()
			t.Logf("infos: %v", infos)
			time.Sleep(3 * time.Second)
		}
	}()

	time.Sleep(11 * time.Second)
}
