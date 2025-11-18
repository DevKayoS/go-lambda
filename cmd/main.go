package main

import (
	"context"
	"log/slog"

	"github.com/DevKayoS/go-lambda/internal/api"
	"github.com/DevKayoS/go-lambda/internal/store/database"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	ctx := context.Background()

	if err := database.Init(ctx); err != nil {
		slog.Error("Erro ao conectar com o banco: ", err.Error())
	}

	slog.Info("Banco Conectado com sucesso!")
}

func main() {
	r := api.SetupAPI()

	ginLambda = ginadapter.New(r)
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
