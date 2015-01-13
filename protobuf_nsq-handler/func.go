package handler

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

//包外层结构规则。
const (
	Version_Head_1 = 0x0001
	Version_Head_2 = 0x0002
	Version_Head_3 = 0x0003
	HEAD_END       = 0x0009
)

//转码字典。
var codemap = map[uint8][]byte{
	Version_Head_3: []byte("|*"),
	Version_Head_2: []byte("|$"),
	Version_Head_1: []byte("|#"),
	HEAD_END:       []byte("|&"), //Version_Head_3: []byte("|*"),

}

//版本解析对照字典。
var VersionMap = map[uint8]func(req *Response, b []byte) error{
	Version_Head_1: HeadV1Unmarsharl, //第一类解析。
	Version_Head_2: HeadV2Unmarsharl, //第二类解析
	Version_Head_3: HeadV3Unmarsharl, //第三类解析。
}

//创建一个消息异常通告。
func NewResponse(buf []byte) (*Response, error) {
	res := new(Response)
	if err := UnmarshalToResponse(res, buf); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

//将字节编译给Resposne
func UnmarshalToResponse(res *Response, buf []byte) error {
	if len(buf) == 0 {
		return errors.New("Buf is to shout!")
	}

	if fun, ok := VersionMap[buf[0]]; ok {
		buf = BDecode(buf)
		if err := fun(res, buf); err != nil {
			return err
		}
	}
	return nil
}

// 传入一个ReadCloser 会自动解包，
func ReadResponseByConnect(buf []byte, conn io.ReadCloser, call func([]byte) bool) error {
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

	buf = append(buf, data...)
	if data[0] != HEAD_END {
		goto RELOAD_DATA
	}

	goon = call(buf)
	if goon {
		buf = buf[len(buf):]
		goto RELOAD_DATA
	}
	return nil
}

//将字节转换为Respoonse ，输出多个。
func ByteToResponse(b []byte) ([]Response, error) {

	bufarr := bytes.Split(b, []uint8{HEAD_END})
	rl := []Response{}

	for _, buf := range bufarr {
		if len(buf) < 4 {
			continue
		}
		buf = append(buf, HEAD_END) //被裁切掉的END 符号补充上去，
		if r, err := NewResponse(buf); err == nil {
			rl = append(rl, *r)
		}
	}

	return rl, nil

}

//20字节包数据类型解析。
func HeadBigFunc(r *Response, b []byte) error {
	return HeadV2Unmarsharl(r, b)
}
func HeadV2Unmarsharl(r *Response, b []byte) error {
	buf := bytes.NewBuffer(b)
	r.V = buf.Next(1)[0] //version  旧版本的头消息。

	tbuf := bytes.NewBuffer(buf.Next(8)) //修正时间参数。
	binary.Read(tbuf, binary.BigEndian, &r.Time)
	tbuf.Reset()
	/*
		//频道id
		tbuf.Write(buf.Next(2))
		binary.Read(tbuf, binary.BigEndian, &r.Cid)
		tbuf.Reset()

		//房间id
		tbuf.Write(buf.Next(4))
		binary.Read(tbuf, binary.BigEndian, &r.Rid)
		tbuf.Reset()*/

	//token 长度
	tbuf.Write(buf.Next(2))
	binary.Read(tbuf, binary.BigEndian, &r.TokenLen)
	tbuf.Reset()
	//token
	tbuf.Write(buf.Next(int(r.TokenLen)))
	binary.Read(tbuf, binary.BigEndian, &r.Token)
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

//6字节包数据解析。
func Head4BFunc(r *Response, b []byte) error {
	return HeadV1Unmarsharl(r, b)
}
func HeadV1Unmarsharl(r *Response, b []byte) error {
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

//包含token的消息头解析。
func HeadV3Unmarsharl(r *Response, b []byte) error {
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

//任意字节数据被传入后，将被转译处理。并开头追加字符长度。
func BEncode(buf []byte) []byte {
	for a, b := range codemap {
		buf = bytes.Replace(buf, []byte{a}, b, -1)
	}
	return buf
}

func BDecode(buf []byte) []byte {
	for a, b := range codemap {
		buf = bytes.Replace(buf, b, []byte{a}, -1)
	}
	return buf
}
