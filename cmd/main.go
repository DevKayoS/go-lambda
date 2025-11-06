package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func main() {
	// Cria o roteador do Gin
	r := gin.Default()

	// Rota simples de health check
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Adiciona mais rotas aqui, ex:
	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, " + name,
		})
	})

	// Cria o adaptador que converte eventos API Gateway -> requests do Gin
	ginLambda = ginadapter.New(r)

	// Inicia a lambda
	lambda.Start(Handler)
}

// Handler que integra com o runtime da AWS Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}
