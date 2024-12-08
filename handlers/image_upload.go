package handlers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PlayersIntroHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Welcome to the players endpoint",
		"data": gin.H{
			"endpoints": gin.H{
				"players": []gin.H{
					{
						"method": "GET",
						"path":   "/api/v1/players",
						"desc":   "Get all players",
					},
					{
						"method": "GET",
						"path":   "/api/v1/players/:id",
						"desc":   "Get player Informations",
					},
					{
						"method": "POST",
						"path":   "/api/v1/players/create",
						"desc":   "Create a player",
					},
					{
						"method": "PUT",
						"path":   "/api/v1/players/update/:d",
						"desc":   "Create a player",
					},
				},
			},
		},
	})
}

type FileJson struct {
	name string          `form:"name" binding:"required"`
	file *multipart.File `form:"file"`
}
