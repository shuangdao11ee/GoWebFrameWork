package gee

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

//H  json like format, used this type to parse into json format
type H map[string]interface{}

//Context extremely important struct
//everytime a request coming, gin framework will
//used this context struct to handle different URL that have been registered
//in other words, everytime a request come in, a new context will be created
//to handle different function that have been registered
type Context struct {
	//origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	//request info
	Path   string
	Method string
	Params map[string]string
	//response info
	StatusCode int
	//middlewares
	handlers []HandlerFunc
	index    int
	//database
	Db *DB
	//accesstoken
	AccessToken *AccessTokenJson
}

//use to get parameter of :name or * in URL
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

//return a pointer context
func newContext(w http.ResponseWriter, req *http.Request, db *DB, accesstoken *AccessTokenJson) *Context {
	return &Context{
		Writer:      w,
		Req:         req,
		Path:        req.URL.Path,
		Method:      req.Method,
		index:       -1,
		Db:          db,
		AccessToken: accesstoken,
	}
}

//stop all the handler actively
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

//Next used to handler functions in c.handlers alternately
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//PostForm parse json like data
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//Query get parameters in url by input corresponding name
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//Status setting status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//SetHeader setting reply header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//String reply string type request body
func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

//JSON reply json body
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//Data reply int only
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//HTML reply html
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

//signature check
//a facticity verify function for all url that include "/wechat"
//a middleware
func (c *Context) CheckSignature() bool {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	//if one of the necessary parameter is empty, stop the function
	if signature == "" || timestamp == "" || nonce == "" {
		return false
	}
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
	return SHA1_after == signature
}
