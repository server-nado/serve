package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ablegao/serve-nado"
	. "github.com/ablegao/serve-nado/lib"
)

var (
	Error = log.New(os.Stderr, "ERROR ", log.Lshortfile|log.LstdFlags)
	Debug = log.New(os.Stderr, "DEBUG ", log.Lshortfile|log.LstdFlags)
)

func init() {
	nado.AddServerHandle(new(HttpServeHandle))
}

type HttpServeHandle struct {
	dataVerify DataVerifyType
	conf       *Configure
}

func (self *HttpServeHandle) Run(conf *Configure) (err error) {
	self.dataVerify = conf.DataVerify
	self.conf = conf
	http.HandleFunc(conf.HttpHandleUrl, self.userRequest)
	return
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

	/*
		hj, ok := w.(http.Hijacker)
		if !ok {
			errorMsg := "The Web Server does not support Hijacking! "
			http.Error(w, errorMsg, http.StatusInternalServerError)
			return
		}
		conn, bufrw, err1 := hj.Hijack()
		if err1 != nil {
			errorMsg := "Internal error!"
			http.Error(w, errorMsg, http.StatusInternalServerError)
			Error.Printf(errorMsg+" Hijacking Error: %s\\n", err)
			return
		}
		defer conn.Close()*/
	rw := HttpResponse{}
	rw.stop = make(chan bool)
	rw.w = make(chan []byte)
	req := new(HttpRequest)
	req.AppKey = self.conf.AppKey
	req.AppSecret = self.conf.AppSecret
	err = req.UnmarshalData(b)
	if err != nil {
		Error.Println(err, string(b))
		w.Write(rw.WriteError(err.Error()))
		return
	}

	nado.WriteToServer <- RequestResponse{Req: req, Res: rw}

	for {
		select {
		case b := <-rw.w:
			w.Write(b)
		case <-time.After(self.conf.MessageTimeout):
			Error.Println("connect timeout")
			http.Error(w, "connect timeout ", http.StatusGatewayTimeout)
			return
		case <-rw.stop:
			return
		}
	}
}
