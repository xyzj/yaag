package yaag

// Spec Spec
type Spec struct {
	APISpecs []APISpec
}

// APISpec APISpec
type APISpec struct {
	Calls      []APICall
	MethodType string
	Path       string
	Idx        int
}

// APICall APICall
type APICall struct {
	PostForm             map[string]string
	RequestHeader        map[string]string
	CommonRequestHeaders map[string]string
	ResponseHeader       map[string]string
	RequestURIParams     map[string]string
	RequestBody          string
	ResponseBody         string
	CallHash             string
	CurrentPath          string
	MethodType           string
	ID                   uint64
	ResponseCode         int
}
