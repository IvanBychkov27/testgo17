// https://zalinux.ru/?p=6293
// https://stackoverflow.com/questions/2886719/unix-sockets-in-go

package main

import (
	"io"
	"log"
	"net"
	"time"
)

func reader(r io.Reader) {
	buf := make([]byte, 128)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			log.Println("error:", err.Error())
			return
		}
		log.Println("client:", string(buf[:n]))
	}
}

func main() {
	//addr := "/home/ivan/tmp/unix.sock"
	addr := "/home/ivan/tmp/echo.sock"

	//client01(addr)
	client02(addr)

}

func client02(addr string) {
	// Get unix socket address based on file path
	uaddr, err := net.ResolveUnixAddr("unix", addr)
	if err != nil {
		log.Println("error", err.Error())
		return
	}

	// Connect server with unix socket
	unixClient, err := net.DialUnix("unix", nil, uaddr)
	if err != nil {
		log.Println("error", err.Error())
		return
	}
	defer unixClient.Close()

	go reader(unixClient)

	for {
		//msg := "time " + time.Now().Format("15:04:05")
		msg := "0123456789"
		//msg := buildMSG()
		_, err := unixClient.Write([]byte(msg))
		if err != nil {
			log.Println("error:", err.Error())
			break
		}
		time.Sleep(time.Second * 2)
	}

}

func client01(addr string) {
	unixClient, err := net.Dial("unix", addr)
	if err != nil {
		log.Println("error:", err.Error())
		return
	}
	defer unixClient.Close()

	go reader(unixClient)

	for {
		//msg := "time " + time.Now().Format("15:04:05")
		msg := buildMSG()
		_, err := unixClient.Write([]byte(msg))
		if err != nil {
			log.Println("error:", err.Error())
			break
		}
		time.Sleep(time.Second * 2)
	}
}

func buildMSG() string {
	msg := `binder_ch_endpoint_request{endpoint_id="1"}
`
	msg += `binder_ch_endpoint_request_bysource{endpoint_id="1",source_id="123"}
`
	msg += `binder_ch_endpoint_request_bycountry{endpoint_id="1",country="RU"}
`
	return msg
}
