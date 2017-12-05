// +build windows

package daemon

import (
	"sort"

	"github.com/hyperhq/hypercli/container"
	"github.com/hyperhq/hypercli/daemon/execdriver"
	derr "github.com/hyperhq/hypercli/errors"
	"github.com/hyperhq/hypercli/volume"
)

// setupMounts configures the mount points for a container by appending each
// of the configured mounts on the container to the execdriver mount structure
// which will ultimately be passed into the exec driver during container creation.
// It also ensures each of the mounts are lexographically sorted.
func (daemon *Daemon) setupMounts(container *container.Container) ([]execdriver.Mount, error) {
	var mnts []execdriver.Mount
	for _, mount := range container.MountPoints { // type is volume.MountPoint
		if err := daemon.lazyInitializeVolume(container.ID, mount); err != nil {
			return nil, err
		}
		// If there is no source, take it from the volume path
		s := mount.Source
		if s == "" && mount.Volume != nil {
			s = mount.Volume.Path()
		}
		if s == "" {
			return nil, derr.ErrorCodeVolumeNoSourceForMount.WithArgs(mount.Name, mount.Driver, mount.Destination)
		}
		mnts = append(mnts, execdriver.Mount{
			Source:      s,
			Destination: mount.Destination,
			Writable:    mount.RW,
		})
	}

	sort.Sort(mounts(mnts))
	return mnts, nil
}

// setBindModeIfNull is platform specific processing which is a no-op on
// Windows.
func setBindModeIfNull(bind *volume.MountPoint) *volume.MountPoint {
	return bind
}
