package http

import (
	gohttp "net/http"
)

func NullHandler() gohttp.Handler {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {
		return
	}

	return gohttp.HandlerFunc(fn)
}
