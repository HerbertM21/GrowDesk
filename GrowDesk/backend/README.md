# GrowDesk Go Backend

This is the Go backend application for GrowDesk, a ticket management system.

## Requirements

- Go 1.20 or higher
- Docker and Docker Compose (for containerized deployment)

## Getting Started

### Local Development

1. Clone this repository
2. Navigate to the project directory:
   ```
   cd GrowDesk-Go
   ```
3. Install dependencies:
   ```
   go mod download
   ```
4. Run the application:
   ```
   go run cmd/server/main.go
   ```

The server will start on port 8080 (or the port specified in the `PORT` environment variable).

### Using Docker

1. Build and run the container using Docker Compose:
   ```
   docker-compose up --build
   ```

2. To run in the background:
   ```
   docker-compose up -d
   ```

3. To stop the services:
   ```
   docker-compose down
   ```

## Environment Variables

- `PORT`: The port on which the server will run (default: 8080)
- `DATA_DIR`: Directory for storing data files (default: ./data)
- `MOCK_AUTH`: Enable mock authentication for development (default: true)

## API Endpoints

The API provides the following endpoints:

### Tickets
- `GET /api/tickets`: Get all tickets
- `GET /api/tickets/:id`: Get a specific ticket
- `POST /api/tickets`: Create a new ticket
- `PUT /api/tickets/:id`: Update a ticket
- `DELETE /api/tickets/:id`: Delete a ticket

### Categories
- `GET /api/categories`: Get all categories
- `GET /api/categories/:id`: Get a specific category
- `POST /api/categories`: Create a new category
- `PUT /api/categories/:id`: Update a category
- `DELETE /api/categories/:id`: Delete a category

### FAQs
- `GET /api/faqs`: Get all FAQs
- `GET /api/faqs/:id`: Get a specific FAQ
- `POST /api/faqs`: Create a new FAQ
- `PUT /api/faqs/:id`: Update a FAQ
- `DELETE /api/faqs/:id`: Delete a FAQ

## License

MIT 