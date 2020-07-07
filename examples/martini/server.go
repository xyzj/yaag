package main

import (
	"github.com/xyzj/yaag/martiniyaag"
	"github.com/xyzj/yaag/yaag"
	"github.com/go-martini/martini"
)

func main() {
	yaag.Init(&yaag.Config{On: true, DocTitle: "Martini", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
	m := martini.Classic()
	m.Use(martiniyaag.Document)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}
