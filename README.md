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