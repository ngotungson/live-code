package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"mindx/internal/app/mindx/models"
	"net/http"
	"strconv"
	"strings"
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
	routes.POST("/login", h.Login)
	routes.POST("/signup", h.CreateUser)
}

func (h *Handler) ListAllUser(c *gin.Context) {
	// Paging
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "2"))

	offset := (page - 1) * perPage
	limit := perPage
	var users []models.User
	result := h.db.Offset(offset).Limit(limit)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
	}

	//Sorting
	sorting := c.DefaultQuery("sorting", "user_name")
	for _, sortedFiled := range strings.Split(sorting, ",") {
		// -username
		if string(sortedFiled[0]) == "-" {
			result = result.Order(fmt.Sprintf("%s desc", sortedFiled[1:]))
		} else {
			// +username || username
			if string(sortedFiled[0]) == "+" {
				result = result.Order(fmt.Sprintf("%s asc", sortedFiled[1:]))
			} else {
				result = result.Order(fmt.Sprintf("%s asc", sortedFiled))
			}
		}
	}

	result = result.Find(&users)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.UserPayload
	err := c.BindJSON(&user)
	if err != nil {
		fmt.Println(err.Error())
	}

	if user.UserName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "User name is not allowed to be blank",
		})
		return
	}

	if user.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Email is not allowed to be blank",
		})
		return
	}

	if user.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Password is not allowed to be blank",
		})
		return
	}

	err = h.db.AutoMigrate(&models.User{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	result := h.db.Create(&models.User{
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
	})

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Added user %s successfully", user.UserName),
		})
	}
}

func (h *Handler) Login(c *gin.Context) {
	var user models.UserPayload
	err := c.BindJSON(&user)
	if err != nil {
		fmt.Println(err.Error())
	}

	var userRecord models.User
	result := h.db.Where("email = ? AND password = ?", user.Email, user.Password).Find(&userRecord)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userRecord.ID,
		"user_name": userRecord.UserName,
	})

	// Sign and get the complete encoded token as a string using the secret
	var hmacSampleSecret []byte
	tokenString, err := token.SignedString(hmacSampleSecret)

	fmt.Println(tokenString, err)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}
}
