package handlers

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"github.com/gin-gonic/gin"
	"net/http"
	"camsystem/db"
	"strconv"

)

type CameraRequest struct {
	ID int `json:"id"`
	URL string `json:"url"`
}

// func CreateEstacionamentos(c *gin.Context) {
// 	db := config.GetDB()

// 	var input schemas.Estacionamentos

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   "JSON inválido",
// 			"details": err.Error(),
// 		})
// 		return
// 	}
// 	if input.Nome == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "O campo 'nome' é obrigatório",
// 		})
// 		return
// 	}

// 	if err := db.Create(&input).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error":   "Erro ao criar estacionamento",
// 			"details": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, gin.H{
// 		"data": input,
// 	})
// }

func SaveCamera(c *gin.Context) {
	var cam CameraRequest

	if err := c.ShouldBindJSON(&cam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON inválido",
			"details": err.Error(),
		})
		return
	}

	if err := SaveCameraToDB(cam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro ao salvar câmera",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "câmera salva com sucesso",
	})
}

func GetCameras(c *gin.Context) {

	data, err := GetAllCameras()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "erro ao acessar banco de dados",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}


func SaveCameraToDB(cam CameraRequest) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cameras"))

		data, err := json.Marshal(cam)
		if err != nil {
			return err
		}

		key := []byte(fmt.Sprintf("%d", cam.ID))

		return b.Put(key, data)
	})
}

func GetAllCameras() (map[int]string, error) {
	cameras := make(map[int]string)

	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cameras"))

		return b.ForEach(func(k, v []byte) error {
			id, err := strconv.Atoi(string(k))
			if err != nil {
				return err
			}

			cameras[id] = string(v)
			return nil
		})
	})

	return cameras, err
}