package renderer

import (
	"fmt"
	"os/exec"
	"errors"
	"saint-angels/shaderbox/pkg/models"
	"io"
	"bytes"
	"time"
	"saint-angels/shaderbox/pkg/models/mysql"
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
    fmt.Printf(">shady+ffmpeg\n %s\n", outputBuffer.String())

	return nil
}

// https://medium.com/@j.d.livni/write-a-go-worker-pool-in-15-minutes-c9b42f640923
type Work struct {
	ID	int
	Job *models.Artwork//Replace with data needed for rendering
}

type Worker struct {
	ID int
	WorkerChannel chan chan Work // used to communicate between dispatcher and workers
	Channel chan Work
	End chan bool
}

//TODO: put it inside StartDispatcher
func (w *Worker) Start() {
	go func() {
		for {
			//Provide worker's channel in a queue to notify collector that it's free
			w.WorkerChannel <-w.Channel
			select {
				case artwork := <-w.Channel: // worker has received job
					Render(artwork.ID)
					// work.DoWork(job.Job, w.ID) // do work
					// time.Sleep(5 * time.Second)
					fmt.Printf("worker %d did some work with artwork %d\n", w.ID, artwork.ID)
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

var WorkerChannel = make(chan chan Work)

type Collector struct {
	Work chan Work // receives jobs to send to workers
	End chan bool // when receives bool stops workers
}

func StartDispatcher(workerCount int, artworkModel *mysql.ArtworkModel) Collector {
	var i int
	var workers []Worker
	input := make(chan Work, 1) // channel to recieve work
	end := make(chan bool) // channel to kill workers
	collector := Collector{Work: input, End: end}

	for i < workerCount {
		i++
		fmt.Println("starting worker: ", i)
		worker := Worker{
				ID: i,
				Channel: make(chan Work),
				WorkerChannel: WorkerChannel,
				End: make(chan bool)}
		worker.Start()
		workers = append(workers, worker) // stores worker
	}

	// start collector
	go func() {
		workId := 0
		for {
			select {
			case <-end:
				for _, w := range workers {
					w.Stop()
				}
				return
			case work := <-input:
				worker := <-WorkerChannel // wait for available channel
				worker <-work // dispatch work to worker
			default:
				artwork, err := artworkModel.GetArtForRender()
				if err != nil {
					if errors.Is(err, models.ErrNoRecord) {
						fmt.Println("no art to render")
						time.Sleep(5 * time.Second)
					} else {
						fmt.Println(err)
						time.Sleep(100 * time.Second)
					}
					continue
				}
				// fmt.Println("got the artwork", artwork.ID)
				input <- Work{ID:workId, Job: artwork}
				workId += 1
			}
		}
	}()
	return collector
}
