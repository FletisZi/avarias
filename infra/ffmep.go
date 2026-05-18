package infra

import (
	"fmt"
	"io"
	"os/exec"
)

func NewFFmpegCommand(url string) *exec.Cmd {
	return exec.Command(
		"ffmpeg",
		"-rtsp_transport", "tcp",
		"-timeout", "10000000", // 500ms
		"-i", url,
		"-c", "copy",
		"-f", "mpegts",
		"pipe:1",
	)
}

func StartFFmpeg(url string) (*exec.Cmd, io.ReadCloser, error) {
	cmd := NewFFmpegCommand(url)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Erro ao iniciar FFmpeg: %v\n", err)
		return nil, nil, err
	}

	return cmd, stdout, nil
}

// func StartFFmpeg(url string) (*exec.Cmd, io.ReadCloser, error) {
// 	cmd := NewFFmpegCommand(url)

// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// logs do ffmpeg
// 	go func() {
// 		buf := make([]byte, 1024)

// 		for {
// 			n, err := stderr.Read(buf)
// 			msg := string(buf[:n])
// 			fmt.Printf("[FFmpeg] %s", msg)

// 			if err != nil {
// 				if err != io.EOF {
// 					fmt.Printf("Erro lendo stderr do FFmpeg: %v\n", err)
// 				}
// 				break
// 			}
// 		}
// 	}()

// 	if err := cmd.Start(); err != nil {
// 		return nil, nil, err
// 	}

// 	return cmd, stdout, nil
// }
