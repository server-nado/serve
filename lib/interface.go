package lib

import (
	"time"

	"github.com/server-nado/go-nsq"
)

type Header func(w ResponseWrite, r Request)
type Headers map[uint16]Header

type Message interface {
	ByteBody() []byte
}

type RequestResponse struct {
	Res ResponseWrite
	Req Request
}

type Request interface {
	Reset()
	Type() uint16
	Unmarshal(interface{}) error
	Marshal(interface{}) error
	Byte() []byte
	BaseByte() []byte
	GetId() uint64
	SetType(uint16)
	//Finish() //表示消息处理完毕。
}

type RequestByNsq interface {
	Reset()
	Type() uint16
	Unmarshal(interface{}) error
	Marshal(interface{}) error
	Byte() []byte
	BaseByte() []byte
	GetId() uint64
	GetRoute() string
	SetRoute(string)
	SetId(uint64)
	Copy() RequestByNsq
	//Finish() //表示消息处理完毕。
}

type ResponseWrite interface {
	Write([]byte) error
	Close() error
}
type ResponseWriteByNsq interface {
	Write([]byte) error
	Close() error
	Copy() ResponseWriteByNsq
}
type Configure struct {
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

//用来验证数据有效性的借口
type DataVerifyType func([]byte) error

type ServeHandle interface {
	Run(conf *Configure) error
	Stop()
}
