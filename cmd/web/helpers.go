package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {trace := fmt.Sprintf("%s\b%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// func (app *application) render(w http.ResponseWriter, r *http.Request, pageName string, td *templateData) {
// 	ts, ok := app.templateCache[pageName]
//     if !ok {
//         app.serverError(w, fmt.Errorf("The template %s does not exist", pageName))
//         return
//     }

// 	// Create intermediate result buffer, so we can check for errors before presenting it
// 	buf := new(bytes.Buffer)
//     // Execute the template set, passing in any dynamic data.
//     err := ts.Execute(buf, app.addDefaultData(td, r))
// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}

// 	buf.WriteTo(w)
// }

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}
