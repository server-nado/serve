package handler

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync"

	"github.com/ablegao/go-nsq"
)

//用来预处理nsq 发过来的数据
type JsonRequest struct {
	sync.RWMutex
	Typ       uint16 //消息类型
	Id        uint64 //用户id
	RouteName string //来源路由服务器 -回写使用
	b         []byte //字节数据
}

func (self *JsonRequest) Type() uint16 {
	self.RLock()
	defer self.RUnlock()
	return self.Typ
}

func (self *JsonRequest) Unmarshal(info interface{}) error {
	self.Lock()
	defer self.Unlock()
	return json.Unmarshal(self.b, info)
}

func (self *JsonRequest) Marshal(data interface{}) (err error) {
	self.Lock()
	defer self.Unlock()
	self.b, err = json.Marshal(data)
	return err
}
func (self *JsonRequest) GetId() uint64 {
	self.RLock()
	defer self.RUnlock()
	return self.Id
}

func (self *JsonRequest) Byte() []byte {
	self.RLock()
	defer self.RUnlock()
	w := bytes.NewBuffer(nil)
	binary.Write(w, binary.BigEndian, self.Id)
	binary.Write(w, binary.BigEndian, uint16(len(self.RouteName)))
	w.Write([]byte(self.RouteName))
	binary.Write(w, binary.BigEndian, self.Typ)
	binary.Write(w, binary.BigEndian, uint16(len(self.b)))
	w.Write(self.b)
	return w.Bytes()
}
func (self *JsonRequest) SetType(id uint16) {
	self.Lock()
	defer self.Unlock()
	self.Typ = id
}
func (self *JsonRequest) BaseByte() []byte {
	self.RLock()
	defer self.RUnlock()
	return self.b
}
func (self *JsonRequest) Reset() {

}

func (self *JsonRequest) UnmarshalData(data []byte) (err error) {
	self.RLock()
	defer self.RUnlock()
	var n uint16

	buf := bytes.NewBuffer(data)
	r := bytes.NewBuffer(buf.Next(8))
	err = binary.Read(r, binary.BigEndian, &self.Id)
	r.Reset()
	r.Write(buf.Next(2))
	err = binary.Read(r, binary.BigEndian, &n)
	r.Reset()
	self.RouteName = string(buf.Next(int(n)))
	r.Write(buf.Next(2))
	err = binary.Read(r, binary.BigEndian, &self.Typ)
	r.Reset()

	r.Write(buf.Next(2))
	err = binary.Read(r, binary.BigEndian, &n)
	r.Reset()
	self.b = buf.Next(int(n))
	return
}

type JsonResponseWrite struct {
	RouteName string
	producer  *nsq.Producer
	Id        uint64
}

func (self *JsonResponseWrite) Write(body []byte) (err error) {
	self.producer.Publish(self.RouteName, body)
	return
}

/**
重点用在nsq 向route 发送数据时。
*/
func (self *JsonResponseWrite) Close() (err error) {
	w := bytes.NewBuffer(nil)
	binary.Write(w, binary.BigEndian, self.Id)
	binary.Write(w, binary.BigEndian, uint16(len(self.RouteName)))
	//binary.Write(w, binary.BigEndian, self.RouteName)
	w.Write([]byte(self.RouteName))
	binary.Write(w, binary.BigEndian, uint16(1))
	binary.Write(w, binary.BigEndian, uint16(0))
	//w.Write([]byte{0})
	Debug.Println(w.Bytes())
	self.producer.Publish(self.RouteName, w.Bytes())

	return
}
