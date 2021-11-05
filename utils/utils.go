package utils

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"net"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

type DeviceStatus int
type DeviceType string



const (
	OK DeviceStatus = iota
	OCCUPIED
	PENDING
	OFFLINE
	DEAD
)

var (
	//K8sClient *kubernetes.Clientset = createK8sInClusterClient()
	K8sClient = createK8sOutClusterClient()
	_, b, _, _ = runtime.Caller(0)
	RootPath = filepath.Join(filepath.Dir(b), "../..")
)

//func Copy(src, dst string) error {
//	in, err := os.Open(src)
//	if err != nil {
//		return err
//	}
//	defer in.Close()
//
//	out, err := os.Create(dst)
//	if err != nil {
//		return err
//	}
//	defer out.Close()
//
//	_, err = io.Copy(out, in)
//	if err != nil {
//		return err
//	}
//	return out.Close()
//}


func createK8sInClusterClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func createK8sOutClusterClient() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func CheckPodStatus(poduid string) v1.PodPhase {
	namespace := "default"
	pod, err := K8sClient.CoreV1().Pods(namespace).Get(context.TODO(), poduid, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", poduid, namespace)
		panic(err)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			poduid, namespace, statusError.ErrStatus.Message)
		panic(err)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", poduid, namespace)
		return pod.Status.Phase
	}

}

func GetPodLimits(poduid, containerid string) {

}

func Dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadFile(filename string) string {
	dat, err := ioutil.ReadFile(filename)
	Check(err)
	return string(dat)
}

func ExtractNumber(str string) (int,error){
	regex := "[-]?\\d[\\d,]*[\\.]?[\\d{2}]*"
	re := regexp.MustCompile(regex)
	gpuLoadStr := re.FindAllString(str, 1)[0]
	gpuLoad, err := strconv.ParseFloat(gpuLoadStr, 64)
	Check(err)
	return int(gpuLoad),nil

}

