package server

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed static/*

var content embed.FS

func Serve(port int, jsonBytes []byte) error {
	r := &Router{&mux.Router{}}

	r.MustResponse("GET", "/", func(res http.ResponseWriter, req *http.Request) {
		data, _ := content.ReadFile("static/index.html")
		res.Header().Set("Content-Type", "text/html")
		fmt.Fprint(res, string(data))
	})

	r.MustResponse("GET", "/normalize.css", func(res http.ResponseWriter, req *http.Request) {
		data, _ := content.ReadFile("static/normalize.css")
		res.Header().Set("Content-Type", "text/css")
		fmt.Fprint(res, string(data))
	})

	r.MustResponse("GET", "/javascript.js", func(res http.ResponseWriter, req *http.Request) {
		data, _ := content.ReadFile("static/javascript.js")
		res.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(res, string(data))
	})

	r.MustResponse("GET", "/data.json", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		fmt.Fprint(res, string(jsonBytes))
	})

	return r.Run(fmt.Sprintf(":%d", port))
}

type Router struct {
	*mux.Router
}

func (r *Router) MustResponse(meth, path string, h http.HandlerFunc) {
	r.HandleFunc(path, h).Methods(meth)
}

func (r *Router) Run(address string) error {
	return http.ListenAndServe(address, r)
}
