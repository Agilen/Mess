package server

import (
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

func SockCreate(CSI ClientSocketInfo) *Socket {
	S := &Socket{
		LifeTime:         10,
		Client:           CSI,
		AcceptedAttempts: 5,
		ErrorChan:        make(chan error),
	}
	S.SockListen()
	return S
}

func (S *Socket) SockBind() {

}

func (S *Socket) SockListen() error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(S.Client.Port))
	if err != nil {
		return err
	}
	S.SockAccept(listener)

	return nil
}

func (S *Socket) SockAccept(listener net.Listener) (err error) {
	S.Connection.SetDeadline(time.Now().Add(time.Second * 2))
	S.Connection, err = listener.Accept()
	S.Connection.SetDeadline(time.Now().Add(time.Minute * time.Duration(S.LifeTime)))
	S.Connection.Write([]byte("connection is successful"))
	return
}

func (S *Socket) SockRecv() {

}

func (S *Socket) SockSend() {

}

func (S *Socket) SockClose() {
	S.Connection.Close()

}
