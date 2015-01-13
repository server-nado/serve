package lib

import (
	"time"
)

type Header func(w ResponseWrite, r Request)
type Headers map[uint16]Header

type Message interface {
	ByteBody() []byte
}

type Handler interface {
	ServeMessage(message Message) error
	HandFunc(typ uint16, fun Header)
	DefaultHandle(w ResponseWrite, r Request)
	SetHeaders(Headers)
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
	MessageTimeout      time.Duration
	AppKey              string
	AppSecret           string
}

//用来验证数据有效性的借口
type DataVerifyType func([]byte) error

type ServeHandle interface {
	Run(conf *Configure) error
}
