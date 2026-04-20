package tools

import (
	"sync"
)


// 1. Definimos o RingBuffer (O nosso "Relógio de Frames")
type RingBuffer struct {
	Data       [][]byte
	size       int
	WritePos   int
	isFull     bool
	mu         sync.RWMutex // RWMutex é melhor para leitura/escrita simultânea
}

// 2. Criamos o "Fabricador" de Buffers
func NewRingBuffer(seconds int, framesPerSecond int) *RingBuffer {
	totalFrames := seconds * framesPerSecond
	return &RingBuffer{
		Data: make([][]byte, totalFrames),
		size: totalFrames,
	}
}

// 3. O Método para adicionar frames (A lógica que não move ninguém de lugar)
func (rb *RingBuffer) Push(frame []byte) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.Data[rb.WritePos] = frame
	rb.WritePos = (rb.WritePos + 1) % rb.size
	
	if rb.WritePos == 0 {
		rb.isFull = true
	}
}


func (rb *RingBuffer) GetAll() [][]byte {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	orderedData := make([][]byte, 0, rb.size)

	if !rb.isFull {
		// Se não encheu, é só ler do zero até onde escreveu
		return rb.Data[:rb.WritePos]
	}

	// Se encheu, a ordem correta é:
	// 1. Do WritePos até o Final (os mais antigos)
	// 2. Do Início até o WritePos (os mais novos)
	orderedData = append(orderedData, rb.Data[rb.WritePos:]...)
	orderedData = append(orderedData, rb.Data[:rb.WritePos]...)

	return orderedData
}