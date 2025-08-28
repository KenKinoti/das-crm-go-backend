# GoFiber AGO CRM Backend

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.9+-00ADD8?style=flat&logo=gin)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-v1.25+-00ADD8?style=flat)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/kenkinoti/gofiber-ago-crm-backend)

A comprehensive **NDIS Care Management CRM System** built with Go, Gin Framework, and GORM. Designed specifically for Australian Disability Services providers to manage participants, staff, shifts, documents, and care plans while ensuring NDIS compliance.

## 🌟 Features

### Core Functionality
- 🔐 **User Authentication & Authorization** - JWT-based auth with role-based access control
- 👥 **Participant Management** - Complete participant profiles with medical info and funding tracking
- 📅 **Staff Scheduling** - Advanced shift management with conflict detection and status tracking
- 📄 **Document Management** - Secure file upload/download with categorization and expiry tracking
- 🏥 **Care Plans** - Comprehensive care planning with approval workflows
- 🆘 **Emergency Contacts** - Priority-based emergency contact management
- 🏢 **Organization Management** - Multi-tenant architecture with data isolation
- 📊 **NDIS Compliance** - Built-in NDIS registration and funding tracking

### Technical Features
- 🚀 **RESTful API** - Clean, standard-compliant REST endpoints
- 🔍 **Advanced Filtering** - Comprehensive search and filter capabilities
- 📄 **Pagination** - Efficient data loading with page-based navigation
- 🔒 **Data Security** - Organization-based data isolation and role permissions
- 📝 **Input Validation** - Comprehensive request validation and sanitization
- 🗄️ **Database Optimization** - Optimized queries with proper indexing
- 📤 **File Handling** - Secure file upload with validation and download

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │◄──►│   REST API       │◄──►│   Database      │
│   (React/Vue)   │    │   (Gin/Go)      │    │   (PostgreSQL)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │   File Storage   │
                       │   (Local/Cloud)  │
                       └──────────────────┘
```

## 📋 API Endpoints

### 🔐 Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `POST /api/v1/auth/logout` - User logout

### 👤 Users Management
- `GET /api/v1/users` - List all users (with filtering & pagination)
- `GET /api/v1/users/me` - Get current user profile
- `POST /api/v1/users` - Create new user (admin only)
- `PUT /api/v1/users/:id` - Update user details
- `DELETE /api/v1/users/:id` - Delete user (admin only)

### 👥 Participants Management
- `GET /api/v1/participants` - List participants (with filtering & search)
- `GET /api/v1/participants/:id` - Get participant details
- `POST /api/v1/participants` - Create new participant
- `PUT /api/v1/participants/:id` - Update participant
- `DELETE /api/v1/participants/:id` - Delete participant

### 📅 Shifts Management
- `GET /api/v1/shifts` - List shifts (with date range & status filtering)
- `GET /api/v1/shifts/:id` - Get shift details
- `POST /api/v1/shifts` - Create new shift
- `PUT /api/v1/shifts/:id` - Update shift
- `PATCH /api/v1/shifts/:id/status` - Update shift status
- `DELETE /api/v1/shifts/:id` - Delete shift

### 📄 Documents Management
- `GET /api/v1/documents` - List documents (with category filtering)
- `GET /api/v1/documents/:id` - Get document metadata
- `POST /api/v1/documents` - Upload new document
- `PUT /api/v1/documents/:id` - Update document metadata
- `GET /api/v1/documents/:id/download` - Download document file
- `DELETE /api/v1/documents/:id` - Delete document

### 🆘 Emergency Contacts
- `GET /api/v1/emergency-contacts` - List emergency contacts
- `GET /api/v1/emergency-contacts/:id` - Get contact details
- `POST /api/v1/emergency-contacts` - Create new contact
- `PUT /api/v1/emergency-contacts/:id` - Update contact
- `DELETE /api/v1/emergency-contacts/:id` - Delete contact

### 🏥 Care Plans
- `GET /api/v1/care-plans` - List care plans (with status filtering)
- `GET /api/v1/care-plans/:id` - Get care plan details
- `POST /api/v1/care-plans` - Create new care plan
- `PUT /api/v1/care-plans/:id` - Update care plan
- `PATCH /api/v1/care-plans/:id/approve` - Approve/reject care plan
- `DELETE /api/v1/care-plans/:id` - Delete care plan

### 🏢 Organization
- `GET /api/v1/organization` - Get organization details
- `PUT /api/v1/organization` - Update organization (admin only)

## 🚀 Quick Start

### Prerequisites
- [Go 1.21+](https://golang.org/doc/install)
- [PostgreSQL 13+](https://www.postgresql.org/download/) (or MySQL/SQLite for development)
- [Air](https://github.com/cosmtrek/air) (optional, for hot reload)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/kenkinoti/gofiber-ago-crm-backend.git
   cd gofiber-ago-crm-backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database and JWT configurations
   ```

4. **Run database migrations**
   ```bash
   go run cmd/migrate/main.go
   ```

5. **Start the server**
   ```bash
   make run
   # Or with hot reload: make dev
   ```

The API will be available at `http://localhost:8080`

## 🛠️ Development

### Available Commands
```bash
make dev          # Run with hot reload (requires air)
make run          # Run the application
make build        # Build the application
make test         # Run tests
make test-cover   # Run tests with coverage
make clean        # Clean build artifacts
make lint         # Run linting
make docs         # Generate API documentation
```

### Project Structure
```
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   └── migrate/           # Database migration tool
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── database/         # Database connection and setup
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/          # Data models and database schemas
│   └── services/        # Business logic services
├── pkg/                  # Public packages
│   ├── utils/           # Utility functions
│   └── validators/      # Custom validators
├── migrations/          # Database migration files
├── uploads/            # File upload directory
├── docs/               # Documentation files
├── scripts/            # Build and deployment scripts
└── web/               # Static web assets
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file with the following variables:

```env
# Server Configuration
PORT=8080
ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=care_crm_db
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

# File Upload Configuration
UPLOAD_MAX_SIZE=10485760  # 10MB in bytes
UPLOAD_PATH=./uploads

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Database Support

- **PostgreSQL** (recommended for production)
- **MySQL** (alternative production option)
- **SQLite** (development/testing only)

## 📊 Database Schema

The system includes the following main entities:

- **Organizations** - Care provider organizations
- **Users** - System users (admin, manager, care_worker, support_coordinator)
- **Participants** - Care recipients with NDIS details
- **Emergency Contacts** - Participant emergency contacts
- **Shifts** - Scheduled care shifts with status tracking
- **Documents** - File attachments with metadata
- **Care Plans** - Participant care planning with approval workflow
- **Refresh Tokens** - JWT refresh token management
- **User Permissions** - Granular permission system

## 🔒 Security Features

- **JWT Authentication** - Secure token-based authentication
- **Role-Based Access Control** - Admin, Manager, Care Worker, Support Coordinator roles
- **Organization Data Isolation** - Multi-tenant architecture
- **Input Validation** - Comprehensive request validation
- **File Upload Security** - File type and size validation
- **SQL Injection Protection** - GORM ORM with prepared statements
- **Password Security** - bcrypt hashing with salt

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover

# Run specific test package
go test ./internal/handlers -v

# Run tests with race detection
go test -race ./...
```

## 📋 API Documentation

API documentation is automatically generated and available at:
- Development: `http://localhost:8080/swagger/index.html`
- Production: `/swagger/index.html`

## 📝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Run `golint` before submitting
- Write tests for new features
- Update documentation as needed

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Support

- 📧 Email: [kennedy@dasyin.com.au](mailto:kennedy@dasyin.com.au)
- 🐛 Issues: [GitHub Issues](https://github.com/kenkinoti/gofiber-ago-crm-backend/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/kenkinoti/gofiber-ago-crm-backend/discussions)

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/) - HTTP web framework
- [GORM](https://gorm.io/) - Go ORM library
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Testify](https://github.com/stretchr/testify) - Testing toolkit

## 📊 Project Status

- ✅ **Authentication System** - Complete
- ✅ **User Management** - Complete  
- ✅ **Participant Management** - Complete
- ✅ **Shift Management** - Complete
- ✅ **Document Management** - Complete
- ✅ **Emergency Contacts** - Complete
- ✅ **Care Plans** - Complete
- ✅ **Organization Management** - Complete
- 🔄 **API Documentation** - In Progress
- 🔄 **Unit Testing** - In Progress
- 📋 **Frontend Integration** - Planned
- 📋 **Advanced Reporting** - Planned

---

**Built with ❤️ by the DASYIN Team**

For more information, visit our [website](https://dasyin.com.au) or check out our [documentation](https://docs.dasyin.com.au).