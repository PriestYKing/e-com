# e-com

An e-commerce platform built with a full-stack architecture.

## Features

- User management (create, read, delete users)
- RESTful API with Go backend
- PostgreSQL database integration
- Modular code structure

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Node.js & npm (for frontend, if applicable)

### Backend Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/PriestYKing/e-com.git
   cd e-com/server
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure your database in `.env` or config file.

4. Run the server:
   ```bash
   go run main.go
   ```

### Frontend Setup

_(If you have a frontend folder)_

1. Navigate to the frontend directory:

   ```bash
   cd ../client
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

## API Endpoints

- `GET /users` - List all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create a new user
- `DELETE /users/{id}` - Delete a user
- `GET /products` - Get all products

## License

MIT

## Author

PriestYKing
