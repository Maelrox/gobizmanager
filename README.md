# GoBizManager Backend

🚧 **Work in Progress** 🚧

A backend service for GoBizManager, a free business management platform built with Go.

## Overview

GoBizManager is a business management platform that helps companies manage their operations, users, and permissions. This repository contains the backend service built with Go.

## Current Status

The project is currently under active development. I am working on:

- 🔄 Refactoring database operations to the repository layer
- 🌐 Implementing internationalization support
- 🔒 Enhancing security and validation
- 📝 Improving error handling and messages
- 🧪 Adding comprehensive testing

## Features (In Progress)

- User authentication and authorization
- Company management
- Role-based access control (RBAC)
- Multi-language support
- Data encryption for sensitive information

## Tech Stack

- Go
- PostgreSQL
- Chi (HTTP router)
- JWT for authentication
- Validator for input validation

## Getting Started

### Prerequisites

- Go 1.21 or later
- PostgreSQL
- Make (optional)

### Environment Variables

Create a `.env` file with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db
ENCRYPTION_KEY=your_encryption_key
```

### Installation

1. Clone the repository
2. Install dependencies
3. Run the application

```bash
git clone https://github.com/yourusername/gobizmanager-backend.git
cd gobizmanager-backend
go mod download
go run cmd/api/main.go
```

## Contributing

We welcome contributions! Since this is a work in progress, please:

1. Check the current issues and pull requests
2. Create a new issue if you find a bug or want to suggest a feature
3. Follow the existing code style and patterns
4. Ensure all tests pass before submitting a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

⚠️ This is a work in progress. The codebase is actively being refactored and improved. Some features may not be fully implemented or may change significantly in future updates. 