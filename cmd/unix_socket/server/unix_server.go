// https://zalinux.ru/?p=6293
// https://stackoverflow.com/questions/2886719/unix-sockets-in-go

package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
)

func echoServer(unixClient net.Conn) {
	defer unixClient.Close()

	for {
		// читаем инфо из канала socket
		buf := make([]byte, 10)
		n, err := unixClient.Read(buf)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("error read:", err.Error())
			}
			return
		}

		data := buf[:n]
		log.Println("server:", string(data))

		// пишем обратную инфо в канал socket
		_, err = unixClient.Write(data)
		if err != nil {
			log.Println("error:", err.Error())
			return
		}
	}
}

func main() {
	addr := "/home/ivan/tmp/echo.sock"
	//server01(addr)
	server02(addr)
}

func server02(addr string) {
	os.Remove(addr)

	// Get unix socket address based on file path
	uaddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		log.Println("error", err.Error())
		return
	}

	// Listen on the address
	unixListener, err := net.ListenUnix("unix", uaddr)
	if err != nil {
		log.Println("error", err.Error())
		return
	}
	defer unixListener.Close()

	for {
		uconn, err := unixListener.AcceptUnix()
		if err != nil {
			log.Println("error", err.Error())
			continue
		}

		//data, errParse := parseRequest(uconn)
		//if errParse != nil {
		//	continue
		//}
		//log.Println("server:", string(data))
		go echoServer(uconn)
	}
}

func server01(addr string) {
	os.Remove(addr)

	ln, err := net.Listen("unix", addr)
	if err != nil {
		log.Println("error:", err.Error())
		return
	}

	for {
		connData, err := ln.Accept()
		if err != nil {
			break
		}

		go echoServer(connData)
	}
}

func parseRequest(conn *net.UnixConn) ([]byte, error) {
	var reqLen uint32
	lenBytes := make([]byte, 4)
	if _, err := io.ReadFull(conn, lenBytes); err != nil {
		return nil, err
	}

	lenBuf := bytes.NewBuffer(lenBytes)
	if err := binary.Read(lenBuf, binary.BigEndian, &reqLen); err != nil {
		return nil, err
	}

	reqBytes := make([]byte, reqLen)
	_, err := io.ReadFull(conn, reqBytes)
	if err != nil {
		return nil, err
	}

	return reqBytes, nil
}
