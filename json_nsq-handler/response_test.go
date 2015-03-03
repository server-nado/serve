package handler

import (
	"testing"

	"github.com/ablegao/serve-nado/lib"
)

func Test_request_copy(t *testing.T) {
	callback := func(w lib.ResponseWrite, r lib.Request) {
		rr := r.(lib.RequestByNsq).Copy()
		rr.SetId(11)

		t.Log(r.GetId(), rr.GetId())
	}

	res := new(JsonResponseWrite)
	req := new(JsonRequest)
	req.Id = 10
	callback(res, req)
}
