package frontend

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kawijayaa/pintas/db"
)

func Routes(r *gin.Engine) {
	r.LoadHTMLGlob("frontend/templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"host": c.Request.Host,
		})
	})
	r.POST("/shorten", func(c *gin.Context) {
		var url db.ShortUrlDto
		var resUrl db.ShortUrl

		c.Bind(&url)
		urlJson, _ := json.Marshal(url)

		res, err := http.Post("http://localhost:8080/api/v1/urls", "application/json", bytes.NewBuffer(urlJson))
		if err != nil || res.StatusCode == 500 {
			c.Data(500, "text/html", []byte("500 Internal Server Error"))
			return
		}
		if res.StatusCode == 409 {
			c.Data(409, "text/html", []byte("409 Conflict"))
			return
		}
		if res.StatusCode == 400 {
			c.Data(400, "text/html", []byte("400 Bad Request"))
			return
		}

		body, err := io.ReadAll(res.Body)
		json.Unmarshal(body, &resUrl)

		c.HTML(200, "success.html", gin.H{
			"host": c.Request.Host,
			"path": resUrl.Path,
		})
	})
	r.GET("/:shortUrl", func(c *gin.Context) {
		var url db.ShortUrl

		conn := db.Connect()
		conn.Where("path = ?", c.Param("shortUrl")).First(&url)

		if url.RedirectUrl != "" {
			c.Redirect(301, url.RedirectUrl)
		} else {
			c.HTML(404, "404.html", gin.H{})
		}
	})
	r.Static("/static", "frontend/static")
}
