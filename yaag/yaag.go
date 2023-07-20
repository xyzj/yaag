// Package yaag 记录接口调用日志
package yaag

/*
 * This is the main core of the yaag package
 */
import (
	"html/template"
	"io"
	"os"

	"github.com/xyzj/gopsu/json"
	"github.com/xyzj/gopsu/loopfunc"
)

var count uint64
var config *Config

// Initial empty spec
var spec *Spec = &Spec{APISpecs: make([]APISpec, 0)}
var htmlTemplate *template.Template
var htmlFile string
var dataFile string
var chanGenHTML = make(chan APICall, 1000)

// IsOn 是否启用
func IsOn() bool {
	return config.On
}

// Init 初始化
func Init(conf *Config) {
	defer func() { recover() }()
	config = conf
	// load the config file
	if conf.DocPath == "" {
		conf.DocPath = "apirecord.html"
	}
	var err error
	htmlTemplate, _ = template.New("apirec").Parse(Template)
	// t, _ := template.New("runtime").Parse(TPLHEAD + TPLCSS + TPLBODY)
	// h := render.HTML{
	// 	Name:     "runtime",
	// 	Data:     statusInfo,
	// 	Template: t,
	// }
	// h.WriteContentType(c.Writer)
	// h.Render(c.Writer)
	// 模板
	// funcs := template.FuncMap{"add": add, "mult": mult}
	// t := template.New("API Documentation").Funcs(funcs)
	// htmlString := TemplateLocal
	// htmlTemplate, err = t.Parse(htmlString)
	// if err != nil {
	// 	println(err.Error())
	// 	return
	// }
	dataFile = conf.DocPath + ".json"
	htmlFile = conf.DocPath
	b, err := os.ReadFile(dataFile)
	if err == nil {
		json.Unmarshal(b, spec)
	}
	go loopfunc.LoopFunc(func(params ...interface{}) {
		for apicall := range chanGenHTML {
			GenerateHTML(apicall)
		}
	}, "genhtml", nil)
	if spec == nil {
		return
	}
	for k, v := range spec.APISpecs {
		if v.Calls == nil {
			continue
		}
		for idx, call := range v.Calls {
			if call.RequestHeader == nil {
				spec.APISpecs[k].Idx = idx
				break
			}
		}
	}
	// f, err := os.Open(dataFile)
	// if err == nil {
	// 	json.NewDecoder(io.Reader(f)).Decode(spec)
	// 	generateHTML()
	// }
	// defer dataFile.Close()
}

// SetGenHTML SetGenHTML
func SetGenHTML(apicall APICall) {
	chanGenHTML <- apicall
}

// GenerateHTML 生成html
func GenerateHTML(apiCall APICall) {
	shouldAddPathSpec := true
	// deleteCommonHeaders(apiCall)
	for k, apiSpec := range spec.APISpecs {
		if apiSpec.Path == apiCall.CurrentPath && apiSpec.MethodType == apiCall.MethodType {
			shouldAddPathSpec = false
			// found := false
			// for _, call := range apiSpec.Calls {
			// 	if call.CallHash == apiCall.CallHash {
			// 		found = true
			// 		break
			// 	}
			// }
			// if found {
			// 	break
			// }
			if spec.APISpecs[k].Idx >= 10 {
				spec.APISpecs[k].Idx = 0
			}
			// apiCall.Id = atomic.AddUint64(&count, 1)
			// avoid := false
			// for _, currentAPICall := range spec.APISpecs[k].Calls {
			// 	if apiCall.RequestBody == currentAPICall.RequestBody &&
			// 		apiCall.ResponseCode == currentAPICall.ResponseCode { // &&
			// 		// apiCall.ResponseBody == currentAPICall.ResponseBody {
			// 		avoid = true
			// 	}
			// }
			// if !avoid {
			// 	spec.APISpecs[k].Calls = append(apiSpec.Calls, *apiCall)
			// } else {

			// 	spec.APISpecs[k].Calls[0].RequestURIParams = apiCall.RequestURIParams
			// 	spec.APISpecs[k].Calls[0].PostForm = apiCall.PostForm
			// 	spec.APISpecs[k].Calls[0].ResponseBody = apiCall.ResponseBody
			// }
			// if len(spec.APISpecs[k].Calls) == 0 {
			// spec.APISpecs[k].Calls = append(apiSpec.Calls, apiCall)
			spec.APISpecs[k].Calls[spec.APISpecs[k].Idx] = apiCall
			spec.APISpecs[k].Idx++
			break
			// } else {
			// 	spec.APISpecs[k].Calls[0] = *apiCall
			// }
		}
	}

	if shouldAddPathSpec {
		apiSpec := APISpec{
			Idx:        0,
			MethodType: apiCall.MethodType,
			Path:       apiCall.CurrentPath,
			Calls:      make([]APICall, 10),
		}
		// apiCall.Id = atomic.AddUint64(&count, 1)
		// apiSpec.Calls = append(apiSpec.Calls, apiCall)
		apiSpec.Calls[0] = apiCall
		spec.APISpecs = append(spec.APISpecs, apiSpec)
	}
	generateHTML()
	if b, err := json.Marshal(spec); err == nil {
		os.WriteFile(dataFile, b, 0664)
	}
}

func generateHTML() {
	homeHTMLFile, err := os.Create(htmlFile)
	if err != nil {
		panic("Error while creating documentation file : " + err.Error())
	}
	defer homeHTMLFile.Close()
	homeWriter := io.Writer(homeHTMLFile)
	htmlTemplate.Execute(homeWriter,
		map[string]interface{}{
			"array":    spec.APISpecs,
			"baseUrls": config.BaseUrls,
			"Title":    config.DocTitle,
		})
}

func deleteCommonHeaders(call *APICall) {
	delete(call.RequestHeader, "Accept")
	delete(call.RequestHeader, "Accept-Encoding")
	delete(call.RequestHeader, "Accept-Language")
	delete(call.RequestHeader, "Cache-Control")
	delete(call.RequestHeader, "Connection")
	delete(call.RequestHeader, "Cookie")
	delete(call.RequestHeader, "Origin")
	delete(call.RequestHeader, "User-Agent")
	delete(call.RequestHeader, "Postman-Token")
	delete(call.RequestHeader, "User-Token")
	delete(call.RequestHeader, "Vary")
	delete(call.ResponseHeader, "Content-Encoding")
	delete(call.ResponseHeader, "Vary")
}

// IsStatusCodeValid 检查状态码
func IsStatusCodeValid(code int) bool {
	if code >= 200 && code < 300 { // || (code >= 400 && code < 500) {
		return true
	}
	return false
}
