package gee

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
)

//H no
type H map[string]interface{}

//Context aaa
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

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

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

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//PostForm ?
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//Query ?
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//Status ?
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//SetHeader ?
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

//String ?
func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

//JSON ?
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

//Data ?
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//HTML ?
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

//signature check
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
