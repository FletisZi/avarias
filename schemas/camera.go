package schemas

import (
	"camsystem/tools"
	"fmt"
	"io"
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
	mu              sync.RWMutex
}

func NewCamera(id int, url string) *Camera {
	return &Camera{
		ID:              id,
		URL:             url,
		Buffer:          tools.NewRingBuffer(10, 1),
		Recording:       false,
		RecordingBuffer: make([][]byte, 0),
	}
}

func (c *Camera) setRecording(status bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.IsRecording = status
}

func (c *Camera) StartCapture() {
	go func() {
		for {
			fmt.Printf("[Câmera %d] Tentando conectar em: %s\n", c.ID, c.URL)

			// 1. Configura o comando
			c.Cmd = exec.Command("ffmpeg",
				"-rtsp_transport", "tcp",
				"-i", c.URL,
				"-c", "copy", "-f", "mpegts", "pipe:1")

			stdout, _ := c.Cmd.StdoutPipe()

			// 2. Inicia o processo
			if err := c.Cmd.Start(); err != nil {
				c.setRecording(false)
				fmt.Printf("[Câmera %d] Erro ao iniciar: %v. Tentando novamente em 5s...\n", c.ID, err)
				time.Sleep(5 * time.Second)
				continue
			}
			c.setRecording(true)
			fmt.Printf("[Câmera %d] ✅ Conectado e Gravando!\n", c.ID)

			// 3. Lê o stream em uma goroutine separada
			// Passamos o stdout para uma função que fará o Push no RingBuffer
			done := make(chan error, 1)
			go func() {
				done <- c.processStream(stdout)
			}()

			// 4. O "Vigia": Ele fica parado aqui até o FFmpeg morrer ou o stream fechar
			err := c.Cmd.Wait()

			if err != nil {
				fmt.Printf("[Câmera %d] FFmpeg caiu! Motivo: %v\n", c.ID, err)
			} else {
				fmt.Printf("[Câmera %d] FFmpeg encerrou normalmente.\n", c.ID)
			}

			// 5. Delay de segurança antes de reiniciar (para não fritar a CPU se a URL estiver errada)
			fmt.Printf("[Câmera %d] Reiniciando captura em 5 segundos...\n", c.ID)
			time.Sleep(5 * time.Second)
		}
	}()
}

func (c *Camera) processStream(stdout io.ReadCloser) error {
	buf := make([]byte, 1024*64)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			return err // Se houver erro na leitura, sai do loop e avisa o "Vigia"
		}

		frame := make([]byte, n)
		copy(frame, buf[:n])
		c.Buffer.Push(frame)

		if c.isStopping {
			// fmt.Printf("[Câmera %d] Gravando frame de %d bytes\n", c.ID, n)
			c.mu.Lock()
			c.RecordingBuffer = append(c.RecordingBuffer, frame)
			c.mu.Unlock()
		} else {
			// fmt.Printf("[Câmera %d] Não está gravando frame de %d bytes\n", c.ID, n)
		}
		// if c.IsRecording {
		// 	c.mu.Lock()
		// 	c.RecordingBuffer = append(c.RecordingBuffer, frame)
		// 	c.mu.Unlock()
		// }
	}
}

type StreamManager struct {
	// O map usa o ID da câmera como chave para busca rápida
	Cameras map[int]*Camera
	mu      sync.RWMutex
}

func NewStreamManager() *StreamManager {
	return &StreamManager{
		Cameras: make(map[int]*Camera),
	}
}

// AddCamera adiciona uma nova câmera ao sistema e já inicia a captura
func (m *StreamManager) AddCamera(id int, url string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Cria a instância usando o construtor que já fizemos
	newCam := NewCamera(id, url)
	m.Cameras[id] = newCam

	// Inicia o processo de auto-healing que criamos antes
	go newCam.StartCapture()

	fmt.Printf("[Manager] Câmera %d adicionada e iniciando...\n", id)
}

func (m *StreamManager) GetCamera(id int) (*Camera, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %d...\n", id)

	cam, exists := m.Cameras[id]

	cam.isStopping = true
	return cam, exists
}

func (m *StreamManager) StartGravação(id int) (*Camera, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %d...\n", id)

	cam, exists := m.Cameras[id]

	cam.isStopping = true
	return cam, exists
}

func (m *StreamManager) StopGravação(id int) (*Camera, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %d...\n", id)

	cam, exists := m.Cameras[id]

	cam.isStopping = false

	return cam, exists
}


func (c *Camera) RecordingBufferSize() int {
    total := 0
    for _, b := range c.RecordingBuffer {
        total += len(b)
    }
    return total
}

// func (c *Camera) Stop() {
// 	c.mu.Lock()
// 	c.isStopping = true // Avisa que a parada é intencional
// 	c.mu.Unlock()

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
// 	m.mu.Lock()
// 	defer m.mu.Unlock()

// 	if cam, ok := m.Cameras[id]; ok {
// 		// Precisaremos criar um método Stop na Camera depois
// 		// cam.Stop()
// 		delete(m.Cameras, id)
// 		fmt.Printf("[Manager] Câmera %d removida.\n", id)
// 	}
// }
