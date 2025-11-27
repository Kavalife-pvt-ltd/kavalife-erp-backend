# Kava Life ERP - Backend

Backend server for the Kava Life ERP system.

## ✨ Tech Stack

- Language: Golang
- Initial Setup: Vanilla Go HTTP server
- Framework: Will add Gin later
- Database: Will connect to Supabase PostgreSQL later

## 🚀 How to Run

```bash
To install:
    go mod tidy
To run code:
    go run cmd/server/main.go
```
## 🚀 How to Create Docker Image
#### Docker Build & Docker run
    docker build -t kavalife-erp-backend .  
    docker run -d  --name kavalife-erp-backend --env-file .env -p 8080:8080 kavalife-erp-backend
```
kavalife-erp-backend
├─ .air.toml
├─ .dockerignore
├─ README.md
├─ cmd
│  └─ server
│     └─ main.go
├─ config
│  └─ config.go
├─ dockerfile
├─ docs
├─ go.mod
├─ go.sum
├─ internal
│  ├─ db
│  │  └─ connection.go
│  ├─ handlers
│  │  ├─ grn.go
│  │  ├─ notification.go
│  │  ├─ products.go
│  │  ├─ qaqc.go
│  │  ├─ sales_po.go
│  │  ├─ users.go
│  │  ├─ vendors.go
│  │  └─ vir.go
│  ├─ models
│  │  ├─ grn.go
│  │  ├─ notification.go
│  │  ├─ product.go
│  │  ├─ qaqc.go
│  │  ├─ sales_po.go
│  │  ├─ sales_po_status.go
│  │  ├─ users.go
│  │  ├─ vendors.go
│  │  └─ vir.go
│  ├─ routes
│  │  ├─ middleware.go
│  │  └─ router.go
│  ├─ services
│  │  ├─ grn.go
│  │  ├─ products.go
│  │  ├─ qaqc.go
│  │  ├─ sales_po.go
│  │  ├─ users.go
│  │  ├─ vendors.go
│  │  └─ vir.go
│  └─ utils
│     ├─ context.go
│     ├─ customRes.go
│     ├─ email.go
│     ├─ email_template.go
│     ├─ hashEncrypt.go
│     ├─ jwt.go
│     ├─ logger.go
│     └─ utils.go
├─ makefile
├─ migrations
└─ tmp

```