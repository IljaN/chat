package main

import (
	"fmt"
	"github.com/IljaN/chat/user"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math/rand"
	"net/http"
)

const baseUrl string = "/chat"
const host string = "golang-VirtualBox.fritz.box:8080"
const pubKeyPath string = "keys/app.rsa.pub"
const privKeyPath string = "keys/app.rsa"

var backend *Backend
var ch *Chat
var persistence *user.Persistence
var manager *user.Manager

func main() {

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")

	if err != nil {
		log.Fatal(err)
	}

	persistence = &user.Persistence{conn}
	a := user.NewAuthenticator(privKeyPath, pubKeyPath)
	manager = user.NewManager(persistence, a)

	backend = NewBackend()
	ch = &Chat{backend, make(map[string]Room)}

	router := gin.Default()
	router.SetHTMLTemplate(html)
	router.BasePath = baseUrl
	router.POST("/login", login)
	router.GET("/rooms", roomIndexGET)
	router.POST("/rooms", createRoom)
	router.GET("/rooms/:roomid", roomGET)
	router.POST("/rooms/:roomid", roomPOST)
	router.DELETE("/rooms/:roomid", roomDELETE)
	router.GET("/streams/:roomid", stream)
	router.GET("/protected", Firewall(manager), protected)

	router.GET("/debug", debug)

	router.Run(host)
}

func stream(c *gin.Context) {
	roomid := c.Param("roomid")
	listener := backend.OpenListener(roomid)
	defer backend.CloseListener(roomid, listener)

	c.Stream(func(w io.Writer) bool {
		c.SSEvent("message", <-listener)
		return true
	})
}

func login(c *gin.Context) {

	token, err := manager.Login(c.Request.FormValue("username"), c.Request.FormValue("password"))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Header("Chat-Auth", token)

}

func protected(c *gin.Context) {
	c.String(http.StatusOK, "INSIDE!")

}

func createRoom(c *gin.Context) {
	r := Room{}
	c.BindJSON(&r)
	locFormat := c.Request.URL.Path + "/%v"
	r = ch.CreateRoom(r, locFormat)

	c.Header("Location", r.Location)
	c.String(http.StatusCreated, "")
}

func roomIndexGET(c *gin.Context) {
	c.JSON(200, ch.Rooms)
}

func roomGET(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := fmt.Sprint(rand.Int31())
	c.HTML(200, "chat_room", gin.H{
		"roomid":  roomid,
		"userid":  userid,
		"baseUrl": baseUrl,
	})
}

func roomPOST(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := c.PostForm("user")
	message := c.PostForm("message")
	backend.Room(roomid).Submit(userid + ": " + message)

	c.JSON(200, gin.H{
		"status":  "success",
		"message": message,
	})

}

func roomDELETE(c *gin.Context) {
	id := c.Param("roomid")
	if ch.DissolveRoom(id) {
		c.String(http.StatusOK, "")
	} else {
		c.String(http.StatusNotFound, "")
	}
}

func debug(c *gin.Context) {
	c.String(200, "#backendRoomChannels: %v | #rooms: %v", len(ch.Backend.roomChannels), len(ch.Rooms))
}
