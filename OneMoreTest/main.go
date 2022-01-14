package main

import (
	"encoding/binary"
	"encoding/json"
	"net"

	"github.com/Agilen/Mess/server/model"
)

func main() {

}

func Conn() (net.Conn, error) {
	con, err := net.Dial("tcp", ":3333")
	if err != nil {
		return nil, err
	}
	return con, nil
}

func Registration(conn net.Conn) error {
	info := model.BaseUserInfo{
		Name:     "Denis",
		Email:    "denis@denis.com",
		Password: "12345678",
	}
	jsoninfo, err := json.Marshal(info)
	if err != nil {
		return err
	}

	size := len(jsoninfo)
	bytesize := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytesize, uint32(size))
	var data []byte
	data = append(data, bytesize...)
	command := make([]byte, 20)
	copy(command, []byte("registration"))
	data = append(data, command...)
	data = append(data, jsoninfo...)

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func Sock(conn net.Conn) error {
	var data []byte
	bytesize := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytesize, 0)
	data = append(data, bytesize...)
	command := make([]byte, 20)
	copy(command, []byte("conntosocket"))
	data = append(data, command...)
	_, err := conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
