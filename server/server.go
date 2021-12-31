package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

var ()

type (
	ServerContext struct {
		ServerName string
		ServerIpv4 string
		Ports      map[int]bool
		Sockets    map[string]*Socket
		ErrorChan  chan error
	}
)

func ListenAndServe() {
	SC := NewServerContext()
	go SC.ListenClientWish()
	HadnleError(<-SC.ErrorChan)
}

func NewServerContext() *ServerContext {
	return &ServerContext{
		ServerName: "SuperPuperServer",
		ServerIpv4: "172.16.9.120",
		Ports:      PrepPortsMap(),
		Sockets:    PrepSocketMap(),
		ErrorChan:  make(chan error),
	}
}

func PrepPortsMap() map[int]bool {
	m := make(map[int]bool)
	for i := 3000; i < 65536; i++ {
		m[i] = false
	}
	return m
}

func PrepSocketMap() map[string]*Socket {
	return make(map[string]*Socket)
}

func TakeFreePort(server *ServerContext) int {
	for i := 3000; i < 65536; i++ {
		if !server.Ports[i] {
			server.Ports[i] = true
			return i
		}
	}
	return 0
}

func (SC *ServerContext) ListenClientWish() {

	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatal(err.Error())
	}
	for {
		//ждём подключения
		connection, err := listener.Accept()
		if err != nil {
			SC.ErrorChan <- err
			continue
		}

		BaseInfo, err := json.Marshal(SC.Ports) //маршалим карту с портами
		if err != nil {
			SC.ErrorChan <- err
			continue
		}

		_, err = connection.Write(BaseInfo) //отправлем карту
		if err != nil {
			SC.ErrorChan <- err
			continue
		}

		buffer := make([]byte, 2048)
		////////////////////////////////////////////////////////
		connection.SetReadDeadline(time.Now().Add(time.Millisecond * 500)) //ждем ответа от клиента, если ответа нет, то рвем подключение и ждем новое

		_, err = connection.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				SC.ErrorChan <- err
				continue
			}
			SC.ErrorChan <- err
			continue
		}
		CSI := new(ClientSocketInfo)
		err = json.Unmarshal(buffer, &CSI)
		if err != nil {
			SC.ErrorChan <- err
			continue
		}
		////////////////////////////////////////////////////////
		CSI.Number = len(SC.Sockets)

		connection.Close()

		SC.Ports[int(buffer[0])] = true
		SC.Sockets[CSI.Name] = SockCreate(*CSI)
	}
}

func HadnleError(ch error) {
	for {
		if ch != nil {
			fmt.Println(ch)
			ch = nil
		}
	}
}
