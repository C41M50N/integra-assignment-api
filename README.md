## Integra Assignment API
API server for managing users.

### 1. Setup database
- start/create postgres service
```bash
docker-compose -p integra -f docker-compose.yml up -d
```

- setup users table
```bash
psql -h localhost -U user -d postgres < setup.sql
```

### 2. Start API server
```bash
go run main.go
```
