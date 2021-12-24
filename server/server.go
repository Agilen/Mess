package server

import (
	"net"
	"strconv"
)

var ()

type (
	serverContext struct {
		Ports map[int]bool
	}

	ConnWithUser struct {
		Port int
		Conn net.Conn
	}

	Socket struct {
	}
)

func NewServerContext() *serverContext {
	return &serverContext{
		Ports: PrepPortsMap(),
	}
}

func (CWU *ConnWithUser) CreateNewSocket() (err error) {

	CWU, err = net.Listen("tcp", ":"+strconv.Itoa(TakeFreePort()))
	return nil
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
