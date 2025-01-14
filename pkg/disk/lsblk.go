package disk

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xib1uvXi/armboxkit/pkg/shutil"
)

const (
	// CmdTmpl adds device name, if add empty string - command will print info about all devices
	CmdTmpl = "lsblk %s --paths --json --bytes --fs " +
		"--output NAME,TYPE,SIZE,ROTA,MOUNTPOINT,FSTYPE,PARTUUID,LABEL,UUID"

	// outputKey is the key to find block devices in lsblk json output
	outputKey = "blockdevices"
)

// CustomInt64 to handle Size lsblk output - 8001563222016 or "8001563222016"
type CustomInt64 struct {
	Int64 int64
}

// UnmarshalJSON customizes string size unmarshalling
func (ci *CustomInt64) UnmarshalJSON(data []byte) error {
	QuotesByte := byte(34)
	if data[0] == QuotesByte {
		err := json.Unmarshal(data[1:len(data)-1], &ci.Int64)
		if err != nil {
			return errors.New("CustomInt64: UnmarshalJSON: " + err.Error())
		}
	} else {
		err := json.Unmarshal(data, &ci.Int64)
		if err != nil {
			return errors.New("CustomInt64: UnmarshalJSON: " + err.Error())
		}
	}
	return nil
}

// MarshalJSON customizes size marshalling
func (ci *CustomInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(ci.Int64)
}

// CustomBool to handle Rota lsblk output - true/false or "1"/"0"
type CustomBool struct {
	Bool bool
}

// UnmarshalJSON customizes string rota unmarshalling
func (cb *CustomBool) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"true"`, `true`, `"1"`, `1`:
		cb.Bool = true
		return nil
	case `"false"`, `false`, `"0"`, `0`, `""`:
		cb.Bool = false
		return nil
	default:
		return errors.New("CustomBool: parsing \"" + string(data) + "\": unknown value")
	}
}

// MarshalJSON customizes rota marshalling
func (cb CustomBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(cb.Bool)
}

type BlockDevice struct {
	Name       string        `json:"name,omitempty"`
	Type       string        `json:"type,omitempty"`
	Size       *CustomInt64  `json:"size,omitempty"`
	Rota       *CustomBool   `json:"rota,omitempty"`
	MountPoint string        `json:"mountpoint,omitempty"`
	FSType     string        `json:"fstype,omitempty"`
	PartUUID   string        `json:"partuuid,omitempty"`
	Label      string        `json:"label,omitempty"`
	UUID       string        `json:"uuid,omitempty"`
	Children   []BlockDevice `json:"children,omitempty"`
}

type Lsblk struct {
	executor *shutil.CmdExecutor
}

func NewLsblk() *Lsblk {
	return &Lsblk{executor: shutil.NewCmdExecutor()}
}

// GetBlockDevices run os lsblk command for device and construct BlockDevice struct based on output
// Receives device path. If device is empty string, info about all devices will be collected
// Returns slice of BlockDevice structs or error if something went wrong
func (c *Lsblk) GetBlockDevices(device string) ([]BlockDevice, error) {
	cmdStr := fmt.Sprintf(CmdTmpl, device)
	stdOut, _, err := c.executor.ExecuteCmd(cmdStr)
	if err != nil {
		return nil, err
	}

	rawOut := make(map[string][]BlockDevice, 1)
	if err := json.Unmarshal([]byte(stdOut), &rawOut); err != nil {
		return nil, fmt.Errorf("unable to unmarshal output to BlockDevice instance, error: %v", err)
	}

	var (
		devs []BlockDevice
		ok   bool
	)
	if devs, ok = rawOut[outputKey]; !ok {
		return nil, fmt.Errorf("unexpected lsblk output format, missing \"%s\" key", outputKey)
	}

	return GetLeafBlockDevices(devs), nil
}

func GetLeafBlockDevices(devs []BlockDevice) []BlockDevice {
	result := make([]BlockDevice, 0)
	for _, dev := range devs {
		if len(dev.Children) == 0 {
			result = append(result, dev)
		} else {
			leafs := GetLeafBlockDevices(dev.Children)
			result = append(result, leafs...)
		}
	}
	return result
}

// SearchBlockDevice searches for a block device by fstype, lable, mountpoint
func (c *Lsblk) SearchBlockDevice(fstype, label string, hasMountPoint bool) ([]BlockDevice, error) {
	devices, err := c.GetBlockDevices("")
	if err != nil {
		return nil, err
	}

	var res []BlockDevice = make([]BlockDevice, 0)
	for _, dev := range devices {
		if dev.FSType == fstype && dev.Label == label {
			if hasMountPoint && dev.MountPoint != "" {
				res = append(res, dev)
			} else if !hasMountPoint {
				res = append(res, dev)
			}
		}
	}

	return res, nil
}
