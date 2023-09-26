package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
  "net/http"
  "os/exec"
  "strings"
)

func RestartDcs(c *gin.Context) {
  log.Println("Restarting DCS")
  exec.Command("sudo", "systemctl", "restart", "dcs")
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func DcsStatus(c *gin.Context) {
  out, err := exec.Command("sudo", "systemctl", "status", "dcs").Output()

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch service status"})
    return
  }

  if strings.Contains(string(out[:]), "active (running)") {
    c.JSON(http.StatusOK, gin.H{"status": "running"})
    return
  } 
  c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}

func RestartSrs(c *gin.Context) {
  log.Println("Restarting SRS")
  exec.Command("sudo", "systemctl", "restart", "srs")
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func SrsStatus(c *gin.Context) {
  out, err := exec.Command("sudo", "systemctl", "status", "srs").Output()

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Can't fetch service status"})
    return
  }

  if strings.Contains(string(out[:]), "active (running)") {
    c.JSON(http.StatusOK, gin.H{"status": "running"})
    return
  } 
  c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}
