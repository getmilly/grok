package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/myheartz/grok/logging"
	uuid "github.com/satori/go.uuid"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

//Logging ...
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer recovery()
		defer c.Request.Body.Close()

		requestID := uuid.NewV4()

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		blw.Header().Set("Request-Id", requestID.String())
		c.Writer = blw

		now := time.Now()
		req := request(c)

		c.Next()

		elapsed := time.Since(now)
		fields := make(map[string]interface{})

		fields["request"] = req
		fields["claims"] = c.Keys
		fields["errors"] = c.Errors
		fields["ip"] = c.ClientIP()
		fields["latency"] = elapsed.Seconds()
		fields["request_id"] = requestID.String()
		fields["response"] = blw

		logging.LogWith(fields).Info(
			"Request incoming from %s elapsed %s completed with %d",
			c.ClientIP(),
			elapsed.String(),
			c.Writer.Status(),
		)
	}
}

func request(context *gin.Context) interface{} {
	r := make(map[string]interface{})

	bodyCopy := new(bytes.Buffer)
	io.Copy(bodyCopy, context.Request.Body)
	bodyData := bodyCopy.Bytes()

	var body map[string]interface{}
	json.Unmarshal(bodyData, &body)

	r["body"] = body
	r["headers"] = context.Request.Header
	r["host"] = context.Request.Host
	r["form"] = context.Request.Form
	r["path"] = context.Request.URL.Path
	r["method"] = context.Request.Method
	r["url"] = context.Request.URL.String()
	r["post_form"] = context.Request.PostForm
	r["remote_addr"] = context.Request.RemoteAddr
	r["query_string"] = context.Request.URL.Query()

	context.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyData))

	return r
}

func response(writer *bodyLogWriter) interface{} {
	r := make(map[string]interface{})

	var body map[string]interface{}
	json.Unmarshal(writer.body.Bytes(), &body)

	r["body"] = body
	r["status"] = writer.Status()
	r["headers"] = writer.Header()

	return r
}

func marshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func unmarshal(str string) interface{} {
	v := make(map[string]interface{})

	json.Unmarshal([]byte(str), &v)

	return v
}

func recovery() {
	if err := recover(); err != nil {
		logging.LogWith(err).Error("Error on logging middleware")
	}
}
