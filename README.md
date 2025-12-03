# Go Lambda API – Serverless API com Go, JWT, SQLC e AWS

Este repositório foi criado com o objetivo de estudar e entender como construir APIs serverless utilizando **Golang**, explorando autenticação com **JWT**, geração de queries com **SQLC**, migrations com **Tern**, além de integrar a aplicação com **AWS Lambda** e **API Gateway**. O projeto simula uma API real com criação de usuários, autenticação e controle de permissões.

---

## 🚀 Tecnologias Utilizadas

* **Golang** – Linguagem principal da aplicação
* **AWS Lambda** – Execução serverless da API
* **API Gateway** – Exposição pública da função Lambda como endpoint HTTP
* **JWT (JSON Web Token)** – Autenticação e autorização de usuários
* **SQLC** – Geração de código Go a partir de queries SQL
* **Tern** – Migrations de banco de dados
* **Makefile** – Automatização de tarefas e deploy
* **Gin Framework** – Criação das rotas HTTP e middlewares

---

## 📚 Objetivo do Projeto

O foco principal deste repositório é o **estudo** da arquitetura serverless com Go. Aqui você encontrará:

* Uma API completa com autenticação JWT
* Criação e gerenciamento de usuários
* Middlewares responsáveis por validar permissões de acesso
* Estrutura organizada para entender como montar um projeto Go escalável
* Processo de deploy inicial usando AWS Lambda
* Scripts automatizados via Makefile

Tudo foi desenvolvido com o intuito de aprender na prática como funciona a construção e deploy de APIs serverless utilizando Go.

---

## 🔐 Autenticação e Autorização

A API utiliza JWT para autenticação. Algumas funcionalidades implementadas:

* **Login** gerando um token JWT válido
* **Rota protegida** para criação de usuários
* **Sistema de permissões** simples
* **Middlewares** que validam se o usuário pode acessar determinada rota

Isso permite simular um ambiente real onde cada usuário tem permissões específicas.

---

## 🏗️ Estrutura do Projeto

A estrutura foi organizada com foco em modularidade, testabilidade e clareza, utilizando **interfaces** para todos os serviços e controllers — o mesmo padrão que você usa no dia a dia. Isso facilita manutenção, mocks em testes e substituição de implementações.

```
.
├── cmd/
│   └── main.go                  → Entrada da aplicação
├── internal/
│   ├── api/                     → Setup geral da API e router
│   ├── controllers/             → Controllers (todos usando interfaces)
│   ├── errors/                  → Estrutura de erros customizados
│   ├── middleware/              → Autenticação, autorização e handler de erros
│   ├── models/                  → Modelos principais (Token, Transaction)
│   ├── pgstore/                 → Código SQLC + migrations + queries
│   ├── routes/                  → Separação organizada das rotas
│   ├── services/                → Serviços com interfaces (user, token, transaction)
│   └── utils/                   → Utilidades gerais (ex: hash de senha)
├── go.mod
├── go.sum
└── Makefile                     → Scripts automatizados (build/deploy)
```

---

## 🛠️ Migrations com Tern

O **Tern** foi utilizado para gerenciar a evolução do banco de dados.

Comandos úteis:

```bash
make migrate-new
make migrate-up
```

---

## 🧬 SQLC – Funções de Banco Geradas Automaticamente

O **SQLC** converte queries SQL em código Go totalmente tipado. Isso reduz erros e aumenta produtividade.

Comando usado:

```bash
sqlc-generate:
```

---

## ☁️ Deploy Serverless na AWS

A API é empacotada e enviada como uma **função Lambda**, utilizando:

* **AWS Lambda** para executar o código Go
* **API Gateway** para receber requisições HTTP

A automação inicial está no `Makefile`, permitindo futuramente evoluir para um **CI/CD** completo.

---

## 📦 Makefile – Scripts Automatizados

O projeto possui um Makefile que ajuda nas tarefas de desenvolvimento e deploy, como:

* Gerar migrations
* Gerar SQLC
* Build da Lambda
* Deploy da função na AWS

Exemplo:

```bash
make deploy
```

---

## 🔮 Próximos Passos

* Configurar CI/CD automático
* Melhorar logs e rastreamento
* Criar novos módulos de domínio (ex: roles, audit logs)
* Adicionar testes automatizados

---

## 🤝 Contribuições

Este projeto é totalmente experimental, mas contribuições e sugestões são bem-vindas!

Fique à vontade para abrir issues, enviar PRs ou mandar um feedback.

---


Feito com ❤️ para estudar Golang e AWS Serverless.
