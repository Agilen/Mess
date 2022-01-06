package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type (
	ServerContext struct {
		ServerName string
		ServerIpv4 string
		Ports      map[int]bool
		Sockets    map[string]*Socket
		ErrorChan  chan error
	}
)

func NewServerContext() *ServerContext {
	return &ServerContext{
		ServerName: "SuperPuperServer",
		ServerIpv4: "172.16.9.120",
		Ports:      PrepPortsMap(),
		Sockets:    PrepSocketMap(),
		ErrorChan:  make(chan error),
	}
}

func ListenAndServe() {
	SC := NewServerContext()
	go SC.ListenClientWish()
	for {
		HadnleError(<-SC.ErrorChan)
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

func (SC *ServerContext) TakeFreePort() int {
	for i := 3000; i < 65536; i++ {
		if !SC.Ports[i] {
			SC.Ports[i] = true
			return i
		}
	}
	return 0
}

func (SC *ServerContext) ListenClientWish() {

	listener, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("сервер запущен")
	for {
		fmt.Println("Новая сессия создана")
		//ждём подключения
		connection, err := listener.Accept()
		if err != nil {
			SC.ErrorChan <- err
			continue
		}
		fmt.Println("подключение создано")
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(SC.TakeFreePort()))
		_, err = connection.Write(bs) //отправлем карту
		if err != nil {
			SC.ErrorChan <- err
			continue
		}

		fmt.Println("список отправлен")

		buffer := make([]byte, 1024)
		////////////////////////////////////////////////////////
		err = connection.SetReadDeadline(time.Now().Add(time.Millisecond * 500)) //ждем ответа от клиента, если ответа нет, то рвем подключение и ждем новое
		if err != nil {
			SC.ErrorChan <- err
			continue
		}
		fmt.Println("Ждем ответа")
		_, err = connection.Read(buffer)
		if err != nil {
			SC.ErrorChan <- err
			continue
		}
		n := bytes.Split(buffer, []byte{0})
		fmt.Println("ответ получен")
		CSI := ClientSocketInfo{}

		err = json.Unmarshal(n[0], &CSI)

		if err != nil {
			SC.ErrorChan <- err
			continue
		}
		fmt.Println("анмаршал ответа")
		////////////////////////////////////////////////////////
		CSI.Number = len(SC.Sockets)

		connection.Close()
		fmt.Println("создаем сокет")
		SC.Ports[int(buffer[0])] = true
		SC.Sockets[CSI.Name] = SockCreate(CSI, SC.ErrorChan)
		go SC.Sockets[CSI.Name].ServeSocket()
	}
}

func (SC *ServerContext) LookForCloseSocket() {

}

//TODO: Придумать нормальный хендлер
func HadnleError(ch error) {
	for {
		if ch != nil {
			fmt.Println(ch)
			ch = nil
		}
	}
}
