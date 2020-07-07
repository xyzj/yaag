package main

import (
	"github.com/xyzj/yaag/yaag"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/xyzj/yaag/middleware"
	"fmt"
	"time"
	"github.com/urfave/negroni"
)

func main() {
	yaag.Init(&yaag.Config{On: true, DocTitle: "Negroni-gorilla", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})

	router := mux.NewRouter()

	router.HandleFunc("/", middleware.HandleFunc(handler))
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":5000")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, time.Now().String())
}
