package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// Listen for incoming connections on port 8001
	ln, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	// Accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Read the message type
	var msgType int8
	err := binary.Read(conn, binary.LittleEndian, &msgType)
	if err != nil {
		log.Fatal(err)
	}

	switch msgType {
	case 0: // file upload
		go handleFileUpload(conn)
	case 1: // file download
		go handleFileDownload(conn)
	default:
		log.Fatalf("Unknown message type %d", msgType)
	}
}

func handleFileDownload(conn net.Conn) {
	defer conn.Close()

	// Read file name length
	var nameLen int64
	err := binary.Read(conn, binary.LittleEndian, &nameLen)
	if err != nil {
		log.Fatal(err)
	}

	// Read file name
	nameBuf := make([]byte, nameLen)
	_, err = io.ReadFull(conn, nameBuf)
	if err != nil {
		log.Fatal(err)
	}
	fileName := string(nameBuf)

	// Create the file on disk
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read file content size
	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		log.Fatal(err)
	}

	// Copy file content to disk
	_, err = io.CopyN(file, conn, size)
	if err != nil {
		log.Fatal(err)
	}
}

func handleFileUpload(conn net.Conn) {
	defer conn.Close()

	// Read file name length
	var nameLen int64
	err := binary.Read(conn, binary.LittleEndian, &nameLen)
	if err != nil {
		log.Fatal(err)
	}

	// Read file name
	nameBuf := make([]byte, nameLen)
	_, err = io.ReadFull(conn, nameBuf)
	if err != nil {
		log.Fatal(err)
	}
	fileName := string(nameBuf)

	// Open the file on disk
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Get file size
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := fi.Size()

	// Send file content size
	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		log.Fatal(err)
	}

	// Send file content
	_, err = io.CopyN(conn, file, size)
	if err != nil {
		log.Fatal(err)
	}
}
