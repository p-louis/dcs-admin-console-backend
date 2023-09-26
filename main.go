package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/p-louis/dcs-admin/controllers"
	"github.com/p-louis/dcs-admin/middlewares"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Uploading file")
}

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	basePath := os.Getenv("STATIC_DIRECTORY")

	router.StaticFile("", basePath+"/static/index.html")
	router.StaticFile("main.js", basePath+"/static/main.js")
	router.Static("css", basePath+"/static/css")
	router.Static("webfonts", basePath+"/static/webfonts")
	router.Static("img", basePath+"/static/img")

	public := router.Group("/api")
	public.POST("/login", controllers.Login)
	public.Static("/tacviews", os.Getenv("TACVIEW_DIRECTORY"))
	public.Static("/liberation", os.Getenv("LIBERATION_DIRECTORY"))

	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	protected.GET("/tacview", controllers.TacViews)
	protected.POST("/upload", controllers.Upload)
	protected.GET("/mission", controllers.Missions)
	protected.GET("/mission/current", controllers.CurrentMission)
	protected.POST("/mission", controllers.MissionChange)
	protected.DELETE("/mission", controllers.MissionRemove)
	protected.POST("/pause", controllers.PauseMission)
	protected.POST("/unpause", controllers.UnpauseMission)
	protected.GET("/pause", controllers.GetPause)
	protected.POST("/chat", controllers.SendChatMessage)
	protected.GET("/dcs", controllers.DcsStatus)
	protected.POST("/dcs", controllers.RestartDcs)
	protected.GET("/srs", controllers.SrsStatus)
	protected.POST("/srs", controllers.RestartSrs)

	router.Run("localhost:8080")
}
