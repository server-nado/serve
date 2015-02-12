package lib

import (
	"time"

	"github.com/ablegao/go-nsq"
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

	//Finish() //表示消息处理完毕。
}

type ResponseWrite interface {
	Write([]byte) error
	Close() error
}

type Configure struct {
	HttpHandleUrl       string
	WebsocketHandlerUrl string
	DataVerify          DataVerifyType
	NadoDefaultHandle   Header

	MessageTimeout time.Duration
	AppKey         string
	AppSecret      string

	NsqConsumerTopic  string
	NsqProducterTopic string
	NsqChannel        string
	NsqDefaultHandle  Header
	NsqMaxConsumer    int
	NsqdLookupds      []string
	NsqdAddress       string
	NsqConfig         *nsq.Config

	OnConnectStop func(w ResponseWrite, r Request) //当链接中断时的回调函数
}

//用来验证数据有效性的借口
type DataVerifyType func([]byte) error

type ServeHandle interface {
	Run(conf *Configure) error
}
