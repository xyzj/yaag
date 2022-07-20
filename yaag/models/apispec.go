package models

type APISpec struct {
	Idx      int
	HttpVerb string
	Path     string
	Calls    []*ApiCall
}
