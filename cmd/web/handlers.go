package main

import (
	"net/http"
	"os/exec"
    "saint-angels/glsl-goback/pkg/renderer"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("go play somewhere else."))
}

func (app *application) generate(w http.ResponseWriter, r *http.Request) {
    err := renderer.Render()
    if err != nil {
       app.serverError(w, err)
       return
    }
    w.Write([]byte("rendering finished!"))
}

func (app *application) generate2(w http.ResponseWriter, r *http.Request) {
    // path, err := exec.LookPath("shady")
    // if err != nil {
    //    app.serverError(w, errors.New("shady not found"))
    //    return
    // }
    // app.infoLog.Printf("shady path: %s",path)
    //
    shaderName := "seascape"
    shaderPath := "./shaders/" + shaderName + ".glsl"
    encoding := "x11"
    resolution := "1080x1080"
    cmdShady := exec.Command("/home/teeth/go/bin/shady", "-i", shaderPath, "-ofmt", encoding, "-g",resolution, "-framerate", "30")
    stdoutStderr, err := cmdShady.CombinedOutput()

    app.infoLog.Printf("shady output: %s",stdoutStderr)
    if err != nil {
       app.serverError(w, err)
       return
    }

    // videoDuration := "10"
    // outputPath := "./renders/" + shaderName
    // cmdffmpeg := exec.Command(
    //     "/usr/bin/ffmpeg",
    //     "-f", "rawvideo",
    //     "-pixel_format", "rgb24",
    //     "-video_size", "1080x1080",
    //     "-framerate", "30",
    //     "-t", videoDuration,
    //     "-i", "-", outputPath,
    //     "-y")

	// app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

    w.Write([]byte("rendering finished!"))
}
