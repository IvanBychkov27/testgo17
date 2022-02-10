package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"math/rand"
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
	fileName := "/home/ivan/projects/testgo17/headless_browsers/proxy_auto_ping/data/servers_020222.json"
	data := getDataIP(getData(fileName))
	log.Println("count ip :", len(data))
	if len(data) == 0 {
		return
	}

	autoPing(data)
}

func autoPing(data []*dataIP) {
	rand.Seed(time.Now().UnixNano())

	uas := []string{
		"Mozilla/5.0 (Unknown; Linux i686) AppleWebKit/534.34 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/534.34",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36",
		"Mozilla/5.0 (X11; U; U; Linux x86_64; vi-vn) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.136 Safari/537.36 Puffin/9.2.0.50581AV",
	}

	interval := 30

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

		ua := uas[rand.Intn(3)]

		s, keyCode := proxyPing(ip, flag, ua)

		code[keyCode]++
		res += s
		log.Print(s)

		select {
		case <-ticker.C:
		case <-ctx.Done():
			s = ""
			for key, val := range code {
				s += fmt.Sprintf("Code: %3d = %d \n", key, val)
			}
			res += s
			date := time.Now().Format("20060102_150405")
			fileNameSaved := "/home/ivan/projects/testgo17/headless_browsers/proxy_auto_ping/log/log" + date + ".txt"
			saveDataFile(fileNameSaved, []byte(res))
			return
		}
	}
}

func proxyPing(ip, flag, ua string) (string, int) {
	proxyURL := "https://" + ip + ":89"
	log.Printf("[%s] %16s \n", flag, ip)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	ctx, cancel = chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = chromedp.NewExecAllocator(
		ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ProxyServer(proxyURL),
			chromedp.Flag("ignore-certificate-errors", "1"),
			chromedp.UserAgent(ua),
		)...,
	)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	defer chromedp.Cancel(ctx)

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		go func() {
			switch ev := ev.(type) {
			case *fetch.EventAuthRequired:
				c := chromedp.FromContext(ctx)
				execCtx := cdp.WithExecutor(ctx, c.Target)

				resp := &fetch.AuthChallengeResponse{
					Response: fetch.AuthChallengeResponseResponseProvideCredentials,
					Username: "ysETCpBC8JvFzJnt7SjsJxJC",
					Password: "qRhXR6k8yW4pcW71c34ReDW3",
				}

				err := fetch.ContinueWithAuth(ev.RequestID, resp).Do(execCtx)
				if err != nil {
					log.Print(err)
				}

			case *fetch.EventRequestPaused:
				c := chromedp.FromContext(ctx)
				execCtx := cdp.WithExecutor(ctx, c.Target)
				err := fetch.ContinueRequest(ev.RequestID).Do(execCtx)
				if err != nil {
					log.Print(err)
				}
			}
		}()
	})

	offerURL := "http://4000.99.adscompass.ru"

	var body, title string
	var code int

	resp, errRunResponse := chromedp.RunResponse(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Navigate(offerURL),
		chromedp.Sleep(time.Second),
		chromedp.Title(&title),
		chromedp.OuterHTML("body", &body),
	)
	if errRunResponse != nil {
		log.Print(errRunResponse)
		return fmt.Sprintf("[%s] %16s -> code: %d \n", flag, ip, code), code
	}
	code = int(resp.Status)

	//log.Println("title:", title)
	//
	//indexKey := "IndexKey=testKey"
	//if strings.Contains(body, indexKey) {
	//	log.Println("key is verified")
	//} else {
	//	log.Println("key not found")
	//}

	return fmt.Sprintf("[%s] %16s -> code: %d  [UA: %s]\n", flag, ip, code, ua), code
}

//-----------------------

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
