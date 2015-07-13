package lib

import "encoding/json"

func NewTestRW() (ResponseWrite, Request, chan []byte, chan bool) {
	w := new(TestResponseWrite)
	w.Stop = make(chan bool)
	w.W = make(chan []byte)

	r := new(TestRequest)
	r.route = "testwrite"
	return w, r, w.W, w.Stop
}

type TestRequest struct {
	TypeId uint16
	Id     uint32
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
func (self *TestRequest) SetId(id uint32) {
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
func (self *TestRequest) GetId() uint32 {
	return self.Id
}
func (self *TestRequest) SetType(id uint16) {
	self.TypeId = id
}
func (self TestRequest) Copy() Request {
	return &self
}

func NewTestW() (w ResponseWrite, write chan []byte, stop chan bool) {
	write = make(chan []byte)
	stop = make(chan bool)
	w = &TestResponseWrite{stop, write}

	return

}

type TestResponseWrite struct {
	Stop chan bool
	W    chan []byte
}

func (self *TestResponseWrite) Write(b []byte) (int, error) {
	self.W <- b
	return len(b), nil
}

func (self *TestResponseWrite) Close() error {
	self.Stop <- true
	return nil
}
