package main

import (
	"context"
	"log"
	"os"

	handlers "mc-dns-image-lambda/handlers"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	middleware "mc-dns-image-lambda/middleware"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// This function loads the `.env` file.
//
// # The file should contain your AWS credentials
//
// AWS_REGION - AWS_ACCESS_KEY - AWS_ACCESS_KEY_ID
func LoadDotEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

var (
	ginLambda *ginadapter.GinLambda
	s3Client  *s3.Client
)

func init() {
	// LoadDotEnv()

	temp := os.Getenv("BUCKET_PUBLIC")
	log.Printf("BUCKET_PUBLIC %s", temp)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:5173", "https://mc-server.kr"}, // Replace with allowed origins
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour, // Cache preflight response duration
	}))

	router.Use(middleware.AuthMiddleware())

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	handlers.SetS3Client(s3Client)

	router.GET("/s3", handlers.ListObjects)

	uploader_route := router.Group("/upload")
	{
		uploader_route.POST("/image", handlers.UploadObject)
		uploader_route.POST("/logo", handlers.UploadIcon)
	}
	// router.POST("/upload/image", handlers.UploadObject)

	ginLambda = ginadapter.New(router)
}

func GinRequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {
	lambda.Start(GinRequestHandler)
}
