package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [username]")
		os.Exit(1)
	}

	username := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = conn.Write([]byte("join:" + username))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Welcome to server ", username)

	// Menangani Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		_, err := conn.Write([]byte("left:exit"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("\nYou have left the chat.")
		os.Exit(0)
	}()

	go handleConn(conn)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("[You]: ") 
		message, _ := inputReader.ReadString('\n')
		message = strings.TrimSpace(message)

		// Jika pengguna ingin keluar
		if message == "" {
			continue
		}

		_, err := conn.Write([]byte("message:" + message))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}


func handleConn(request *net.UDPConn) {
	for {
		var buf [512]byte
		n, _, err := request.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		// Menampilkan pesan yang diterima
		incomingMessage := string(buf[:n])
		if incomingMessage != "" {
			fmt.Printf("\r%s\n[You]: ", incomingMessage)
		}

	}
}
