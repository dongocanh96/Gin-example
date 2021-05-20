package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	r := gin.Default()

	r.Static("/static-file", "./assets")
	// r.Use(ACustomMiddleware)
	r.Use(ACustomMiddleware())

	r.GET("/ping", getPing)
	r.POST("/ping", postPing)
	r.GET("/detail/:id", getDetail)

	//Upload single file
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", func(context *gin.Context) {
		// single file
		file, _ := context.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		context.SaveUploadedFile(file, "./assets/uploads/"+uuid.New().String()+filepath.Ext(file.Filename))

		context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	//Upload multiple files
	r.POST("/upload-multiple", func(context *gin.Context) {
		// Multipart form
		form, _ := context.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			context.SaveUploadedFile(file, "./assets/uploads/"+uuid.New().String()+filepath.Ext(file.Filename))
		}
		context.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "Hello from v1 ping",
				})
			})
			v1.GET("/pong", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "Hello from v1 pong",
				})
			})
		}
		v2 := api.Group("/v2")
		{
			v2.GET("/a", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "Hello from v2 a",
				})
			})
			v2.GET("/b", func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{
					"message": "Hello from v2 b",
				})
			})
		}
	}
	r.Run(":8080")
}

func getPing(context *gin.Context) {

	log.Println("I am in get ping handler")
	name := context.DefaultQuery("name", "Guest")
	var data = map[string]interface{}{
		"message": "Hello " + name + " from get ping",
	}
	context.JSON(http.StatusOK, data)
}

func postPing(context *gin.Context) {

	log.Println("I am in post ping handler")
	add := context.DefaultPostForm("add", "Viet Nam")
	context.JSON(http.StatusOK, gin.H{
		"message": "Hello from " + add + " post ping",
	})
}

func getDetail(context *gin.Context) {

	log.Println("I am in post ping handler")
	id := context.Param("id")
	context.JSON(http.StatusOK, gin.H{
		"message": id,
	})
}

func ACustomMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Println("I'm in a global middleware")
		if true {
			context.Next()
		}
	}
}

// func ACustomMiddleware(context *gin.Context) {
// 	log.Println("Hello, i am a middleware")
// 	context.Next()
// }
