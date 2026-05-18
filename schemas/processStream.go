package schemas

import (
	"fmt"
	"io"
	"time"
)

func (c *Camera) ProcessStream(stdout io.ReadCloser) error {
	buf := make([]byte, 1024*64)
	const maxFrames = 30000000

	for {

		n, err := stdout.Read(buf)
		fmt.Printf("[Câmera %d] Recebendo %d bytes do stream...\n", c.ID, n)

		if err != nil {
			fmt.Printf("[Câmera %d] Erro na leitura do stream: %v\n", c.ID, err)
			return err // Se houver erro na leitura, sai do loop e avisa o "Vigia"
		}

		if n > 0 {
			c.Mu.Lock()
			c.IsRecording = true
			c.LastData = time.Now()
			c.Mu.Unlock()
		}

		frame := make([]byte, n)
		copy(frame, buf[:n])
		c.Buffer.Push(frame)

		if c.isStopping {
			c.Mu.Lock()

			if len(c.RecordingBuffer) >= maxFrames {
				c.RecordingBuffer = c.RecordingBuffer[1:]
			}
			c.RecordingBuffer = append(c.RecordingBuffer, frame)
			c.Mu.Unlock()
		}
	}
}
