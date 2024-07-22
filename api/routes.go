package api

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kawijayaa/pintas/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func generateHash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func verifyHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(user db.User) (string, error) {
	expiry := time.Now().Add(time.Minute * 15)
	jti := uuid.New()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": jti,
		"sub": user.Id,
		"exp": expiry.Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	conn := db.Connect()
	result := conn.Create(&db.UserToken{
		Id:      jti,
		UserId:  user.Id,
		User:    user,
		Token:   tokenString,
		Expires: expiry.Unix(),
	})

	if result.Error != nil {
		return "", fmt.Errorf("Failed to create token: %s", result.Error.Error())
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		conn := db.Connect()

		var fetchedToken db.UserToken
		result := conn.Where("id = ?", claims["jti"].(string)).First(&fetchedToken)
		if result.Error != nil {
			return nil, false
		}

		expiry, _ := claims.GetExpirationTime()
		if expiry.Unix() < time.Now().Unix() {
			return nil, false
		}

		var fetchedUser db.User
		result = conn.Where("id = ?", claims["sub"].(string)).First(&fetchedUser)
		if result.Error != nil {
			return nil, false
		}

		return token.Claims.(jwt.MapClaims), true
	}

	return nil, false
}

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
	r.POST("/api/v1/login", func(c *gin.Context) {
		var user db.User
		var fetchedUser db.User
		c.BindJSON(&user)

		conn := db.Connect()
		result := conn.Where("username = ?", user.Username).First(&fetchedUser)
		if result.Error != nil {
			c.JSON(401, gin.H{
				"error": "Invalid username or password",
			})
			return
		}

		if !verifyHash(user.Password, fetchedUser.Password) {
			c.JSON(401, gin.H{
				"error": "Invalid username or password",
			})
			return
		}

		token, err := generateToken(fetchedUser)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Header("Authorization", "Bearer "+token)
		c.JSON(200, gin.H{
			"token": token,
		})
	})
	r.POST("/api/v1/register", func(c *gin.Context) {
		var user db.User
		c.BindJSON(&user)

		conn := db.Connect()

		user.Id = uuid.New()
		user.Password = generateHash(user.Password)

		result := conn.Create(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				c.JSON(409, gin.H{
					"error": "Username already exists",
				})
				return
			} else {
				c.JSON(500, gin.H{
					"error": result.Error.Error(),
				})
			}
			return
		}
		c.JSON(201, user)
	})
	r.GET("/api/v1/urls", func(c *gin.Context) {
		var urls []db.Url
		conn := db.Connect()
		conn.Find(&urls)
		c.JSON(200, urls)
	})
	r.POST("/api/v1/urls", func(c *gin.Context) {
		var urlDto db.UrlDto

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

		url := db.Url{
			Id:          uuid.New(),
			Path:        urlDto.Path,
			RedirectUrl: urlDto.RedirectUrl,
		}

		conn := db.Connect()

		if c.GetHeader("Authorization") != "" {
			token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
			if claims, ok := verifyToken(token); !ok {
				c.JSON(401, gin.H{
					"error": "Invalid token",
				})
				return
			} else {
				var user db.User
				conn.Where("id = ?", claims["sub"].(string)).First(&user)
				url.Owner = user
				url.OwnerId = &user.Id
			}
		}

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
