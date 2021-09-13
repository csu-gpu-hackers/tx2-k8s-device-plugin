package devices

import (
	"github.com/fsnotify/fsnotify"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
)

type Device interface {
	GetDeviceParts() []*plugin.Device
	GetDeviceStatus() utils.DeviceStatus
	GetDeviceLoads() int
	WatchDevice() error
	GetDeviceChangeNotifier() *fsnotify.Watcher
}
