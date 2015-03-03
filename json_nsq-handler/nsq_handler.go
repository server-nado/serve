package handler

import (
	"errors"
	"log"
	"os"

	"github.com/ablegao/serve-nado"

	"github.com/ablegao/go-nsq"
	. "github.com/ablegao/serve-nado/lib"
)

var (
	Error            = log.New(os.Stderr, "ERROR ", log.Lshortfile|log.LstdFlags)
	Debug            = log.New(os.Stderr, "DEBUG ", log.Lshortfile|log.LstdFlags)
	writeToTop       = make(chan []byte)
	NewReaderHandler *ReaderHandler
)

func init() {
	NewReaderHandler = new(ReaderHandler)
	nado.RouterToConsumer = NewReaderHandler.routerToNsq
	nado.AddServerHandle(NewReaderHandler)
	log.Println("=======  import test ")
}

type NsqReaderHandler struct {
	Producer *nsq.Producer
	conf     *Configure
}

func (self *NsqReaderHandler) HandleMessage(message *nsq.Message) error {
	if len(message.Body) < 1 {
		return errors.New("Error message")
	}
	r := new(JsonRequest)

	if err := r.UnmarshalData(message.Body); err == nil {
		w := new(JsonResponseWrite)
		w.RouteName = r.RouteName
		w.Id = r.Id
		w.Producer = self.Producer
		//self.conf.NsqDefaultHandle(&w, &r)
		//Debug.Println(r.RouteName)
		nado.WriteToServer <- RequestResponse{Req: r, Res: w}

	} else {
		Error.Println(err)
	}
	return nil
}

type ReaderHandler struct {
	conf     *Configure
	Consumer *nsq.Consumer
	Producer *nsq.Producer
}

func (self *ReaderHandler) Stop() {
	self.Consumer.Stop()
	self.Producer.Stop()
}
func (self *ReaderHandler) Run(conf *Configure) error {
	var err error
	self.conf = conf

	handler := new(NsqReaderHandler)
	handler.Producer, err = nsq.NewProducer(conf.NsqdAddress, conf.NsqConfig)
	if err != nil {
		panic(err)
	}
	self.Producer = handler.Producer
	self.Consumer, err = nsq.NewConsumer(conf.NsqConsumerTopic, conf.NsqChannel, conf.NsqConfig)
	if err != nil {
		panic(err)
	}

	self.Consumer.AddConcurrentHandlers(handler, conf.NsqMaxConsumer)
	err = self.Consumer.ConnectToNSQLookupds(conf.NsqdLookupds)
	if err != nil {
		panic(err)
	}

	return nil
}
func (self *ReaderHandler) routerToNsq(w ResponseWrite, r Request) {

}
