package vDevice_manager

import (
	container "github.com/csu-gpu-hackers/tx2-k8s-device-plugin/pod_config"
	"github.com/csu-gpu-hackers/tx2-k8s-device-plugin/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"time"
)

type VDeviceManager struct{
	VDeviceType string
	ContainerConfigs map[string]*container.ContainerConfig
	grpcServer *grpc.Server
	socketPath string
	ctx context.Context
	cancel context.CancelFunc
}

func NewVDeviceManager(VDeviceType string) *VDeviceManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &VDeviceManager{
		VDeviceType: VDeviceType,
		ContainerConfigs: make(map[string]*container.ContainerConfig),
		grpcServer: grpc.NewServer(grpc.EmptyServerOption{}),
		socketPath: "/var/lib/kubelet/device-plugins/vdmgr.sock",
		cancel: cancel,
		ctx: ctx,
	}
}

func (vdm *VDeviceManager) Run() error {
	/* Start grpc service here so that vdevices can be registered */
	lis, err := net.Listen("unix", vdm.socketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go func() {
		err := vdm.grpcServer.Serve(lis)
		utils.Check(err)
	}()
	conn, err := utils.Dial(vdm.socketPath, 5 * time.Second)
	utils.Check(err)
	err = conn.Close()
	utils.Check(err)
	return nil
}



// RegisterVDevice Receiving register request from
// virtual device library running inside docker
// containers.
func (VDeviceManager) RegisterVDevice(c interface{}, request *VDeviceRequest) (*VDeviceResponse, error) {
	log.Println("RegisterVDevice called")
	containerConfig := container.NewContainer(request.ContainerID)
}

