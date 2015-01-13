package handler

import (
	"net/http"
	"time"

	. "github.com/ablegao/serve-nado/lib"

	"github.com/ablegao/serve-nado"
	"golang.org/x/net/websocket"
)

type WebSocketServeHandle struct {
	conf          *Configure
	writeCallback chan []byte
}

func (self *WebSocketServeHandle) Run(conf *Configure) (err error) {
	self.conf = conf
	http.HandleFunc(conf.WebsocketHandlerUrl, self.http)
	return
}
func (self *WebSocketServeHandle) http(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(self.wxHandler)}
	s.ServeHTTP(w, r)
}
func (self WebSocketServeHandle) wxHandler(ws *websocket.Conn) {
	var err error
	var in []byte
	r := RequestResponse{}
	for {
		if err = websocket.Message.Receive(ws, &in); err != nil {
			Error.Println("error ......", err)
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
		req.Unmarshal(in)

		res := HttpResponse{}
		res.w = self.writeCallback
		r.Req = req
		r.Res = res
		nado.WriteToServer <- r
		for {
			select {
			case b := <-self.writeCallback:
				websocket.Message.Send(ws, b)

			case <-time.After(self.conf.MessageTimeout):
				Error.Println("Message time out.")
				websocket.Message.Send(ws, []byte("Message time out."))
				return
			case <-res.stop:

				return
			}
		}

	}
}
