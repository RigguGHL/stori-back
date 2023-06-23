package main

import (
	"fmt"
	"net/http"
	"os"
	"stori-go/email"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	r.POST("/upload", posting)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// Bind data to body structure, save file on tmp folder and send to email.send function,
// then clean tmp folder and send response to client
func posting(c *gin.Context) {
	body := email.Body{}

	err := c.ShouldBind(&body)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	c.SaveUploadedFile(body.File, "tmp/"+body.File.Filename)

	err = email.Send(body)
	if err != nil {
		os.RemoveAll("tmp/")
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	os.RemoveAll("tmp/")

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent",
	})
}
