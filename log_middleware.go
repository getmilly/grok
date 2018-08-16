package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

//LogMiddleware ...
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		blw.Header().Set("Request-Id", requestID)
		c.Writer = blw

		now := time.Now()

		c.Next()

		elapsed := time.Since(now)
		fields := make(map[string]interface{})

		fields["ip"] = c.ClientIP()
		fields["request"] = request(c)
		fields["request_id"] = requestID
		fields["response"] = response(blw)
		fields["latency"] = elapsed.Seconds()
		fields["claims"] = c.Keys

		LogWith(fields).Info(
			"Request incoming from %s elapsed %s completed with %d",
			c.ClientIP(),
			elapsed.String(),
			c.Writer.Status(),
		)
	}
}

func request(context *gin.Context) interface{} {
	r := make(map[string]interface{})

	var body interface{}
	requestBody(context.Request, &body)

	r["body"] = body
	r["host"] = context.Request.Host
	r["form"] = context.Request.Form
	r["path"] = context.Request.URL.Path
	r["method"] = context.Request.Method
	r["headers"] = context.Request.Header
	r["url"] = context.Request.URL.String()
	r["post_form"] = context.Request.PostForm
	r["remote_addr"] = context.Request.RemoteAddr
	r["query_string"] = context.Request.URL.Query()

	return r
}

func response(writer *bodyLogWriter) interface{} {
	r := make(map[string]interface{})

	r["status"] = writer.Status()
	r["headers"] = writer.Header()

	var body interface{}
	json.Unmarshal(writer.body.Bytes(), &body)

	r["body"] = body

	return r
}

func requestBody(request *http.Request, v interface{}) {
	body, err := ioutil.ReadAll(request.Body)

	if err == nil {
		return
	}

	var rBody interface{}
	json.Unmarshal(body, &rBody)
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
