package utils

import (
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"regexp"
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

