package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"log"
)

func main() {
	url := "http://4000.99.adscompass.ru"
	//url := "http://3000.99.adscompass.ru"
	//url := "http://eth0.me"

	fromProxy(url)
}

func fromProxy(url string) {
	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//defer cancel()

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = chromedp.NewExecAllocator(
		ctx,
		append(chromedp.DefaultExecAllocatorOptions[:],
			//chromedp.ProxyServer("https://172.83.40.219:89"), // CA
			chromedp.ProxyServer("https://37.19.218.139:89"), // UA
			chromedp.Flag("ignore-certificate-errors", "1"),
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

	var body string
	if err := chromedp.Run(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Navigate(url),
		chromedp.OuterHTML("body", &body),
	); err != nil {
		log.Fatal(err)
	}

	log.Println(body)
}
