package schemas

import (
	"camsystem/infra"
	"camsystem/tools"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type Camera struct {
	ID              int
	URL             string
	Cmd             *exec.Cmd
	Buffer          *tools.RingBuffer
	Recording       bool
	IsRecording     bool
	isStopping      bool
	LastData        time.Time
	RecordingBuffer [][]byte
	Placa           string
	Mu              sync.RWMutex
}

func NewCamera(id int, url string) *Camera {
	return &Camera{
		ID:              id,
		URL:             url,
		Buffer:          tools.NewRingBuffer(12, 20),
		Recording:       false,
		RecordingBuffer: make([][]byte, 0),
	}
}

func (c *Camera) setRecording(status bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.IsRecording = status
}

func (c *Camera) StartCapture() {

	go c.MonitorStream()

	go func() {
		for {

			fmt.Printf("[Câmera %d] Tentando conectar em: %s\n", c.ID, c.URL)

			cmd, stdout, err := infra.StartFFmpeg(c.URL)

			if err != nil {
				fmt.Printf("[Câmera %d] erro ao iniciar ffmpeg: %v\n", c.ID, err)
				time.Sleep(5 * time.Second)
				continue
			}

			err = c.ProcessStream(stdout)

			fmt.Printf("[Câmera %d] stream encerrado: %v\n", c.ID, err)

			c.Mu.Lock()
			c.IsRecording = false
			c.Mu.Unlock()

			cmd.Wait()

			time.Sleep(5 * time.Second)
		}
	}()
}

// func (c *Camera) StartCapture() {
// 	go func() {
// 		for {
// 			fmt.Printf("[Câmera %d] Tentando conectar em: %s\n", c.ID, c.URL)

// 			cmd, stdout, _ := infra.StartFFmpeg(c.URL)

// 			done := make(chan error, 1)

// 			go func() {
// 				done <- c.ProcessStream(stdout)
// 			}()

// 			// c.setRecording(true)
// 			fmt.Printf("[Câmera %d] O que ta em ISRecording: %v\n", c.ID, c.IsRecording)
// 			fmt.Printf("[Câmera %d] ✅ Conectado e Gravando!\n", c.ID)

// 			erro := cmd.Wait()

// 			if erro != nil {
// 				fmt.Printf("[Câmera %d] FFmpeg caiu! Motivo: %v\n", c.ID, erro)
// 			} else {
// 				fmt.Printf("[Câmera %d] FFmpeg encerrou normalmente.\n", c.ID)
// 			}

// 			// 5. Delay de segurança antes de reiniciar (para não fritar a CPU se a URL estiver errada)
// 			fmt.Printf("[Câmera %d] Reiniciando captura em 5 segundos...\n", c.ID)
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }

func (c *Camera) MonitorStream() {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {

		c.Mu.Lock()

		// Se passou mais de 10 segundos sem dados
		if time.Since(c.LastData) > 10*time.Second {
			c.IsRecording = false

			fmt.Printf("[Câmera %d] ⚠️ Sem dados do stream!\n", c.ID)
		}

		c.Mu.Unlock()
	}
}

// func (c *Camera) StartCapture() {
// 	go func() {
// 		for {
// 			fmt.Printf("[Câmera %d] Tentando conectar em: %s\n", c.ID, c.URL)

// 			// 1. Configura o comando
// 			// c.Cmd = exec.Command("ffmpeg",
// 			// 	"-rtsp_transport", "tcp",
// 			// 	"-i", c.URL,
// 			// 	"-c", "copy", "-f", "mpegts", "pipe:1")

// 			// c.Cmd = infra.NewFFmpegCommand(c.URL)

// 			cmd, stdout, _ := infra.StartFFmpeg(c.URL)

// 			// stdout, _ := c.Cmd.StdoutPipe()

// 			// 2. Inicia o processo
// 			// if err := cmd.Start(); err != nil {
// 			// 	c.setRecording(false)
// 			// 	fmt.Printf("[Câmera %d] Erro ao iniciar: %v. Tentando novamente em 5s...\n", c.ID, err)
// 			// 	time.Sleep(5 * time.Second)
// 			// 	continue
// 			// }
// 			c.setRecording(true)
// 			fmt.Printf("[Câmera %d] ✅ Conectado e Gravando!\n", c.ID)

// 			// 3. Lê o stream em uma goroutine separada
// 			// Passamos o stdout para uma função que fará o Push no RingBuffer
// 			done := make(chan error, 1)
// 			go func() {
// 				done <- c.ProcessStream(stdout)
// 			}()

// 			// 4. O "Vigia": Ele fica parado aqui até o FFmpeg morrer ou o stream fechar
// 			erro := cmd.Wait()

// 			if erro != nil {
// 				fmt.Printf("[Câmera %d] FFmpeg caiu! Motivo: %v\n", c.ID, erro)
// 			} else {
// 				fmt.Printf("[Câmera %d] FFmpeg encerrou normalmente.\n", c.ID)
// 			}

// 			// 5. Delay de segurança antes de reiniciar (para não fritar a CPU se a URL estiver errada)
// 			fmt.Printf("[Câmera %d] Reiniciando captura em 5 segundos...\n", c.ID)
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }

// func (c *Camera) processStream(stdout io.ReadCloser) error {
// 	buf := make([]byte, 1024*64)
// 	const maxFrames = 30000000

// 	for {
// 		n, err := stdout.Read(buf)
// 		if err != nil {
// 			return err // Se houver erro na leitura, sai do loop e avisa o "Vigia"
// 		}

// 		frame := make([]byte, n)
// 		copy(frame, buf[:n])
// 		c.Buffer.Push(frame)

// 		if c.isStopping {
// 			// fmt.Printf("[Câmera %d] Gravando frame de %d bytes\n", c.ID, n)
// 			c.Mu.Lock()
// 			// c.RecordingBuffer = append(c.RecordingBuffer, frame)
// 			if len(c.RecordingBuffer) >= maxFrames {
// 				c.RecordingBuffer = c.RecordingBuffer[1:]
// 			}
// 			c.RecordingBuffer = append(c.RecordingBuffer, frame)
// 			c.Mu.Unlock()
// 		} else {
// 			// fmt.Printf("[Câmera %d] Não está gravando frame de %d bytes\n", c.ID, n)
// 		}
// 		// if c.IsRecording {
// 		// 	c.Mu.Lock()
// 		// 	c.RecordingBuffer = append(c.RecordingBuffer, frame)
// 		// 	c.Mu.Unlock()
// 		// }
// 	}
// }

func (c *Camera) SaveBufferToFile(buffer [][]byte, filename string) error {
	cmd := exec.Command("ffmpeg",
		"-f", "mpegts",
		"-i", "pipe:0",
		"-c", "copy",
		"-movflags", "+faststart",
		filename,
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		for _, frame := range buffer {
			stdin.Write(frame)
		}
	}()
	return cmd.Wait()
}

func (c *Camera) SaveRecording(placa string) error {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	fmt.Println("Placa é", placa)
	type Payload struct {
		Placa string `json:"placa"`
	}

	var data Payload

	err := json.Unmarshal([]byte(placa), &data)
	if err != nil {
		fmt.Println("Erro ao converter:", err)
		return err
	}

	filename, err := tools.GenerateVideoFilePath(data.Placa)
	if err != nil {
		fmt.Println("Erro ao gerar caminho do vídeo:", err)
		return err
	}

	fmt.Println("Salvando em:", filename)

	cmd := exec.Command("ffmpeg",
		"-f", "mpegts",
		"-i", "pipe:0",
		"-c", "copy",
		"-movflags", "+faststart",
		filename,
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		for _, frame := range c.RecordingBuffer {
			stdin.Write(frame)
		}
	}()

	return cmd.Wait()
}

// func (c *Camera) RecordingBufferSize() int {
// 	total := 0
// 	for _, b := range c.RecordingBuffer {
// 		total += len(b)
// 	}
// 	return total
// }

// func (c *Camera) BufferSize() int {
// 	total := 0
// 	for _, b := range c.Buffer.GetAll() {

// 		total += len(b)
// 	}
// 	return total
// }

// func (c *Camera) Stop() {
// 	c.Mu.Lock()
// 	c.isStopping = true // Avisa que a parada é intencional
// 	c.Mu.Unlock()

// 	if c.Cmd != nil && c.Cmd.Process != nil {
// 		fmt.Printf("[Câmera %d] Encerrando processo FFmpeg...\n", c.ID)

// 		// No Windows, o Kill() encerra o processo imediatamente
// 		err := c.Cmd.Process.Kill()
// 		if err != nil {
// 			fmt.Printf("[Câmera %d] Erro ao matar processo: %v\n", c.ID, err)
// 		}
// 	}
// }

// RemoveCamera para a captura e remove do mapa
// func (m *StreamManager) RemoveCamera(id int) {
// 	m.Mu.Lock()
// 	defer m.Mu.Unlock()

// 	if cam, ok := m.Cameras[id]; ok {
// 		// Precisaremos criar um método Stop na Camera depois
// 		// cam.Stop()
// 		delete(m.Cameras, id)
// 		fmt.Printf("[Manager] Câmera %d removida.\n", id)
// 	}
// }
