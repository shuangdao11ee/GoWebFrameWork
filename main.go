package main

import (
	"log"
	"net/http"

	"gee"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gin</h1>")
	})
	r.GET("/hello", func(c *gee.Context) {
		//except /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	log.Fatal(r.Run(":9999"))
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
