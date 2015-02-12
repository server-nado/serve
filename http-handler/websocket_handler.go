package handler

import (
	"encoding/json"
	"net/http"
	"time"

	. "github.com/ablegao/serve-nado/lib"

	"github.com/ablegao/serve-nado"
	"golang.org/x/net/websocket"
)

func init() {
	nado.AddServerHandle(new(WebSocketServeHandle))
}

type WebSocketServeHandle struct {
	conf *Configure
}

func (self *WebSocketServeHandle) Run(conf *Configure) (err error) {
	self.conf = conf
	Debug.Println("websocket address ", conf.WebsocketHandlerUrl)
	http.HandleFunc(conf.WebsocketHandlerUrl, self.http)
	return
}
func (self *WebSocketServeHandle) Stop() {

}
func (self *WebSocketServeHandle) http(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(self.wxHandler)}
	s.ServeHTTP(w, r)
}
func (self WebSocketServeHandle) wxHandler(ws *websocket.Conn) {
	var err error
	var in []byte
	var Id uint64
	var writeCallback chan []byte
	var closeMesage chan bool
	var timeSleep *time.Timer
	r := RequestResponse{}
	stop := make(chan bool)
	connectRun := func() {
		for {
			select {
			case b := <-writeCallback:
				info := map[string]interface{}{}
				if err := json.Unmarshal(b, &info); err == nil {
					md := NewSignInData(info)
					info["sign"] = md.GetSign(self.conf.AppSecret)
					b, _ = json.Marshal(info)
				}
				if err := websocket.Message.Send(ws, b); err != nil {

					Error.Println("write error", err)
					return
				}
			case <-closeMesage:

			case <-stop:
				return
			}
		}

	}
	for {
		if err = websocket.Message.Receive(ws, &in); err != nil {
			Error.Println("read error , ", err.Error())
			select {
			case stop <- true:
			default:
			}
			return
		}

		err = self.conf.DataVerify(in)
		if err != nil {
			Error.Println(err)
			websocket.Message.Send(ws, "Data verify error ")
			return
		}

		req := new(HttpRequest)
		req.AppKey = self.conf.AppKey
		req.AppSecret = self.conf.AppSecret
		err = req.UnmarshalData(in)
		if err != nil {
			Error.Println(err)
			continue
		}
		Debug.Println("========== websocket .... sonc ", string(in))
	ADD_INDEX_LINE:
		if req.GetId() == 0 {
			add_index <- true
			Id = <-back_index
			info := struct {
				Code uint16 `json:"code"`
				Nid  uint64 `json:"_nid"`
			}{2, Id}
			if b, e := json.Marshal(&info); e == nil {
				websocket.Message.Send(ws, b)
			}
			writeCallback, closeMesage, timeSleep = nado.WaitId(Id, func() {
				func(id uint64, callback func(w ResponseWrite, r Request)) {
					w := new(HttpResponse)
					rr := new(HttpRequest)
					rr.Id = id
					callback(w, rr)
				}(Id, self.conf.OnConnectStop)
			})
			go connectRun()
			req.Id = Id
		}

		if writeCallback == nil || timeSleep == nil {
			if writeCallback, closeMesage, timeSleep, err = nado.GetWait(Id); err != nil {
				Error.Println(err)
				req.Id = 0
				goto ADD_INDEX_LINE
			}
			go connectRun()
		}

		switch req.Type() {
		case 2:
		case 3:
			timeSleep.Reset(time.Second * 40)

		default:

			res := HttpResponse{}
			res.w = writeCallback
			res.stop = closeMesage
			r.Req = req
			r.Res = res
			nado.WriteToServer <- r
		}

	}
}
