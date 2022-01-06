package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	server "github.com/Agilen/Mess/server"
)

func main() {

	// Подключаемся к сокету
	fmt.Println("Подключаемя к сокету")
	conn, err := net.Dial("tcp", "127.0.0.1:3333")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Читаем")

	buffer := make([]byte, 4)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	port := binary.LittleEndian.Uint32(buffer)

	fmt.Println("получили порт")

	CSI := server.ClientSocketInfo{
		Port: int(port),
		Name: "TestClient",
		Ipv4: "127.0.0.1",
	}
	fmt.Println("маршал информации для создания сокета")
	jsonCSI, err := json.Marshal(CSI)
	if err != nil {
		log.Fatal(err)
	}
	// time.Sleep(time.Second * 1)
	fmt.Println("Отправляем информацию")
	_, err = conn.Write(jsonCSI)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Подключаемя к сокету", strconv.Itoa(int(port)))
	for {
		conn, err = net.Dial("tcp", ":"+strconv.Itoa(int(port)))
		if err != nil {
			log.Fatal(err)
		}
		if conn != nil {
			break
		}
	}
	fmt.Println("ждем хороших новостей")
	for {
		_, err = conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(buffer))
	}
	fmt.Println(string(buffer))

}
