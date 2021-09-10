package conf
var PluginManagerConfig map[string]string = map[string]string {
	"plugin_socket_path": "/var/lib/kubelet/device-plugins/plugin-manager.socket",
	"max_device_types_num": "16",
}