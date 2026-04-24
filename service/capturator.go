package service

import (
	. "camsystem/handlers"
	"camsystem/schemas"
)

func Capturarador(manager *schemas.StreamManager) {
	data, err := GetAllCameras()

	if err != nil {
		panic("Erro ao acessar banco de dados: " + err.Error())
	}

	for _, cam := range data {
		manager.AddCamera(cam.ID, cam.URL)
	}

	// O programa principal fica rodando...
	select {}
}
