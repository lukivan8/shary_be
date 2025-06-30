# Shary BE - Renting Platform

A Go REST API for a renting platform with a clean architecture using:
- `net/http` for HTTP server
- `sqlx` with PostgreSQL for database operations
- `zap` for structured logging
- `go-playground/validator/v10` for request validation
- `golang-migrate/migrate` for database migrations

## Project Structure

```
shary_be/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── README.md              # This file
├── internal/              # Internal application code
│   ├── config/           # Configuration management
│   ├── database/         # Database utilities and migrations
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models and validation
│   ├── repository/       # Database operations
│   └── service/          # Business logic
├── migrations/           # Database migrations (golang-migrate format)
│   ├── 000001_create_items_table.up.sql
│   ├── 000001_create_items_table.down.sql
│   ├── 000002_insert_sample_data.up.sql
│   └── 000002_insert_sample_data.down.sql
└── scripts/             # Utility scripts
    └── setup_db.sh      # Database setup script
```

## Architecture

This project follows a clean architecture pattern:

- **Handlers**: Handle HTTP requests and responses
- **Services**: Contain business logic
- **Repositories**: Handle database operations
- **Models**: Define data structures and validation
- **Middleware**: Provide cross-cutting concerns like logging and recovery

## Data Model

### Item Entity
The core entity represents items available for rent:

```go
type Item struct {
    ID          int       `json:"id" db:"id"`
    Title       string    `json:"title" db:"title" validate:"required,min=1,max=200"`
    Description string    `json:"description" db:"description" validate:"required,min=10,max=2000"`
    ImageURL    string    `json:"image_url" db:"image_url" validate:"required,url"`
    PricePerDay float64   `json:"price_per_day" db:"price_per_day" validate:"required,min=0.01"`
    Location    string    `json:"location" db:"location" validate:"required,min=1,max=500"`
    IsAvailable bool      `json:"is_available" db:"is_available"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

## Prerequisites

- Go 1.21 or later
- PostgreSQL database

## Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd shary_be
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Run the setup script (it will install golang-migrate if needed)
   ./scripts/setup_db.sh
   ```

4. **Set environment variables** (optional)
   ```bash
   export DATABASE_URL="postgres://username:password@localhost:5432/shary_be?sslmode=disable"
   export PORT=4000
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on port 4000 (or the port specified in the PORT environment variable).

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Items
- `GET /api/items` - Get all items (with filtering)
- `POST /api/items` - Create a new item
- `GET /api/items/{id}` - Get item by ID
- `PUT /api/items/{id}` - Update item
- `DELETE /api/items/{id}` - Delete item
- `GET /api/items/location/{location}` - Get items by location

### Query Parameters for Filtering
- `min_price` - Minimum price per day
- `max_price` - Maximum price per day
- `location` - Filter by location (partial match)
- `available` - Filter by availability (true/false)
- `search` - Search in title and description
- `limit` - Number of items to return (default: 20)
- `offset` - Number of items to skip (default: 0)

## Example Usage

### Create an item
```bash
curl -X POST http://localhost:4000/api/items \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Mountain Bike - Trek Marlin 7",
    "description": "High-quality mountain bike perfect for trail riding. Features 21-speed Shimano drivetrain.",
    "image_url": "https://example.com/images/mountain-bike.jpg",
    "price_per_day": 25.00,
    "location": "San Francisco, CA"
  }'
```

### Get all items
```bash
curl http://localhost:4000/api/items
```

### Get items with filters
```bash
# Get items under $50 per day in San Francisco
curl "http://localhost:4000/api/items?max_price=50&location=San%20Francisco"

# Search for bikes
curl "http://localhost:4000/api/items?search=bike"

# Get only available items
curl "http://localhost:4000/api/items?available=true"

# Pagination
curl "http://localhost:4000/api/items?limit=5&offset=10"
```

### Get item by ID
```bash
curl http://localhost:4000/api/items/1
```

### Get items by location
```bash
curl http://localhost:4000/api/items/location/San%20Francisco
```

### Update an item
```bash
curl -X PUT http://localhost:4000/api/items/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price_per_day": 30.00,
    "is_available": false
  }'
```

### Delete an item
```bash
curl -X DELETE http://localhost:4000/api/items/1
```

## Database Migrations

This project uses `golang-migrate/migrate` for database migrations:

### Running migrations manually
```bash
# Run all migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/shary_be?sslmode=disable" up

# Rollback last migration
migrate -path migrations -database "postgres://user:pass@localhost:5432/shary_be?sslmode=disable" down 1

# Check migration status
migrate -path migrations -database "postgres://user:pass@localhost:5432/shary_be?sslmode=disable" version
```

### Creating new migrations
```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq add_new_table

# This creates:
# migrations/000003_add_new_table.up.sql
# migrations/000003_add_new_table.down.sql
```

## Configuration

The application can be configured using environment variables:

- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://postgres:password@localhost:5432/shary_be?sslmode=disable`)
- `PORT`: Server port (default: `4000`)

## Development

### Running tests
```bash
go test ./...
```

### Building
```bash
go build -o shary_be main.go
```

## How to Add New Requests (For New Go Developers)

As a new Go developer, here's how to add new API endpoints to this project:

### 1. **Define the Data Model** (if needed)
Create or update models in `internal/models/`:

```go
// internal/models/rental.go
type Rental struct {
    ID        int       `json:"id" db:"id"`
    ItemID    int       `json:"item_id" db:"item_id"`
    UserID    int       `json:"user_id" db:"user_id"`
    StartDate time.Time `json:"start_date" db:"start_date"`
    EndDate   time.Time `json:"end_date" db:"end_date"`
    Status    string    `json:"status" db:"status"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### 2. **Create Repository Methods**
Add database operations in `internal/repository/`:

```go
// internal/repository/rental.go
func (r *RentalRepository) Create(rental *models.Rental) error {
    query := `INSERT INTO rentals (item_id, user_id, start_date, end_date, status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
    // Implementation...
}

func (r *RentalRepository) GetByID(id int) (*models.Rental, error) {
    // Implementation...
}
```

### 3. **Add Service Layer Logic**
Implement business logic in `internal/service/`:

```go
// internal/service/rental.go
func (s *RentalService) CreateRental(req *models.CreateRentalRequest) (*models.Rental, error) {
    // Validate request
    if err := req.Validate(); err != nil {
        return nil, err
    }
    
    // Business logic (check availability, calculate price, etc.)
    // Call repository
    // Return result
}
```

### 4. **Create HTTP Handler**
Add request handling in `internal/handlers/`:

```go
// internal/handlers/rental.go
func (h *RentalHandler) HandleRentals(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.getAllRentals(w, r)
    case http.MethodPost:
        h.createRental(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *RentalHandler) createRental(w http.ResponseWriter, r *http.Request) {
    var req models.CreateRentalRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    rental, err := h.rentalService.CreateRental(&req)
    if err != nil {
        // Handle error
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(rental)
}
```

### 5. **Register Routes in main.go**
Add the new routes to the main function:

```go
// main.go
func main() {
    // ... existing setup ...
    
    // Initialize new components
    rentalRepo := repository.NewRentalRepository(db)
    rentalService := service.NewRentalService(rentalRepo, logger)
    rentalHandler := handlers.NewRentalHandler(rentalService, logger)
    
    // Add routes
    mux.HandleFunc("/api/rentals", rentalHandler.HandleRentals)
    mux.HandleFunc("/api/rentals/", rentalHandler.HandleRentalByID)
}
```

### 6. **Create Database Migration**
Add migration files:

```sql
-- migrations/000003_create_rentals_table.up.sql
CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    item_id INTEGER REFERENCES items(id),
    user_id INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 7. **Run Migration**
```bash
migrate -path migrations -database "your_db_url" up
```

### Key Principles:
- **Separation of Concerns**: Each layer has a specific responsibility
- **Dependency Injection**: Services depend on repositories, handlers depend on services
- **Error Handling**: Always handle errors appropriately at each layer
- **Validation**: Validate input data using struct tags and custom validation
- **Logging**: Log important events and errors for debugging

## Features

- ✅ Clean architecture with separation of concerns
- ✅ Structured logging with Zap
- ✅ Request validation with go-playground/validator
- ✅ Database operations with sqlx and PostgreSQL
- ✅ Database migrations with golang-migrate
- ✅ Graceful shutdown handling
- ✅ Middleware for logging and panic recovery
- ✅ RESTful API design with filtering and pagination
- ✅ Proper error handling and HTTP status codes
- ✅ Sample data for testing

## Next Steps

Consider adding these features for a production application:

- Authentication and authorization (JWT, OAuth)
- User management system
- Rental booking system
- Payment processing
- Image upload and storage
- Email notifications
- Rate limiting
- CORS middleware
- API documentation (Swagger/OpenAPI)
- Unit and integration tests
- Docker configuration
- CI/CD pipeline
- Environment-specific configurations
- Database connection pooling
- Metrics and monitoring 