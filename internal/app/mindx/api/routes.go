package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mindx/internal/app/mindx/api/comment"
	"mindx/internal/app/mindx/api/post"
	"mindx/internal/app/mindx/api/user"
	database2 "mindx/internal/app/mindx/database"
)

func SetupDefault() {
	// All defaults are set here
	viper.SetDefault("api.listening", ":8080")
}

// InitAPI configure and start API
func InitAPI() {
	SetupDefault()
	apiListening := viper.GetString("api.listening")
	go func() {
		router := gin.Default()

		initRoutes(router)
		router.Run(apiListening)
	}()
}

func NewOpenAPIMiddleware() gin.HandlerFunc {
	validator := OpenapiInputValidator("./openapi.yaml")
	return validator
}

func initRoutes(router *gin.Engine) {

	openapiMiddleware := NewOpenAPIMiddleware()

	database := database2.GetDB()

	userRoute := router.Group("/users")
	{
		userHandler := user.NewHandler(database)
		userRoute.Use(openapiMiddleware)
		userHandler.DefineEndpoints(userRoute)
	}
	postRoute := router.Group("/posts")
	{
		postHandler := post.NewHandler(database)
		postRoute.Use(openapiMiddleware)
		postHandler.DefineEndpoints(postRoute)
	}

	commentRoute := router.Group("/posts/:id/comments")
	{
		commentHandler := comment.NewHandler(database)
		commentRoute.Use(openapiMiddleware)
		commentHandler.DefineEndpoints(commentRoute)
	}
}

var InitRoutes = initRoutes
