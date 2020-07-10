package gin

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	"github.com/xyzj/yaag/middleware"
	"github.com/xyzj/yaag/yaag"
	"github.com/xyzj/yaag/yaag/models"
)

// Document 生成api文档中间件
func Document() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !yaag.IsOn() || strings.Index(c.Request.RequestURI, "/api") > 0 {
			return
		}
		apiCall := &models.ApiCall{
			RequestHeader:    make(map[string]string),
			PostForm:         make(map[string]string),
			RequestUrlParams: make(map[string]string),
			ResponseHeader:   make(map[string]string),
		}
		// header
		b := bytes.NewBuffer([]byte(""))
		err := c.Request.Header.WriteSubset(b, middleware.ReqWriteExcludeHeaderDump)
		if err != nil {
			apiCall.RequestHeader = make(map[string]string)
		}
		for _, header := range strings.Split(b.String(), "\n") {
			values := strings.Split(header, ":")
			if strings.EqualFold(values[0], "") {
				continue
			}
			apiCall.RequestHeader[values[0]] = values[1]
		}
		// request params
		ct := c.Request.Header.Get("Content-Type")
		switch ct {
		case "", "application/x-www-form-urlencoded":
			apiCall.RequestBody = "?" + strings.Split(c.Request.RequestURI, "?")[1]
		// x, _ := url.ParseQuery(c.Request.URL.RawQuery)
		// for k, v := range x {
		// 	apiCall.RequestUrlParams[k] = v[0]
		// }
		// apiCall.RequestBody = gjson.Parse(c.Param("_body")).String()
		// case "application/json, application/x-www-form-urlencoded":
		// 	gjson.Parse(c.Param("_body")).ForEach(func(key gjson.Result, value gjson.Result) bool {
		// 		apiCall.PostForm[key.String()] = value.String()
		// 		return true
		// 	})
		default:
			apiCall.RequestBody = gjson.Parse(c.Param("_body")).String()
		}
		c.Next()
		if yaag.IsStatusCodeValid(c.Writer.Status()) {
			apiCall.MethodType = c.Request.Method
			apiCall.CurrentPath = strings.Split(c.Request.RequestURI, "?")[0]
			var body string
			if len(c.Keys) > 0 {
				jsonBytes, err := json.Marshal(c.Keys)
				if err != nil {
					body = ""
				} else {
					body = string(jsonBytes)
				}
			}
			apiCall.ResponseBody = body
			apiCall.ResponseCode = c.Writer.Status()
			headers := map[string]string{}
			for k, v := range c.Writer.Header() {
				headers[k] = strings.Join(v, " ")
			}
			apiCall.ResponseHeader = headers
			go yaag.GenerateHtml(apiCall)
		}
	}
}
