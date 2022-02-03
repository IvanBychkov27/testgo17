package main

/*
99.adscompass.ru
2000.99.adscompass.ru
3000.99.adscompass.ru
4000.99.adscompass.ru
*/

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type dataIP struct {
	IP   string
	Flag string
}

func main() {
	//serverNordVPN()

	fileName := "/home/ivan/projects/testgo17/cmd/vpn/data/servers_020222.json"
	data := getDataIP(getData(fileName))
	log.Println("count ip :", len(data))
	if len(data) == 0 {
		return
	}

	autoPing(data)
}

func autoPing(data []*dataIP) {
	rand.Seed(time.Now().UnixNano())

	interval := 10

	res := fmt.Sprintf("start auto ping interval [%d sec] \n", interval)
	log.Print(res)
	defer log.Println("done")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	ticker := time.NewTicker(time.Second * time.Duration(interval))
	defer ticker.Stop()

	n := len(data)
	code := make(map[int]int)
	for {
		i := rand.Intn(n)
		ip := data[i].IP
		flag := data[i].Flag

		s, keyCode := pingIP(ip, flag)

		code[keyCode]++
		res += s
		log.Print(s)

		select {
		case <-ticker.C:
		case <-ctx.Done():
			s = ""
			for key, val := range code {
				s += fmt.Sprintf("Code: %d = %3d \n", key, val)
			}
			res += s
			date := time.Now().Format("20060102_150405")
			fileNameSaved := "/home/ivan/projects/testgo17/cmd/vpn/log/log" + date + ".txt"
			saveDataFile(fileNameSaved, []byte(res))
			return
		}
	}
}

func pingIP(ip, flag string) (string, int) {
	user := "ysETCpBC8JvFzJnt7SjsJxJC"
	password := "qRhXR6k8yW4pcW71c34ReDW3"

	protocol := "https"
	port := "89"

	addrNordVPN := fmt.Sprintf("%s://%s:%s@%s:%s", protocol, user, password, ip, port)

	urlAddrNordVPN, err := url.Parse(addrNordVPN)
	if err != nil {
		log.Println(err)
	}

	transport := &http.Transport{
		Proxy:           http.ProxyURL(urlAddrNordVPN),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	urlData := "http://4000.99.adscompass.ru/" + flag

	request, err := http.NewRequest("GET", urlData, nil)
	if err != nil {
		return fmt.Sprintf("[%s] %16s -> error: %s \n", flag, ip, err.Error()), 0
	}

	resp, err := client.Do(request)
	if err != nil {
		return fmt.Sprintf("[%s] %16s -> error: %s \n", flag, ip, err.Error()), 0
	}

	return fmt.Sprintf("[%s] %16s -> code: %d \n", flag, ip, resp.StatusCode), resp.StatusCode
}

func saveDataFile(fileName string, data []byte) {
	df, errCreateFile := os.Create(fileName)
	if errCreateFile != nil {
		fmt.Errorf("error create file, %w", errCreateFile)
		return
	}
	defer df.Close()

	_, errWrite := df.Write(data)
	if errWrite != nil {
		fmt.Errorf("error write data %w", errWrite)
		return
	}
	log.Printf("file saved: %s\n", fileName)
}

//-----------------

type serversNordVPNData struct {
	ID             int          `json:"id"`
	IPAddress      string       `json:"ip_address"`
	SearchKeywords []string     `json:"search_keywords"`
	Categories     []categories `json:"categories"`
	Name           string       `json:"name"`
	Domain         string       `json:"domain"`
	Price          int          `json:"price"`
	Location       location     `json:"location"`
	Flag           string       `json:"flag"`
	Country        string       `json:"country"`
	Load           int          `json:"load"`
	Features       features     `json:"features"`
}

type location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type categories struct {
	Name string `json:"name"`
}

type features struct {
	Ikev2               bool `json:"ikev2"`                 // 1
	OpenvpnUDP          bool `json:"openvpn_udp"`           // 2
	OpenvpnTCP          bool `json:"openvpn_tcp"`           // 3
	Socks               bool `json:"socks"`                 // 4
	Proxy               bool `json:"proxy"`                 // 5
	Pptp                bool `json:"pptp"`                  // 6
	L2tp                bool `json:"l2tp"`                  // 7
	OpenvpnXorUDP       bool `json:"openvpn_xor_udp"`       // 8
	OpenvpnXorTCP       bool `json:"openvpn_xor_tcp"`       // 9
	ProxyCybersec       bool `json:"proxy_cybersec"`        // 10
	ProxySSL            bool `json:"proxy_ssl"`             // 11
	ProxySSLCybersec    bool `json:"proxy_ssl_cybersec"`    // 12
	Ikev2V6             bool `json:"ikev2v6"`               // 13
	OpenvpnUDPv6        bool `json:"openvpn_udp_v6"`        // 14
	OpenvpnTCPv6        bool `json:"openvpn_tcp_v6"`        // 15
	WireguardUDP        bool `json:"wireguard_udp"`         // 16
	OpenvpnUdpTlsCrypt  bool `json:"openvpn_udp_tls_crypt"` // 17
	OpenvpnTcpTlsCrypt  bool `json:"openvpn_tcp_tls_crypt"` // 18
	OpenvpnDedicatedUdp bool `json:"openvpn_dedicated_udp"` // 19
	OpenvpnDedicatedTcp bool `json:"openvpn_dedicated_tcp"` // 20
	Skylark             bool `json:"skylark"`               // 21
	MeshRelay           bool `json:"mesh_relay"`            // 22
}

func getData(fileName string) []serversNordVPNData {
	file, errOpen := ioutil.ReadFile(fileName)
	if errOpen != nil {
		fmt.Println("error open file stamp data", errOpen.Error())
		return nil
	}

	data := make([]serversNordVPNData, 0)

	errJSON := json.Unmarshal(file, &data)
	if errJSON != nil {
		fmt.Println("error unmarshal data", errJSON.Error())
		return nil
	}

	return data
}

func getDataIP(data []serversNordVPNData) []*dataIP {
	res := make([]*dataIP, 0, len(data))

	for _, d := range data {
		if !d.Features.ProxySSL {
			continue
		}
		r := &dataIP{
			IP:   d.IPAddress,
			Flag: d.Flag,
		}
		res = append(res, r)
	}
	return res
}

//==================

func serverNordVPN() {
	log.Println("start...")

	user := "ysETCpBC8JvFzJnt7SjsJxJC"
	password := "qRhXR6k8yW4pcW71c34ReDW3"

	// Socks
	//protocol, ip, port, flag := "socks5", "109.202.99.35", "1080", "NL"
	//protocol, ip, port, flag := "socks5", "165.231.210.171", "1080", "US"
	//protocol, ip, port, flag := "socks5", "196.196.244.3", "1080", "SE"
	//protocol, ip, port, flag := "socks5", "196.196.192.11", "1080", "IE"

	// ProxySSL
	//protocol, ip, port, flag := "https", "172.83.40.219", "89", "CA"
	//protocol, ip, port, flag := "https", "194.99.105.100", "89", "PL"
	//protocol, ip, port, flag := "https", "82.102.20.236", "89", "DK"
	//protocol, ip, port, flag := "https", "82.102.19.137", "89", "BE"
	//protocol, ip, port, flag := "https", "185.206.225.196", "89", "NO"
	//protocol, ip, port, flag := "https", "185.189.114.28", "89", "HU"
	//protocol, ip, port, flag := "https", "185.245.87.59", "89", "US"
	//protocol, ip, port, flag := "https", "185.216.34.100", "89", "AT"
	//protocol, ip, port, flag := "https", "81.92.202.11", "89", "GB"
	//protocol, ip, port, flag := "https", "89.238.186.244", "89", "CZ"
	protocol, ip, port, flag := "https", "82.102.18.252", "89", "FR"

	addrNordVPN := fmt.Sprintf("%s://%s:%s@%s:%s", protocol, user, password, ip, port)
	urlAddrNordVPN, err := url.Parse(addrNordVPN)
	if err != nil {
		log.Println(err)
	}
	log.Println("urlAddrNordVPN =", urlAddrNordVPN)

	transport := &http.Transport{
		Proxy:           http.ProxyURL(urlAddrNordVPN),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	//generating the HTTP GET request
	//urlData := "http://eth0.me"
	urlData := "http://4000.99.adscompass.ru/" + flag

	request, err := http.NewRequest("GET", urlData, nil)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("StatusCode:", resp.StatusCode)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error body:", err)
		return
	}

	log.Println("Body:")
	log.Println(string(data))

	log.Println("Done...")
}
