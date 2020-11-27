package gee

type AccessTokenJson struct {
	Access_Token string
	Expires_In   int
}

type InfoXML struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	Content      string
	MediaId      string
}

const (
	Token          = "gouqunzhu"
	AppID          = "wx94214222423759cc"
	Appserect      = "a76e60a158fc7c306b410652a6e90601"
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	MsgXml         = "<xml>\n  <ToUserName>%s</ToUserName>\n  <FromUserName>%s</FromUserName>\n  <CreateTime>%d</CreateTime>\n  <MsgType>%s</MsgType>\n  <Content>%s</Content>\n</xml>"
	ImgXml         = "<xml>\n  <ToUserName>%s</ToUserName>\n  <FromUserName>%s</FromUserName>\n  <CreateTime>%d</CreateTime>\n  <MsgType>%s</MsgType>\n  <Image>\n    <MediaId>%s</MediaId>\n  </Image>\n</xml>"
)
