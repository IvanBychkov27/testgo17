package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func main() {
	log.Println("start...")
	timeOut := time.Second * 5

	//url := "http://4000.99.adscompass.ru"
	//url := "http://3000.99.adscompass.ru"
	url := "http://eth0.me"

	for {
		redURL, errCheckRedirect := checkRedirect(url, timeOut)
		if errCheckRedirect != nil {
			log.Println("error check redirect:", errCheckRedirect.Error())
			os.Exit(1)
		}

		if redURL != url {
			//log.Printf("redirect to url: %s\n", redURL)
			url = redURL
			continue
		}
		//log.Printf("no redirect")
		break
	}

	title, body, errGetData := getData(url)
	if errGetData != nil {
		log.Println("error get data:", errGetData.Error())
		os.Exit(1)
	}

	log.Println("title:", title)
	log.Println("len body:", len(body))
	log.Println("body:", body)

	log.Println("done")
}

func checkRedirect(url string, timeOut time.Duration) (string, error) {
	log.Println("check redirect url", url)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	chURL := make(chan string, 1)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventExceptionThrown:
			chURL <- ev.ExceptionDetails.URL
			ev = nil
		}
	})

	errRun := chromedp.Run(ctx, chromedp.Navigate(url))
	if errRun != nil {
		return url, errRun
	}

	select {
	case url = <-chURL:
	case <-time.After(timeOut):
	}
	return url, nil
}

func getData(url string) (title, body string, err error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err = chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.Title(&title),
		chromedp.OuterHTML("body", &body),
		//chromedp.OuterHTML("html", &html),
	)
	return
}
