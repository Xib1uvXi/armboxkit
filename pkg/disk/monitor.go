package disk

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v4/disk"
	"sync"
	"time"
)

type Info struct {
	Device     string `json:"device"`
	ReadData   uint64 `json:"read_data"`
	WriteData  uint64 `json:"write_data"`
	ReadSpeed  uint64 `json:"read_speed"`
	WriteSpeed uint64 `json:"write_speed"`
}

func (i Info) String() string {
	return fmt.Sprintf("Device: %s, ReadData: %s, WriteData: %s, ReadSpeed: %s/s, WriteSpeed: %s/s", i.Device,
		humanize.Bytes(i.ReadData), humanize.Bytes(i.WriteData), humanize.Bytes(i.ReadSpeed), humanize.Bytes(i.WriteSpeed))
}

type Monitor struct {
	locker    sync.RWMutex
	infos     map[string]*Info
	devices   []string
	closeChan chan struct{}
}

func NewMonitor() (*Monitor, error) {
	mt := &Monitor{
		devices:   make([]string, 0),
		infos:     make(map[string]*Info),
		closeChan: make(chan struct{}),
	}

	ps, err := Ext4Partitions()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		mt.devices = append(mt.devices, p.Device)
	}

	currStats, err := disk.IOCounters(mt.devices...)
	if err != nil {
		return nil, err
	}

	for _, stat := range currStats {
		device := fmt.Sprintf("/dev/%s", stat.Name)
		mt.infos[device] = &Info{
			Device:    device,
			ReadData:  stat.ReadBytes,
			WriteData: stat.WriteBytes,
		}
	}

	return mt, nil
}

func (m *Monitor) GetInfos() []*Info {
	m.locker.RLock()
	defer m.locker.RUnlock()

	infos := make([]*Info, 0, len(m.infos))
	for _, info := range m.infos {
		infos = append(infos, info)
	}

	return infos
}

func (m *Monitor) GetInfo(name string) *Info {
	m.locker.RLock()
	defer m.locker.RUnlock()

	if info, ok := m.infos[name]; ok {
		return info
	}

	return nil
}

func (m *Monitor) Close() {
	close(m.closeChan)
}

func (m *Monitor) loop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.closeChan:
			return

		case <-ticker.C:
			m.update()
		}
	}
}

func (m *Monitor) update() {
	m.locker.Lock()
	defer m.locker.Unlock()

	currStats, err := disk.IOCounters(m.devices...)
	if err != nil {
		return
	}

	for _, stat := range currStats {
		device := fmt.Sprintf("/dev/%s", stat.Name)
		if info, ok := m.infos[device]; ok {
			readSpeed := (stat.ReadBytes - info.ReadData) / 2
			writeSpeed := (stat.WriteBytes - info.WriteData) / 2

			info.ReadSpeed = readSpeed
			info.WriteSpeed = writeSpeed
			info.ReadData = stat.ReadBytes
			info.WriteData = stat.WriteBytes
		}
	}
}
