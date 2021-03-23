package renderer

import (
	"fmt"
	"os/exec"
	"errors"
	"saint-angels/glsl-goback/pkg/models"
	"io"
	"bytes"
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
    // cmdShady := fmt.Sprintf(
    //     "%s -i %s -ofmt %s -g %s -f 30",
    //     shadyPath,
    //     shaderPath,
    //     encodingShady,
    //     resolution,
    // )
    // cmdffmpeg := fmt.Sprintf(
    //     "%s -hide_banner -loglevel error -f rawvideo -pixel_format %s -video_size %s -framerate 30 -t %s -i - %s -y",
    //     ffmpegPath,
    //     ffmpegEncoding,
    //     resolution,
    //     videoDuration,
    //     outputPath,
    // )

	c1 := exec.Command(
		shadyPath,
		"-i", shaderPath,
		"-ofmt", encodingShady,
		"-g", resolution,
		"-f", "30",
	)
	c2 := exec.Command(
		ffmpegPath,
		// "-hide_banner", "-loglevel error", //At the moment using these options break the renderer?!!
		"-f", "rawvideo",
		"-pixel_format", ffmpegEncoding,
		"-video_size", resolution,
		"-framerate", "30",
		"-t", videoDuration,
		"-i", "-", outputPath, "-y",
	)
    fmt.Println("shady cmd:" + c1.String())
    fmt.Println("ffmpeg cmd:" + c2.String())
	readPipe, writePipe := io.Pipe()
	c1.Stdout = writePipe
	c2.Stdin = readPipe
    var outputBuffer bytes.Buffer
    c2.Stdout = &outputBuffer
	err = c1.Start()
	if err != nil {
		return err
	}
    err = c2.Start()
	if err != nil {
		return err
	}
    go func() {
        defer writePipe.Close()
        c1.Wait()
    }()
    c2.Wait()
    fmt.Printf(">shady+ffmpeg\n  %s\n", outputBuffer.String())

	return nil
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

