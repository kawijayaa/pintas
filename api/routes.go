package api

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kawijayaa/pintas/db"
	"gorm.io/gorm"
)

func isValidUrl(str string) bool {
	url, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	address := net.ParseIP(url.Host)

	if address == nil {
		return strings.Contains(url.Host, ".")
	}

	return true
}

func hasSchema(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func Routes(r *gin.Engine) {
	r.GET("/api/v1/urls", func(c *gin.Context) {
		var urls []db.ShortUrl
		conn := db.Connect()
		conn.Find(&urls)
		c.JSON(200, urls)
	})
	r.POST("/api/v1/urls", func(c *gin.Context) {
		var urlDto db.ShortUrlDto

		c.BindJSON(&urlDto)

		if !hasSchema(urlDto.RedirectUrl) {
			urlDto.RedirectUrl = "https://" + urlDto.RedirectUrl
		}

		if !isValidUrl(urlDto.RedirectUrl) {
			c.JSON(400, gin.H{
				"error": "Invalid URL",
			})
			return
		}

		url := db.ShortUrl{
			Id:          uuid.New(),
			Path:        urlDto.Path,
			RedirectUrl: urlDto.RedirectUrl,
		}

		conn := db.Connect()
		result := conn.Create(&url)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				c.JSON(409, gin.H{
					"error": "Path already exists",
				})
				return
			} else {
				c.JSON(500, gin.H{
					"error": result.Error.Error(),
				})
				return
			}
		}

		c.JSON(200, url)
	})
}
