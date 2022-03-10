package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func main() {
	getBody()
}

func getBody() {
	ctxWT, cancelWT := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelWT()

	ctx, cancel := chromedp.NewContext(ctxWT)
	defer cancel()
	defer chromedp.Cancel(ctx)

	url := "http://4000.99.adscompass.ru"

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		go func() {
			switch ev := ev.(type) {
			case *fetch.EventAuthRequired:
				c := chromedp.FromContext(ctx)
				execCtx := cdp.WithExecutor(ctx, c.Target)

				resp := &fetch.AuthChallengeResponse{
					Response: fetch.AuthChallengeResponseResponseProvideCredentials,
				}

				err := fetch.ContinueWithAuth(ev.RequestID, resp).Do(execCtx)
				if err != nil {
					fmt.Println("error continue with auth:", err.Error())
				}

			case *fetch.EventRequestPaused:
				c := chromedp.FromContext(ctx)
				execCtx := cdp.WithExecutor(ctx, c.Target)
				err := fetch.ContinueRequest(ev.RequestID).Do(execCtx)
				if err != nil {
					log.Println("error continue request:", err.Error())
				}
			}
		}()
	})

	var ids []cdp.NodeID
	var body string
	resp, errRunResponse := chromedp.RunResponse(ctx,
		fetch.Enable().WithHandleAuthRequests(true),
		chromedp.Navigate(url),
		chromedp.Sleep(time.Second*1),
		//chromedp.OuterHTML("body", &body),
		chromedp.NodeIDs(`document`, &ids, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			body, err = dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
			return err
		}),
	)
	if errRunResponse != nil {
		log.Println("error:", errRunResponse.Error())
		return
	}

	code := int(resp.Status)

	log.Println("Outer HTML:")
	log.Println(body)
	log.Println("status code:", code)
}

func getBody1() {
	ctxWT, cancelWT := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelWT()

	ctx, cancel := chromedp.NewContext(ctxWT)
	defer cancel()

	//url := "http://poker-iv.herokuapp.com/"
	//url := "https://ya.ru/"
	url := "http://4000.99.adscompass.ru"

	var ids []cdp.NodeID
	var body string
	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.NodeIDs(`document`, &ids, chromedp.ByJSPath),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			body, err = dom.GetOuterHTML().WithNodeID(ids[0]).Do(ctx)
			return err
		}),
	); err != nil {
		log.Println("error:", err.Error())
	}

	log.Println("Outer HTML:")
	log.Println(body)
}
