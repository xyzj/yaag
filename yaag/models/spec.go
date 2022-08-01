package models

// Spec Spec
type Spec struct {
	APISpecs []APISpec
}

// APISpec APISpec
type APISpec struct {
	Idx      int
	HttpVerb string
	Path     string
	Calls    []ApiCall
}

// ApiCall ApiCall
type ApiCall struct {
	Id uint64

	CurrentPath string
	MethodType  string

	PostForm map[string]string

	RequestHeader        map[string]string
	CommonRequestHeaders map[string]string
	ResponseHeader       map[string]string
	RequestUrlParams     map[string]string

	RequestBody  string
	ResponseBody string
	ResponseCode int
	CallHash     string
}
