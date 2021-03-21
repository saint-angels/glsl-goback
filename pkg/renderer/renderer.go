package renderer

import (
	"fmt"
	"os/exec"
	"errors"
)

func Render() (error) {
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
    outputPath := "./renders/" + shaderName + ".avi"
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
    fmt.Printf("shady+ffmpeg output: %s\n",out)

	return err
}
