/*
Copyright 2020,2021 Avi Zimmerman

This file is part of kvdi.

kvdi is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

kvdi is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with kvdi.  If not, see <https://www.gnu.org/licenses/>.
*/

package v1

import (
	"path"
	"strconv"

	appv1 "github.com/tinyzimmer/kvdi/apis/app/v1"
	v1 "github.com/tinyzimmer/kvdi/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
)

// TODO: Move
var (
	qemuBootImagePath  = path.Join(v1.DesktopRunPath, "boot.img")
	qemuCloudImagePath = path.Join(v1.DesktopRunPath, "cloud.img")

	qemuBootImageEnvVar  = "BOOT_IMAGE"
	qemuCloudImageEnvVar = "CLOUD_IMAGE"
	qemuCPUsEnvVar       = "CPUS"
	qemuMemoryEnvVar     = "MEMORY"
)

// IsQEMUTemplate returns true if this template is for a QEMU vm.
func (t *Template) IsQEMUTemplate() bool { return t.Spec.QEMUConfig != nil }

// GetQEMUContainer returns the container for launching the QEMU vm.
func (t *Template) GetQEMUContainer(cluster *appv1.VDICluster, instance *Session) corev1.Container {
	return corev1.Container{
		Name:            "qemu-kvm",
		Image:           t.GetQEMUImage(),
		ImagePullPolicy: t.GetQEMUImagePullPolicy(),
		VolumeMounts:    t.GetDesktopVolumeMounts(cluster, instance),
		Env: []corev1.EnvVar{
			{
				Name:  v1.VNCSockEnvVar,
				Value: t.GetDisplaySocketURI(),
			},
			{
				Name:  v1.UserEnvVar,
				Value: instance.GetUser(),
			},
			{
				Name:  v1.UIDEnvVar,
				Value: strconv.Itoa(int(v1.DefaultUser)),
			},
			{
				Name:  qemuBootImageEnvVar,
				Value: qemuBootImagePath,
			},
			{
				Name:  qemuCloudImageEnvVar,
				Value: qemuCloudImagePath,
			},
			{
				Name:  qemuCPUsEnvVar,
				Value: strconv.Itoa(t.GetQEMUNumCPUs()),
			},
			{
				Name:  qemuMemoryEnvVar,
				Value: strconv.Itoa(t.GetQEMUMemory()),
			},
		},
		SecurityContext: &corev1.SecurityContext{
			Privileged: &v1.True,
			RunAsUser:  &v1.DefaultUser,
		},
		Lifecycle: &corev1.Lifecycle{
			PreStop: &corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"/bin/sh", "-c", "echo system_powerdown | socat - UNIX-CONNECT:/run/qemu-monitor.sock"},
				},
			},
		},
	}
}

// GetQEMUImage returns the qemu utility image to use.
func (t *Template) GetQEMUImage() string {
	if t.Spec.QEMUConfig != nil && t.Spec.QEMUConfig.QEMUImage != "" {
		return t.Spec.QEMUConfig.QEMUImage
	}
	return "ghcr.io/tinyzimmer/kvdi:qemu-latest"
}

// GetQEMUImagePullPolicy returns the pull policy for the qemu utility image.
func (t *Template) GetQEMUImagePullPolicy() corev1.PullPolicy {
	if t.Spec.QEMUConfig != nil && t.Spec.QEMUConfig.QEMUImagePullPolicy != "" {
		return t.Spec.QEMUConfig.QEMUImagePullPolicy
	}
	return corev1.PullIfNotPresent
}

// GetQEMUDiskImage returns the qemu disk image to use.
func (t *Template) GetQEMUDiskImage() string {
	if t.Spec.QEMUConfig != nil {
		return t.Spec.QEMUConfig.DiskImage
	}
	return ""
}

// GetQEMUDiskImagePullPolicy returns the pull policy for the qemu disk image.
func (t *Template) GetQEMUDiskImagePullPolicy() corev1.PullPolicy {
	if t.Spec.QEMUConfig != nil && t.Spec.QEMUConfig.DiskImagePullPolicy != "" {
		return t.Spec.QEMUConfig.QEMUImagePullPolicy
	}
	return corev1.PullIfNotPresent
}

// GetQEMUDiskPath returns the path to the boot image inside the disk image.
func (t *Template) GetQEMUDiskPath() string {
	if t.Spec.QEMUConfig != nil && t.Spec.QEMUConfig.DiskPath != "" {
		return t.Spec.QEMUConfig.DiskPath
	}
	return "/disk/boot.img"
}

// GetQEMUCloudInitPath returns the path to the cloud image inside the disk image.
// An empty string means to generate one.
func (t *Template) GetQEMUCloudInitPath() string {
	if t.Spec.QEMUConfig != nil {
		return t.Spec.QEMUConfig.CloudInitPath
	}
	return ""
}

// GetQEMUNumCPUs returns the number of CPUs to use for the vm.
func (t *Template) GetQEMUNumCPUs() int {
	if t.Spec.QEMUConfig != nil && t.Spec.QEMUConfig.CPUs != 0 {
		return t.Spec.QEMUConfig.CPUs
	}
	return 1
}

// GetQEMUMemory returns the amount of memory to use for the vm.
func (t *Template) GetQEMUMemory() int {
	if t.Spec.QEMUConfig != nil && t.Spec.QEMUConfig.Memory != 0 {
		return t.Spec.QEMUConfig.Memory
	}
	return 1024
}
