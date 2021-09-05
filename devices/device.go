package devices

import (
	"dev-play/utils"
	"github.com/fsnotify/fsnotify"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type Device interface {
	GetDeviceParts() []*plugin.Device
	GetDeviceStatus() utils.DeviceStatus
	GetDeviceLoads() int
	WatchDevice() error
	GetDeviceChangeNotifier() *fsnotify.Watcher
}
