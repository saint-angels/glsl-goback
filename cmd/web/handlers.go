package main

import (
	"net/http"
    "saint-angels/glsl-goback/pkg/renderer"
    "fmt"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("go play somewhere else."))
}

func (app *application) generate(w http.ResponseWriter, r *http.Request) {
    id, err := app.artworks.Insert()
    if err != nil {
        app.serverError(w, err)
        return
    }

    err = renderer.Render(id)
    if err != nil {
       app.serverError(w, err)
       return
    }

    reply := fmt.Sprintf("art id:%d", id)
    w.Write([]byte(reply))
}
