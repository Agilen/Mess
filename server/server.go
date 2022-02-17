package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

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

	Response struct {
		Response string
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
	go SC.ListenClientWish("3333")
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

func (SC *ServerContext) ListenClientWish(port string) {

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Run server")
	for {
		fmt.Println("Listen")
		SC.Conn, err = listener.Accept()
		if err != nil {
			SC.ErrorChan <- err
			continue
		}
		SC.Conn.SetDeadline(time.Now().Add(time.Second * 5))

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

func (SC *ServerContext) DataProcessing(message []byte) error {
	//4byte - len
	//16byte - command
	//else - info

	len := binary.LittleEndian.Uint32(message[:4])
	command := string(message[4:20])
	data := message[20+4 : len+4]

	g := string(bytes.Split([]byte(command), []byte{0})[0])

	err := SC.Command(g, data)
	if err != nil {
		return err
	}

	return nil
}

func (SC *ServerContext) Command(command string, data []byte) error {
	switch command {
	case "registration":
		err := SC.NewUser(data)
		if err != nil {
			return err
		}
	case "auth":

		reg, err := SC.CheckUser(data)
		if err != nil {
			return err
		}

		if reg {
			bs, err := SC.GiveClientFreePort()
			if err != nil {
				return err
			}

			err = SC.Sock(data, bs)
			if err != nil {
				return err
			}
		} else {
			response, err := json.Marshal(Response{Response: NOT_REGISTER})
			if err != nil {
				return err
			}
			SC.Conn.Write(response)
		}

	default:
		return errors.New("unknown command")
	}

	return nil
}

func (SC *ServerContext) GiveClientFreePort() (int, error) {
	port := SC.TakeFreePort()

	response, err := json.Marshal(Response{Response: strconv.Itoa(port)})
	if err != nil {
		return 0, err
	}
	_, err = SC.Conn.Write(response)
	if err != nil {
		return 0, err
	}
	return port, nil
}

func (SC *ServerContext) Sock(data []byte, port int) error {
	CSI := ClientSocketInfo{}

	err := json.Unmarshal(data, &CSI)
	if err != nil {
		return err
	}

	CSI.Port = port
	SC.Ports[CSI.Port] = true
	SC.Sockets[CSI.Name] = SockCreate(CSI, SC.ErrorChan, &SC.Sockets)

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

func (SC *ServerContext) CheckUser(data []byte) (bool, error) {
	BUI := model.BaseUserInfo{}

	err := json.Unmarshal(data, &BUI)
	if err != nil {
		return false, err
	}

	return SC.Store.User().FindUser(BUI.Login, BUI.Password)

}

func (SC *ServerContext) CheckPassword(password string) {

}
