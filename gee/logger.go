package gee

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		//Start timer
		t := time.Now()
		//Process request
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func VerifySignature() HandlerFunc {
	return func(c *Context) {
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		//if one of the necessary parameter is empty, stop the function
		//Starting sha1 crypto
		//get token
		token := Token
		//sort 3 of the parameters
		SHA1_before := []string{token, timestamp, nonce}
		sort.Strings(SHA1_before)
		//[]string to string
		sha1_string := ""
		for _, v := range SHA1_before {
			sha1_string += v
		}
		//get hash.Hash struct
		sha1 := sha1.New()
		io.WriteString(sha1, sha1_string)
		SHA1_after := fmt.Sprintf("%x", sha1.Sum(nil)) //finishing crypto
		//verify that result and signature are same or not, if yes, return true
		if SHA1_after != signature {
			log.Println(c.Req.RemoteAddr, "====>signature verification failed")
			c.Fail(http.StatusInternalServerError, "status internal")
		}
	}
}

func OnlyForV2() HandlerFunc {
	return func(c *Context) {
		//start timer
		t := time.Now()
		//if a server error occurred
		//c.Fail(500, "Internal Server Error")
		c.Next()
		//Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
