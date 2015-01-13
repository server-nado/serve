package handler

import (
	"sync"

	. "github.com/ablegao/serve-nado/lib"
	"github.com/bitly/go-nsq"
)

type Handler struct {
	mu      *sync.RWMutex
	r       ResponseInterface
	handler map[uint16]func()
}

func (self *Handler) HandleMessage(message *nsq.Message) error {
	nsqMessage{message.Body}

}

func (self *Handler) ServerMessage(message Message) {

}

func (self *Handler) HandFunc(typ uint16, fun func(r Response)) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.handler[typ] = fun
}
func (self *Handler) DefaultHandle(r Response) {

}
func (self *Handler) SetResponse(r Response) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.r = r
}

func (self *Handler) Route(r Response) {

}

func (self *Handler) Splite(msg []byte) {

}
