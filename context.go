package gos

import (
	"bytes"
	"encoding/json"
	"math"
	"net"
	"net/http"
	"strings"
)

const (
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-IP"

	charset       = "; charset=UTF-8"
	contentType   = "Content-Type"
	jsonContext   = "application/json" + charset
	stringContext = "text/plain" + charset
	htmlContext   = "text/html" + charset

	acceptLanguage = "Accept-Language"
	abortIndex     = math.MaxInt8 / 2
)

//C is context for every goroutine
type Context struct {
	// params  httprouter.Params
	Request *http.Request
	Reponse *Response
	handles []HandlerFunc
	index   int8
	gos     *Gos
}

func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.index = -1
	c.Reponse.Writer = w
	c.Request = r
}

func (c *Context) Next() {

	len := int8(len(c.handles))
	c.index++
	if c.index < len {
		c.handles[c.index](c)
	}
}

func (c *Context) Json(statusCode int, obj interface{}) {
	c.Reponse.Header().Set(contentType, jsonContext)
	c.Reponse.WriteHeader(statusCode)
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(obj)
	if err != nil {
		panic(err)
	}
	c.Reponse.Write(buf.Bytes())
}
func (c *Context) String(statusCode int, s string) {
	c.Reponse.Header().Set(contentType, jsonContext)
	c.Reponse.WriteHeader(statusCode)

	_, err := c.Reponse.Write([]byte(s))

	if err != nil {
		panic(err)
	}
}

func (c *Context) Html(statusCode int, name string, data interface{}) {
	c.Reponse.Header().Set(contentType, htmlContext)
	c.gos.Template.ExecuteTemplate(c.Reponse.Writer, name, data)
}

func (c *Context) BindJsonBody(i interface{}) error {
	defer c.Request.Body.Close()
	return json.NewDecoder(c.Request.Body).Decode(i)
}

func (c *Context) RealIP() string {
	if ip := c.Request.Header.Get(HeaderXForwardedFor); ip != "" {
		return strings.Split(ip, ", ")[0]
	}
	if ip := c.Request.Header.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	return ra

	return c.Request.RemoteAddr
}
