package comment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mindx/internal/app/mindx/models"
	"net/http"
	"strconv"
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
	routes.POST("/", h.CreateComment)
	routes.GET("/", h.GetComments)
}

func (h *Handler) CreateComment(c *gin.Context) {
	var commentPayload models.CommentPayload
	err := c.BindJSON(&commentPayload)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = h.db.AutoMigrate(&models.Comment{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	postID := c.Param("id")
	u64, err := strconv.ParseUint(postID, 10, 32)
	wd := uint(u64)
	result := h.db.Create(&models.Comment{
		Content: commentPayload.Content,
		PostID:  wd,
	})

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Added comment successfully"),
		})
	}

}

func (h *Handler) GetComments(c *gin.Context) {
	postID := c.Param("id")
	u64, _ := strconv.ParseUint(postID, 10, 32)
	wd := uint(u64)

	var comments []models.Comment
	result := h.db.Where("post_id = ?", wd).Find(&comments)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
		"count":    len(comments),
	})

}
