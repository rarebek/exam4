package api

import (
	_ "4microservice/api_gateway/api/docs"
	v1 "4microservice/api_gateway/api/handlers/v1"
	"4microservice/api_gateway/config"
	"4microservice/api_gateway/pkg/logger"
	"4microservice/api_gateway/services"
	"4microservice/api_gateway/storage/repo"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AddCorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Option Struct
type Option struct {
	Cfg            config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Reds           repo.RedisStorageI
	NatsCon        *nats.Conn
}

func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Log:            option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Cfg,
		Reds:           option.Reds,
	})

	api := router.Group("/v1")
	api.Use(AddCorsMiddleware())

	//users
	api.POST("/admin/user/create", handlerV1.CreateUser)
	api.GET("/user/get/:id", handlerV1.GetUser)
	api.PUT("/user/update/:id", handlerV1.UpdateUser)
	api.DELETE("/user/delete/:id", handlerV1.DeleteUser)
	api.GET("/users/:page/:limit", handlerV1.GetAllUser)
	api.POST("/user/register", handlerV1.RegisterUser)
	api.POST("/user/verify/:email/:code", handlerV1.Verify)
	api.POST("/user/login/:email/:password", handlerV1.Login)

	//posts
	api.POST("/post/create", handlerV1.CreatePost)
	api.GET("/post/get/:id", handlerV1.GetPost)
	api.PUT("/post/update/:id", handlerV1.UpdatePost)
	api.GET("/posts/:page/:limit", handlerV1.GetAllPosts)
	api.DELETE("/post/delete/:id", handlerV1.DeletePost)

	//comments
	api.POST("/comment/create", handlerV1.CreateComment)
	api.GET("/comment/get/:id", handlerV1.GetComment)
	api.PUT("/comment/update/:id", handlerV1.UpdateComment)
	api.GET("/comments/:post_id", handlerV1.GetAllComments)
	api.DELETE("/comment/delete/:id", handlerV1.DeleteComment)
	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" && strings.HasPrefix(origin, "http://localhost") && strings.HasPrefix(origin, "http://rarebek.jprq.app") {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
