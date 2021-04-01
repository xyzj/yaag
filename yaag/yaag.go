/*
 * This is the main core of the yaag package
 */
package yaag

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/xyzj/yaag/yaag/models"
)

var count int
var config *Config

// Initial empty spec
var spec *models.Spec = &models.Spec{}

func IsOn() bool {
	return config.On
}

func Init(conf *Config) {
	config = conf
	// load the config file
	if conf.DocPath == "" {
		conf.DocPath = "apidoc.html"
	}

	filePath, _ := filepath.Abs(conf.DocPath + ".json")
	dataFile, err := os.Open(filePath)
	if err == nil {
		json.NewDecoder(io.Reader(dataFile)).Decode(spec)
		generateHtml()
	}
	defer dataFile.Close()
}

func add(x, y int) int {
	return x + y
}

func mult(x, y int) int {
	return (x + 1) * y
}

func GenerateHtml(apiCall *models.ApiCall) {
	shouldAddPathSpec := true
	for k, apiSpec := range spec.ApiSpecs {
		if apiSpec.Path == apiCall.CurrentPath && apiSpec.HttpVerb == apiCall.MethodType {
			shouldAddPathSpec = false
			apiCall.Id = count
			count++
			deleteCommonHeaders(apiCall)
			// avoid := false
			// for _, currentAPICall := range spec.ApiSpecs[k].Calls {
			// 	if apiCall.RequestBody == currentAPICall.RequestBody &&
			// 		apiCall.ResponseCode == currentAPICall.ResponseCode { // &&
			// 		// apiCall.ResponseBody == currentAPICall.ResponseBody {
			// 		avoid = true
			// 	}
			// }
			// if !avoid {
			// 	spec.ApiSpecs[k].Calls = append(apiSpec.Calls, *apiCall)
			// } else {

			// 	spec.ApiSpecs[k].Calls[0].RequestUrlParams = apiCall.RequestUrlParams
			// 	spec.ApiSpecs[k].Calls[0].PostForm = apiCall.PostForm
			// 	spec.ApiSpecs[k].Calls[0].ResponseBody = apiCall.ResponseBody
			// }
			if len(spec.ApiSpecs[k].Calls) == 0 {
				spec.ApiSpecs[k].Calls = append(apiSpec.Calls, *apiCall)
			} else {
				spec.ApiSpecs[k].Calls[0] = *apiCall
			}
		}
	}

	if shouldAddPathSpec {
		apiSpec := models.ApiSpec{
			HttpVerb: apiCall.MethodType,
			Path:     apiCall.CurrentPath,
		}
		apiCall.Id = count
		count++
		deleteCommonHeaders(apiCall)
		apiSpec.Calls = append(apiSpec.Calls, *apiCall)
		spec.ApiSpecs = append(spec.ApiSpecs, apiSpec)
	}
	filePath, _ := filepath.Abs(config.DocPath)
	dataFile, err := os.Create(filePath + ".json")
	if err != nil {
		log.Println(err)
		return
	}
	defer dataFile.Close()
	data, err := json.Marshal(spec)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = dataFile.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
	generateHtml()
}

func generateHtml() {
	funcs := template.FuncMap{"add": add, "mult": mult}
	t := template.New("API Documentation").Funcs(funcs)
	htmlString := TemplateLocal
	t, err := t.Parse(htmlString)
	if err != nil {
		log.Println(err)
		return
	}
	filePath, err := filepath.Abs(config.DocPath)
	if err != nil {
		panic("Error while creating file path : " + err.Error())
	}
	homeHTMLFile, err := os.Create(filePath)
	if err != nil {
		panic("Error while creating documentation file : " + err.Error())
	}
	defer homeHTMLFile.Close()
	homeWriter := io.Writer(homeHTMLFile)
	t.Execute(homeWriter, map[string]interface{}{"array": spec.ApiSpecs,
		"baseUrls": config.BaseUrls, "Title": config.DocTitle})
}

func deleteCommonHeaders(call *models.ApiCall) {
	delete(call.RequestHeader, "Accept")
	delete(call.RequestHeader, "Accept-Encoding")
	delete(call.RequestHeader, "Accept-Language")
	delete(call.RequestHeader, "Cache-Control")
	delete(call.RequestHeader, "Connection")
	delete(call.RequestHeader, "Cookie")
	delete(call.RequestHeader, "Origin")
	delete(call.RequestHeader, "User-Agent")
	delete(call.RequestHeader, "Vary")
}

// IsStatusCodeValid 检查状态码
func IsStatusCodeValid(code int) bool {
	if code >= 200 && code < 300 {
		return true
	} else {
		return false
	}
}
