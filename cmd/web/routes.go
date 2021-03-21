package main

import (
	"net/http"
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest)


	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home));
	mux.Get("/generate", http.HandlerFunc(app.generate));

    fileServer := http.FileServer(http.Dir("./renders/"))
	mux.Get("/renders/", http.StripPrefix("/renders", fileServer))

	return standardMiddleware.Then(mux)
}
