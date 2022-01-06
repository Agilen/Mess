package server

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Socket struct {
	LifeTime         int
	Client           ClientSocketInfo
	Connection       net.Conn
	AcceptedAttempts int
	ErrorChan        chan error
}

type ClientSocketInfo struct {
	Port   int
	Name   string
	Ipv4   string
	Number int
}

func SockCreate(CSI ClientSocketInfo, errorchan chan error) *Socket {
	S := &Socket{
		LifeTime:         10,
		Client:           CSI,
		AcceptedAttempts: 5,
		ErrorChan:        errorchan,
	}
	fmt.Println("сокет создан")

	return S
}

func (S *Socket) SockBind() {

}

func (S *Socket) SockListen() error {
	fmt.Println("ждем-с")
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(S.Client.Port))
	if err != nil {
		return err
	}
	S.SockAccept(listener)

	return nil
}

func (S *Socket) SockAccept(listener net.Listener) (err error) {
	// err := S.Connection.SetDeadline(time.Now().Add(time.Second * 2))
	// if err != nil {
	// 	S.ErrorChan <- err
	// 	return
	// }

	fmt.Println("Еще ждем")
	S.Connection, err = listener.Accept()
	if err != nil {
		return err
	}
	fmt.Println("дождались")
	err = S.Connection.SetDeadline(time.Now().Add(time.Minute * time.Duration(S.LifeTime)))
	if err != nil {
		return err
	}

	_, err = S.Connection.Write([]byte("connection is successful"))
	if err != nil {
		return err
	}
	return
}

func (S *Socket) SockRecv() error {
	var buffer []byte
	for {
		_, err := S.Connection.Read(buffer)
		if err != nil {
			return err
		}
		_, err = S.Connection.Write(buffer)
		if err != nil {
			S.ErrorChan <- err
			return err
		}
	}
}

func (S *Socket) SockSend(message []byte) error {
	_, err := S.Connection.Write(message)
	if err != nil {
		return err
	}
	return nil
}

func (S *Socket) SockClose() error {
	err := S.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (S *Socket) ServeSocket() error {
	S.SockListen()
	go S.SockRecv()
	return nil
}
