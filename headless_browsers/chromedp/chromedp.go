// https://github.com/chromedp/chromedp
// https://itnext.io/scrape-the-web-faster-in-go-with-chromedp-c94e43f116ce
package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func main() {
	log.Println("start...")
	//ExampleTitle_Redirect()
	//ExampleChromeEDP_StatusCode()

	//ExampleChromeEDP()
	//ExampleChromedp_Proxy()

	//getBody()
	listenTarget()
}

func listenTarget() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	url := "http://3000.99.adscompass.ru"

	gotException := make(chan bool, 1)
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		//case *runtime.EventConsoleAPICalled:
		//	log.Printf("* console.%s call:\n", ev.Type)
		//	for _, arg := range ev.Args {
		//		log.Printf("%s - %s\n", arg.Type, arg.Value)
		//	}
		case *runtime.EventExceptionThrown:
			log.Printf("redirect to url: %s\n", ev.ExceptionDetails.URL)
			gotException <- true
		}
	})

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		log.Fatal(err)
	}
	<-gotException
}

func getBody() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	//url := "http://poker-iv.herokuapp.com/"
	//url := "https://yandex.ru/"
	//url := "http://4000.99.adscompass.ru"
	url := "http://3000.99.adscompass.ru"

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

func ExampleChromedp_Proxy() {
	o := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("--proxy-server=https://ysETCpBC8JvFzJnt7SjsJxJC:qRhXR6k8yW4pcW71c34ReDW3@172.83.40.219:89"),
	)
	cx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(cx)
	defer cancel()

	url := "http://4000.99.adscompass.ru"
	//url := "http://eth0.me"

	var title string

	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.Title(&title),
	)

	if err != nil {
		log.Println("error:", err.Error())
	}

	log.Println("title:", title)
}

func ExampleChromeEDP_StatusCode() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	//url := "http://poker-iv.herokuapp.com/"
	//url := "https://yandex.ru/"
	//url := "http://4000.99.adscompass.ru"
	url := "http://3000.99.adscompass.ru"

	resp, errRunResponse := chromedp.RunResponse(ctx, chromedp.Navigate(url))
	if errRunResponse != nil {
		log.Fatal(errRunResponse)
	}
	log.Println("status code:", resp.Status)
}

func ExampleChromeEDP() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	//url := "http://poker-iv.herokuapp.com/"
	//url := "https://yandex.ru/"
	url := "http://3000.99.adscompass.ru" // redirect in poker-iv.herokuapp.com

	var title string

	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.Title(&title),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("title:", title)

	//resp, errRunResponse := chromedp.RunResponse(ctx, chromedp.Navigate(url))
	//if errRunResponse != nil {
	//	log.Fatal(errRunResponse)
	//}
	//log.Println("status code:", resp.Status)
}

// =====================================================
func writeHTML(content string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, strings.TrimSpace(content))
	})
}

func ExampleTitle() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(writeHTML(`
<head>
	<title>fancy website title</title>
</head>
<body>
	<div id="content"></div>
</body>
	`))
	defer ts.Close()

	var title string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.Title(&title),
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println(title)

	// Output:
	// fancy website title
}

func ExampleTitle_Redirect() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(writeHTML(`
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
</head>
<body>
    Hello World

    <script>
        (function(){
            document.location.href = "https://ya.ru"
        }())
    </script>
</body>
</html>
	`))
	defer ts.Close()

	var title string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.Title(&title),
	); err != nil {
		log.Fatal(err)
	}
	log.Println("title:", title)

	// Output:
	// fancy website title
}
