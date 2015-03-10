package serve

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/server-nado/go-nsq"
	. "github.com/server-nado/serve/lib"
)

/***
默认处理方式
**/
var DefaultServer *NadoServer

func init() {
}

func NewConfig(jsonFile string) *Configure {
	Config := new(Configure)
	if jsonFile == "" {

		Config.HttpHandleUrl = "/w"
		Config.WebsocketHandlerUrl = "/s"
		Config.Host = ":8080"
		Config.Fastcgi = ""
		Config.NsqProducterTopic = ""
		Config.NsqConsumerTopic = ""
		Config.NsqChannel = "default"

		Config.NsqMaxConsumer = 1
		Config.NsqdLookupds = nil
		Config.NsqdAddress = ""

	} else {
		file, err := os.OpenFile(jsonFile, os.O_RDONLY, 0660)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(b, Config)
		if err != nil {
			panic(err)
		}

	}
	Config.NsqConfig = nsq.NewConfig()
	Config.OnConnectStop = func(w ResponseWrite, r Request) {
	}
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
	Config.OnServeStop = func() {

	}
	Config.OnServeStart = func() {}

	if Config.Databases == nil {
		Config.Databases = [][3]string{}
	}
	if Config.RedisAddress == nil {
		Config.RedisAddress = []string{"127.0.0.1:6379"}
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
func RunHand(w ResponseWrite, r Request, fun Header) {
	if DefaultServer == nil {
		DefaultServer = new(NadoServer)
	}
	DefaultServer.RunHandler(w, r, fun)
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
