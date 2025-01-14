package net

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v4/net"
	"sync"
	"time"
)

type Info struct {
	Interface string `json:"interface"`
	RecvData  uint64 `json:"recv_data"`
	SendData  uint64 `json:"send_data"`
	RecvSpeed uint64 `json:"recv_speed"`
	SendSpeed uint64 `json:"send_speed"`
}

func (i Info) String() string {
	return fmt.Sprintf("Interface: %s, RecvData: %s, SendData: %s, RecvSpeed: %s/s, SendSpeed: %s/s", i.Interface,
		humanize.Bytes(i.RecvData), humanize.Bytes(i.SendData), humanize.Bytes(i.RecvSpeed), humanize.Bytes(i.SendSpeed))
}

type Monitor struct {
	locker    sync.RWMutex
	infos     map[string]*Info
	closeChan chan struct{}
}

func NewMonitor() (*Monitor, error) {
	mt := &Monitor{
		infos:     make(map[string]*Info),
		closeChan: make(chan struct{}),
	}

	currStats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	for _, stat := range currStats {
		mt.infos[stat.Name] = &Info{
			Interface: stat.Name,
			RecvData:  stat.BytesRecv,
			SendData:  stat.BytesSent,
		}
	}

	go mt.loop()

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

	currStats, err := net.IOCounters(true)
	if err != nil {
		return
	}

	for _, stat := range currStats {
		if info, ok := m.infos[stat.Name]; ok {
			info.RecvSpeed = (stat.BytesRecv - info.RecvData) / 2
			info.SendSpeed = (stat.BytesSent - info.SendData) / 2
			info.RecvData = stat.BytesRecv
			info.SendData = stat.BytesSent
		}
	}
}
