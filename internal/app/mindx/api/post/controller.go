package post

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"mindx/internal/app/mindx/models"
	"net/http"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) DefineEndpoints(routes *gin.RouterGroup) {
	routes.POST("/", h.CreatePost)
	routes.GET("/", h.GetPosts)
}

func (h *Handler) CreatePost(c *gin.Context) {
	var postPayload models.PostPayload
	err := c.BindJSON(&postPayload)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = h.db.AutoMigrate(&models.Post{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	tokenInfo := GetTokenInfo(c)
	userID := tokenInfo["user_id"].(float64)
	userName := tokenInfo["user_name"]
	result := h.db.Create(&models.Post{
		Content:  postPayload.Content,
		UserID:   uint(userID),
		UserName: userName.(string),
	})

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Added post successfully"),
		})
	}

}

func (h *Handler) GetPosts(c *gin.Context) {
	var posts []models.Post
	result := h.db.Find(&posts)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
		"count": len(posts),
	})

}

var hmacSampleSecret []byte

func GetTokenInfo(c *gin.Context) jwt.MapClaims {
	// sample token string taken from the New example
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	tokenString := c.Request.Header.Get("Token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

	return claims
}
