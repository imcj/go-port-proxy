package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	argv := os.Args
	if len(argv) < 4 {
		fmt.Println("Argument err, use port proxy example:")
		fmt.Println("root@host:~$./proxy.exe tcp 127.0.0.1:80 192.168.1.102:8080")
		fmt.Println("Source from https://github.com/JamesWone/go-port-proxy")
		return
	}
	listen, err := net.Listen(argv[1], argv[2])
	if err != nil {
		fmt.Println("Port listen err:", err.Error())
		if err := listen.Close(); err != nil {
			panic(err)
		}
		return
	}
	for {
		client, err := listen.Accept()
		if err != nil {
			if addr := client.RemoteAddr(); addr != nil {
				log.Println("Client->Proxy Connect Fail:", addr)
			}

			if err := client.Close(); err != nil {
				panic(err)
			}
			break
		}

		server, err := net.Dial(argv[1], argv[3])
		if err != nil {
			fmt.Println("Proxy->Server Connect Fail:", err.Error())
			if err := server.Close(); err != nil {
				panic(err)
			}
			if err := client.Close(); err != nil {
				panic(err)
			}
			break
		}

		log.Println("New connection:", client.RemoteAddr().Network(), client.RemoteAddr(), "<->", argv[2], "<->", server.RemoteAddr())
		go proxy(client, server)
		go proxy(server, client)
	}

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
}

func proxy(client net.Conn, server net.Conn) {
	reader := bufio.NewReader(client)
	_, err := reader.WriteTo(server)
	if err != nil {
		fmt.Println("Proxy err,", err.Error())
		if err := client.Close(); err != nil {
			panic(err)
		}
		if err := server.Close(); err != nil {
			panic(err)
		}
	}
}
