# Ledger System - API Financeira de Alta Concorrência

> 🏆 **Projeto Focado em Engenharia de Software, Consistência ACID e Sistemas Distribuídos.**

## 📖 Sobre
Este projeto é uma implementação de referência de um **Sistema de Ledger (Contabilidade)** capaz de processar transações financeiras com **Zero Race Conditions**.

O objetivo principal não é apenas mover dinheiro, mas garantir que o dinheiro nunca seja criado ou destruído, mesmo quando milhares de requisições tentam acessar a mesma conta simultaneamente. Para isso, utilizamos estratégias avançadas de Locking no banco de dados.

## 🏗️ Arquitetura & Tech Stack

O projeto segue os princípios da **Clean Architecture** para garantir desacoplamento e testabilidade.

-   **Linguagem:** Go 1.22+ (Standard Lib + Chi router).
-   **Banco de Dados:** PostgreSQL 15.
-   **Infraestrutura:** Docker & Docker Compose.
-   **Design Pattern:** Repository Pattern, Dependency Injection.

---

## 🛡️ O Diferencial: Estratégia de Concorrência & Segurança

A parte crítica de qualquer sistema financeiro é evitar **Race Conditions** e **Deadlocks**.

### 1. Pessimistic Locking (ACID)
Utilizamos `SELECT ... FOR UPDATE` dentro de uma transação SQL. Isso garante que, ao ler o saldo de uma conta, nenhuma outra transação possa modificá-la até que a transação atual termine.

### 2. Prevenção de Deadlocks
Para evitar que a Transação A trave a conta 1 esperando a 2, e a Transação B trave a conta 2 esperando a 1, implementamos uma **ordenação determinística de IDs**. Sempre bloqueamos o ID menor primeiro.

### 🔬 Prova de Implementação (Snippet)

```go
// internal/repository/postgres/account.go

// 1. Prevenção de Deadlocks: Ordenação Determinística
firstID, secondID := transfer.FromAccountID, transfer.ToAccountID
if firstID > secondID {
    firstID, secondID = secondID, firstID
}

// 2. Lock Pessimista (SELECT ... FOR UPDATE)
if err := r.lockAccount(ctx, tx, firstID); err != nil { return err }
if err := r.lockAccount(ctx, tx, secondID); err != nil { return err }

// ... Helper function
func (r *AccountRepository) lockAccount(ctx, tx, id) error {
    // "FOR UPDATE" trava a linha no Postgres
    return tx.QueryRowContext(ctx, "SELECT id FROM accounts WHERE id = $1 FOR UPDATE", id).Scan(&dummy)
}
```

---

## 🚀 Como Rodar

### Pré-requisitos
-   Docker & Docker Compose

### 1. Iniciar Aplicação
Suba o banco de dados e a API com um único comando:
```bash
make up
# ou
docker-compose up -d --build
```
A API estará disponível em `http://localhost:8080`.

### 2. Rodar Testes de Estresse (Prova de Fogo)
Execute o script que simula **50 workers simultâneos** realizando transferências:
```bash
go run cmd/stress_test/main.go
```
✅ **Resultado Esperado:** O teste deve finalizar confirmando que o saldo total do sistema permanece inalterado (Consistência Contábil).

### 3. Rodar Testes de Ponta a Ponta (Edge Cases)
Valide se a API rejeita transações inválidas (ex: saldo negativo):
```bash
go run cmd/e2e_test/main.go
```

---

## ✅ Resultados Alcançados
-   [x] **Robustez:** 50 Workers concorrentes testados sem falhas.
-   [x] **Segurança:** Tratamento de erros HTTP 400/404/422 implementado.
-   [x] **Qualidade:** Código validado com `go vet` e testes automatizados.

---
Desenvolvido por Gabriel Cau.
