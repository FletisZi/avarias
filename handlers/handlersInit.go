package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PageVagasEstacionamento(c *gin.Context) {
	c.HTML(http.StatusOK, "vagasestacionamento.html", gin.H{})
}