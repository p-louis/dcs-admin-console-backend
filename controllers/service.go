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
  if GetDcsStatus() == 1 {
    exec.Command("sudo", "systemctl", "restart", "dcs").Run()
  } else {
    exec.Command("sudo", "systemctl", "start", "dcs").Run()
  }
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func DcsStatus(c *gin.Context) {
  if GetDcsStatus() == 1 {
    c.JSON(http.StatusOK, gin.H{"status": "running"})
    return
  }
  c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}


func RestartSrs(c *gin.Context) {
  log.Println("Restarting SRS")

  if GetSrsStatus() == 1 {
    exec.Command("sudo", "systemctl", "restart", "srs").Run()
  } else {
    exec.Command("sudo", "systemctl", "start", "srs").Run()
  }
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func SrsStatus(c *gin.Context) {
  if GetSrsStatus() == 1 {
    c.JSON(http.StatusOK, gin.H{"status": "running"})
    return
  } 
  c.JSON(http.StatusOK, gin.H{"status": "stopped"})
}

func GetDcsStatus() int {
  out, _ := exec.Command("sudo", "systemctl", "status", "dcs").Output()

  if strings.Contains(string(out[:]), "active (running)") {
    return 1
  } 
  return 0
}

func GetSrsStatus() int {
  out, _ := exec.Command("sudo", "systemctl", "status", "srs").Output()

  if strings.Contains(string(out[:]), "active (running)") {
    return 1
  } 
  return 0
}
