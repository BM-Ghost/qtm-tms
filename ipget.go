package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func handleUpdateMessage(parts []string, conn net.Conn) {
	if len(parts) != 5 {
		log.Println("Invalid message format for update message: expected 4 sections")
		return
	}

	//transactionType := parts[0]
	fullSerialNumber := parts[1]
	//lastDownloadDate := parts[2]
	//batteryPercentage := parts[3]

	// Log the request with date, time, serial number, and request type
	log.Printf("[%s] Serial Number: %s - Received Update Message Request\n", time.Now().Format("2006-01-02 15:04:05"), fullSerialNumber)

	// Handle the update message here, for example, you can log or process the data.

	// Construct the response
	serverPublicIP := "3.6.122.107"
	thisPort := "12784"
	numberOfFiles := "1"
	dirPath := "/dummy"
	log.Printf("1|%s|%s|%s|%s|%s|test\n", serverPublicIP, thisPort, numberOfFiles, dirPath, fullSerialNumber)
	response := fmt.Sprintf("1|%s|%s|%s|%s|%s|test", serverPublicIP, thisPort, numberOfFiles, dirPath, fullSerialNumber)

	// Send the response to the client
	_, err := conn.Write([]byte(response + "\n"))
	if err != nil {
		log.Println("Failed to send response to the client:", err)
	}
}

func handleUpdateNotification(parts []string) {
	log.Print(len(parts))
	if len(parts) < 5 {
		log.Println("Invalid message format for update notification: expected 5 sections")
		return
	}

	transactionType := parts[0]
	countOfDownloadedFiles := parts[1]
	status := parts[2]
	timeTaken := parts[3]
	terminalSerial := parts[4]
	log.Printf("Trac type" + transactionType + "with a status of " + status + "files downloaded " + countOfDownloadedFiles + "timetaken " + timeTaken)
	// Log the request with date, time, serial number, and request type
	log.Printf("[%s] Serial Number: %responses - Received Update Notification Request\n", time.Now().Format("2023-11-07 10:30:05"), terminalSerial)

	// Handle the update notification here, for example, you can log or process the data.
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
		}
	}(conn)

	reader := bufio.NewReader(conn)
	data, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to read from client:", err)
		return
	}

	// Split the received data by the pipe character '|'
	parts := strings.Split(strings.TrimSpace(data), "|")

	// Check the transaction type (first token)
	transactionType := parts[0]

	if transactionType == "1" {
		// If the transaction type is 1, handle it as an update message.
		handleUpdateMessage(parts, conn)

	} else if transactionType == "2" {
		// If the transaction type is 2, handle it as an update notification.
		handleUpdateNotification(parts)
	} else {
		log.Println("Unsupported transaction type:", transactionType)
	}

}

func main() {
	listener, err := net.Listen("tcp", "localhost:8099")
	if err != nil {
		log.Fatalln("Failed to start server:", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
		}
	}(listener)

	log.Println("Server started on localhost:8099")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		println("Handling socker...")
		go handleConnection(conn)
	}
}
