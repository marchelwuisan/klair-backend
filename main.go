package main

import (
	"context"
	"golanglearn/configs"
	"golanglearn/routes"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"

	"github.com/gofiber/fiber/v2"
)

var fiberLambda *fiberadapter.FiberLambda

func main() {
	lambda.Start(Handler)
}

func init() {

	log.Printf("Fiber cold start")

	app := fiber.New()

	configs.ConnectDB()

	configs.SetupFirebase()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	routes.Route(app)

	app.Listen(":6000")

	fiberLambda = fiberadapter.New(app)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fiberLambda.ProxyWithContext(ctx, req)
}
