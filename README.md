## Integra Assignment API
API server for managing users. Usernames and emails are auto-generated. Usernames cannot be modified. All other fields can be modified.

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

### Endpoints
- get all users
```http
GET /users
```

- add user
```http
POST /users
Body:
{
  first_name: string,
  last_name: string,
  department: string,
}
```

- update user
```http
PUT /users/{id:int}
Body:
{
  id: int,
  first_name: string,
  last_name: string,
  email: string,
  user_status: string,
  department: string,
}
```

- delete user
```http
DELETE /users/{id:int}
```
