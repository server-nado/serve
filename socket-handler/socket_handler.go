package handler

import (
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/websocket"

	"github.com/server-nado/go-nsq"
	"github.com/server-nado/serve"
	. "github.com/server-nado/serve/lib"
)

var (
	Error = log.New(os.Stderr, "ERROR ", log.Lshortfile|log.LstdFlags)
	Debug = log.New(os.Stderr, "DEBUG ", log.Lshortfile|log.LstdFlags)
)

func init() {
	serve.AddServerHandle(new(SocketServer))
}

type SocketServer struct {
	conf        Configure
	producer    *nsq.Producer
	consumer    *nsq.Consumer
	socker_type string
}

func (self *SocketServer) connectNsq() {
	var err error
	self.producer, err = nsq.NewProducer(self.conf["nsqd_address"].(string), self.conf["nsqConf"].(*nsq.Config))
	if err != nil {
		panic(err)
	}

	handler := new(NsqHandler)
	handler.config = self.conf
	self.consumer, err = nsq.NewConsumer(self.conf["master_topic"].(string), "default", self.conf["nsqConf"].(*nsq.Config))
	self.consumer.AddConcurrentHandlers(handler, 2048)
	self.consumer.ConnectToNSQLookupd(self.conf["nsq_lookupd"].(string))
	self.consumer.SetLogger(Debug, nsq.LogLevelInfo)
	self.producer.SetLogger(Debug, nsq.LogLevelInfo)
}

func (self *SocketServer) Run(conf Configure) error {

	self.conf = conf
	self.connectNsq()

	self.socker_type = "socket"
	if socketType, ok := self.conf["socket_type"]; ok {
		self.socker_type = socketType.(string)
	}
	switch self.socker_type {
	case "socket":

		self.listenSocker()
	case "websocket":
		self.listenWebSocket()
	}

	return nil
}

func (self *SocketServer) readByte(buf []byte, w *RouteResponseWrite, r *RouteRequest) bool {
	if buf[0] != serve.HEAD_1 {
		Error.Println("error hand")
		return true
	}
	r.UnmarshalData(buf)
	if fun, ok := self.conf["route_callback"]; ok {
		w.Id = r.Id
		if route, err := fun.(func(ResponseWrite, Request) (bool, error))(w, r); err == nil {
			if route {
				self.producer.Publish(self.conf["consumer_topic"].(string), r.Byte())

			}
		} else {

			return false
		}
	} else {
		Debug.Println("MESSAGE", r.Id, r.Typ)
	}

	return true

}

func (self *SocketServer) listenSocker() {
	var err error
	var listen net.Listener
	listen, err = net.Listen("tcp", self.conf["address"].(string))
	if err != nil {
		panic(err)
	}
	Debug.Println("START SERVER:", self.conf["address"])
	for {
		conn, err := listen.Accept()
		if err != nil {
			Error.Println(err)
			return
		}
		go self.sockerRead(conn)
		Debug.Println("New Connect ===========")

	}

}

func (self SocketServer) sockerRead(conn net.Conn) error {
	replay := []byte{}
	var err error
	r := RouteRequest{}
	r.RouteName = self.conf["master_topic"].(string)
	r.conn = conn
	w := RouteResponseWrite{}
	w.conn = conn

	err = serve.ReadResponseByConnect(replay, conn, func(buf []byte) bool {
		return self.readByte(buf, &w, &r)
	})
	if conn != nil {
		r.conn = nil
		w.conn = nil
		conn.Close()

	}
	if fun, ok := self.conf["on_connect_close"]; ok {
		fun.(func(uid uint32))(r.Id)
	}
	if err != nil {
		Error.Println(err.Error())
		return err
	}
	return nil
}

func (self *SocketServer) listenWebSocket() {
	var err error
	//var listen net.Listener
	lookpath := "/ws"
	if path, ok := self.conf["websocket_path"]; ok {
		lookpath = path.(string)
	}
	Debug.Println("START SERVER:", "ws://"+self.conf["address"].(string)+lookpath)

	http.HandleFunc(lookpath, self.websocket)
	err = http.ListenAndServe(self.conf["address"].(string), nil)
	if err != nil {
		panic(err)
	}
}

func (self *SocketServer) websocket(w http.ResponseWriter, r *http.Request) {
	s := websocket.Server{Handler: websocket.Handler(self.websocketRead)}
	s.ServeHTTP(w, r)
}

func (self SocketServer) websocketRead(conn *websocket.Conn) {
	var replay []byte
	var err error
	r := RouteRequest{}
	r.RouteName = self.conf["master_topic"].(string)
	r.conn = conn
	w := RouteResponseWrite{}
	w.conn = conn

	//read := bytes.NewBuffer(nil)

	for {
		err = websocket.Message.Receive(conn, &replay)
		if err != nil {
			Error.Println(err)
			return
		}
		Debug.Println(" message : ", replay)
		//read.Write(replay)
		//err = serve.ReadResponseByConnect(replay, read, func(buf []byte) bool {
		//return self.readByte(buf, &w, &r)
		//})
		//read.Reset()
		if ok := self.readByte(replay, &w, &r); !ok {
			conn.Close()
			return
		}

	}
}

func (self *SocketServer) Stop() {
	self.producer.Stop()
	self.consumer.Stop()
}
