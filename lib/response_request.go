package lib

import "encoding/json"

func NewTestRW() (ResponseWrite, Request, chan []byte, chan bool) {
	w := new(TestResponseWrite)
	w.stop = make(chan bool)
	w.w = make(chan []byte)
	r := new(TestRequest)
	r.route = "testwrite"
	return w, r, w.w, w.stop
}

type TestRequest struct {
	TypeId uint16
	Id     uint64
	b      []byte
	route  string
}

func (self *TestRequest) Reset() {

}
func (self *TestRequest) Type() uint16 {
	return self.TypeId
}

func (self *TestRequest) Unmarshal(but interface{}) error {
	return json.Unmarshal(self.b, but)
}

func (self *TestRequest) Marshal(info interface{}) (err error) {
	self.b, err = json.Marshal(info)
	return
}
func (self *TestRequest) SetRoute(s string) {
	self.route = s
}
func (self *TestRequest) SetId(id uint64) {
	self.Id = id
}
func (self *TestRequest) GetRoute() string {
	return self.route
}
func (self *TestRequest) Byte() []byte {
	return self.b
}
func (self *TestRequest) BaseByte() []byte {
	return self.b
}
func (self *TestRequest) GetId() uint64 {
	return self.Id
}
func (self *TestRequest) SetType(id uint16) {
	self.TypeId = id
}

type TestResponseWrite struct {
	stop chan bool
	w    chan []byte
}

func (self *TestResponseWrite) Write(b []byte) error {
	self.w <- b
	return nil
}

func (self *TestResponseWrite) Close() error {
	self.stop <- true
	return nil
}
