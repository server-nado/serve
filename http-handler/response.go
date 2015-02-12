package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	. "github.com/ablegao/serve-nado/lib"
)

type HttpResponse struct {
	stop chan bool
	w    chan []byte
}

func (self HttpResponse) Write(b []byte) (err error) {
	select {
	case self.w <- b:
	default:
		Debug.Println("Write timeout .. ")
	}
	return
}

func (self HttpResponse) Close() (err error) {
	select {
	case self.stop <- true:
	default:
	}
	return
}

func (self HttpResponse) WriteError(msg string) []byte {
	b, _ := json.Marshal(map[string]interface{}{"code": 9000, "msg": msg})
	return b
}

type HttpRequest struct {
	sync.RWMutex
	TypeId    uint16 `json:"code"`
	Id        uint64
	b         []byte
	AppKey    string
	AppSecret string
}

func (self *HttpRequest) UnmarshalData(data []byte) (err error) {
	self.Lock()
	defer self.Unlock()
	info := map[string]interface{}{}
	err = json.Unmarshal(data, &info)
	if err != nil {
		Error.Println(err)
		return
	}
	sd := NewSignInData(info)
	if sd.VerifySign(self.AppSecret) == false {
		err = errors.New("sign error")
		return
	}
	self.b = data

	switch info["code"].(type) {
	case float64:
		self.TypeId = uint16(info["code"].(float64))
	default:
		err = errors.New(fmt.Sprintf("code is %v", info["code"]))
	}
	switch info["_nid"].(type) {
	case float64:
		self.Id = uint64(info["_nid"].(float64))
	default:
		err = errors.New(fmt.Sprintf("_nid is %v", info["_nid"]))
	}

	return
}

func (self *HttpRequest) Unmarshal(data interface{}) (err error) {
	self.RLock()
	defer self.RUnlock()
	err = json.Unmarshal(self.b, data)
	return
}
func (self HttpRequest) encodeSign(b []byte) string {
	md := md5.New()
	md.Write(append(b, []byte(self.AppSecret)...))
	return hex.EncodeToString(md.Sum(nil))
}
func (self *HttpRequest) Marshal(info interface{}) (err error) {
	self.Lock()
	defer self.Unlock()
	self.b, err = json.Marshal(info)
	return
}

func (self *HttpRequest) Type() uint16 {
	return self.TypeId
}
func (self *HttpRequest) Reset() {

}
func (self *HttpRequest) GetId() uint64 {
	self.RLock()
	defer self.RUnlock()
	return self.Id
}
func (self *HttpRequest) Byte() []byte {
	return self.b
}
func (self *HttpRequest) BaseByte() []byte {
	return self.b
}
func (self *HttpRequest) SetType(id uint16) {
	self.Lock()
	defer self.Unlock()
	self.TypeId = id
}
