package gos

import (
	"net/http"
)

type (
	Response struct {
		Writer http.ResponseWriter
		status int
		size   int
	}
)

func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

func (r *Response) WriteHeader(n int) {
	r.status = n
	r.Writer.WriteHeader(n)

}

func (r *Response) Write(b []byte) (n int, err error) {
	n, err = r.Writer.Write(b)
	r.size += n
	return n, err
}

func (r *Response) Status() int {
	return r.status
}
