package plugin_manager

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/conf"
	dev_plg "github.com/csu-gpu-hackers/tx2-k8s-device-plugin/device-plugin"
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type PluginManager struct {
	regisSvr *grpc.Server
	devicePlugins []*dev_plg.DevPlg
	socketPath string
}

func NewPluginManager() *PluginManager {
	pluginPath := conf.PluginManagerConfig["plugin_socket_path"]
	log.Println(pluginPath)
	maxDeviceTypesNum, err := strconv.Atoi(conf.PluginManagerConfig["max_device_types_num"])
	utils.Check(err)
	return &PluginManager{
		regisSvr:      grpc.NewServer(),
		devicePlugins: make([]*dev_plg.DevPlg, maxDeviceTypesNum),
		socketPath:    pluginPath,
	}
}

func (pmg *PluginManager) Run() error {
	lis, err := net.Listen("unix", pmg.socketPath)
	utils.Check(err)
	go func() {
		err := pmg.regisSvr.Serve(lis)
		utils.Check(err)
	}()
	conn, err := utils.Dial(pmg.socketPath, 5 * time.Second)
	utils.Check(err)
	err = conn.Close()
	utils.Check(err)
	return nil
}

func (pmg *PluginManager) RegisterHandler(context context.Context,
	deviceRequest *DeviceRegisterRequest) (*DeviceRegisterReply, error) {
	devicePlugin := dev_plg.NewDevPlg(deviceRequest.DeviceType, deviceRequest.SocketPath)
	pmg.devicePlugins = append(pmg.devicePlugins, devicePlugin)
	go devicePlugin.Run()
	go func() {
		err := devicePlugin.RegisterToKubelet()
		utils.Check(err)
		log.Println("Register finished")
	}()
	return &DeviceRegisterReply{RegisterResult: true}, nil
}
