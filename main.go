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

	public := router.Group("/api")
	public.POST("/login", controllers.Login)
	public.Static("/tacviews", os.Getenv("TACVIEW_DIRECTORY"))

	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	protected.GET("/tacview", controllers.TacViews)
	protected.POST("/upload", controllers.Upload)
	protected.GET("/mission", controllers.Missions)
	protected.POST("/chat", controllers.SendChatMessage)

	router.Run("localhost:8080")
}
