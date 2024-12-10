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
		AllowOrigins:     []string{"https://localhost:5173"}, // Replace with allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight response duration
	}))

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	handlers.SetS3Client(s3Client)

	router.GET("/", handlers.PlayersIntroHandler)
	router.GET("/s3", handlers.ListObjects)
	router.POST("/upload", handlers.UploadObject)
	temp_bucket_route := router.Group("/temp")
	{
		temp_bucket_route.POST("/upload", handlers.UploadObject)
		temp_bucket_route.DELETE("/delete", handlers.DeleteObject)
	}
	ginLambda = ginadapter.New(router)
}

func GinRequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {
	lambda.Start(GinRequestHandler)
}
