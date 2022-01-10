package server

import (
	"encoding/binary"
	"encoding/json"
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
	Buffer           []byte
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
		Buffer:           make([]byte, 2048),
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
	go S.SockClose()

	_, err = S.Connection.Write([]byte("connection is successful"))
	if err != nil {
		return err
	}
	return
}

func (S *Socket) SockRecv() error {
	buffer := make([]byte, 4)
	for {
		_, err := S.Connection.Read(buffer)
		if err != nil {
			return err
		}
		mes, data, err := S.ReadFullMessage(buffer)
		if err != nil {
			return err
		}
		b := *S.Sockets //?????
		err = b[mes.To].SockSend(data)
		if err != nil {
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
	time.Sleep(time.Minute * time.Duration(S.LifeTime))
	err := S.Connection.Close()
	if err != nil {
		return err
	}
	delete(*S.Sockets, S.Client.Name)

	return nil
}

func (S *Socket) ServeSocket() error {
	err := S.SockListen()
	if err != nil {
		return err
	}
	go S.SockRecv()
	return nil
}

func (S *Socket) ReadFullMessage(message []byte) (Message, []byte, error) {
	var data []byte
	var i int
	mes := Message{}
	len := binary.LittleEndian.Uint32(message[:4])

	for i < int(len) {
		_, err := S.Connection.Read(S.Buffer)
		if err != nil {
			return mes, nil, err
		}
		data = append(data, S.Buffer...)
		i += 2048
	}

	err := json.Unmarshal(data, &mes)
	if err != nil {
		return mes, nil, err
	}

	return mes, data, nil
}
