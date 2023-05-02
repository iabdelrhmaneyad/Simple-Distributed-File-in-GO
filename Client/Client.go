package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func sendFile(fileName string) error {
	// Check if file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return err
	}

	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	size := fileInfo.Size()

	// Connect to the server
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the message type to the server
	err = binary.Write(conn, binary.LittleEndian, int8(0))
	if err != nil {
		return err
	}

	// Send the file name length and file name to the server
	name := []byte(filepath.Base(fileName))
	nameLen := int64(len(name))
	err = binary.Write(conn, binary.LittleEndian, nameLen)
	if err != nil {
		return err
	}
	_, err = conn.Write(name)
	if err != nil {
		return err
	}

	// send the file size to the server
	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		return err
	}

	// Send the file content over the connection
	var totalSent int64
	for totalSent < size {
		n, err := io.CopyN(conn, file, size-totalSent)
		if err != nil {
			return err
		}
		totalSent += n
	}

	fmt.Printf("Written %d bytes over connection\n", totalSent)
	return nil
}

func receiveFile(fileName string) error {
	// Connect to the server
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send the message type to the server
	err = binary.Write(conn, binary.LittleEndian, int8(1))
	if err != nil {
		return err
	}

	// Send the file name length and file name to the server
	name := []byte(filepath.Base(fileName))
	nameLen := int64(len(name))
	err = binary.Write(conn, binary.LittleEndian, nameLen)
	if err != nil {
		return err
	}
	_, err = conn.Write(name)
	if err != nil {
		return err
	}

	// Read the file size from the server
	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		return err
	}

	// Create a new file to write to
	newFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// Receive the file content from the server
	var totalReceived int64
	for totalReceived < size {
		n, err := io.CopyN(newFile, conn, size-totalReceived)
		if err != nil {
			return err
		}
		totalReceived += n
	}

	fmt.Printf("Received %d bytes over connection\n", totalReceived)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client <upload/download> <filename>")
		return
	}

	serviceType := os.Args[1]
	// Check if the service type is valid
if serviceType != "upload" && serviceType != "download" {
	fmt.Println("Invalid service type")
	return
	}
	
	fileName := os.Args[2]
	if serviceType == "upload" {
	err := sendFile(fileName)
	if err != nil {
	fmt.Printf("Error sending file: %v\n", err)
	return
	}
	fmt.Println("File sent successfully")
	} else if serviceType == "download" {
	err := receiveFile(fileName)
	if err != nil {
	fmt.Printf("Error receiving file: %v\n", err)
	return
	}
	fmt.Println("File received successfully")
	}
	}