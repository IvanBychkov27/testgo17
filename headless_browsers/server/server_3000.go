package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	address := ":3000"
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

// --- 1 ---
func handler(rw http.ResponseWriter, req *http.Request) {
	page := `<!doctype html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport"
        content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Document3000</title>
</head>
<body>

 Redirect in poker-iv

  <script>
      setTimeout (function(){
          document.location.href = "http://poker-iv.herokuapp.com/"
      }, 3000)
  </script>
</body>
</html>
`
	log.Println("UA:", req.UserAgent())
	rw.Write([]byte(page))
}

// --- 2 ---
//func handler(rw http.ResponseWriter, req *http.Request) {
//	page := `<!DOCTYPE html>
//<html>
//<head>
//  <title></title>
//  <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
//  <meta name="googlebot" content="noindex">
//  <meta http-equiv="x-dns-prefetch-control" content="on">
//  <link rel="dns-prefetch" href="http://poker-iv.herokuapp.com/">
//  <meta name="google" content="notranslate">
//  <meta name="robots" content="noindex, nofollow">
//  <meta name="format-detection" content="telephone=no">
//  <meta http-equiv="refresh" content="1; url=http://poker-iv.herokuapp.com/" />
//  <meta name="referrer" content="no-referrer">
//  <script type="text/javascript">
//      var submitted = false;
//      function s() {
//          if (!submitted) {
//              window.setTimeout('_submit()', 3000);
//          }
//          submitted = true;
//      }
//      function _submit() {
//          document.location.href = "http://poker-iv.herokuapp.com/";
//      }
//  </script>
//</head>
//<body onload="s();">
//<form id="dest" method="post" action="http://poker-iv.herokuapp.com/"></form>
//<br/>
//
//<script type="text/javascript">(function() {s();})();</script>
//</body>
//</html>
//`
//
//	rw.Write([]byte(page))
//}
