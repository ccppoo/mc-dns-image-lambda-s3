package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	handlers "mc-dns-image-lambda/handlers"

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

	temp := os.Getenv("BUCKET_TEMP")
	log.Printf("temp %s", temp)
	static := os.Getenv("BUCKET_STATIC")
	log.Printf("static %s", static)
	user := os.Getenv("BUCKET_USER")
	log.Printf("user %s", user)

	router := gin.Default()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
	handlers.SetS3Client(s3Client)

	router.GET("/", handlers.PlayersIntroHandler)
	router.GET("/s3", handlers.ListObjects)
	router.POST("/upload", handlers.UploadObject)
	// api := router.Group("/api")
	// {
	// 	players := api.Group("/players")

	// 	// Mapping Player routes to their handlers
	// 	players.GET("/", handlers.PlayersIntroHandler)
	// 	// players.GET("/get/:id", PlayerRequestsHandler.GetPlayerHandler)
	// 	// players.GET("/getAll", PlayerRequestsHandler.GetAllPlayersHandler)
	// 	// players.POST("/create", PlayerRequestsHandler.CreatePlayerHandler)
	// 	// players.PATCH("/update/:id", PlayerRequestsHandler.UpdatePlayerHandler)
	// 	// players.DELETE("/delete/:id", PlayerRequestsHandler.DeletePlayerHandler)
	// }
	ginLambda = ginadapter.New(router)
}

func GinRequestHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func main() {
	lambda.Start(GinRequestHandler)
}
