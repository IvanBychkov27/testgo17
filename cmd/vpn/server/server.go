package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	address := ":4000"
	log.Println("start server", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("error listen address, %v", err)
		os.Exit(1)
	}
	server := &http.Server{
		Handler: http.HandlerFunc(handler),
	}
	errServe := server.Serve(ln)
	if errServe != nil {
		log.Printf("error serve, %v", errServe)
		os.Exit(1)
	}
	log.Printf("done")
}
func handler(rw http.ResponseWriter, req *http.Request) {
	//log.Printf("%s %s [%s]", req.Method, req.RequestURI, req.RemoteAddr)
	log.Printf("%s %s [%s]", req.Method, req.RequestURI, req.Header["X-Forwarded-For"][0])

	d := `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Server 4000</title>
</head>
<body>
    Server :4000
	IndexKey=testKey
</body>
</html>
`
	rw.Write([]byte(d))
}
