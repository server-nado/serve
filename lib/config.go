package lib

type Configure map[string]interface{}

/* type Configure struct  {
	HttpHandleUrl       string `json:"HttpHandleUrl"`
	WebsocketHandlerUrl string `json:"WebsocketHandlerUrl"`
	Host                string `json:"Host"`
	Fastcgi             string `json:"Fastcgi"`
	DataVerify          DataVerifyType
	NadoDefaultHandle   Header

	MessageTimeout time.Duration `json:"MessageTimeout,timeunit:s"`
	AppKey         string        `json:"AppKey"`
	AppSecret      string        `json:"AppSecret"`

	NsqConsumerTopic  string `json:"NsqConsumerTopic"`
	NsqProducterTopic string `json:"NsqProducterTopic"`
	NsqChannel        string `json:"NsqChannel"`
	NsqDefaultHandle  Header
	NsqMaxConsumer    int      `json:"NsqMaxConsumer"`
	NsqdLookupds      []string `json:"NsqdLookupds"`
	NsqdAddress       string   `json:"NsqdAddress"`
	NsqConfig         *nsq.Config

	RedisDb      int         `json:"RedisDB"`
	RedisAddress []string    `json:"RedisAddress"`
	Databases    [][3]string `json:"Databases"`

	OnServeStop   func()
	OnServeStart  func()
	OnConnectStop func(w ResponseWrite, r Request) //当链接中断时的回调函数
}
*/
