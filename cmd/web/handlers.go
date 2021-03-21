package main

import (
	"net/http"
	"os/exec"
    // "errors"
    "fmt"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("go play somewhere else."))
}

func (app *application) generate(w http.ResponseWriter, r *http.Request) {
    shadyPath := "/home/teeth/go/bin/shady"
    shaderName := "seascape"
    shaderPath := "./shaders/" + shaderName + ".glsl"
    encodingShady := "rgb24"
    resolution := "1080x1080"
    ffmpegPath := "/usr/bin/ffmpeg"
    ffmpegEncoding := "rgb24"
    videoDuration := "10"
    outputPath := "./renders/" + shaderName + ".avi"
    cmdShady := fmt.Sprintf(
        "%s -i %s -ofmt %s -g %s -f 30",
        shadyPath,
        shaderPath,
        encodingShady,
        resolution,
    )
    cmdffmpeg := fmt.Sprintf(
        "%s -f rawvideo -pixel_format %s -video_size %s -framerate 30 -t %s -i - %s -y",
        ffmpegPath,
        ffmpegEncoding,
        resolution,
        videoDuration,
        outputPath,
    )
    cmdCombined := fmt.Sprintf("%s | %s", cmdShady, cmdffmpeg)
    fmt.Println("cmd:" + cmdCombined)
    out, err := exec.Command("bash","-c", cmdCombined).CombinedOutput()
    app.infoLog.Printf("shady output: %s",out)
    if err != nil {
       app.serverError(w, err)
       return
    }
    w.Write([]byte("rendering finished!"))
}
func (app *application) generate2(w http.ResponseWriter, r *http.Request) {
// shady -i fire_storm_cube.glsl -ofmt rgb24  -g 1024x768 -f 30 |
// ffmpeg -f rawvideo -pixel_format rgb24 -video_size 1024x768 -framerate 30 -t 12 -i - example.avi -y
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
