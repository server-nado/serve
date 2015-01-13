package handler

type nsqMessage struct {
	Body []byte
}

func (self *nsqMessage) ByteBody() []byte {
	return self.Body
}
