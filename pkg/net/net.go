package net

import (
	"errors"
	"fmt"
	"github.com/ccding/go-stun/stun"
	"net"
	"time"
)

type BaseInfoOption func(*BaseInfo)

func SetNatCheckInterval(interval time.Duration) BaseInfoOption {
	return func(s *BaseInfo) {
		s.natCheckInterval = interval
	}
}

type BaseInfo struct {
	stunServer string

	publicIP  string
	privateIP string
	natType   int

	natCheckInterval time.Duration
	closeCh          chan struct{}
}

func NewBaseInfo(stunServer string, opts ...BaseInfoOption) (*BaseInfo, error) {
	s := &BaseInfo{
		stunServer:       stunServer,
		closeCh:          make(chan struct{}),
		natCheckInterval: 4 * time.Hour,
	}

	privateIP, err := s.initPrivateIP()
	if err != nil {
		return nil, err
	}

	s.privateIP = privateIP

	err = s.getNATType()
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(s)
	}

	go s.natCheckLoop()

	return s, nil
}

func (s *BaseInfo) Close() {
	close(s.closeCh)
}

func (s *BaseInfo) natCheckLoop() {
	ticker := time.NewTicker(s.natCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_ = s.getNATType()
		case <-s.closeCh:
			return
		}
	}
}

func (s *BaseInfo) GetPublicIP() string {
	return s.publicIP
}

func (s *BaseInfo) GetPrivateIP() string {
	return s.privateIP
}

func (s *BaseInfo) GetNATType() int {
	return s.natType
}

func (s *BaseInfo) initPrivateIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("private ip not found")
}

func (s *BaseInfo) getNATType() error {
	client := stun.NewClient()
	client.SetServerAddr(s.stunServer)

	nat, host, err := client.Discover()
	if err != nil {
		return err
	}

	s.publicIP = host.IP()
	nresult, err := s.dtoNATType(nat)
	if err != nil {
		return err
	}

	s.natType = int(nresult)
	return nil
}

func (s *BaseInfo) dtoNATType(sNatType stun.NATType) (NATType, error) {
	switch sNatType {
	case stun.NATNone:
		return None, nil
	case stun.NATFull:
		return FullCone, nil
	case stun.NATRestricted:
		return RestrictedCone, nil
	case stun.NATPortRestricted:
		return PortRestrictedCone, nil
	case stun.NATSymetric:
		return Symmetric, nil
	default:
		return UnKnown, fmt.Errorf("unsupported NAT type: %v", sNatType)
	}
}
