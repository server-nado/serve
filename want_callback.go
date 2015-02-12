package nado

import (
	"errors"
	"time"
)

var (
	idLists     = map[uint64]waitItem{}
	addWaitItem = make(chan waitItem)
	delWaitItem = make(chan uint64)
	getWaitItem = make(chan getWaitId)
)

func init() {
	go _initWaitId()
}

func _initWaitId() {
	for {
		select {
		case item := <-addWaitItem:
			if _, ok := idLists[item.id]; !ok {
				idLists[item.id] = item
			}
		case id := <-delWaitItem:
			if item, ok := idLists[id]; ok {
				close(item.b)
				close(item.c)
				delete(idLists, id)
				if item.autoClose != nil {
					item.autoClose.Stop()
				}
				item.onStop()
				Debug.Println("del wait id ", item.id)
			}
		case g := <-getWaitItem:
			if item, ok := idLists[g.id]; ok {
				select {
				case g.item <- item:
				case <-time.After(time.Second * 10):
					Error.Println("timeout ")
				}

			} else {
				select {
				case g.err <- errors.New("Id not exist!"):
				case <-time.After(time.Second * 10):
					Error.Println("timeout ")
				}
			}
		}
	}
}
func WaitId(id uint64, callback func()) (b chan []byte, c chan bool, t *time.Timer) {
	b = make(chan []byte, 100)
	c = make(chan bool)
	item := waitItem{b, c, id, nil, callback}
	t = time.AfterFunc(time.Second*40, item.Close)
	item.autoClose = t
	select {
	case addWaitItem <- item:
	case <-time.After(time.Second * 10):
		Error.Println("timeout")
	}
	return
}
func StopId(id uint64) {
	select {
	case delWaitItem <- id:
	case <-time.After(time.Second * 10):
		Error.Println("timeout")
	}
}

func GetWait(id uint64) (b chan []byte, c chan bool, t *time.Timer, e error) {
	item := getWaitId{}
	item.id = id
	item.item = make(chan waitItem)
	item.err = make(chan error)
	select {
	case getWaitItem <- item:
	case <-time.After(time.Second):
		Error.Println("timeout")
	}

	select {
	case i := <-item.item:
		b = i.b
		c = i.c
		t = i.autoClose

	case e = <-item.err:
		return
	case <-time.After(time.Second):
		Error.Println("timeout")
	}
	return
}

type waitItem struct {
	b         chan []byte
	c         chan bool
	id        uint64
	autoClose *time.Timer
	onStop    func()
}

func (self *waitItem) Close() {
	StopId(self.id)
}

type getWaitId struct {
	id   uint64
	item chan waitItem
	err  chan error
}
