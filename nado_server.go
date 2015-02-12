package nado

import (
	"log"
	"os"
	"sync"

	. "github.com/ablegao/serve-nado/lib"
)

var RouterToConsumer Header
var WriteToServer = make(chan RequestResponse)
var Stop = make(chan bool)
var (
	Error = log.New(os.Stderr, "ERROR ", log.Lshortfile|log.LstdFlags)
	Debug = log.New(os.Stderr, "DEBUG ", log.Lshortfile|log.LstdFlags)
)

type UserRoute struct {
	Uid   uint64
	Route chan []byte
}

type NadoServer struct {
	sync.RWMutex
	headers     Headers
	defaultHead Header
	serve       []ServeHandle
	config      *Configure
	AddRoute    chan *UserRoute
	DelRoute    chan uint64
	Routes      map[uint64]chan []byte
}

//启动一个通道服务
func (self *NadoServer) Run() {

	r := RequestResponse{}
	var fun Header
	var ok bool
	for _, r := range self.serve {
		go r.Run(self.config)
	}
	for {
		select {
		case res := <-self.AddRoute:
			self.Lock()
			self.Routes[res.Uid] = res.Route
			self.Unlock()
		case uid := <-self.DelRoute:
			self.Lock()
			delete(self.Routes, uid)
			self.Unlock()
		case r = <-WriteToServer:
			self.RLock()
			if fun, ok = self.headers[r.Req.Type()]; ok {
				go fun(r.Res, r.Req)
			} else {
				go self.config.NadoDefaultHandle(r.Res, r.Req)
			}
			self.RUnlock()

		case <-Stop:
			return
		}
	}
}
func (self *NadoServer) AddServeHandle(s ServeHandle) {
	self.Lock()
	defer self.Unlock()
	self.serve = append(self.serve, s)
}

func (self *NadoServer) HandFunc(typ uint16, fun Header) {
	self.Lock()
	defer self.Unlock()
	self.headers[typ] = fun
}
