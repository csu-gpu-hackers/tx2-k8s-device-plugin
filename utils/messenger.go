package utils

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
)

type Messenger struct {
	addr string
	Listener net.Listener
	connection net.Conn
	Actions map[string]func(string, string)string
}

type Message struct {
	MessageType string
	PodUID		string
	ContainerID string
	MessageContent string
}

func InitMessenger(socketPath string, actions map[string]func(string, string)string) *Messenger {
	listener, err := net.Listen("unix", socketPath)
	for true {
		if err  != nil {
			err = os.Remove(socketPath)
			if err != nil {
				log.Fatalln("Lack permission")
			}
			continue 
		} else {
			break
		}
	}


	return &Messenger{
		addr:       socketPath,
		Listener:   listener,
		connection: nil,
		Actions: actions,
	}
}


func (messenger *Messenger) Serve()  {
	var message Message
	log.Println("Start serving, waiting for connection")
	//conn, err := messenger.Listener.Accept()


	for true {
		var buffer []byte = make([]byte, 2048)
		conn, err := messenger.Listener.Accept()
		Check(err)
		log.Println("received connection!", conn.RemoteAddr().String())
		messenger.connection = conn
		conn.Read(buffer)
		log.Println("received message:", string(buffer))
		err = json.Unmarshal(buffer[:bytes.IndexByte(buffer, 0)], &message)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("No error while unmarshlling")
		}
		//buffer = []
		log.Println("received struct from %s: %s", message.ContainerID, message.MessageContent)
		conn.Write([]byte(messenger.Actions[message.MessageType](message.ContainerID, message.MessageContent)))
	}

}

func (messenger *Messenger) Request(socketPath string, content string) {

}


//
//func (messenger *Messenger) receive()  {
//
//}



//func main() {
//	listener, err := net.Listen("unix", "a.sock")
//	Check(err)
//	conn, err := listener.Accept()
//}
