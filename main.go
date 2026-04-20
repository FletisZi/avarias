// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"time"
// )

// func main() {
// 	http.HandleFunc("/record", RecordHandler)

// 	fmt.Println("Server running on :8080")
// 	http.ListenAndServe(":8080", nil)
// }


// func RecordHandler(w http.ResponseWriter, r *http.Request) {
// 	rtspURL := "rtsp://admin:Ferrasa321@10.100.150.32:554/Streaming/Channels/101"

// 	filename := fmt.Sprintf("video_%d.mp4", time.Now().Unix())
// 	outputPath := "./videos/" + filename
	
// 	os.MkdirAll("./videos", os.ModePerm)
// 	// FFmpeg command
// 	cmd := exec.Command(
// 		`C:\ffmpeg\bin\ffmpeg.exe`,
// 		"-rtsp_transport", "tcp",
// 		"-i", rtspURL,
// 		"-t", "30",              // duração de 30 segundos
// 		"-c:v", "copy",         // não re-encode (mais rápido)
// 		"-an",                  // sem áudio (opcional)
// 		outputPath,
// 	)

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(string(output))
// 		http.Error(w, string(output), 500)
// 		return
// 	}
	

// 	w.Write([]byte("Vídeo gravado: " + filename))
// }


// package main

// import (
// 	"log"
// 	"os"
// 	"os/exec"
// )

// package main

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"sort"
// 	"sync"
// 	"time"
// )

// var recording bool
// var mu sync.Mutex
// var stopChan chan bool
// var currentRecordingDir string

// func main() {
// 	// cria pastas
// 	os.MkdirAll("buffer", os.ModePerm)
// 	os.MkdirAll("recordings", os.ModePerm)

// 	// inicia FFmpeg em background
// 	go startFFmpeg()

// 	// rotas HTTP
// 	http.HandleFunc("/start", startHandler)
// 	http.HandleFunc("/stop", stopHandler)

// 	log.Println("Servidor rodando em :8080")
// 	http.ListenAndServe(":8080", nil)
// }

// func startFFmpeg() {
// 	cmd := exec.Command(
// 		`C:\ffmpeg\bin\ffmpeg.exe`,
// 		"-rtsp_transport", "tcp",
// 		"-i", "rtsp://admin:Ferrasa321@10.100.150.32:554/Streaming/Channels/101",
// 		"-f", "segment",
// 		"-segment_time", "10",
// 		"-segment_wrap", "1",
// 		"-reset_timestamps", "1",
// 		"buffer/out%03d.ts",
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	log.Println("Iniciando FFmpeg (buffer)...")

// 	err := cmd.Start()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func startHandler(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	if recording {
// 		fmt.Fprintln(w, "Já está gravando")
// 		return
// 	}

// 	recording = true
// 	stopChan = make(chan bool)

// 	go record()

// 	fmt.Fprintln(w, "Gravação iniciada")
// }

// func stopHandler(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	if !recording {
// 		fmt.Fprintln(w, "Não está gravando")
// 		return
// 	}

// 	stopChan <- true
// 	recording = false

// 	// gerar mp4 final
// 	go generateMP4(currentRecordingDir)

// 	fmt.Fprintln(w, "Gravação parada")
// }

// func record() {
// 	timestamp := time.Now().Format("20060102_150405")
// 	outputDir := filepath.Join("recordings", timestamp)
// 	currentRecordingDir = outputDir

// 	os.MkdirAll(outputDir, os.ModePerm)

// 	// 🔥 copia buffer inicial
// 	files, _ := filepath.Glob("buffer/*.ts")
// 	sort.Strings(files)

// 	for _, f := range files {
// 		copyFile(f, filepath.Join(outputDir, filepath.Base(f)))
// 	}

// 	log.Println("Buffer inicial copiado")

// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-stopChan:
// 			log.Println("Parando gravação...")
// 			return

// 		case <-ticker.C:
// 			files, _ := filepath.Glob("buffer/*.ts")
// 			sort.Strings(files)

// 			for _, f := range files {
// 				dst := filepath.Join(outputDir, filepath.Base(f))

// 				if _, err := os.Stat(dst); os.IsNotExist(err) {
// 					copyFile(f, dst)
// 				}
// 			}
// 		}
// 	}
// }

// func copyFile(src, dst string) {
// 	in, err := os.Open(src)
// 	if err != nil {
// 		return
// 	}
// 	defer in.Close()

// 	out, err := os.Create(dst)
// 	if err != nil {
// 		return
// 	}
// 	defer out.Close()

// 	io.Copy(out, in)
// }

// func generateMP4(dir string) {
// 	log.Println("Gerando MP4 final...")

// 	files, _ := filepath.Glob(filepath.Join(dir, "*.ts"))
// 	sort.Strings(files)

// 	listFile := filepath.Join(dir, "list.txt")
// 	f, err := os.Create(listFile)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer f.Close()

// 	for _, file := range files {
// 		fmt.Fprintf(f, "file '%s'\n", filepath.Base(file))
// 	}

// 	output := filepath.Join(dir, "output.mp4")

// 	cmd := exec.Command(
// 		`C:\ffmpeg\bin\ffmpeg.exe`,
// 		"-f", "concat",
// 		"-safe", "0",
// 		"-i", listFile,
// 		"-c", "copy",
// 		output,
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	err = cmd.Run()
// 	if err != nil {
// 		log.Println("Erro ao gerar MP4:", err)
// 		return
// 	}

// 	log.Println("MP4 gerado com sucesso:", output)
// }


// package main

// import (	
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/exec"
// 	"sync"
// 	"time"
// )

// const (
// 	bufferSeconds = 10
// )

// var (
// 	buffer      [][]byte
// 	mu          sync.Mutex
// 	recording   bool
// 	recordData  [][]byte
// 	ffmpegCmd   *exec.Cmd
// 	ffmpegPipe  io.ReadCloser
// 	stopFFmpeg  chan bool
// )

// func main() {
// 	os.MkdirAll("recordings", os.ModePerm)

// 	go startFFmpeg()

// 	http.HandleFunc("/start", startHandler)
// 	http.HandleFunc("/stop", stopHandler)

// 	log.Println("Servidor rodando em :8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func startFFmpeg() {
// 	stopFFmpeg = make(chan bool)

// 	ffmpegCmd = exec.Command(
// 		`C:\ffmpeg\bin\ffmpeg.exe`,
// 		"-rtsp_transport", "tcp",
// 		"-i", "rtsp://admin:Ferrasa321@10.100.150.32:554/Streaming/Channels/101",
// 		"-c", "copy",
// 		"-f", "mpegts",
// 		"pipe:1",
// 	)

// 	stdout, err := ffmpegCmd.StdoutPipe()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ffmpegPipe = stdout

// 	ffmpegCmd.Stderr = os.Stderr

// 	log.Println("FFmpeg iniciado...")

// 	err = ffmpegCmd.Start()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	go readStream(stdout)
// }

// func readStream(stdout io.ReadCloser) {
// 	buf := make([]byte, 1024*64)

// 	for {
// 		n, err := stdout.Read(buf)
// 		if err != nil {
// 			log.Println("FFmpeg finalizado")
// 			return
// 		}

// 		frame := make([]byte, n)
// 		copy(frame, buf[:n])

// 		mu.Lock()

// 		buffer = append(buffer, frame)

// 		if len(buffer) > bufferSeconds {
// 			buffer = buffer[1:]
// 		}

// 		if recording {
// 			recordData = append(recordData, frame)
// 		}

// 		mu.Unlock()
// 	}
// }

// func startHandler(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	if recording {
// 		fmt.Fprintln(w, "Já está gravando")
// 		return
// 	}

// 	recording = true

// 	// copia buffer (pré 10s)
// 	recordData = make([][]byte, len(buffer))
// 	copy(recordData, buffer)

// 	fmt.Fprintln(w, "Gravação iniciada (com buffer de 10s)")
// }

// func stopHandler(w http.ResponseWriter, r *http.Request) {
// 	mu.Lock()
// 	// defer mu.Unlock()

// 	if !recording {
// 		fmt.Fprintln(w, "Não está gravando")
// 		return
// 	}

// 	recording = false

// 	filename := fmt.Sprintf("recordings/%d.ts", time.Now().Unix())

// 	file, err := os.Create(filename)
// 	if err != nil {
// 		fmt.Fprintln(w, "Erro ao criar arquivo")
// 		return
// 	}
// 	defer file.Close()

// 	for _, chunk := range recordData {
// 		file.Write(chunk)
// 	}

// 	mu.Unlock()

// 	// converter para mp4 (fora do lock)
// 	go convertToMP4(filename)

// 	fmt.Fprintln(w, "Gravação finalizada:", filename)
// }

// func convertToMP4(input string) {
// 	out := input + ".mp4"

// 	cmd := exec.Command(
// 		`C:\ffmpeg\bin\ffmpeg.exe`,
// 		"-y",
// 		"-i", input,
// 		"-c", "copy",
// 		out,
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	log.Println("Convertendo para MP4...")

// 	err := cmd.Run()
// 	if err != nil {
// 		log.Println("Erro ao converter:", err)
// 		return
// 	}

// 	log.Println("MP4 gerado:", out)
// }



package main

import (
	"camsystem/schemas" // Ajuste para o nome do seu modulo no go.mod
	
)

// func main() {
// 	// 1. Criamos uma instância única da câmera (A nossa "Casa nº 1")
// 	// Vamos configurar um buffer minúsculo de 5 posições para facilitar o teste
// 	cam1 := schemas.NewCamera(1, "rtsp://10.100.150.32:554/ch1")
	
// 	fmt.Printf("Câmera [%d] iniciada no endereço: %p\n", cam1.ID, cam1)

// 	// 2. Simulando a entrada de 15 frames de vídeo
// 	fmt.Println("--- Iniciando captura de dados simulada ---")
// 	for i := 1; i <= 15; i++ {
// 		frame := []byte(fmt.Sprintf("Frame Video #%d", i))
		
// 		// Usamos o ponteiro para o Buffer e damos um Push
// 		cam1.Buffer.Push(frame)
		
// 		fmt.Printf("Inserindo: %s | Posição no Relógio: %d\n", string(frame), cam1.Buffer.WritePos)
// 		time.Sleep(100 * time.Millisecond) // Pequena pausa para vermos o log
// 	}

// 	// 3. Vamos ver o que sobrou no buffer? 
// 	// Como o tamanho é 5, ele deve ter apenas os frames de 11 a 15.
// 	fmt.Println("\n--- Dados na Ordem Correta (Temporal) ---")
// 	framesOrdenados := cam1.Buffer.GetAll()
// 	for i, frame := range framesOrdenados {
// 		fmt.Printf("Frame %d: %s\n", i, string(frame))
// 	}
// }

func main() {
	manager := schemas.NewStreamManager()

	// Você pode carregar isso de um banco de dados ou arquivo JSON no futuro
	listaDeCameras := map[int]string{
		1: "rtsp://admin:Ferrasa321@192.168.0.40:8000/Streaming/Channels/101",
	}

	for id, url := range listaDeCameras {
		manager.AddCamera(id, url)
	}
	// O programa principal fica rodando...
	select {
		
	} 
}



// package main

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"camsystem/schemas" // Ajuste para o seu módulo
// )

// var manager *schemas.StreamManager

// func main() {
// 	// 1. Iniciamos o nosso Maestro
// 	manager = schemas.NewStreamManager()

// 	// 2. Rota para Adicionar/Conectar uma Câmera
// 	// Exemplo: http://localhost:8080/add?id=1&url=rtsp://...
// 	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
// 		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
// 		url := r.URL.Query().Get("url")
		
// 		manager.AddCamera(id, url)
// 		fmt.Fprintf(w, "Câmera %d adicionada com sucesso!", id)
// 	})

// 	// 3. Rota para Iniciar Gravação (Pegando os 10s anteriores)
// 	// Exemplo: http://localhost:8080/start-record?id=1
// 	http.HandleFunc("/start-record", func(w http.ResponseWriter, r *http.Request) {
// 		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
// 		manager.StartRecording(id)
// 		fmt.Fprintf(w, "Gravação da câmera %d iniciada!", id)
// 	})

// 	// 4. Rota para Parar e Salvar MP4
// 	// Exemplo: http://localhost:8080/stop-record?id=1
// 	http.HandleFunc("/stop-record", func(w http.ResponseWriter, r *http.Request) {
// 		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
// 		path := manager.StopRecording(id)
// 		fmt.Fprintf(w, "Gravação encerrada. Arquivo sendo processado em: %s", path)
// 	})

// 	fmt.Println("Servidor de Monitoramento rodando em http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }
