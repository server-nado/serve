package handler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"sync"

	"github.com/server-nado/go-nsq"
	"github.com/server-nado/serve"
	. "github.com/server-nado/serve/lib"
	"golang.org/x/net/websocket"
)

var ErrHand = errors.New("HEAD ERROR ")

//用来预处理nsq 发过来的数据
type RouteRequest struct {
	sync.RWMutex
	Id        uint32 //用户id
	Typ       uint16 //消息类型
	RouteName string //来源路由服务器 -回写使用
	b         []byte //字节数据

}

func (self *RouteRequest) Type() uint16 {
	self.RLock()
	defer self.RUnlock()
	return self.Typ
}

func (self *RouteRequest) Unmarshal(info interface{}) (err error) {
	self.Lock()
	defer self.Unlock()
	return
}

func (self *RouteRequest) Marshal(data interface{}) (err error) {
	self.Lock()
	defer self.Unlock()
	self.b = data.([]byte)
	return
}

func (self RouteRequest) GetId() uint32 {
	self.RLock()
	defer self.RUnlock()
	return self.Id
}
func (self *RouteRequest) SetType(id uint16) {
	self.Lock()
	defer self.Unlock()
	self.Typ = id
}

func (self RouteRequest) Byte() []byte {
	self.RLock()
	defer self.RUnlock()
	w := bytes.NewBuffer(nil)
	binary.Write(w, byteOrder, self.Id)
	binary.Write(w, byteOrder, uint16(len(self.RouteName)))
	w.Write([]byte(self.RouteName))
	binary.Write(w, byteOrder, self.Typ)
	binary.Write(w, byteOrder, uint32(len(self.b)))
	w.Write(self.b)
	buf := []byte{serve.HEAD_2}
	buf = append(buf, serve.ByteEncode(w.Bytes())...)
	buf = append(buf, serve.HEAD_END)
	return buf
}

func (self RouteRequest) BaseByte() []byte {
	self.RLock()
	defer self.RUnlock()
	return self.b
}
func (self *RouteRequest) Reset() {

}

func (self RouteRequest) Copy() RequestByRoute {

	return &self
}

func (self *RouteRequest) UnmarshalData(data []byte) (err error) {
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

	if err != nil {
		return
	}
	r.Reset()
	b := buf.Next(int(n1))

	self.RouteName = string(b)

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

func (self *RouteRequest) MarshalData() []byte {
	buf := bytes.NewBuffer(nil)

	//buf.WriteByte(serve.HEAD_1)
	serve.WriteToByte(buf, byteOrder, self.Id, self.Typ, uint32(len(self.b)))
	buf.Write(self.b)
	//buf.WriteByte(serve.HEAD_END)
	b := []byte{serve.HEAD_2}
	b = append(b, serve.ByteEncode(buf.Bytes())...)
	b = append(b, serve.HEAD_END)
	return b

}

func (self *RouteRequest) SetId(id uint32) {
	self.Lock()
	defer self.Unlock()
	self.Id = id
}

func (self *RouteRequest) SetRoute(s string) {
	self.Lock()
	defer self.Unlock()
	self.RouteName = s
}

func (self RouteRequest) GetRoute() string {
	self.Lock()
	defer self.Unlock()
	return self.RouteName
}

type RouteResponseWrite struct {
	sync.RWMutex
	RouteName string
	conn      io.ReadWriteCloser
	Id        uint32
	producer  *nsq.Producer
}

func (self *RouteResponseWrite) Write(body []byte) (n int, err error) {
	self.Lock()
	defer self.Unlock()

	switch self.conn.(type) {
	case *websocket.Conn:
		err = websocket.Message.Send(self.conn.(*websocket.Conn), body)
		if err != nil {
			n = 0
		} else {
			n = len(body)
		}
	default:
		n, err = self.conn.Write(body)
	}

	return
}

/**
重点用在nsq 向route 发送数据时。
*/
func (self *RouteResponseWrite) Close() (err error) {

	return
}
