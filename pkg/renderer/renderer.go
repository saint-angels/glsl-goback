package renderer

import (
	"fmt"
	"os/exec"
	"errors"
	"saint-angels/glsl-goback/pkg/models"
)

func Render(artId int) (error) {
    pathShady, err := exec.LookPath("shady")
    if err != nil {
       return errors.New("renderer: shady not found")
    }
    pathffmpeg, err := exec.LookPath("ffmpeg")
    if err != nil {
       return errors.New("renderer: ffmpeg not found")
    }

    shadyPath := pathShady
    shaderName := "seascape"
    shaderPath := "./shaders/" + shaderName + ".glsl"
    encodingShady := "rgb24"
    resolution := "1080x1080"
    ffmpegPath := pathffmpeg
    ffmpegEncoding := "rgb24"
    videoDuration := "10"
    outputPath := fmt.Sprintf("./renders/%d_%s.avi", artId, shaderName)
    cmdShady := fmt.Sprintf(
        "%s -i %s -ofmt %s -g %s -f 30",
        shadyPath,
        shaderPath,
        encodingShady,
        resolution,
    )
    cmdffmpeg := fmt.Sprintf(
        "%s -hide_banner -loglevel error -f rawvideo -pixel_format %s -video_size %s -framerate 30 -t %s -i - %s -y",
        ffmpegPath,
        ffmpegEncoding,
        resolution,
        videoDuration,
        outputPath,
    )
    cmdCombined := fmt.Sprintf("%s | %s", cmdShady, cmdffmpeg)
    fmt.Println("cmd:" + cmdCombined)
    out, err := exec.Command("bash","-c", cmdCombined).CombinedOutput()
    fmt.Printf(">shady+ffmpeg\n  %s\n",out)

	return err
}
// https://medium.com/@j.d.livni/write-a-go-worker-pool-in-15-minutes-c9b42f640923
type Work struct {
	ID	int
	Job models.Artwork//Replace with data needed for rendering
}

type Worker struct {
	ID int
	WorkerChannel chan chan Work // used to communicate between dispatcher and workers
	Channel chan Work
	End chan bool
}

func (w *Worker) Start() {
	go func() {
		for {
			//Provide worker's channel in a queue to notify that it's free
			w.WorkerChannel <-w.Channel
			select {
			case job := <-w.Channel: // worker has received job
				// work.DoWork(job.Job, w.ID) // do work
				fmt.Printf("worker %d did some work with artwork %d", w.ID, job.ID)
			case <-w.End:
				return
			}
		}
	}()
}

// end worker
func (w *Worker) Stop() {
	fmt.Printf("worker [%d] is stopping", w.ID)
	w.End <- true
}

//TODO: Rewrite cmd calls to that
    // path, err := exec.LookPath("shady")
    // if err != nil {
    //    app.serverError(w, errors.New("shady not found"))
    //    return
    // }
    // app.infoLog.Printf("shady path: %s",path)
    //
    // shaderName := "seascape"
    // shaderPath := "./shaders/" + shaderName + ".glsl"
    // encoding := "x11"
    // resolution := "1080x1080"
    // cmdShady := exec.Command("/home/teeth/go/bin/shady", "-i", shaderPath, "-ofmt", encoding, "-g",resolution, "-framerate", "30")
    // stdoutStderr, err := cmdShady.CombinedOutput()

    // app.infoLog.Printf("shady output: %s",stdoutStderr)
    // if err != nil {
    //    app.serverError(w, err)
    //    return
    // }

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
