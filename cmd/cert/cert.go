// Получение Ключа и сертификата от Let's Encrypt

package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("a.letestx.ru"), // Наш домен
		Cache:      autocert.DirCache("cmd/cert/certs"),    // Папка для хранения сертификатов
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Let's Encrypt"))
	})

	server := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))

	log.Fatal(server.ListenAndServeTLS("certFile", "keyFile")) // Ключ и сертификат поступают от Let's Encrypt
}
