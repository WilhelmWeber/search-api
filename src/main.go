package main

import (
	"net/http"

	"github.com/WilhelmWeber/search-api/src/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("src/views/*.html")
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "This is IIIF Search API ver.1",
		})
	})
	r.GET("/service/manifest/search", controllers.Search)
	r.Run(":3000")
}
