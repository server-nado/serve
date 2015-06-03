package handler

import (
	"bytes"
	"encoding/binary"

	"github.com/server-nado/go-nsq"
	"github.com/server-nado/serve"
	"github.com/server-nado/serve/lib"
)

var byteOrder = binary.BigEndian

type NsqHandler struct {
	config   lib.Configure
	producer *nsq.Producer
}

func (self *NsqHandler) HandleMessage(message *nsq.Message) error {
	conn := bytes.NewBuffer(message.Body)

	replay := []byte{}
	Debug.Println("a message", message.Body)
	serve.ReadResponseByConnect(replay, conn, func(replay []byte) bool {
		r := new(RouteRequest)

		w := new(RouteResponseWrite)
		w.producer = self.producer

		r.UnmarshalData(replay)
		w.RouteName = r.RouteName

		serve.RunHand(w, r, self.config["on_concumer_default_callback"].(lib.Header))
		return true
	})

	return nil
}
