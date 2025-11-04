package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/paaart/kavalife-erp-backend/config"
	"github.com/paaart/kavalife-erp-backend/internal/db"
	"github.com/paaart/kavalife-erp-backend/internal/routes"
	util "github.com/paaart/kavalife-erp-backend/internal/utils"
)

var (
	ginLambda *ginadapter.GinLambda
)

func main() {
	Log := util.InitLogger()
	appConfig := config.ConfigLoader()

	// Connect to DB
	pool, err := db.Connect(appConfig)
	if err != nil {
		Log.Fatal("Database connection failed:", err)
	}
	defer pool.Close()

	Log.Info("Starting Kava Life ERP Backend on Lambda...")

	// Initialize Gin routes
	r := routes.RunApp()
	ginLambda = ginadapter.New(r)

	// Graceful shutdown (useful for local dev / SAM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		Log.Info("Shutting down server...")
		pool.Close()
		os.Exit(0)
	}()

	// Start the Lambda handler
	lambda.Start(Handler)
}

// Lambda handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
