package disk

import "github.com/shirou/gopsutil/v4/disk"

func Ext4Partitions() ([]disk.PartitionStat, error) {
	raw, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var ext4s []disk.PartitionStat
	for _, r := range raw {
		if r.Fstype == "ext4" {
			ext4s = append(ext4s, r)
		}
	}

	return ext4s, nil
}

func Ext4PartitionsByLabel(label string) ([]disk.PartitionStat, error) {
	lsblk := NewLsblk()

	bds, err := lsblk.SearchBlockDevice("ext4", label, true)
	if err != nil {
		return nil, err
	}

	raw, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var ext4s []disk.PartitionStat
	for _, r := range raw {
		if r.Fstype == "ext4" {
			for _, bd := range bds {
				if r.Device == bd.Name {
					ext4s = append(ext4s, r)
				}
			}
		}
	}

	return ext4s, nil
}
