# Medicine Reminder API

A RESTful API for managing medicine reminders. This API allows you to create, read, update, and delete medicine records with their dosage, frequency, and timing information.

## Features

- CRUD operations for medicine records
- Structured medicine information including dosage, frequency, and timing
- Input validation
- PostgreSQL database storage
- RESTful API design
- CORS support
- Comprehensive unit tests

## Prerequisites

- Go 1.16 or later
- PostgreSQL 12 or later
- Git

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd medicine-reminder
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Configure PostgreSQL:
   - Create a database named `medicine_reminder`
   - Update database configuration in `database/database.go` if needed

4. Run the application:
   ```bash
   go run main.go
   ```

The server will start on port 8080.

## Project Structure

```
medicine-reminder/
├── main.go                 # Application entry point
├── database/
│   └── database.go        # Database connection and initialization
├── handlers/
│   ├── medicine_handler.go      # HTTP handlers
│   └── medicine_handler_test.go # Unit tests
├── models/
│   └── medicine.go        # Data models
└── go.mod                 # Dependencies
```

## API Endpoints

### GET /api/medicines
Returns a list of all medicines.

Response:
```json
[
  {
    "id": 1,
    "name": "Paracetamol",
    "dosage": "500mg",
    "frequency": "3 times a day",
    "time_of_day": ["08:00", "14:00", "20:00"],
    "start_date": "2024-03-20T00:00:00Z",
    "end_date": "2024-04-20T00:00:00Z",
    "notes": "Take after meals",
    "created_at": "2024-03-20T15:55:14.721116Z",
    "updated_at": "2024-03-20T15:55:14.721116Z"
  }
]
```

### GET /api/medicines/{id}
Returns a specific medicine by ID.

Response:
```json
{
  "id": 1,
  "name": "Paracetamol",
  "dosage": "500mg",
  "frequency": "3 times a day",
  "time_of_day": ["08:00", "14:00", "20:00"],
  "start_date": "2024-03-20T00:00:00Z",
  "end_date": "2024-04-20T00:00:00Z",
  "notes": "Take after meals",
  "created_at": "2024-03-20T15:55:14.721116Z",
  "updated_at": "2024-03-20T15:55:14.721116Z"
}
```

### POST /api/medicines
Creates a new medicine record.

Request:
```json
{
  "name": "Paracetamol",
  "dosage": "500mg",
  "frequency": "3 times a day",
  "time_of_day": ["08:00", "14:00", "20:00"],
  "start_date": "2024-03-20T00:00:00Z",
  "end_date": "2024-04-20T00:00:00Z",
  "notes": "Take after meals"
}
```

Response: Returns the created medicine with status 201 Created.

### PUT /api/medicines/{id}
Updates an existing medicine record.

Request:
```json
{
  "name": "Paracetamol",
  "dosage": "1000mg",
  "frequency": "2 times a day",
  "time_of_day": ["08:00", "20:00"],
  "start_date": "2024-03-20T00:00:00Z",
  "end_date": "2024-04-20T00:00:00Z",
  "notes": "Take after meals, updated dosage"
}
```

Response: Returns the updated medicine with status 200 OK.

### DELETE /api/medicines/{id}
Deletes a medicine record.

Response: Returns status 204 No Content on success.

## Testing

Run the unit tests:
```bash
go test ./handlers -v
```

## cURL Examples

1. Get All Medicines:
```bash
curl -X GET http://localhost:8080/api/medicines
```

2. Get Single Medicine:
```bash
curl -X GET http://localhost:8080/api/medicines/1
```

3. Create Medicine:
```bash
curl -X POST http://localhost:8080/api/medicines \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Aspirin",
    "dosage": "100mg",
    "frequency": "Once daily",
    "time_of_day": ["08:00"],
    "start_date": "2024-03-20T00:00:00Z",
    "end_date": "2024-04-20T00:00:00Z",
    "notes": "Take with food"
  }'
```

4. Update Medicine:
```bash
curl -X PUT http://localhost:8080/api/medicines/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Aspirin",
    "dosage": "200mg",
    "frequency": "Twice daily",
    "time_of_day": ["08:00", "20:00"],
    "start_date": "2024-03-20T00:00:00Z",
    "end_date": "2024-04-20T00:00:00Z",
    "notes": "Take with food, updated dosage"
  }'
```

5. Delete Medicine:
```bash
curl -X DELETE http://localhost:8080/api/medicines/1
```

## License

This project is licensed under the MIT License. 