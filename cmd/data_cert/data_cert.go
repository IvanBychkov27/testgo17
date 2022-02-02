package main

import (
	"crypto/tls"
	"log"
	"time"
)

// получение даты окончания срока сертификата домена
func main() {
	conn, err := tls.Dial("tcp", "yandex.com:443", &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	log.Printf("expiry: %s", time.Now().Format(time.RFC822))
	log.Printf("expiry: %d", time.Now().Unix())

	// в нулевом элементе массива находится нужная нам дата окончания срока сертификата
	for idx, cert := range conn.ConnectionState().PeerCertificates {
		log.Printf("Cert #%d", idx)
		//log.Printf("issuer: %s", cert.Issuer.String())
		log.Printf("expiry: %s", cert.NotAfter.Format(time.RFC822))
		log.Printf("expiry: %d", cert.NotAfter.Unix())
		//log.Printf("dns names: %s", strings.Join(cert.DNSNames, ","))
	}
}
