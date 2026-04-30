package schemas

import (
	"fmt"
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

func (m *StreamManager) StopGravação(id int) (*Camera, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	fmt.Printf("[Manager] Buscando câmera %d...\n", id)

	cam, exists := m.Cameras[id]

	cam.isStopping = false

	return cam, exists
}
