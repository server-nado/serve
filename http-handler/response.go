package handler

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

type dataItem struct {
	k   string
	val interface{}
}

func NewSignInData(m interface{}) signInData {

	switch m.(type) {
	case map[string]interface{}:
		ms := make(signInData, 0, len(m.(map[string]interface{})))
		for k, v := range m.(map[string]interface{}) {
			ms = append(ms, dataItem{k, v})
		}
		return ms
	default:
		typ := reflect.TypeOf(m).Elem()
		val := reflect.ValueOf(m).Elem()
		ms := make(signInData, 0, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			ms = append(ms, dataItem{typ.Field(i).Tag.Get("json"), val.Field(i).Interface()})
		}
		return ms
	}
}

// 用来生成sign
type signInData []dataItem

func (s signInData) Len() int {
	return len(s)
}

func (s signInData) Less(i, j int) bool {
	return s[i].k < s[j].k
}
func (s signInData) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s signInData) GetSign(scre string) string {
	sort.Sort(s)
	b := []byte{}
	for _, v := range s {
		if v.k == "sign" {
			continue
		}

		typ := reflect.TypeOf(v.val)
		val := reflect.ValueOf(v.val)
		switch typ.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendUint(b, val.Uint(), 10)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendInt(b, val.Int(), 10)
		case reflect.String:
			b = append(b, []byte(v.k)...)
			b = append(b, []byte(val.String())...)
		case reflect.Bool:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendBool(b, val.Bool())
		case reflect.Float32, reflect.Float64:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendFloat(b, val.Float(), 'f', 0, 64)
		}
	}
	b = append(b, []byte(scre)...)
	md := md5.New()
	md.Write(b)
	return hex.EncodeToString(md.Sum(nil))
}

func (s signInData) VerifySign(scre string) bool {
	sort.Sort(s)
	sign := ""
	b := []byte{}
	for _, v := range s {
		if v.k == "sign" {
			sign = v.val.(string)
			continue
		}
		b = append(b, []byte(v.k)...)
		typ := reflect.TypeOf(v.val)
		val := reflect.ValueOf(v.val)
		switch typ.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b = strconv.AppendUint(b, val.Uint(), 10)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b = strconv.AppendInt(b, val.Int(), 10)
		case reflect.String:
			b = append(b, []byte(val.String())...)
		case reflect.Bool:
			b = strconv.AppendBool(b, val.Bool())
		case reflect.Float32, reflect.Float64:
			b = strconv.AppendFloat(b, val.Float(), 'f', 0, 64)
		}
	}
	//Debug.Println(string(b), len(b), scre)
	b = append(b, []byte(scre)...)

	md := md5.New()
	md.Write(b)
	mysign := hex.EncodeToString(md.Sum(nil))
	//Debug.Println("get", mysign, "send", sign)
	return mysign == sign
}

type HttpResponse struct {
	stop chan bool
	w    chan []byte
}

func (self HttpResponse) Write(b []byte) (err error) {
	self.w <- b
	return
}

func (self HttpResponse) Close() (err error) {
	self.stop <- true

	return
}

func (self HttpResponse) WriteError(msg string) []byte {
	b, _ := json.Marshal(map[string]interface{}{"code": 9000, "msg": msg})
	return b
}

type HttpRequest struct {
	sync.RWMutex
	FromTopic string
	Uid       uint32
	TypeId    uint16 `json:"code"`
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

	md := NewSignInData(info)
	switch info.(type) {
	case map[string]interface{}:
		info.(map[string]interface{})["sign"] = md.GetSign(self.AppSecret)
	default:
		val := reflect.ValueOf(info).Elem()
		typ := reflect.TypeOf(info).Elem()
		if _, ok := typ.FieldByName("Sign"); ok {
			//val.FieldByName("Sign").SetString(md.GetSign(self.AppSecret))
			val.FieldByName("Sign").SetString(md.GetSign(self.AppSecret))
		}

	}
	self.b, err = json.Marshal(info)
	return
}

func (self *HttpRequest) Type() uint16 {
	return self.TypeId
}
func (self *HttpRequest) Reset() {

}
func (self HttpRequest) Byte() []byte {
	return self.b
}
