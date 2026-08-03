package router

import (
	"context"
	"net/http"

	"github.com/ao9911/go-matrix/transport/httpserver"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ao9911/bluebell-new/conf"
	"github.com/ao9911/bluebell-new/service"
)

var srv *service.Service

func startHttp(c *conf.Config) {
	srv = service.New(c)
	r := gin.New()
	r.UseH2C = true
	// r.Use(gin.Recovery(), RateLimitMiddleware(c.RateLimit.MaxQPS, c.RateLimit.MaxBurst))
	r.Use(gin.Recovery())
	r.GET("/health", health)
	r.StaticFile("/openapi.yaml", "./api/openapi.yaml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/openapi.yaml")))

	v1 := r.Group("/api/v1")
	v1.POST("/signup", SignUp)
	v1.POST("/login", Login)
	v1.POST("/refresh_token", RefreshToken)

	v1.Use(JWTAuthMiddleware())
	v1.GET("/community", GetCommunityList)
	v1.GET("/community/:community_id", GetCommunityDetail)

	v1.POST("/post", CreatePost)
	v1.GET("/post/:post_id", GetPostDetail)
	v1.GET("/post", GetPostList)
	v1.GET("/post2", GetPostList2)

	v1.POST("/vote", PostVote)

	c.HttpServer.Handler = r
	httpserver.Run(c.HttpServer)
}

func stopHttp(ctx context.Context) {
	httpserver.Stop(ctx)
	srv.Close()
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
