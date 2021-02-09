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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DesktopInit represents the init system that the desktop container uses.
// +kubebuilder:validation:Enum=supervisord;systemd
type DesktopInit string

const (
	// InitSupervisord signals that the image uses supervisord.
	InitSupervisord = "supervisord"
	// InitSystemd signals that the image uses systemd.
	InitSystemd = "systemd"
)

// TemplateSpec defines the desired state of Template
type TemplateSpec struct {
	// The docker repository and tag to use for desktops booted from this template.
	Image string `json:"image"`
	// The pull policy to use when pulling the container image.
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`
	// Any pull secrets required for pulling the container image.
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Additional environment variables to pass to containers booted from this template.
	Env []corev1.EnvVar `json:"env,omitempty"`
	// Optionally map additional information about the user (and potentially extended further
	// in the future) into the environment of desktops booted from this template. The keys in the
	// map are the environment variable to set inside the desktop, and the values are go templates
	// or strings to set to the value. Currently the go templates are only passed a `Session` object
	// containing the information in the claims for the user that created the desktop. For more information
	// see the [JWTCLaims object](./metav1.md#JWTClaims) and corresponding go types.
	EnvTemplates map[string]string `json:"envTemplates,omitempty"`
	// Resource requirements to apply to desktops booted from this template.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// Configuration options for the instances. These are highly dependant on using
	// the Dockerfiles (or close derivitives) provided in this repository.
	Config *DesktopConfig `json:"config,omitempty"`
	// Volume configurations for the instances. These can be used for mounting custom
	// volumes at arbitrary paths in desktops.
	VolumeConfig *DesktopVolumeConfig `json:"volumeConfig,omitempty"`
	// Arbitrary tags for displaying in the app UI.
	Tags map[string]string `json:"tags,omitempty"`
}

// DesktopConfig represents configurations for the template and desktops booted
// from it.
type DesktopConfig struct {
	// Extra system capabilities to add to desktops booted from this template.
	Capabilities []corev1.Capability `json:"capabilities,omitempty"`
	// AllowRoot will pass the ENABLE_ROOT envvar to the container. In the Dockerfiles
	// in this repository, this will add the user to the sudo group and ability to
	// sudo with no password.
	AllowRoot bool `json:"allowRoot,omitempty"`
	// The address the VNC server listens on inside the image. This defaults to the
	// UNIX socket /var/run/kvdi/display.sock. The kvdi-proxy sidecar will forward
	// websockify requests validated by mTLS to this socket.
	// Must be in the format of `tcp://{host}:{port}` or `unix://{path}`.
	SocketAddr string `json:"socketAddr,omitempty"`
	// Override the address of the PulseAudio server that the proxy will try to connect to
	// when serving audio. This defaults to what the ubuntu/arch desktop images are configured
	// to do during init. The value is assumed to be a unix socket.
	PulseServer string `json:"pulseServer,omitempty"`
	// AllowFileTransfer will mount the user's home directory inside the kvdi-proxy image.
	// This enables the API endpoint for exploring, downloading, and uploading files to
	// desktop sessions booted from this template.
	AllowFileTransfer bool `json:"allowFileTransfer,omitempty"`
	// The image to use for the sidecar that proxies mTLS connections to the local
	// VNC server inside the Desktop. Defaults to the public kvdi-proxy image
	// matching the version of the currrently running manager.
	ProxyImage string `json:"proxyImage,omitempty"`
	// The type of init system inside the image, currently only `supervisord` and `systemd`
	// are supported. Defaults to `systemd`. `systemd` containers are run privileged and
	// downgrading to the desktop user must be done within the image's init process. `supervisord`
	// containers are run with minimal capabilities and directly as the desktop user.
	Init DesktopInit `json:"init,omitempty"`
}

// DesktopVolumeConfig represents configurations for volumes attached to pods booted from
// a template.
type DesktopVolumeConfig struct {
	// Additional volumes to attach to pods booted from this template. To mount them there
	// must be cooresponding `volumeMounts` or `volumeDevices` specified.
	Volumes []corev1.Volume `json:"volumes,omitempty"`
	// Volume mounts for the desktop container.
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	// Volume devices for the desktop container.
	VolumeDevices []corev1.VolumeDevice `json:"volumeDevices,omitempty"`
}

// TemplateStatus defines the observed state of Template
type TemplateStatus struct{}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=templates,scope=Cluster
//+kubebuilder:subresource:status

// Template is the Schema for the templates API
type Template struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateSpec   `json:"spec,omitempty"`
	Status TemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TemplateList contains a list of Template
type TemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Template `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Template{}, &TemplateList{})
}