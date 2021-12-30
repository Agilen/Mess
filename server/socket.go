package server

import (
	"net"
	"strconv"
)

type Socket struct {
	Port       int
	LifeTime   int
	Client     ClientSocketInfo
	Connection net.Conn
}

type ClientSocketInfo struct {
	Name   string
	Ipv4   string
	Number int
}

func SockCreate(CSI ClientSocketInfo, port int) *Socket {
	return &Socket{
		Port:     port,
		LifeTime: 5000000,
		Client:   CSI,
	}
}

func (S *Socket) SockBind() {

}

func (S *Socket) SockListen() error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(S.Port))
	if err != nil {
		return err
	}
	S.SockAccept()

	return nil
}

func (S *Socket) SockAccept() {

}

func (S *Socket) SockRecv() {

}

func (S *Socket) SockSend() {

}

func (S *Socket) SockClose() {

}
