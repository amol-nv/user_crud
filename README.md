# amol-nv/user_crud

This repository contains a Go implementation of a Product CRUD service.

## Endpoints (REST)
- `POST /products` - create a product
- `GET /products/{id}` - get a product by id
- `PUT /products/{id}` - update a product
- `DELETE /products/{id}` - delete a product

## Run

```bash
go run ./cmd/server
```

## Test

```bash
go test ./...
```
