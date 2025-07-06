# GoPass - Proper Go Project Structure

```
gopass/
├── cmd/
│   └── gopass/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── crypto/
│   │   └── crypto.go              # Encryption/decryption utilities
│   ├── database/
│   │   └── database.go            # Database operations
│   ├── models/
│   │   └── models.go              # Data structures
│   ├── handlers/
│   │   ├── auth.go                # Authentication handlers
│   │   ├── dashboard.go           # Dashboard handlers
│   │   ├── entries.go             # Password entry handlers
│   │   └── api.go                 # API handlers
│   ├── middleware/
│   │   └── auth.go                # Authentication middleware
│   ├── service/
│   │   └── password.go            # Business logic
│   └── templates/
│       ├── login.html             # Login template
│       ├── dashboard.html         # Dashboard template
│       ├── form.html              # Add/Edit form template
│       └── base.html              # Base template
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css          # CSS styles
│   │   └── js/
│   │       └── app.js             # JavaScript functionality
│   └── templates/                 # Alternative template location
├── pkg/
│   └── utils/
│       └── generator.go           # Password generator utility
├── scripts/
│   ├── build.sh                   # Build script
│   └── install.sh                 # Installation script
├── docs/
│   ├── README.md                  # Project documentation
│   └── SECURITY.md               # Security considerations
├── .gitignore                     # Git ignore file
├── go.mod                         # Go module file
├── go.sum                         # Go module checksums
├── Makefile                       # Build automation
└── LICENSE                        # License file
```

## Key Benefits of This Structure:

### 1. **Separation of Concerns**
- **cmd/**: Application entry points (main packages)
- **internal/**: Private application code that can't be imported by other projects
- **pkg/**: Public library code that can be imported by other projects
- **web/**: Web assets (templates, static files)

### 2. **Scalability**
- Each component has its own package
- Easy to add new features without cluttering
- Clear boundaries between layers

### 3. **Testability**
- Each package can be tested independently
- Business logic separated from HTTP handlers
- Database operations isolated

### 4. **Security**
- Crypto operations in dedicated package
- Authentication middleware separate from handlers
- Configuration management centralized

## Example Package Contents:

### internal/models/models.go
```go
package models

import "time"

type PasswordEntry struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Username  string    `json:"username"`
    Password  string    `json:"password"`
    Website   string    `json:"website"`
    Notes     string    `json:"notes"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Database struct {
    Entries []PasswordEntry `json:"entries"`
}

type Session struct {
    ID        string
    CreatedAt time.Time
}
```

### internal/config/config.go
```go
package config

import (
    "os"
    "path/filepath"
)

type Config struct {
    DataDir  string
    DataFile string
    Port     string
}

func New() *Config {
    homeDir, _ := os.UserHomeDir()
    dataDir := filepath.Join(homeDir, ".gopass")
    
    return &Config{
        DataDir:  dataDir,
        DataFile: filepath.Join(dataDir, "passwords.enc"),
        Port:     getEnv("PORT", "8080"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### internal/service/password.go
```go
package service

import (
    "gopass/internal/crypto"
    "gopass/internal/database"
    "gopass/internal/models"
)

type PasswordService struct {
    db     *database.Database
    crypto *crypto.Manager
}

func NewPasswordService(db *database.Database, crypto *crypto.Manager) *PasswordService {
    return &PasswordService{
        db:     db,
        crypto: crypto,
    }
}

func (s *PasswordService) AddEntry(entry models.PasswordEntry) error {
    // Business logic here
    return s.db.SaveEntry(entry)
}

// ... other business logic methods
```

### Makefile
```makefile
.PHONY: build run test clean install

BINARY_NAME=gopass
BUILD_DIR=build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/gopass/main.go

run:
	go run cmd/gopass/main.go

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

dev:
	go run cmd/gopass/main.go

fmt:
	go fmt ./...

lint:
	golint ./...

deps:
	go mod tidy
	go mod download
```

## Additional Recommendations:

1. **Use dependency injection** for better testability
2. **Add configuration validation** in the config package
3. **Implement proper logging** with structured logging (e.g., logrus, zap)
4. **Add graceful shutdown** handling
5. **Use context** for request handling and cancellation
6. **Add health check endpoints**
7. **Implement rate limiting** for security
8. **Add comprehensive tests** for each package

This structure makes your codebase more maintainable, testable, and follows Go community standards.
