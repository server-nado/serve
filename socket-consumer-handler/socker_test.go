package handler

import (
	"net"
	"testing"

	"github.com/server-nado/go-nsq"
	"github.com/server-nado/serve"

	_ "net/http/pprof"

	"github.com/server-nado/serve/lib"
)

var start = make(chan bool)
var conn net.Conn
var config lib.Configure

func init() {
	config = serve.NewConfig("./test.json")
	config["nsqConf"] = nsq.NewConfig()
	config["OnServeStart"] = func() {
		//start <- true
	}

}

/*
func Test_socket_connect(t *testing.T) {
	serve.NewServer(config)
	go serve.ServerListen()
	<-start
	var err error
	conn, err = net.Dial("tcp", "127.0.0.1"+config["address"].(string))
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println("connect is create ", conn)
	}

	buf := bytes.NewBuffer(nil)
	var a, b, c, d, e, f uint32 = 1, 2, 3, 4, 5, 6
	buf.WriteByte(serve.HEAD_1)
	serve.WriteToByte(buf, a, b, c, d, e, f)
	buf.WriteByte(serve.HEAD_END)
	n, err := conn.Write(buf.Bytes())

	if err != nil {
		t.Error(err)
	} else {
		t.Log(n)
	}
	<-time.After(time.Second * 1)

}*/

func Test_websocket_connect(t *testing.T) {
	/*
		config["socket_type"] = "websocket"
		serve.NewServer(config)
		go serve.ServerListen()
		<-start
		conn, err := websocket.Dial("ws://127.0.0.1:7776/ws", "", "http://127.0.0.1")
		if err != nil {
			t.Error(err)
		} else {
			fmt.Println("websocket connect is create ", conn)
		}

		buf := bytes.NewBuffer(nil)
		var a, b, c, d, e, f uint32 = 1, 2, 3, 4, 5, 6
		buf.WriteByte(serve.HEAD_1)
		serve.WriteToByte(buf, byteOrder, a, b, c, d, e, f)
		buf.WriteByte(serve.HEAD_END)
		n, err := conn.Write(buf.Bytes())

		if err != nil {
			t.Error(err)
		} else {
			t.Log(n)
		}

		<-time.After(time.Second * 1)*/
}
