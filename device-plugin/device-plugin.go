package device_plugin

import (
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/devices"
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/vDeviceManager"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	plugin "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"net"
	"path"
	"strings"
	"time"
)
//type DeviceType string
const (
	kubeletSocket = "kubelet.sock"
	DevicePluginDir = "/var/lib/kubelet/device-plugins/"
)

type DevPlg struct {
	deviceType    string
	srv           *grpc.Server
	ctx           context.Context
	DeviceManager devices.Device
	cancel		  context.CancelFunc
	devSocketPath string
	vDevMgr		  *vDeviceManager.VDeviceManager
	//devUpdate	  chan bool
}

func NewDevPlg(deviceType string, devSocketPath string) *DevPlg {
	devmgr := devices.NewGPUManager()
	vDevMgr := &vDeviceManager.VDeviceManager{}
	go vDevMgr.Serve()
	ctx, cancel := context.WithCancel(context.Background())
	log.Println("Construction of DevPlg starting: ", deviceType)
	return &DevPlg{
		deviceType: deviceType,
		srv:     grpc.NewServer(grpc.EmptyServerOption{}),
		ctx:     ctx,
		DeviceManager: devmgr,
		cancel:  cancel,
		devSocketPath: devSocketPath,
		vDevMgr: &vDeviceManager.VDeviceManager{},
		//devUpdate: make(chan bool),
	}


}

func (dp *DevPlg) Run() error {
	log.Println("dp start running")
	go dp.DeviceManager.WatchDevice()

	plugin.RegisterDevicePluginServer(dp.srv,dp)
	lis, err := net.Listen("unix", dp.devSocketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	utils.Check(err)

	go func() {
		log.Println("Start serving listener")
		err = dp.srv.Serve(lis)
		utils.Check(err)
	}()

	conn, err := utils.Dial(dp.devSocketPath, 5 * time.Second)
	utils.Check(err)
	err = conn.Close()
	utils.Check(err)
	return nil
}

func (dp *DevPlg) RegisterToKubelet() error {
	log.Println("Start registering to kubelet")
	var kubeletSocketFile = DevicePluginDir + kubeletSocket
	conn, err := utils.Dial(kubeletSocketFile, 5*time.Second)
	utils.Check(err)
	defer conn.Close()

	dpClient := plugin.NewRegistrationClient(conn)
	req := &plugin.RegisterRequest{
		Version:      plugin.Version,
		Endpoint:     path.Base(dp.devSocketPath),
		ResourceName: dp.deviceType,
	}
	_, err = dpClient.Register(context.Background(), req)
	utils.Check(err)
	log.Println("Registering finished")
	return nil
}



func (dp *DevPlg) ListAndWatch(empty *plugin.Empty, server plugin.DevicePlugin_ListAndWatchServer) error {
	//panic("implement me")
	log.Infoln("ListAndWatch called")
	devs := make([]*plugin.Device, len(dp.DeviceManager.GetDeviceParts()))
	i := 0
	for _, dev := range dp.DeviceManager.GetDeviceParts() {
		devs[i] = dev
		i++
	}
	err := server.Send(&plugin.ListAndWatchResponse{Devices: devs})
	utils.Check(err)

	for true {
		log.Println("waiting for device change")
		select {
		case event, ok := <- dp.DeviceManager.GetDeviceChangeNotifier().Events:

			if !ok {
				continue
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Println("File has been removed")
				for _, part := range dp.DeviceManager.GetDeviceParts() {
					part.Health = plugin.Unhealthy
				}
			} else if event.Op&fsnotify.Write == fsnotify.Write {
				load := dp.DeviceManager.GetDeviceLoads()
				log.Printf("used cores: %v\n", load)
				devs := make([]*plugin.Device, 100)
				for i, dev := range dp.DeviceManager.GetDeviceParts(){
					devs[i] = dev
					i++
				}
				err := server.Send(&plugin.ListAndWatchResponse{Devices: devs})
				utils.Check(err)
			}

		}

		//select {
		//case <- dp.DeviceManager.DeviceChangeNotifier:
		//	load := dp.DeviceManager.GetDeviceLoads()
		//	log.Printf("Available cores: %v\n", load)
		//	devs := make([]*plugin.Device, 100 - load)
		//	for i, dev := range dp.DeviceManager.DeviceParts[load:]{
		//		devs[i] = dev
		//		i++
		//	}
		//	err := server.Send(&plugin.ListAndWatchResponse{Devices: devs})
		//	utils.Check(err)
		//case <-dp.ctx.Done():
		//	log.Println("ListAndWatch exit")
		//	return nil
		//}

	}
	return nil
}

func (dp *DevPlg) Allocate(ctx context.Context, requests *plugin.AllocateRequest) (*plugin.AllocateResponse, error) {
	log.Println("Allocate called")
	resps := &plugin.AllocateResponse{}

	for _, req := range requests.ContainerRequests {
		log.Printf("received request: %s\n", strings.Join(req.DevicesIDs, ","))
		//req.
		resp := plugin.ContainerAllocateResponse{
			Envs: map[string]string{
				dp.deviceType: strings.Join(req.DevicesIDs, ","),
				"LD_LIBRARY_PATH": "/etc/vcuda/:$LD_LIBRARY_PATH",
			},
		}
		vlmMgr := vDeviceManager.NewVolumeManager()
		// mounting hijacking library path, which
		// also will be containing config file
		dp.vDevMgr.NewDevice(dp.deviceType, vlmMgr)
		resp.Mounts = append(resp.Mounts, &plugin.Mount{
			ContainerPath:        "/etc/vcuda/",
			HostPath:             vlmMgr.VCudaLibHostPath,
			ReadOnly:             false,
		})
		//resp.Mounts =



		//dp.vDevMgr.
		resps.ContainerResponses = append(resps.ContainerResponses, &resp)
	}
	return resps, nil
}

func (dp *DevPlg) GetDevicePluginOptions(ctx context.Context, empty *plugin.Empty) (*plugin.DevicePluginOptions, error) {
	//panic("implement me")
	log.Infoln("GetDevicePluginOptions called")
	return &plugin.DevicePluginOptions{PreStartRequired: true}, nil
}


func (dp *DevPlg) GetPreferredAllocation(ctx context.Context, request *plugin.PreferredAllocationRequest) (*plugin.PreferredAllocationResponse, error) {
	log.Infoln("GetPreferredAllocation called")
	panic("implement me")
}

func (dp *DevPlg) PreStartContainer(ctx context.Context, request *plugin.PreStartContainerRequest) (*plugin.PreStartContainerResponse, error) {
	//panic("implement me")
	log.Println("PreStartContainer called")
	return &plugin.PreStartContainerResponse{}, nil
}
