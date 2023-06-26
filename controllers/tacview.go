package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/p-louis/dcs-admin/models"

	"net/http"
	"os"
)

func TacViews(c *gin.Context) {
	src := os.Getenv("TACVIEW_DIRECTORY")

	contents, err := os.ReadDir(src)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error fetching tacview files"})
		return
	}

	var tacviews []models.Mission

	for i := range contents {
		if !contents[i].IsDir() {
			mis := models.Mission{}
			mis.Filename = contents[i].Name()
			tacviews = append(tacviews, mis)
		}
	}

	c.JSON(http.StatusOK, tacviews)
}
