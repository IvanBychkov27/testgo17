package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func main() {
	timeOut := time.Second * 5

	//url := "http://4000.99.adscompass.ru"
	url := "http://3000.99.adscompass.ru"
	//url := "http://eth0.me"

	fromProxy(url, timeOut)
}

func fromProxy(url string, timeOut time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(context.Background())
	defer cancel()

	//ua := "Mozilla/5.0 (Unknown; Linux i686) AppleWebKit/534.34 (KHTML, like Gecko) Chrome/20.0.1132.57 Safari/534.34"
	//ua := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36"
	//ua := "Mozilla/5.0 (X11; U; U; Linux x86_64; vi-vn) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.136 Safari/537.36 Puffin/9.2.0.50581AV"
	ctx, cancel = chromedp.NewExecAllocator(
		ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ProxyServer("https://172.83.40.219:89"), // CA
			//chromedp.ProxyServer("https://37.19.218.139:89"), // UA
			chromedp.Flag("ignore-certificate-errors", "1"),
			//chromedp.UserAgent(ua),
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

	var body, title string
	if err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Navigate(url),
		chromedp.Sleep(timeOut),
		chromedp.Title(&title),
		chromedp.OuterHTML("body", &body),
	); err != nil {
		log.Fatal(err)
	}

	log.Println("title:", title)

	if len(body) < 200 {
		log.Println(body)
	} else {
		log.Println("len body:", len(body))
	}

}
