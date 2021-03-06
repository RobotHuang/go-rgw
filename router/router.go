package router

import (
	"github.com/gin-gonic/gin"
	"go-rgw/auth"
)

func SetupRouter(a auth.Auth) *gin.Engine {
	r := gin.Default()
	authorized := registerAuthMiddleware(r, a)
	{
		authorized.GET("/createbucket/:bucket", createBucket)

		authorized.POST("/upload/:bucket/:object", putObject)
		authorized.GET("/download/:bucket/:object", getObject)

		authorized.POST("/uploads/create/:bucket/:object", createMultipartUpload)
		authorized.POST("/uploads/upload/:bucket/:object", uploadPart)
		authorized.POST("/uploads/complete/:bucket/:object", completeMultipartUpload)
		authorized.POST("/uploads/abort/:bucket/:object", abortMultipartUpload)

		authorized.GET("/image/blur/:bucket/:object", imageBlur)
		authorized.GET("/image/resize/:bucket/:object", imageResize)
		authorized.GET("/image/cropAnchor/:bucket/:object", imageCropAnchor)
	}
	return r
}

func registerAuthMiddleware(e *gin.Engine, auth auth.Auth) *gin.RouterGroup {
	e.POST("/register", auth.CreateUser)
	e.POST("/login", auth.Login)
	g := e.Group("/")
	g.Use(auth.Auth())
	return g
}
