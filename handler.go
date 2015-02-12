package nado

import (
	"time"

	"github.com/ablegao/go-nsq"
	. "github.com/ablegao/serve-nado/lib"
)

/***
默认处理方式
**/
var DefaultServer *NadoServer

func init() {
}

func NewConfig() *Configure {
	Config := new(Configure)
	Config.HttpHandleUrl = "/w"
	Config.WebsocketHandlerUrl = "/s"
	Config.MessageTimeout = time.Second * 30
	Config.DataVerify = func(b []byte) error {
		return nil
	}
	Config.NadoDefaultHandle = func(w ResponseWrite, r Request) {
		defer w.Close()
		return
	}

	Config.NsqDefaultHandle = func(w ResponseWrite, r Request) {
		defer w.Close()
		return
	}

	Config.NsqProducterTopic = ""
	Config.NsqConsumerTopic = ""
	Config.NsqChannel = "default"
	Config.NsqConfig = nsq.NewConfig()
	Config.NsqMaxConsumer = 1
	Config.NsqdLookupds = nil
	Config.NsqdAddress = ""
	Config.OnConnectStop = func(w ResponseWrite, r Request) {

	}

	return Config
}
func NewServer(conf *Configure) {
	if DefaultServer == nil {
		DefaultServer = new(NadoServer)
	}

	DefaultServer.config = conf
}

func HandFunc(typ uint16, fun Header) {
	if DefaultServer == nil {
		DefaultServer = new(NadoServer)
	}
	if DefaultServer.headers == nil {
		DefaultServer.headers = make(Headers)
	}
	DefaultServer.HandFunc(typ, fun)

}

func ServerListen() {
	if DefaultServer == nil {
		DefaultServer = new(NadoServer)
	}
	DefaultServer.Run()
}

func AddServerHandle(handle ServeHandle) {
	if DefaultServer == nil {
		DefaultServer = new(NadoServer)
	}

	DefaultServer.AddServeHandle(handle)
}

/*
func SendToUser(uid uint64, r Response) {
	DefaultServer.SendToUid(uid, r)
}*/
