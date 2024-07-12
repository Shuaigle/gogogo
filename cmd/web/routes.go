package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("/snippet/view/{id}", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	customHandler := &customHandler{}
	mux.Handle("/custom", customHandler)

	return commonHeaders(mux)
}
