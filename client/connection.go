package client

import (
	"fmt"
	"net"
)

type (
	ClientContext struct {
		Conn Connection
	}

	Connection struct {
		serverAddr string
		port       int
		conn       net.Conn
	}
)

func NewContext() *ClientContext {
	return &ClientContext{
		Conn: Connection{
			serverAddr: "127.0.0.1:10000",
		},
	}
}

func (clientCtx *ClientContext) ConnectToServer() (err error) {

	clientCtx.Conn.conn, err = net.Dial("tcp", clientCtx.Conn.serverAddr)
	if err != nil {
		return err
	}
	fmt.Println("connected")

	return nil
}
