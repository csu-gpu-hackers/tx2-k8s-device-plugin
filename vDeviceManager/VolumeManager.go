package vDeviceManager
import (
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

const (
	vcudalibcontainerbase = "/etc/vcuda/"
)

type VolumeManager struct {
	allocationIdentifier string
	containerID string
	VCudaLibHostPath string
	VCudaLibContainerPath string
	VCudaConfigHostPath string
	VCudaConfigContainerPath string
}

func NewVolumeManager() *VolumeManager {
	allocationIdentifier := uuid.NewV4().String()
	VCudaLibHostPath := path.Join(vcudalibcontainerbase, allocationIdentifier)
	vlmMgr := &VolumeManager{
		allocationIdentifier:     allocationIdentifier,
		containerID:                   "",
		VCudaLibHostPath:         VCudaLibHostPath,
		VCudaLibContainerPath:    vcudalibcontainerbase,
		VCudaConfigHostPath:      path.Join(VCudaLibHostPath, "vcuda.config"),
		VCudaConfigContainerPath: path.Join(vcudalibcontainerbase, "vcuda.config"),
	}
	vlmMgr.prepareDirectories()
	return vlmMgr
}

func (v *VolumeManager) UpdateContainerID(containerID string) {
	v.containerID = containerID
}

func (v *VolumeManager) prepareDirectories() {
	err := os.MkdirAll(v.VCudaLibHostPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Directory exists")
		panic(err)
	}
}


func (v *VolumeManager) WriteConfig() {
	dl := []byte("")
	os.WriteFile(v.VCudaConfigContainerPath, dl, os.ModePerm)
}

func (v *VolumeManager) ReleaseConfig() {
	os.Remove(v.VCudaLibHostPath)
}




