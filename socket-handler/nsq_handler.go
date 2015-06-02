package handler

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"

	"github.com/server-nado/go-nsq"
	"github.com/server-nado/serve"
	"github.com/server-nado/serve/lib"
)

var byteOrder = binary.BigEndian

//用来预处理nsq 发过来的数据
type NsqRouteRequest struct {
	sync.RWMutex
	Id        uint32 //用户id
	Typ       uint16 //消息类型
	RouteName string //来源路由服务器 -回写使用
	b         []byte //字节数据
	conn      net.Conn
}

func (self *NsqRouteRequest) Type() uint16 {
	self.RLock()
	defer self.RUnlock()
	return self.Typ
}

func (self *NsqRouteRequest) Unmarshal(info interface{}) (err error) {
	self.Lock()
	defer self.Unlock()
	return
}

func (self *NsqRouteRequest) Marshal(data interface{}) (err error) {
	self.Lock()
	defer self.Unlock()

	return
}

func (self *NsqRouteRequest) GetId() uint32 {
	self.RLock()
	defer self.RUnlock()
	return self.Id
}
func (self *NsqRouteRequest) SetType(id uint16) {
	self.Lock()
	defer self.Unlock()
	self.Typ = id
}

func (self *NsqRouteRequest) Byte() []byte {
	self.RLock()
	defer self.RUnlock()
	w := bytes.NewBuffer(nil)
	binary.Write(w, byteOrder, self.Id)
	binary.Write(w, byteOrder, self.Typ)
	binary.Write(w, byteOrder, uint32(len(self.b)))
	w.Write(self.b)

	buf := []byte{serve.HEAD_1}
	buf = append(buf, serve.ByteEncode(w.Bytes())...)
	buf = append(buf, serve.HEAD_END)
	return buf

}

func (self *NsqRouteRequest) BaseByte() []byte {
	self.RLocker()
	defer self.RUnlock()
	return self.b
}
func (self *NsqRouteRequest) Reset() {

}
func (self NsqRouteRequest) Copy() *NsqRouteRequest {

	return &self
}

func (self *NsqRouteRequest) UnmarshalData(data []byte) (err error) {
	self.RLock()
	defer self.RUnlock()
	var n uint32
	if data[0] != serve.HEAD_2 {
		err = ErrHand
		return
	}
	data = data[1 : len(data)-1]
	data = serve.ByteDecode(data)
	buf := bytes.NewBuffer(data)
	defer buf.Reset()
	r := bytes.NewBuffer(buf.Next(4))
	err = binary.Read(r, byteOrder, &self.Id)
	if err != nil {
		return
	}
	r.Reset()

	var n1 uint16

	r.Write(buf.Next(2))

	err = binary.Read(r, byteOrder, &n1)
	r.Reset()

	self.RouteName = string(buf.Next(int(n1)))

	r.Write(buf.Next(2))
	err = binary.Read(r, byteOrder, &self.Typ)
	if err != nil {
		return
	}
	r.Reset()

	r.Write(buf.Next(4))
	err = binary.Read(r, byteOrder, &n)
	if err != nil {
		return
	}
	r.Reset()

	self.b = buf.Next(int(n))

	return
}

type NsqHandler struct {
	config lib.Configure
}

func (self *NsqHandler) HandleMessage(message *nsq.Message) error {
	conn := bytes.NewBuffer(message.Body)

	replay := []byte{}
	r := NsqRouteRequest{}

	err := serve.ReadResponseByConnect(replay, conn, func(replay []byte) bool {
		if fun, ok := self.config["on_consumer_to_client"]; ok {
			r.UnmarshalData(replay)
			fun.(func(lib.Request))(&r)
		}
		return true
	})
	if err != nil {
		Error.Println(err)
	}

	return nil
}
