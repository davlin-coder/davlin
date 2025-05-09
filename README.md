# Davlin

Davlin is a modern Go web application that aims to create an intelligent chat application based on deepresearch technology. It integrates with document analysis capabilities to answer user questions by leveraging relevant context from their documents. The project is currently under active development.

## Features

- **User Management**
  - Registration with email verification
  - Login with password or verification code
  - JWT-based authentication

- **Document-Assisted Chat System**
  - Document upload and analysis
  - Contextual message handling based on document content
  - Chat history with document references

- **Integrations**
  - Deep research-based LLM integration
  - Document parsing and indexing
  - Intelligent search capabilities

- **Infrastructure**
  - MySQL database for persistent storage
  - Redis for caching and session management
  - Email service for notifications

## Project Structure

```
davlin/
├── cmd/                 # Application entry points
├── internal/            # Internal application code
│   ├── config/          # Configuration handling
│   ├── controller/      # API controllers
│   ├── middleware/      # HTTP middleware
│   ├── model/           # Data models
│   ├── resource/        # External resources
│   ├── router/          # API routing
│   └── service/         # Business logic
├── docs/                # Documentation
│   └── swagger.yaml     # API specification
└── scripts/             # Utility scripts
```

## Getting Started

### Prerequisites

- Go 1.23.5 or higher
- Docker and Docker Compose
- MySQL
- Redis

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/davlin-coder/davlin.git
   cd davlin
   ```

2. Create a configuration file:
   ```
   cp config.template.yaml config.yaml
   ```

3. Edit `config.yaml` to set up your environment-specific configurations.

4. Start the database services:
   ```
   docker-compose up -d
   ```

5. Run the application:
   ```
   go run cmd/main.go
   ```

### Building and Running with Docker

```
docker build -t davlin .
docker run -p 8080:8080 davlin
```

## Current Development Status

This project is in active development. The following components are being implemented:

- Document parsing and indexing system
- Deep research-based LLM integration
- Document-context aware chat functionality
- User interface for document upload and management

## Development

### Dependency Management

```
go mod tidy
```

### Testing

```
go test ./...
```

### Configuration

The application uses Viper for configuration management, supporting both YAML files and environment variables.

## Vendors

### 七牛云
- 提供商网站: [七牛云](https://www.qiniu.com/)
- 为本项目提供对象存储、CDN加速及文档处理服务

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

