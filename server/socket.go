package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

type Socket struct {
	LifeTime         int
	Client           ClientSocketInfo
	Connection       net.Conn
	AcceptedAttempts int
	ErrorChan        chan error
	Sockets          *map[string]*Socket
}

type Message struct {
	From    string
	To      string
	Message []byte
}

type ClientSocketInfo struct {
	Port   int
	Name   string
	Ipv4   string
	Number int
}

func SockCreate(CSI ClientSocketInfo, errorchan chan error, Sockets *map[string]*Socket) *Socket {
	S := &Socket{
		LifeTime:         10,
		Client:           CSI,
		AcceptedAttempts: 5,
		ErrorChan:        errorchan,
		Sockets:          Sockets,
	}

	return S
}

func (S *Socket) SockListen() error {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(S.Client.Port))
	if err != nil {
		return err
	}
	err = S.SockAccept(listener)
	if err != nil {
		return err
	}

	return nil
}

func (S *Socket) SockAccept(listener net.Listener) (err error) {

	S.Connection, err = listener.Accept()
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

	for {

		mes, data, err := S.ReadFullMessage()
		if err != nil {
			return err
		}

		b := *S.Sockets
		if b[mes.To] != nil {
			err = b[mes.To].SockSend(data)
			if err != nil {
				return err
			}
		} else {
			err := S.SockSend([]byte("no such socket"))
			if err != nil {
				return err
			}
		}

	}
}

func (S *Socket) SockSend(message []byte) error {
	fmt.Println(string(message))
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

	delete(*S.Sockets, S.Client.Name)
	fmt.Println("Close socket", S.Client.Name)
	return nil
}

func (S *Socket) ServeSocket() error {

	err := S.SockListen()
	if err != nil {
		S.ErrorChan <- err
		S.SockClose()
	}

	go S.SockRecv()

	return nil
}

func (S *Socket) ReadFullMessage() (Message, []byte, error) {
	buf := make([]byte, 4)
	_, err := S.Connection.Read(buf)
	if err != nil {
		return Message{}, nil, err
	}

	mes := Message{}

	data := make([]byte, binary.LittleEndian.Uint32(buf[:4]))

	_, err = S.Connection.Read(data)
	if err != nil {
		return mes, nil, err
	}

	err = json.Unmarshal(data, &mes)
	if err != nil {

		return mes, nil, err
	}

	return mes, data, nil
}
