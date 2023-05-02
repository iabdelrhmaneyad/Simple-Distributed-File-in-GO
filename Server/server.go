package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// Constants to represent the 4 storage servers
const (
	StorageServer1 = "localhost:8001"
	StorageServer2 = "localhost:8002"
	StorageServer3 = "localhost:8003"
	StorageServer4 = "localhost:8004"
)

type FileServer struct{}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("server listening on 8000")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.readloop(conn)
	}
}

func (fs *FileServer) readloop(conn net.Conn) {

	// Read the message type
	var msgType int8
	err := binary.Read(conn, binary.LittleEndian, &msgType)
	if err != nil {
		log.Fatal(err)
	}

	switch msgType {
	case 0: // file upload
		fs.handleFileUpload(conn)
	case 1: // file download
		fs.handleFileDownload(conn)
	default:
		log.Fatalf("Unknown message type %d", msgType)
	}
}

func (fs *FileServer) handleFileUpload(conn net.Conn) {
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

	// Read file content
	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		log.Fatal(err)
	}

	contentBuf := new(bytes.Buffer)
	_, err = io.CopyN(contentBuf, conn, size)
	if err != nil {
		log.Fatal(err)
	}

	// Divide the received file into 4 chunks
	chunkSize := size / 4
	chunk1 := contentBuf.Next(int(chunkSize))
	chunk2 := contentBuf.Next(int(chunkSize))
	chunk3 := contentBuf.Next(int(chunkSize))
	chunk4 := contentBuf.Next(int(size - chunkSize*3))

	// Send each chunk to the respective storage server
	go fs.sendToServer(chunk1, "ch1-"+fileName, StorageServer1)
	go fs.sendToServer(chunk2, "ch2-"+fileName, StorageServer2)
	go fs.sendToServer(chunk3, "ch3-"+fileName, StorageServer3)
	go fs.sendToServer(chunk4, "ch4-"+fileName, StorageServer4)
}

func (fs *FileServer) handleFileDownload(conn net.Conn) {
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

	// Send file chunks to the client
	chunk1, err := fs.getFromServer("ch1-"+fileName, StorageServer1)
	if err != nil {
		log.Fatal(err)
	}
	chunk2, err := fs.getFromServer("ch2-"+fileName, StorageServer2)
	if err != nil {
		log.Fatal(err)
	}
	chunk3, err := fs.getFromServer("ch3-"+fileName, StorageServer3)
	if err != nil {
		log.Fatal(err)
	}
	chunk4, err := fs.getFromServer("ch4-"+fileName, StorageServer4)
	if err != nil {
		log.Fatal(err)
	}
	// Concatenate the chunks into a single file
	fileContent := bytes.Join([][]byte{chunk1, chunk2, chunk3, chunk4}, []byte(""))

	// Send file size
	fileSize := int64(len(fileContent))
	err = binary.Write(conn, binary.LittleEndian, fileSize)
	if err != nil {
		log.Fatal(err)
	}

	// Send file content
	_, err = conn.Write(fileContent)
	if err != nil {
		log.Fatal(err)
	}

}

func (fs *FileServer) sendToServer(content []byte, fileName string, server string) {
	// Connect to the storage server
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send the message type to the server
	err = binary.Write(conn, binary.LittleEndian, int8(1))
	if err != nil {
		log.Fatal(err)
	}
	// Send file name length
	nameLen := int64(len(fileName))
	err = binary.Write(conn, binary.LittleEndian, nameLen)
	if err != nil {
		log.Fatal(err)
	}

	// Send file name
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		log.Fatal(err)
	}

	// Send file content size
	size := int64(len(content))
	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		log.Fatal(err)
	}

	// Send file content
	_, err = conn.Write(content)
	if err != nil {
		log.Fatal(err)
	}
}
func (fs *FileServer) getFromServer(fileName string, server string) ([]byte, error) {
	// Connect to the storage server
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Send the message type to the server
	err = binary.Write(conn, binary.LittleEndian, int8(0))
	if err != nil {
		log.Fatal(err)
	}
	// Send file name length
	nameLen := int64(len(fileName))
	err = binary.Write(conn, binary.LittleEndian, nameLen)
	if err != nil {
		return nil, err
	}

	// Send file name
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		return nil, err
	}

	// Read file content size
	var size int64
	err = binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		return nil, err
	}

	// Read file content
	contentBuf := new(bytes.Buffer)
	_, err = io.CopyN(contentBuf, conn, size)
	if err != nil {
		return nil, err
	}

	return contentBuf.Bytes(), nil
}

func main() {
	server := &FileServer{}
	server.start()
}
