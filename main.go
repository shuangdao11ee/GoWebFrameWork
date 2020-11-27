package main

import (
	"gee"
	"log"
)

func main() {
	accesstoken := gee.AccessTokenJson{}
	go accesstoken.CountAndReget()
	DB := gee.DB{}
	DB.Init("./user.db")
	r := gee.New(&DB, &accesstoken)
	r.Use(gee.Logger()) //global middleware, calculate response time
	wechat := r.Group("/wechat")
	wechat.Use(gee.VerifySignature()) //wechat middleware
	{
		infoxml := gee.InfoXML{}
		wechat.GET("/", infoxml.GetWechat)
		wechat.POST("/", infoxml.PostWechat)
	}
	{
		msgjson := gee.InfoJSON{}
		r.POST("/msg", msgjson.SendMsg)
	}
	log.Fatal(r.Run("0.0.0.0:80"))
}

//Engine is the uni handler for all requests
// type Engine struct{}

// func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
// 	switch req.URL.Path {
// 	case "/":
// 		fmt.Fprintf(w, "URL.PATH = %q\n", req.URL.Path)
// 	case "/hello":
// 		for k, v := range req.Header {
// 			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
// 		}
// 	default:
// 		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
// 	}
// }

// func main() {
// 	engine := &Engine{}
// 	//log.Fatal(http.ListenAndServe(":9999", engine))
// }

// func main() {
// 	http.HandleFunc("/", indexHandler)
// 	http.HandleFunc("/hello", helloHandler)
// 	log.Fatal(http.ListenAndServe(":9999", nil))
// }

// func indexHandler(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
// }

// func helloHandler(w http.ResponseWriter, req *http.Request) {
// 	for k, v := range req.Header {
// 		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
// 	}
// }
