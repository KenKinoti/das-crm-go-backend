# Care CRM Backend

A comprehensive care management CRM system built with Go, Gin, and GORM.

## Features

- User authentication and authorization
- Participant management
- Staff scheduling
- Document management
- NDIS compliance tracking
- RESTful API design

## Quick Start

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Install dependencies: `go mod tidy`
4. Run the application: `make run`

## Development

- `make dev` - Run with hot reload (requires air)
- `make test` - Run tests
- `make build` - Build the application
- `make clean` - Clean build artifacts

## API Documentation

The API documentation is available at `/swagger/index.html` when running in development mode.

## Environment Variables

See `.env` file for configuration options.

## Database Support

- PostgreSQL (recommended for production)
- MySQL
- SQLite (for development/testing)

EOF

echo "✅ Project structure created successfully!"
echo ""
echo "📁 Project structure:"
tree -I 'bin|tmp|vendor'
echo ""
echo "🚀 Next steps:"
echo "1. Update the import paths in the files to match your actual module name"
echo "2. Configure your .env file with your database credentials"
echo "3. Install development tools: make install-dev"
echo "4. Run the project: make dev"
echo ""
echo "💡 Don't forget to:"
echo "- Replace 'github.com/yourusername/care-crm-backend' with your actual module path"
echo "- Set up your database (PostgreSQL recommended)"
echo "- Configure your environment variables in .env"