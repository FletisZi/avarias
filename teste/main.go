package main

import (
    "log"
    "os"
    "fmt"
)

func main() {
    arquivo, err := os.ReadFile("texto.txt") // Para acesso de leitura. 
    if err != nil { 
        log.Fatal(err) 
    }
    message := fmt.Sprintf("Hi, %v. Welcome!", arquivo)
    fmt.Println(message)
}