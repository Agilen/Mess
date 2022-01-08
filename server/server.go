package server

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/Agilen/Mess/server/store"

	"github.com/Agilen/Mess/server/model"
)

type (
	ServerContext struct {
		Store      store.Store
		ServerName string
		ServerIpv4 string
		Ports      map[int]bool
		Sockets    map[string]*Socket
		ErrorChan  chan error
		Conn       net.Conn
	}
)

func NewServerContext(store store.Store) *ServerContext {
	return &ServerContext{
		Store:      store,
		ServerName: "SuperPuperServer",
		ServerIpv4: "172.16.9.120",
		Ports:      PrepPortsMap(),
		Sockets:    PrepSocketMap(),
		ErrorChan:  make(chan error),
	}
}

func ListenAndServe(store store.Store) {
	SC := NewServerContext(store)
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
	fmt.Println("Cервер запущен")
	for {
		fmt.Println("Слушаем")
		SC.Conn, err = listener.Accept()
		if err != nil {
			SC.ErrorChan <- err
			continue
		}

		message := make([]byte, 1024)
		_, err = SC.Conn.Read(message)
		if err != nil {
			SC.ErrorChan <- err
			continue
		}
		err = SC.DataProcessing(message)
		if err != nil {
			SC.ErrorChan <- err
			continue
		}

		err = SC.Conn.Close()
		if err != nil {
			SC.ErrorChan <- err
			continue
		}

	}
}

func (SC *ServerContext) LookForCloseSocket() {
}

func (SC *ServerContext) DataProcessing(message []byte) error {
	//4byte - len
	//16byte - command
	//else - info

	len := binary.LittleEndian.Uint32(message[:4])
	command := string(message[4:20])
	data := message[20:len]

	err := SC.Command(command, data)
	if err != nil {
		return err
	}

	return nil
}

func (SC *ServerContext) Command(command string, data []byte) error {
	switch command {
	case "registration":
		err := SC.NewUser(data) //исправлю
		if err != nil {
			return err
		}
	case "authorization":
		//here will be DH
		fmt.Println("auth")
	case "conntosocket":
		err := SC.GiveClientFreePort()
		if err != nil {
			return err
		}
		err = SC.Sock(data)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown command")
	}

	return nil
}

func (SC *ServerContext) GiveClientFreePort() error {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(SC.TakeFreePort()))
	_, err := SC.Conn.Write(bs)
	if err != nil {
		return err
	}
	return nil

}

func (SC *ServerContext) Sock(data []byte) error {
	CSI := ClientSocketInfo{}

	err := json.Unmarshal(data, &CSI)
	if err != nil {
		return err
	}

	SC.Ports[CSI.Port] = true
	SC.Sockets[CSI.Name] = SockCreate(CSI, SC.ErrorChan, SC.Sockets)
	_, err = SC.Conn.Write([]byte("socket is ready"))
	if err != nil {
		return err
	}
	go SC.Sockets[CSI.Name].ServeSocket()

	return nil //add errchan
}

func (SC *ServerContext) NewUser(data []byte) error {
	user := model.User{}
	err := json.Unmarshal(data, &user.BaseUserInfo)
	if err != nil {
		return err
	}
	err = SC.Store.User().Create(&user)
	if err != nil {
		return err
	}
	return nil
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
