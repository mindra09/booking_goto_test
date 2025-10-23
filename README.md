# Go PostgreSQL Docker Application

A complete Go application with PostgreSQL database, Docker containerization, and hot reload for development.

## 🚀 Quick Start

Follow these 3 simple steps to get started:

### 1. Clone the Repository
```bash
git clone <your-repository-url>
```

### 2. Go to Project Directory
```bash
cd go-postgres-app
```

### 3. Start with Docker Compose
```bash
docker-compose up
```

**Done!** Your application is now running with:
- ✅ Go API server with hot reload
- ✅ PostgreSQL database
- ✅ Ready-to-use REST API

## 🌐 Access Your Application

Once running, you can access:

- **API Server**: http://localhost:8080
- **PostgreSQL Database**: localhost:5432

## 🛠️ What You Get

### Built-in Features
- **Hot Reload**: Automatic restart on code changes
- **RESTful API**: Complete CRUD operations
- **Data Validation**: Input validation with Ozzo
- **Structured Logging**: Beautiful logs with Logrus
- **Error Handling**: Comprehensive error management

### Default API Endpoints
```http
GET    /user           # Get all users
POST   /user           # Create new user
GET    /user/{id}      # Get user by ID
PUT    /user/{id}      # Update user
DELETE /user/{id}      # Delete user
DELETE /user/{id}/family/{family_id}  # Delete user family
```



## 🧪 Test Your API

### Create a new user
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "dob": "1990-05-15",
    "email": "john@example.com",
    "nationality_id": 1
  }'
```

### Get all users
```bash
curl http://localhost:8080/users
```



