package handlers

import (
	"camsystem/db"
	"camsystem/schemas"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
)

type CameraRequest struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func SaveCamera(c *gin.Context) {
	var cam CameraRequest

	if err := c.ShouldBindJSON(&cam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "JSON inválido",
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

func GetCamera(manager *schemas.StreamManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID int `json:"id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON inválido",
			})
			return
		}

		cam, exists := manager.GetCamera(req.ID)

		fmt.Printf("[Handler] Requisição para câmera %d. Existe? %v\n", req.ID, cam.URL)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "câmera não encontrada",
			})
			return
		}

		c.JSON(http.StatusOK, cam.URL)
	}
}

func StartCamera(manager *schemas.StreamManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID int `json:"id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON inválido",
			})
			return
		}

		cam, exists := manager.StartGravação(req.ID)

		fmt.Printf("[Handler] Requisição para câmera %d. Existe? %v\n", req.ID, cam.URL)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "câmera não encontrada",
			})
			return
		}

		c.JSON(http.StatusOK, cam.URL)
	}
}

func StopCamera(manager *schemas.StreamManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID int `json:"id"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON inválido",
			})
			return
		}

		cam, exists := manager.StopGravação(req.ID)

		fmt.Printf("[Handler] Requisição para câmera %d. Existe? %v\n", req.ID, cam.URL)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "câmera não encontrada",
			})
			return
		}

		c.JSON(http.StatusOK, cam.RecordingBuffer)
		cam.RecordingBuffer = make([][]byte, 0) // Limpa o buffer de gravação
	}
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

func GetAllCameras() (map[int]CameraRequest, error) {
	cameras := make(map[int]CameraRequest)

	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cameras"))

		return b.ForEach(func(k, v []byte) error {
			id, err := strconv.Atoi(string(k))
			if err != nil {
				return err
			}

			var cam CameraRequest
			if err := json.Unmarshal(v, &cam); err != nil {
				return err
			}
			cameras[id] = cam
			return nil
		})
	})

	return cameras, err
}
