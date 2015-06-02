package serve

import (
	"errors"
	"sync"
)

var Users *UserList

func init() {
	Users = NewUsers()
}

func NewUsers() (u *UserList) {
	u = new(UserList)
	u.L = map[interface{}]UserConn{}
	return
}

func UserGet(key interface{}) (UserConn, error) {
	return Users.Get(key)
}

func UserDel(key interface{}) error {
	return Users.Del(key)
}

func UserAdd(key interface{}, info UserConn) error {
	return Users.Add(key, info)
}

func UserSet(key interface{}, info UserConn) error {
	return Users.Set(key, info)
}

type UserConn interface {
	Write([]byte) (int, error)
	Close() error
}

type UserList struct {
	sync.RWMutex

	L map[interface{}]UserConn
}

func (self *UserList) Add(key interface{}, val UserConn) (err error) {
	self.Lock()
	defer self.Unlock()

	if _, ok := self.L[key]; ok {
		err = errors.New("Key exists!")
	} else {
		self.L[key] = val
	}

	return
}

func (self *UserList) Set(key interface{}, val UserConn) (err error) {
	self.Lock()
	defer self.Unlock()

	self.L[key] = val

	return
}

func (self *UserList) Del(key interface{}) (err error) {
	self.Lock()
	defer self.Unlock()

	if _, ok := self.L[key]; ok {
		delete(self.L, key)
	} else {
		err = errors.New("Key not exists!")
	}

	return
}

func (self *UserList) Get(key interface{}) (info UserConn, err error) {
	self.RLock()
	defer self.RUnlock()
	var ok bool
	if info, ok = self.L[key]; ok {
		delete(self.L, key)
	} else {
		err = errors.New("Key not exists!")
	}

	return
}
