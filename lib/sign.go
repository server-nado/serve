package lib

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
	"sort"
	"strconv"
)

func NewSignInData(m interface{}) signInData {

	switch m.(type) {
	case map[string]interface{}:
		ms := make(signInData, 0, len(m.(map[string]interface{})))
		for k, v := range m.(map[string]interface{}) {
			ms = append(ms, dataItem{k, v})
		}
		return ms
	default:
		typ := reflect.TypeOf(m).Elem()
		val := reflect.ValueOf(m).Elem()
		ms := make(signInData, 0, val.NumField())
		for i := 0; i < val.NumField(); i++ {
			ms = append(ms, dataItem{typ.Field(i).Tag.Get("json"), val.Field(i).Interface()})
		}
		return ms
	}
}

type dataItem struct {
	k   string
	val interface{}
}

// 用来生成sign
type signInData []dataItem

func (s signInData) Len() int {
	return len(s)
}

func (s signInData) Less(i, j int) bool {
	return s[i].k < s[j].k
}
func (s signInData) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s signInData) GetSign(scre string) string {
	sort.Sort(s)
	b := []byte{}
	for _, v := range s {
		if v.k == "sign" {
			continue
		}

		typ := reflect.TypeOf(v.val)
		val := reflect.ValueOf(v.val)
		switch typ.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendUint(b, val.Uint(), 10)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendInt(b, val.Int(), 10)
		case reflect.String:
			b = append(b, []byte(v.k)...)
			b = append(b, []byte(val.String())...)
		case reflect.Bool:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendBool(b, val.Bool())
		case reflect.Float32, reflect.Float64:
			b = append(b, []byte(v.k)...)
			b = strconv.AppendFloat(b, val.Float(), 'f', 0, 64)
		}
	}
	b = append(b, []byte(scre)...)
	md := md5.New()
	md.Write(b)
	return hex.EncodeToString(md.Sum(nil))
}

func (s signInData) VerifySign(scre string) bool {
	sort.Sort(s)
	sign := ""
	b := []byte{}
	for _, v := range s {
		if v.k == "sign" {
			sign = v.val.(string)
			continue
		}
		b = append(b, []byte(v.k)...)
		typ := reflect.TypeOf(v.val)
		val := reflect.ValueOf(v.val)
		switch typ.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b = strconv.AppendUint(b, val.Uint(), 10)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b = strconv.AppendInt(b, val.Int(), 10)
		case reflect.String:
			b = append(b, []byte(val.String())...)
		case reflect.Bool:
			b = strconv.AppendBool(b, val.Bool())
		case reflect.Float32, reflect.Float64:
			b = strconv.AppendFloat(b, val.Float(), 'f', 0, 64)
		}
	}
	//Debug.Println(string(b), len(b), scre)
	b = append(b, []byte(scre)...)

	md := md5.New()
	md.Write(b)
	mysign := hex.EncodeToString(md.Sum(nil))
	//Debug.Println("get", mysign, "send", sign)
	return mysign == sign
}
