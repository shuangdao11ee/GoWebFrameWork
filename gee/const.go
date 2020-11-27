package gee

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	//basic information for wechat server
	Token     = "gouqunzhu"
	AppID     = "wx94214222423759cc"
	Appserect = "a76e60a158fc7c306b410652a6e90601"
	//obviously, this URL is used for accessing access token, it response a json format information
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	//this two xml is used to reply user' message passively
	MsgXml = "<xml>\n  <ToUserName>%s</ToUserName>\n  <FromUserName>%s</FromUserName>\n  <CreateTime>%d</CreateTime>\n  <MsgType>%s</MsgType>\n  <Content>%s</Content>\n</xml>"
	ImgXml = "<xml>\n  <ToUserName>%s</ToUserName>\n  <FromUserName>%s</FromUserName>\n  <CreateTime>%d</CreateTime>\n  <MsgType>%s</MsgType>\n  <Image>\n    <MediaId>%s</MediaId>\n  </Image>\n</xml>"
	//send message actively
	MsgUrl = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"
)

//this struct is used to save the access token and remain time
type AccessTokenJson struct {
	Access_Token string
	Expires_In   int
}

//this function is used to get access token and save it to AccessTokenJson.Access_Token
func (accesstoken *AccessTokenJson) GetAccessToken() {
	URL := fmt.Sprintf(AccessTokenURL, AppID, Appserect)
	resp, err := http.Get(URL)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, accesstoken)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", accesstoken)
}

//a son goroutine will do this function, for counting the remain time and again, get token when time is up
func (accesstoken *AccessTokenJson) CountAndReget() {
	for true {
		if accesstoken.Expires_In > 3600 {
			accesstoken.Expires_In -= 1
			time.Sleep(time.Second)
		} else {
			accesstoken.GetAccessToken()
			time.Sleep(10 * time.Second)
		}
	}
}

//this xml liked struct is used to save data from wechat server
type InfoXML struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	Content      string
	MediaId      string
}

//this is a handler, to test whether this program can verify the facticity of the message or not
//localhost/wechat method-GET
func (infoxml *InfoXML) GetWechat(c *Context) {
	if c.CheckSignature() {
		c.String(http.StatusOK, c.Query("echostr"))
	} else {
		c.Fail(http.StatusInternalServerError, "INTEL ERROR!")
	}
}

//this handler is used to handle the messages from wechat server
//the messages are mainly created because user send a message to wechat server
//the format is mainly xml
//localhost/wechat method-POST
func (infoxml *InfoXML) PostWechat(c *Context) {
	body, _ := ioutil.ReadAll(c.Req.Body)
	xml.Unmarshal(body, infoxml)
	switch infoxml.MsgType {
	case "text":
		infoxml.Reply(c)
		c.String(http.StatusOK, MsgXml, infoxml.FromUserName, infoxml.ToUserName, time.Now().Unix(), infoxml.MsgType, infoxml.Content)
	case "image":
		c.String(http.StatusOK, ImgXml, infoxml.FromUserName, infoxml.ToUserName, time.Now().Unix(), infoxml.MsgType, infoxml.MediaId)
	}
}

//my logic judge to users' message
//when user send a 'id' message, program will check whether users' id exist or not,
//if not , it will create a new corresponding id for user
//if yes, it will return id that exist
//change msg content based on previous content
func (infoxml *InfoXML) Reply(c *Context) {
	infoxml.Content = strings.ToLower(infoxml.Content)
	switch infoxml.Content {
	case "id":
		ID := c.Db.GetID(infoxml.FromUserName)
		if ID != "" {
			infoxml.Content = "your id is" + ID
		} else {
			c.Db.IDCreated(infoxml.FromUserName)
			ID = c.Db.GetID(infoxml.FromUserName)
			infoxml.Content = fmt.Sprintf("your new id is %s", ID)
		}
	default:
		for i, v := range []byte(infoxml.Content) {
			if v == 'i' && i < len(infoxml.Content)-1 && []byte(infoxml.Content)[i+1] == 'd' {
				infoxml.Content = "your content has id inside"
			}
		}
	}
}

//this struct is used for outer api access
type InfoJSON struct {
	Id      string
	Msg     string
	Img_str string
}

//parse []byte data into  InfoJSON
func (infojson *InfoJSON) ParseJson(c *Context) {
	body, _ := ioutil.ReadAll(c.Req.Body)
	err := json.Unmarshal(body, infojson)
	if err != nil {
		log.Println(err)
		return
	}
}

//for outer api access, when their post request.body is complete
//it will post a json data to wechat api server to send a message to appointed user
func (infojson *InfoJSON) SendMsg(c *Context) {
	infojson.ParseJson(c)
	OpenId := c.Db.GetOPENID(infojson.Id)
	if OpenId == "" {
		log.Println("infojson's id doesn't exist")
		c.Fail(http.StatusInternalServerError, "Internal Error")
		return
	}
	msgjson := infojson.JsonReplyFormat(OpenId, "text")
	msgjson["text"] = map[string]string{
		"content": infojson.Msg,
	}
	bytedatas, _ := json.Marshal(msgjson)
	reader := bytes.NewReader(bytedatas)
	request, _ := http.NewRequest("POST", fmt.Sprintf(MsgUrl, c.AccessToken.Access_Token), reader)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	client.Do(request)
	c.String(http.StatusOK, "")
}

//return a json like format
func (infojson *InfoJSON) JsonReplyFormat(openid, msgtype string) map[string]interface{} {
	reply := make(map[string]interface{})
	reply["touser"] = openid
	reply["msgtype"] = msgtype
	return reply
}
