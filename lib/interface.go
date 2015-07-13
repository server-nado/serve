package lib

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

	Unmarshal(interface{}) error
	Marshal(interface{}) error
	Byte() []byte
	BaseByte() []byte
	GetId() uint32
	Type() uint16
	SetId(uint32)
	SetType(uint16)
	GetRoute() string
	SetRoute(string)
	Copy() Request
	//Finish() //表示消息处理完毕。
}

type RequestByRoute interface {
	Reset()
	Type() uint16
	Unmarshal(interface{}) error
	Marshal(interface{}) error
	Byte() []byte
	BaseByte() []byte
	GetId() uint32
	GetRoute() string
	SetRoute(string)
	SetId(uint32)
	//Copy() RequestByRoute
	MarshalData() []byte

	//Finish() //表示消息处理完毕。
}

type ResponseWrite interface {
	Write([]byte) (int, error)
	Close() error
}
type ResponseWriteByRoute interface {
	Write([]byte) error
	Close() error
	Copy() ResponseWriteByRoute
}

//用来验证数据有效性的借口
type DataVerifyType func([]byte) error

type ServeHandle interface {
	Run(conf Configure) error
	Stop()
}
