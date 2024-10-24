package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var clients = make(map[string]string)

func main() {
	fmt.Println("Running Server....")
	udpAddr, err := net.ResolveUDPAddr("udp4", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	for {
		handleConn(conn)
	}
}

func handleConn(conn *net.UDPConn) {
	var buf [512]byte
	n, addr, err := conn.ReadFromUDP(buf[:])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	message := string(buf[:n])
	parts := strings.SplitN(message, ":", 2) // Pisahkan tipe pesan dan isinya
	if len(parts) < 2 {
		fmt.Println("Invalid message format")
		return
	}

	msgType := parts[0]
	content := parts[1]

	switch msgType {
	case "join":
		handleJoin(conn, addr, content)
	case "left":
		handleLeft(conn, addr, content)
	default:
		handleMessage(conn, addr, parts[1])
	}
}

func handleJoin(conn *net.UDPConn, addr *net.UDPAddr, username string) {
	clients[addr.String()] = username
	fmt.Printf("%s joined\n", username)
	go broadcastMessage(conn, fmt.Sprintf("\n%s has joined the chat\n", username), addr)
}

func handleLeft(conn *net.UDPConn, addr *net.UDPAddr, content string) {
	if content == "exit" {
		username := clients[addr.String()]
		delete(clients, addr.String())
		fmt.Printf("%s left\n", username)
	go broadcastMessage(conn, fmt.Sprintf("\r %s has left the chat\n", username), addr)
	}
}

func handleMessage(conn *net.UDPConn, addr *net.UDPAddr, message string) {
	username, ok := clients[addr.String()]
	if ok {
	 go broadcastMessage(conn, fmt.Sprintf("[%s]: %s", username, message), addr)
	}
}

func broadcastMessage(conn *net.UDPConn, message string, senderAddr *net.UDPAddr) {
	for addr, username := range clients {
		if addr != senderAddr.String() {
			recipientAddr, err := net.ResolveUDPAddr("udp", addr)
			if err != nil {
				fmt.Println("Error resolving address:", err)
				continue
			}
			_, err = conn.WriteToUDP([]byte(message), recipientAddr)
			if err != nil {
				fmt.Println("Error broadcasting to", username, ":", err)
			}
		}
	}
}
