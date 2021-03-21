package renderer

import (
	"fmt"
	"os/exec"
)

func Render() (error) {
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
    fmt.Printf("shady output: %s\n",out)

	return err
}
