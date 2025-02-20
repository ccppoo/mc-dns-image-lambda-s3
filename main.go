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

	middleware "mc-dns-image-lambda/middleware"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	ginLambda *ginadapter.GinLambda
	s3Client  *s3.Client
)

func init() {

	origin := os.Getenv("ALLOW_ORIGINS")

	setUpAWS()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	router.Use(middleware.AuthMiddleware())

	router.GET("/s3", handlers.ListObjects)

	uploader_route := router.Group("/upload")
	{
		uploader_route.POST("/image", handlers.UploadObject)
		uploader_route.POST("/logo", handlers.UploadIcon)
	}

	ginLambda = ginadapter.New(router)
}

func setUpAWS() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	handlers.SetS3Client(s3Client)
}

func GinRequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {
	lambda.Start(GinRequestHandler)
}
