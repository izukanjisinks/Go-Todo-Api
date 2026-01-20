# Deployment Guide

This guide explains how to update the Todo API application on your server.

---

## 📋 Prerequisites

- Server with Docker and Docker Compose installed
- Git repository cloned on the server
- `.env` file configured with production credentials
- Access to the server via SSH

---

## 🚀 Standard Update Workflow

### 1. Pull Latest Changes from Git

```bash
git pull origin main
```

This fetches and merges the latest code changes from the main branch.

---

### 2. Rebuild API Container

```bash
docker compose build api
```

This rebuilds the Go application container with the updated code.

**Note:** Only rebuild what changed:
- If only Go code changed → `docker compose build api`
- If Dockerfile changed → `docker compose build api`
- If docker-compose.yml changed → No rebuild needed, just restart

---

### 3. Run Database Migrations (if needed)

```bash
docker compose up migrate
```

This applies any new database migrations.

**When to run:**
- ✅ When new migration files are added
- ✅ After pulling changes that include schema updates
- ❌ Not needed if only application code changed

---

### 4. Restart the API

```bash
docker compose up -d api
```

This restarts the API container with the new code.

The `-d` flag runs it in detached mode (background).

---

## 📝 Complete Update Script

Copy and paste this entire block:

```bash
# Navigate to project directory
cd ~/Go-Todo-Api

# Pull latest changes
git pull origin main

# Rebuild API container
docker compose build api

# Run migrations (if any)
docker compose up migrate

# Restart API
docker compose up -d api

# View logs to verify
docker compose logs -f api
```

---

## 🔍 Verification Steps

### Check Container Status

```bash
docker compose ps
```

Expected output:
```
NAME            STATUS          PORTS
todo-postgres   Up (healthy)    5432/tcp
todo-migrate    Exited (0)
todo-api        Up              0.0.0.0:8080->8080/tcp
```

### View API Logs

```bash
# Follow logs in real-time
docker compose logs -f api

# View last 50 lines
docker compose logs --tail=50 api
```

### Test API Health

```bash
curl http://localhost:8080/health
```

Expected response: `{"status": "healthy"}` or similar.

---

## 🔄 Quick Restart (No Code Changes)

If you just need to restart the API without updating code:

```bash
docker compose restart api
```

---

## 🛑 Full Restart (All Services)

To restart everything (PostgreSQL + API):

```bash
docker compose down
docker compose up -d
```

⚠️ **Warning:** This will briefly interrupt database connections.

---

## 🧹 Clean Rebuild (Fresh Start)

If you encounter issues, do a clean rebuild:

```bash
# Stop all containers
docker compose down

# Remove old images (optional)
docker rmi todo-api:latest

# Pull latest code
git pull origin main

# Rebuild and start everything
docker compose up -d --build

# View logs
docker compose logs -f
```

---

## 📊 Database Management

### View Database Size

```bash
docker exec -it todo-postgres psql -U postgres -d rbca_system -c "SELECT pg_size_pretty(pg_database_size('rbca_system'));"
```

### Backup Database

```bash
docker exec -it todo-postgres pg_dump -U postgres rbca_system > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Restore Database

```bash
cat backup_20240115_120000.sql | docker exec -i todo-postgres psql -U postgres -d rbca_system
```

---

## 🐛 Troubleshooting

### Migration Fails

If migration fails with "dirty" state:

```bash
# Force migration to a specific version
docker compose run --rm migrate -path /migrations \
  -database "postgres://postgres:your-password@postgres:5432/rbca_system?sslmode=disable" \
  force 18
```

Replace `18` with the last successful migration version.

### API Won't Start

Check logs for errors:

```bash
docker compose logs api
```

Common issues:
- Database connection failed → Check `.env` credentials
- Port already in use → Change `SERVER_PORT` in `.env`
- Migration not completed → Run `docker compose up migrate`

### Container Won't Stop

Force remove:

```bash
docker compose down --remove-orphans
docker rm -f todo-api
```

---

## 📦 Environment Variables

Ensure your `.env` file is configured:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-secure-password
DB_NAME=rbca_system
DB_SSLMODE=disable
SERVER_PORT=8080
```

**Never commit `.env` to Git!**

---

## 🔐 Security Best Practices

1. **Change default passwords** in production
2. **Use SSL/TLS** for database connections (`DB_SSLMODE=require`)
3. **Don't expose PostgreSQL port** to the internet
4. **Use secrets management** for sensitive data (e.g., Docker secrets, AWS Secrets Manager)
5. **Keep Docker images updated**: `docker compose pull`

---

## 📞 Support

If you encounter issues:

1. Check logs: `docker compose logs`
2. Verify environment: `docker compose config`
3. Check GitHub issues
4. Contact the development team

---

## 🎯 Quick Reference

| Command | Description |
|---------|-------------|
| `git pull origin main` | Get latest code |
| `docker compose build api` | Rebuild API |
| `docker compose up migrate` | Run migrations |
| `docker compose up -d api` | Start/restart API |
| `docker compose logs -f api` | View logs |
| `docker compose ps` | Check status |
| `docker compose down` | Stop all |
| `docker compose up -d` | Start all |

---

**Last Updated:** 2026-01-20
