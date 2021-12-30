package server

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

var ()

type (
	serverContext struct {
		ServerName string
		ServerIpv4 string
		Ports      map[int]bool
		Sockets    []*Socket
	}

	ConnWithUser struct {
		Port int
		Conn net.Conn
	}
)

func ListenAndServe() error {
	SC := NewServerContext()
	go SC.ListenClientWish()

	return nil
}

func NewServerContext() *serverContext {
	return &serverContext{
		ServerName: "SuperPuperServer",
		ServerIpv4: "172.16.9.120",
		Ports:      PrepPortsMap(),
	}
}

func PrepPortsMap() map[int]bool {
	m := make(map[int]bool)
	for i := 3000; i < 65536; i++ {
		m[i] = false
	}
	return m
}

func TakeFreePort(server *serverContext) int {
	for i := 3000; i < 65536; i++ {
		if !server.Ports[i] {
			server.Ports[i] = true
			return i
		}
	}
	return 0
}

func (SC *serverContext) ListenClientWish() {

	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			break
		}

		BaseInfo, err := json.Marshal(SC)
		if err != nil {
			log.Println(err.Error())
			break
		}

		_, err = connection.Write(BaseInfo)
		if err != nil {
			connection.Close()
			log.Println(err.Error())
			break
		}
		//////////////////////////////////////////////////////
		buffer := make([]byte, 2048)

		connection.SetReadDeadline(time.Now().Add(time.Second * 1))

		_, err = connection.Read(buffer)
		if err != nil {
			connection.Close()
			log.Println(err.Error())
			break
		}
		CSI := new(ClientSocketInfo)
		err = json.Unmarshal(buffer, &CSI)
		if err != nil {
			connection.Close()
			log.Println(err.Error())
			break
		}
		////////////////////////////////////////////////////////
		CSI.Number = len(SC.Sockets)

		connection.Close()

		SC.Ports[int(buffer[0])] = true
		SC.Sockets = append(SC.Sockets, SockCreate(*CSI, int(buffer[0])))
	}
}
