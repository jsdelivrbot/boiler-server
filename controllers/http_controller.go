package controllers

import (
	"github.com/AzureTech/goazure"
	//"net/http"
	//"fmt"
	//"golang.org/x/net/html"
	//"time"
	//"log"
)

type HttpController struct {
	goazure.Controller
}

func init() {
	//server := &http.Server{
	//	Addr:	":12060",
	//	Handler:func(w http.ResponseWriter, r *http.Request) {
	//		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//	},
	//	ReadTimeout:	10 * time.Second,
	//	WriteTimeout:	10 * time.Second,
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//log.Fatal(server.ListenAndServe())
}