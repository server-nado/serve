package serve

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io"
	"net/url"
	"sort"
)

//包外层结构规则。
const (
	HEAD_1   = 0x01
	HEAD_2   = 0x02
	HEAD_3   = 0x03
	HEAD_END = 0x09
)

//转码字典。
var codemap = map[uint8][]byte{
	HEAD_2:   []byte("|$"),
	HEAD_3:   []byte("|*"),
	HEAD_1:   []byte("|#"),
	HEAD_END: []byte("|&"), //
}

//版本解析对照字典。
/*
var VersionMap = map[uint8]func(req *Request, b []byte) error{
	Version_Head_1: HeadV1Unmarsharl, //第一类解析。
	Version_Head_2: HeadV2Unmarsharl, //第二类解析
	Version_Head_3: HeadV3Unmarsharl, //第三类解析。
}
*/
//拼凑字节
func WriteToByte(buf_y *bytes.Buffer, order binary.ByteOrder, params ...interface{}) (err error) {

	for _, param := range params {
		switch param.(type) {
		case string:
			_, err = buf_y.Write([]byte(param.(string)))
			if err != nil {
				return err
			}
		default:

			err = binary.Write(buf_y, order, param)
			if err != nil {
				return err
			}
		}
	}
	return
}

//转译编码
func ByteEncode(buf []byte) []byte {
	for a, b := range codemap {
		buf = bytes.Replace(buf, []byte{a}, b, -1)
	}
	return buf
}

//翻译编码
func ByteDecode(buf []byte) []byte {
	for a, b := range codemap {
		buf = bytes.Replace(buf, b, []byte{a}, -1)
	}
	return buf
}

/*

func HeadV1Unmarsharl(r Request, b []byte) error {
	buf := bytes.NewBuffer(b)
	r.V = buf.Next(1)[0] //消息版本
	tbuf := bytes.NewBuffer(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.Typ) //类型
	tbuf.Reset()
	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.Mlen) //后续消息长度
	r.Msg = buf.Next(int(r.Mlen))                //消息body。
	return nil
}

func HeadV2Unmarsharl(r Request, b []byte) error {
	buf := bytes.NewBuffer(b)
	r.V = buf.Next(1)[0] //version  旧版本的头消息。

	tbuf := bytes.NewBuffer(buf.Next(8)) //修正时间参数。
	binary.Read(tbuf, binary.BigEndian, &r.Time)
	tbuf.Reset()

	tbuf = bytes.NewBuffer(buf.Next(2)) //消息版本号， 或客户端版本号
	binary.Read(tbuf, binary.BigEndian, &r.Version)
	tbuf.Reset()

	//token 长度
	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.TokenLen)
	tbuf.Reset()
	r.Token = string(buf.Next(int(r.TokenLen)))
	//channel id  长度
	tbuf.Write(buf.Next(2))
	var l uint16
	binary.Read(tbuf, binary.BigEndian, &l)
	tbuf.Reset()
	r.Master = string(buf.Next(int(l)))

	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.Cid)
	tbuf.Reset()
	//用户id
	tbuf.Write(buf.Next(4))
	binary.Read(tbuf, binary.BigEndian, &r.Uid)
	tbuf.Reset()

	//消息类型
	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.Typ)
	tbuf.Reset()
	//消息长度
	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.Mlen)
	//消息。
	r.Msg = buf.Next(int(r.Mlen))
	return nil
}

//包含token的消息头解析。
func HeadV3Unmarsharl(r Request, b []byte) error {
	buf := bytes.NewBuffer(b)
	r.V = buf.Next(1)[0]                 //版本
	tbuf := bytes.NewBuffer(buf.Next(2)) //
	binary.Read(tbuf, binary.BigEndian, &r.TokenLen)
	tbuf.Reset()
	r.Token = string(buf.Next(int(r.TokenLen))) //用户token

	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.Typ) //消息类型
	tbuf.Reset()
	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.Mlen) //长度
	r.Msg = buf.Next(int(r.Mlen))                //消息body
	return nil
}
*/
// 传入一个ReadCloser 会自动解包，
func ReadResponseByConnect(replay []byte, conn io.Reader, call func(replay []byte) bool) error {
	//var buf = []byte{}
	var data = make([]byte, 1) //一个一个字节读取
	var goon bool
	var err error
RELOAD_DATA:
	//一个一个字节的提取， 直到找到HEAD_END
	_, err = conn.Read(data) //没有建立连接对象。 这个地方将会终止socket .
	if err != nil {
		return err
	}

	replay = append(replay, data...)
	if data[0] != HEAD_END {
		goto RELOAD_DATA
	}

	goon = call(replay)
	if goon {
		replay = replay[len(replay):]
		goto RELOAD_DATA
	}
	return nil
}

func GetSign(vals url.Values, source string) string {
	keys := []string{}
	for k, _ := range vals {
		if k != "sign" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	outstr := ""
	for _, k := range keys {
		outstr = outstr + k + vals.Get(k)
	}
	outstr = outstr + source

	md := md5.New()
	md.Write([]byte(outstr))
	return hex.EncodeToString(md.Sum(nil))
}
