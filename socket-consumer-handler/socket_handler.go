package handler

import (
	"log"
	"os"

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
	conf     Configure
	producer *nsq.Producer
	consumer *nsq.Consumer
}

func (self *SocketServer) Run(conf Configure) error {
	self.conf = conf
	var err error
	self.producer, err = nsq.NewProducer(self.conf["nsqd_address"].(string), self.conf["nsqConf"].(*nsq.Config))
	if err != nil {
		return err
	}
	Debug.Println(self.conf["consumer_topic"].(string))
	handler := new(NsqHandler)
	handler.config = self.conf
	self.consumer, err = nsq.NewConsumer(self.conf["consumer_topic"].(string), "default", self.conf["nsqConf"].(*nsq.Config))
	handler.producer = self.producer
	self.consumer.AddConcurrentHandlers(handler, 1)
	self.consumer.ConnectToNSQLookupd(self.conf["nsq_lookupd"].(string))
	self.consumer.SetLogger(Debug, nsq.LogLevelInfo)
	self.producer.SetLogger(Debug, nsq.LogLevelInfo)

	Debug.Println("const =====================")
	return nil
}

func (self *SocketServer) Stop() {
	self.consumer.Stop()
	self.producer.Stop()
}
