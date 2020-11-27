package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gee"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//access token func
func GetAccessToken() {
	URL := fmt.Sprintf(gee.AccessTokenURL, gee.AppID, gee.Appserect)
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
	accesstokenjson := gee.AccessTokenJson{}
	err = json.Unmarshal(body, &accesstokenjson)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("%+v\n", accesstokenjson)
}

//localhost/wechat
//Get func
//URL /
func GetWechat(c *gee.Context) {
	if c.CheckSignature() {
		c.String(http.StatusOK, c.Query("echostr"))
	} else {
		c.Fail(http.StatusInternalServerError, "INTEL ERROR!")
	}
}

//Post func
//URL /
func PostWechat(c *gee.Context) {
	body, _ := ioutil.ReadAll(c.Req.Body)
	infoxml := gee.InfoXML{}
	xml.Unmarshal(body, &infoxml)
	switch infoxml.MsgType {
	case "text":
		c.String(http.StatusOK, gee.MsgXml, infoxml.FromUserName, infoxml.ToUserName, time.Now().Unix(), infoxml.MsgType, infoxml.Content)
	case "image":
		c.String(http.StatusOK, gee.ImgXml, infoxml.FromUserName, infoxml.ToUserName, time.Now().Unix(), infoxml.MsgType, infoxml.MediaId)
	}
}
