package main

import (
	"gee"
	"log"
)

func main() {
	accesstoken := gee.AccessTokenJson{} //initialize token access program
	go accesstoken.CountAndReget()       //this goroutine is used to count the remain time of token
	DB := gee.DB{}                       //initialize database program
	DB.Init("./user.db")                 //connect to user.db
	r := gee.New(&DB, &accesstoken)      //initialize engine
	r.Use(gee.Logger())                  //global middleware, calculate response time
	wechat := r.Group("/wechat")         //url/wechat URL, used to access information from wechat server
	wechat.Use(gee.VerifySignature())    //wechat middleware
	{
		infoxml := gee.InfoXML{}             //because format of information from wechat is mainly xml so i create a InfoXML struct
		wechat.GET("/", infoxml.GetWechat)   //for security
		wechat.POST("/", infoxml.PostWechat) //handle the msg from wechat server
	}
	{
		msgjson := gee.InfoJSON{}       //becauser wechat api mainly access json format data. so i create a InfoJSON struct
		r.POST("/msg", msgjson.SendMsg) //handler outer api to send msg to user
	}
	log.Fatal(r.Run("0.0.0.0:80")) //starting server
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
