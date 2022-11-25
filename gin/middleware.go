package yaaggin

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohae/deepcopy"
	"github.com/xyzj/gopsu/json"
	"github.com/xyzj/gopsu/tools"
	"github.com/xyzj/yaag/yaag"
)

var (
	// hashWorker = gopsu.GetNewCryptoWorker(gopsu.CryptoMD5)
	skipPath = []string{"/plain/", "/rd/", "/api", "/root/", "/xgame/"}
)

// Document 生成api文档中间件
func Document(skip ...string) gin.HandlerFunc {
	// b := &bytes.Buffer{}
	emptycall := yaag.APICall{
		RequestHeader:    make(map[string]string),
		PostForm:         make(map[string]string),
		RequestURIParams: make(map[string]string),
		ResponseHeader:   make(map[string]string),
		MethodType:       "",
		CurrentPath:      "",
	}
	for _, v := range skip {
		skipPath = append(skipPath, v)
	}
	return func(c *gin.Context) {
		if !yaag.IsOn() {
			return
		}
		for _, v := range skipPath {
			if strings.Contains(c.Request.RequestURI, v) {
				return
			}
		}
		// b.Reset()
		apiCall := deepcopy.Copy(emptycall).(yaag.APICall)
		apiCall.MethodType = c.Request.Method
		apiCall.CurrentPath = strings.Split(c.Request.RequestURI, "?")[0]
		// header
		apiCall.RequestHeader["Content-Type"] = c.Request.Header.Get("Content-Type")
		if v := c.Request.Header.Get("X-Forwarded-For"); v != "" {
			apiCall.RequestHeader["X-Forwarded-For"] = v
		}
		// err := c.Request.Header.WriteSubset(b, middleware.ReqWriteExcludeHeaderDump)
		// if err != nil {
		// 	apiCall.RequestHeader = make(map[string]string)
		// }
		// for _, header := range strings.Split(b.String(), "\n") {
		// 	values := strings.Split(header, ":")
		// 	if strings.EqualFold(values[0], "") {
		// 		continue
		// 	}
		// 	apiCall.RequestHeader[values[0]] = values[1]
		// }
		// apiCall.CallHash = hashWorker.Hash([]byte(apiCall.CurrentPath + apiCall.MethodType + apiCall.RequestBody))
		c.Next()
		apiCall.RequestHeader["request_time"] = time.Now().Format("2006-01-02 15:04:05")
		if c.Param("_userTokenName") != "" {
			apiCall.RequestHeader["token_name"] = c.Param("_userTokenName")
		}
		// request params
		if v, ok := c.Params.Get("_body"); ok {
			apiCall.RequestBody = v
		} else {
			if strings.Contains(c.Request.RequestURI, "?") {
				apiCall.RequestBody = "?" + strings.Split(c.Request.RequestURI, "?")[1]
			}
		}
		// ct := c.Request.Header.Get("Content-Type")
		// switch ct {
		// case "", "application/x-www-form-urlencoded":
		// 	if !strings.Contains(c.Request.RequestURI, "?") {
		// 		apiCall.RequestBody = ""
		// 	} else {
		// 		apiCall.RequestBody = "?" + strings.Split(c.Request.RequestURI, "?")[1]
		// 	}
		// default:
		// 	apiCall.RequestBody = c.Param("_body")
		// }
		if yaag.IsStatusCodeValid(c.Writer.Status()) {
			var body string
			if len(c.Keys) > 0 {
				jsonBytes, err := json.Marshal(c.Keys)
				if err != nil {
					body = ""
				} else {
					body = tools.String(jsonBytes)
				}
			}
			apiCall.ResponseBody = body
			apiCall.ResponseCode = c.Writer.Status()
			apiCall.ResponseHeader["Content-Type"] = c.Writer.Header().Get("Content-Type")
			// headers := map[string]string{}
			// for k, v := range c.Writer.Header() {
			// 	headers[k] = strings.Join(v, " ")
			// }
			// apiCall.ResponseHeader = headers
			yaag.SetGenHTML(apiCall)
			// go yaag.GenerateHTML(apiCall)
		}
	}
}
