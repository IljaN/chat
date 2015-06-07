package main

import (
	"fmt"
	"io"
	"math/rand"
	"github.com/gin-gonic/gin"


)

const baseUrl string ="/chat"
var backend *Backend

func main() {

	backend = NewBackend()
	router := gin.Default()
	router.SetHTMLTemplate(html)
	router.BasePath = baseUrl

	router.GET("/rooms", roomIndexGET)
	router.GET("/rooms/:roomid", roomGET)
	router.POST("/rooms/:roomid", roomPOST)
	router.DELETE("/rooms/:roomid", roomDELETE)
	router.GET("/streams/:roomid", stream)


	router.Run("golang-VirtualBox.fritz.box:8080")
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

func roomIndexGET(c *gin.Context) {
	channelList := backend.GetRoomChannels()
	c.JSON(200, channelList)
}

func roomGET(c *gin.Context) {
	roomid := c.Param("roomid")
	userid := fmt.Sprint(rand.Int31())
	c.HTML(200, "chat_room", gin.H{
		"roomid": roomid,
		"userid": userid,
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
	roomid := c.Param("roomid")
	backend.DeleteBroadcast(roomid)
}
