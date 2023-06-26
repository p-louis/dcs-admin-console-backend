package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/p-louis/dcs-admin/models"

	// "github.com/p-louis/dcs-admin/utils/token"
	"net/http"
	"os"
	"strings"
)

func Upload(c *gin.Context) {
	dst := os.Getenv("MISSION_DIRECTORY")
	//dst := "/home/dcs/wine/DCSWorld/drive_c/users/dcs/Saved Games/DCS.openbeta_server/Missions/"
	// single file
	file, _ := c.FormFile("file")
	if !strings.HasSuffix(file.Filename, ".miz") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not a mission file"})
		return
	}

	log.Printf("Saving to %s/%s", dst, file.Filename)

	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, dst+"/"+file.Filename)
	log.Print(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error saving mission file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Missions(c *gin.Context) {
	src := os.Getenv("MISSION_DIRECTORY")

	contents, err := os.ReadDir(src)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error fetching mission files"})
		return
	}

	var missions []models.Mission

	for i := range contents {
		if !contents[i].IsDir() && strings.HasSuffix(contents[i].Name(), ".miz") {
			mis := models.Mission{}
			mis.Filename = contents[i].Name()
			missions = append(missions, mis)
		}
	}

	c.JSON(http.StatusOK, missions)
}