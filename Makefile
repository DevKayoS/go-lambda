.PHONY: build clean deploy test run help create-role create-lambda create-api delete-all

# Variaveis
FUNCTION_NAME=minha-api-go
REGION=us-east-1
BINARY_NAME=bootstrap
ZIP_FILE=lambda.zip
ROLE_NAME=lambda-execution-role
API_NAME=minha-api-go
STAGE_NAME=dev
ACCOUNT_ID=$(shell aws sts get-caller-identity --query Account --output text)


ifneq (,$(wildcard .env))
    include .env
endif

# Exporta todas as variáveis
export


help:
	@echo "Comandos disponiveis:"
	@echo "  make build         - Compila o binario para Lambda"
	@echo "  make zip           - Cria o arquivo zip"
	@echo "  make create        - Cria a Lambda na AWS (primeira vez)"
	@echo "  make create-api    - Cria API Gateway e retorna a URL"
	@echo "  make deploy        - Faz deploy na AWS Lambda"
	@echo "  make logs          - Mostra logs da Lambda"
	@echo "  make test          - Roda os testes"
	@echo "  make clean         - Remove arquivos gerados"
	@echo "  make delete-all    - Remove Lambda e API Gateway"
	@echo "  make migrate-new   - Cria uma nova migration"

build:
	@echo "Building for Lambda..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd
	@echo "Build concluido!"

zip: build
	@echo "Creating zip file..."
	zip -j $(ZIP_FILE) $(BINARY_NAME)
	@echo "Zip criado: $(ZIP_FILE)"

create-role:
	@echo "Criando IAM Role..."
	@aws iam create-role \
		--role-name $(ROLE_NAME) \
		--assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}' \
		2>/dev/null || echo "Role ja existe, continuando..."
	@aws iam attach-role-policy \
		--role-name $(ROLE_NAME) \
		--policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole \
		2>/dev/null || echo "Policy ja anexada, continuando..."
	@echo "Aguardando role propagar (10s)..."
	@sleep 10

create-lambda: zip create-role
	@echo "Criando Lambda function..."
	aws lambda create-function \
		--function-name $(FUNCTION_NAME) \
		--runtime provided.al2023 \
		--role arn:aws:iam::$(ACCOUNT_ID):role/$(ROLE_NAME) \
		--handler $(BINARY_NAME) \
		--zip-file fileb://$(ZIP_FILE) \
		--timeout 30 \
		--memory-size 256 \
		--region $(REGION) \
		--environment "Variables={DATABASE_URL=$(DATABASE_URL)}"
	@echo "Lambda criada com sucesso!"


update-env:
	@echo "Atualizando variaveis de ambiente..."
	aws lambda update-function-configuration \
		--function-name $(FUNCTION_NAME) \
		--environment "Variables={DATABASE_URL=$(DATABASE_URL)}" \
		--region $(REGION)
	@echo "Variaveis atualizadas!"

create: create-lambda
	@echo ""
	@echo "Pronto! Agora rode 'make create-api' para criar a URL publica"

create-api:
	@echo "Criando API Gateway..."
	@API_ID=$$(aws apigateway create-rest-api \
		--name $(API_NAME) \
		--region $(REGION) \
		--endpoint-configuration types=REGIONAL \
		--query 'id' \
		--output text); \
	echo "API ID: $$API_ID"; \
	ROOT_ID=$$(aws apigateway get-resources \
		--rest-api-id $$API_ID \
		--region $(REGION) \
		--query 'items[0].id' \
		--output text); \
	echo "Root Resource ID: $$ROOT_ID"; \
	echo "Criando recurso proxy..."; \
	RESOURCE_ID=$$(aws apigateway create-resource \
		--rest-api-id $$API_ID \
		--parent-id $$ROOT_ID \
		--path-part '{proxy+}' \
		--region $(REGION) \
		--query 'id' \
		--output text); \
	echo "Resource ID: $$RESOURCE_ID"; \
	echo "Criando metodo ANY..."; \
	aws apigateway put-method \
		--rest-api-id $$API_ID \
		--resource-id $$RESOURCE_ID \
		--http-method ANY \
		--authorization-type NONE \
		--region $(REGION) > /dev/null; \
	echo "Integrando com Lambda..."; \
	aws apigateway put-integration \
		--rest-api-id $$API_ID \
		--resource-id $$RESOURCE_ID \
		--http-method ANY \
		--type AWS_PROXY \
		--integration-http-method POST \
		--uri arn:aws:apigateway:$(REGION):lambda:path/2015-03-31/functions/arn:aws:lambda:$(REGION):$(ACCOUNT_ID):function:$(FUNCTION_NAME)/invocations \
		--region $(REGION) > /dev/null; \
	echo "Dando permissao para API Gateway invocar Lambda..."; \
	aws lambda add-permission \
		--function-name $(FUNCTION_NAME) \
		--statement-id apigateway-access-$$(date +%s) \
		--action lambda:InvokeFunction \
		--principal apigateway.amazonaws.com \
		--source-arn "arn:aws:execute-api:$(REGION):$(ACCOUNT_ID):$$API_ID/*/*" \
		--region $(REGION) > /dev/null 2>&1 || true; \
	echo "Fazendo deploy..."; \
	aws apigateway create-deployment \
		--rest-api-id $$API_ID \
		--stage-name $(STAGE_NAME) \
		--region $(REGION) > /dev/null; \
	echo ""; \
	echo "========================================"; \
	echo "API criada com sucesso!"; \
	echo "========================================"; \
	echo "URL: https://$$API_ID.execute-api.$(REGION).amazonaws.com/$(STAGE_NAME)"; \
	echo ""; \
	echo "Teste no Postman ou curl:"; \
	echo "  GET https://$$API_ID.execute-api.$(REGION).amazonaws.com/$(STAGE_NAME)/api/v1/health"; \
	echo "========================================"; \
	echo "$$API_ID" > .api-id


deploy: zip
	@echo "Deploying to AWS Lambda..."
	aws lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://$(ZIP_FILE) \
		--region $(REGION)
	@echo "Deploy concluido!"

logs:
	@echo "Fetching logs..."
	aws logs tail /aws/lambda/$(FUNCTION_NAME) --follow --region $(REGION)

test:
	@echo "Running tests..."
	go test -v ./...

run:
	@echo "Running locally..."
	go run ./cmd

clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME) $(ZIP_FILE)
	@echo "Clean concluido!"

delete-all:
	@echo "Deletando API Gateway..."
	@if [ -f .api-id ]; then \
		API_ID=$(cat .api-id); \
		aws apigateway delete-rest-api --rest-api-id $API_ID --region $(REGION) 2>/dev/null || true; \
		rm .api-id; \
	fi
	@echo "Deletando Lambda..."
	@aws lambda delete-function --function-name $(FUNCTION_NAME) --region $(REGION) 2>/dev/null || true
	@echo "Tudo deletado!"

install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencias instaladas!"

migrate-up:
	@echo "Running migrations with pgx driver..."
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "ERRO: DATABASE_URL não definida!"; \
		exit 1; \
	fi
	tern migrate -m ./internal/store/migrations --conn-string $(DATABASE_URL)

migrate-down:
	@echo "Rolling back last migration..."
	@tern migrate -m ./internal/pgstore/migrations --conn-string $(DATABASE_URL) -d -1

migrate-status:
	@echo "Migration status..."
	@tern status -m ./internal/pgstore/migrations --conn-string $(DATABASE_URL)

migrate-new:
	@read -p "Migration name: " name; \
	tern new -m ./internal/pgstore/migrations $$name

# SQLC
sqlc-generate:
	@echo "Generating code from SQL..."
	@sqlc generate -f ./internal/pgstore/sqlc.yaml

test-db-connection:
	@echo "Testing database connection..."
	@psql "$(DATABASE_URL)" -c "SELECT version();" || (echo "Connection failed!"; exit 1)


create-mock:
	@echo "Criando pasta de mocks"
	@mockery --all --output=./internal/mocks --case=underscore

.DEFAULT_GOAL := help
