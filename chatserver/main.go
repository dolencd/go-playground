package main

import (
	"net/http"

	"github.com/dolencd/go-playground/chatserver/common"
	"github.com/dolencd/go-playground/chatserver/rooms"
	"github.com/dolencd/go-playground/chatserver/users"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	godotenv.Load("../.env")
	conn, err := common.InitializeConnection()
	if err != nil {
		panic(err)
	}
	r := gin.Default()

	public := r.Group("/api/")
	public.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	userRepo := users.NewUserRepo(conn)
	roomRepo := rooms.NewRoomRepo(conn)
	private := r.Group("/api/")
	private.Use(common.PopulateUserMiddleware(&userRepo))
	private.Use(common.RequireUserMiddleware())
	users.NewUserController(public, &userRepo)
	rooms.NewRoomController(public, &roomRepo)

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
