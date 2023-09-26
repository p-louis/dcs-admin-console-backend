package controllers

import (
	"encoding/json"
	"log"
	"net"
	"time"
    "strconv"

	"github.com/gin-gonic/gin"
	"github.com/p-louis/dcs-admin/models"

	// "github.com/p-louis/dcs-admin/utils/token"
	"bufio"
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

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return
	}

  _, err = conn.Write([]byte("{\"command\":\"append_mission\", \"mission_name\":\"C:\\\\users\\\\dcs\\\\Saved Games\\\\DCS.openbeta_server\\\\Missions\\\\" + file.Filename + "\"}\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Missions(c *gin.Context) {
  missions := FetchMissionList(c)

	c.JSON(http.StatusOK, missions)
}

func FetchMissionList(c *gin.Context) []models.Mission {
	var missions []models.Mission

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return missions
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return missions
	}

  _, err = conn.Write([]byte("{\"command\":\"get_missionlist\"}\n"))
  if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return missions
	}

	reader := bufio.NewReader(conn)
	reply, err := reader.ReadString('\n')
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading mission-data from DCS"})
		conn.Close()
		return missions
	}
	conn.Close()

	log.Printf("Received Missionlist %s",reply)

	var result models.MissionListResult
	json.Unmarshal([]byte(reply), &result)

  for i := 0; i < len(result.MissionList.Missions); i++ {
    var mis models.Mission
    mis.Index = i+1
    mis.Filename = result.MissionList.Missions[i]
    missions = append(missions, mis)
  }

  return missions
}

func CurrentMission(c *gin.Context) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return
	}
	conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	_, err = conn.Write([]byte("{\"command\":\"get_mission\"}\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return
	}

	reader := bufio.NewReader(conn)
	reply, err := reader.ReadString('\n')
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading mission-data from DCS"})
		conn.Close()
		return
	}

	conn.Close()

	var result map[string]any
	json.Unmarshal([]byte(reply), &result)

	c.JSON(http.StatusOK, gin.H{"filename": result["filename"]})
}

func MissionChange(c *gin.Context) {
	var requestBody models.ChangeMissionBody
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected request body"})
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return
	}
	conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	_, err = conn.Write([]byte("{\"command\":\"run_mission\", \"index\":" + strconv.Itoa(requestBody.MissionIndex) +"}\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func MissionRemove(c *gin.Context) {
	var requestBody models.ChangeMissionBody
	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unexpected request body"})
		return
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return
	}
	conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	_, err = conn.Write([]byte("{\"command\":\"delete_mission\", \"index\":" + strconv.Itoa(requestBody.MissionIndex) +"}\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func PauseMission(c *gin.Context) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return
	}
	conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	_, err = conn.Write([]byte("{\"command\":\"pause_mission\"}\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func UnpauseMission(c *gin.Context) {

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return
	}
	conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	_, err = conn.Write([]byte("{\"command\":\"unpause_mission\"}\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func GetPause(c *gin.Context) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error resolving TCP address"})
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error connecting to DCS"})
		return
	}
	conn.SetWriteDeadline(time.Now().Add(20 * time.Second))

	_, err = conn.Write([]byte("{\"command\":\"get_pause\"}\n"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error writing command to DCS"})
		conn.Close()
		return
	}

	reader := bufio.NewReader(conn)
	reply, err := reader.ReadString('\n')
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error reading pause-state from DCS"})
		conn.Close()
		return
	}

	conn.Close()

	var result map[string]any
	json.Unmarshal([]byte(reply), &result)

	c.JSON(http.StatusOK, gin.H{"pause_state": result["pause_state"]})
}
