package schemas

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type StreamManager struct {
	// O map usa o ID da câmera como chave para busca rápida
	Cameras map[int]*Camera
	Mu      sync.RWMutex
}

func NewStreamManager() *StreamManager {
	return &StreamManager{
		Cameras: make(map[int]*Camera),
	}
}

// AddCamera adiciona uma nova câmera ao sistema e já inicia a captura
func (m *StreamManager) AddCamera(id int, url string) {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	// Cria a instância usando o construtor que já fizemos
	newCam := NewCamera(id, url)
	m.Cameras[id] = newCam

	// Inicia o processo de auto-healing que criamos antes
	go newCam.StartCapture()

	fmt.Printf("[Manager] Câmera %d adicionada e iniciando...\n", id)
}

func (m *StreamManager) GetCamera(id int) (*Camera, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %d...\n", id)

	cam, exists := m.Cameras[id]

	cam.isStopping = true
	return cam, exists
}

func (m *StreamManager) StartGravação(id int) (*Camera, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %d...\n", id)

	cam, exists := m.Cameras[id]

	data := cam.Buffer.GetAll()
	cam.RecordingBuffer = append(cam.RecordingBuffer, data...)
	fmt.Printf("[Manager] Câmera %d: Copiados %d frames do buffer para o buffer de gravação.\n", id, len(cam.RecordingBuffer))

	cam.isStopping = true
	return cam, exists
}

func (m *StreamManager) BuscarPlaca(cameraID string, id int) string {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %s...\n", cameraID)
	cam := m.Cameras[id]
	resp, err := http.Post("http://localhost:9000/detect?cancela=cancela1", "", nil)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	cam.Placa = string(body)

	return cam.Placa
}

func (m *StreamManager) StopGravação(id int) (*Camera, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %d...\n", id)

	cam, exists := m.Cameras[id]

	cam.isStopping = false

	return cam, exists
}
