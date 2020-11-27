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
	Token          = "gouqunzhu"
	AppID          = "wx94214222423759cc"
	Appserect      = "a76e60a158fc7c306b410652a6e90601"
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	MsgXml         = "<xml>\n  <ToUserName>%s</ToUserName>\n  <FromUserName>%s</FromUserName>\n  <CreateTime>%d</CreateTime>\n  <MsgType>%s</MsgType>\n  <Content>%s</Content>\n</xml>"
	ImgXml         = "<xml>\n  <ToUserName>%s</ToUserName>\n  <FromUserName>%s</FromUserName>\n  <CreateTime>%d</CreateTime>\n  <MsgType>%s</MsgType>\n  <Image>\n    <MediaId>%s</MediaId>\n  </Image>\n</xml>"
	MsgUrl         = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"
)

type AccessTokenJson struct {
	Access_Token string
	Expires_In   int
}

//access token func
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

type InfoXML struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	Content      string
	MediaId      string
}

//localhost/wechat method-GET
func (infoxml *InfoXML) GetWechat(c *Context) {
	if c.CheckSignature() {
		c.String(http.StatusOK, c.Query("echostr"))
	} else {
		c.Fail(http.StatusInternalServerError, "INTEL ERROR!")
	}
}

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

type InfoJSON struct {
	Id      string
	Msg     string
	Img_str string
}

func (infojson *InfoJSON) ParseJson(c *Context) {
	body, _ := ioutil.ReadAll(c.Req.Body)
	err := json.Unmarshal(body, infojson)
	if err != nil {
		log.Println(err)
		return
	}
}

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

func (infojson *InfoJSON) JsonReplyFormat(openid, msgtype string) map[string]interface{} {
	reply := make(map[string]interface{})
	reply["touser"] = openid
	reply["msgtype"] = msgtype
	return reply
}
