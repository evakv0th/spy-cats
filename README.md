# Spy Cats API - Swagger Documentation

This project now includes complete Swagger UI documentation for easy API testing and exploration.

## 🚀 Getting Started

### Prerequisites

1. **Setup Environment Variables**
   
   Create a `.env` file in the project root directory:
   ```bash
   # Copy the example file
   cp .env.example .env
   ```

### Running the Server

1. Build and run the server:
   ```bash
   docker compose up 
   ```

2. The server will start on port 8080 with the following message:
   ```
   Server started on port 8080
   Swagger UI available at: http://localhost:8080/swagger/index.html
   ```


### Accessing Swagger UI

Open your browser and navigate to:
```
http://localhost:8080/swagger/index.html
```

## 📚 API Documentation

### Cats Endpoints

- **POST** `/api/cats` - Create a new spy cat
- **GET** `/api/cats` - List all spy cats
- **GET** `/api/cats/{id}` - Get a specific cat by ID
- **PATCH** `/api/cats/{id}/salary` - Update a cat's salary
- **DELETE** `/api/cats/{id}` - Delete a cat

### Missions Endpoints

- **POST** `/api/missions` - Create a new mission
- **GET** `/api/missions` - List all missions
- **GET** `/api/missions/{id}` - Get a specific mission by ID
- **PUT** `/api/missions/{id}/assign` - Assign a cat to a mission
- **DELETE** `/api/missions/{id}` - Delete a mission
- **PATCH** `/api/missions/{id}/complete` - Mark mission as complete

### Target Endpoints

- **POST** `/api/missions/{id}/targets` - Add a target to a mission
- **PATCH** `/api/missions/targets/{targetId}` - Update a target
- **DELETE** `/api/missions/targets/{targetId}` - Delete a target

## 🧪 Testing with Swagger UI

The Swagger UI provides:

1. **Interactive API Testing** - Click "Try it out" on any endpoint
2. **Request/Response Examples** - See example data for all models
3. **Schema Documentation** - Complete model definitions with validation rules
4. **Authentication Support** - Ready for future auth implementation

## 📋 Example Requests

### Create a Cat
```json
{
  "name": "Whiskers",
  "years_of_experience": 5,
  "breed": "Siamese",
  "salary": 50000.0
}
```

### Create a Mission
```json
{
  "cat_id": 5,
  "name": "Operation Stealth",
  "targets": [
    {
      "name": "Agent Smith",
      "country": "Russia",
      "notes": "High priority target",
      "is_complete": false
    }
  ],
  "is_complete": false
}
```

## 🔄 Updating Documentation

To regenerate Swagger docs after making changes:

```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go -o docs
```

## 📁 Generated Files

The Swagger generation creates:
- `docs/docs.go` - Go source file with embedded docs
- `docs/swagger.json` - JSON specification
- `docs/swagger.yaml` - YAML specification

These files are automatically imported and served by the application.