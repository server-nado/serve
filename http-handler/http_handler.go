package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	nado "github.com/server-nado/serve"
	. "github.com/server-nado/serve/lib"
)

var (
	Error             = log.New(os.Stderr, "ERROR ", log.Lshortfile|log.LstdFlags)
	Debug             = log.New(os.Stderr, "DEBUG ", log.Lshortfile|log.LstdFlags)
	index      uint64 = 0
	add_index         = make(chan bool)
	back_index        = make(chan uint64)
)

func init() {
	go getUintId()
	nado.AddServerHandle(new(HttpServeHandle))
}
func getUintId() {
	for {
		select {
		case <-add_index:
			index = index + 1
			back_index <- index
		}
	}
}

type HttpServeHandle struct {
	dataVerify DataVerifyType
	conf       *Configure
}

func (self *HttpServeHandle) Run(conf *Configure) (err error) {
	Debug.Println("http address  ", conf.HttpHandleUrl)
	self.dataVerify = conf.DataVerify
	self.conf = conf
	http.HandleFunc(conf.HttpHandleUrl, self.userRequest)
	return
}
func (self *HttpServeHandle) Stop() {

}

func (self *HttpServeHandle) userRequest(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println(err)
		return
	}
	var t = time.Now().UnixNano()

	defer func(r *http.Request, timestart int64) {
		r.Body.Close()
		Debug.Println("user send over 毫秒:", (time.Now().UnixNano()-timestart)/1000000)

	}(r, t)

	defer r.Body.Close()

	err = self.dataVerify(b)
	if err != nil {
		Error.Println(err)
		return
	}
	add_index <- true
	id := <-back_index

	bytechan, closechan, _ := nado.WaitId(id, func() {
		func(Id uint64, callback func(w ResponseWrite, r Request)) {
			w := new(HttpResponse)
			rr := new(HttpRequest)
			rr.Id = Id
			callback(w, rr)
		}(id, self.conf.OnConnectStop)
	})

	rw := HttpResponse{}
	rw.stop = closechan
	rw.w = bytechan
	req := new(HttpRequest)
	req.Id = id
	req.AppKey = self.conf.AppKey
	req.AppSecret = self.conf.AppSecret

	defer func() {
		self.conf.OnConnectStop(rw, req)
	}()
	err = req.UnmarshalData(b)
	if err != nil {
		Error.Println(err, string(b))
		w.Write(rw.WriteError(err.Error()))
		return
	}

	nado.WriteToServer <- RequestResponse{Req: req, Res: rw}

	for {
		select {
		case bytes := <-rw.w:
			info := map[string]interface{}{}
			if err := json.Unmarshal(bytes, &info); err == nil {
				md := NewSignInData(info)
				info["sign"] = md.GetSign(self.conf.AppSecret)
				bytes, _ = json.Marshal(info)
			}
			w.Write(bytes)

		case <-time.After(self.conf.MessageTimeout):
			Error.Println("connect timeout")
			http.Error(w, "connect timeout ", http.StatusGatewayTimeout)
			return
		case <-rw.stop:
			return
		}
	}
}
