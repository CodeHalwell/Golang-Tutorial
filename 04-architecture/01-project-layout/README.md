# Chapter 1: Standard Go Project Layout

## Learning Objectives

1. Understand the standard Go project structure
2. Learn the purpose of each directory
3. Understand package organization
4. Know when to create new packages
5. Learn naming conventions
6. Understand the internal package
7. Organize projects for scalability

---

## Standard Project Layout

The official Go project layout is documented in [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

```
myproject/
├── cmd/                      # Command-line applications
│   ├── myapp/
│   │   └── main.go          # Entry point
│   └── mycli/
│       └── main.go
│
├── internal/                 # Private application code
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   └── storage/
│       ├── storage.go
│       ├── storage_test.go
│       └── sqlite.go
│
├── pkg/                      # Public libraries (can be imported by others)
│   ├── cache/
│   │   ├── cache.go
│   │   ├── cache_test.go
│   │   └── types.go
│   └── client/
│       ├── client.go
│       ├── client_test.go
│       └── options.go
│
├── api/                      # API definitions
│   └── v1/
│       └── service.proto    # gRPC proto definitions
│
├── config/                   # Configuration files
│   ├── config.yaml
│   └── database.yaml
│
├── docs/                     # Documentation
│   ├── architecture.md
│   ├── api.md
│   └── contributing.md
│
├── examples/                 # Example usage
│   ├── basic_example.go
│   └── advanced_example.go
│
├── scripts/                  # Build/automation scripts
│   ├── build.sh
│   ├── test.sh
│   └── deploy.sh
│
├── testdata/                 # Test data files
│   ├── sample.json
│   └── sample.csv
│
├── third_party/              # Third-party code (vendor)
│   └── (go mod handles this)
│
├── tools/                    # Development tools
│   └── tools.go            # // +build tools
│
├── vendor/                   # Vendored dependencies (optional)
│   └── (go mod uses go.mod)
│
├── .gitignore
├── .gitattributes
├── .github/                  # GitHub-specific files
│   └── workflows/
│       └── ci.yaml
├── Makefile                  # Build automation
├── Dockerfile               # Container image
├── docker-compose.yaml      # Multi-container setup
├── README.md                # Project README
├── go.mod                   # Module definition
├── go.sum                   # Dependency checksums
└── LICENSE                  # License file
```

---

## Directory Purposes

### `/cmd` - Commands

Entry points for executable applications.

```
cmd/
├── myapp/
│   └── main.go              # Main application
├── worker/
│   └── main.go              # Background worker
└── cli/
    └── main.go              # CLI tool
```

**Rules**:
- One main.go per executable
- Put application logic in `/internal` and `/pkg`
- Keep main.go minimal

**Example main.go**:
```go
package main

import (
    "log"
    "myproject/internal/app"
)

func main() {
    app := app.New()
    if err := app.Run(); err != nil {
        log.Fatalf("Error: %v", err)
    }
}
```

### `/internal` - Private Code

Code that should NOT be imported by external packages.

```
internal/
├── auth/              # Authentication
├── storage/           # Data storage
├── config/            # Configuration
└── handler/           # HTTP handlers
```

**Key Feature**: Go compiler prevents imports outside this package tree.

```go
// ✓ Valid: importing internal package from cmd
import "myproject/internal/auth"

// ✗ Invalid: external package cannot import
// This will fail to compile
import "myproject/internal/auth"  // Error: use of internal package
```

### `/pkg` - Public Libraries

Code that CAN be imported by external packages (and third-party projects).

```
pkg/
├── cache/             # Cacheing library
├── client/            # HTTP/gRPC client
└── utils/             # Utilities
```

**Usage**:
```go
// External packages CAN import from /pkg
import "myproject/pkg/cache"

// But NOT from /internal
import "myproject/internal/auth"  // Won't compile if outside module
```

### `/api` - API Definitions

API specifications (gRPC proto files, OpenAPI specs, etc.)

```
api/
├── v1/
│   ├── users.proto
│   ├── products.proto
│   └── common.proto
└── v2/
    └── users.proto
```

### `/config` - Configuration Files

Configuration templates and defaults.

```
config/
├── config.example.yaml
├── database.yaml
└── server.yaml
```

### `/docs` - Documentation

Project documentation.

```
docs/
├── architecture.md       # Architecture decisions
├── api.md                # API documentation
├── deployment.md         # Deployment guide
└── contributing.md       # Contribution guide
```

### `/examples` - Examples

Example code showing how to use the library.

```
examples/
├── basic_usage.go
├── advanced_features.go
└── integration_example.go
```

### `/scripts` - Scripts

Build and deployment scripts.

```
scripts/
├── build.sh
├── test.sh
├── lint.sh
└── deploy.sh
```

### `/testdata` - Test Data

Fixtures and test data files.

```
testdata/
├── sample.json
├── fixture.xml
└── test_config.yaml
```

### `/tools` - Development Tools

Development tool definitions (for go:generate).

```go
// tools/tools.go
// +build tools

package tools

import (
    _ "github.com/golangci/golangci-lint/cmd/golangci-lint"
    _ "github.com/google/wire/cmd/wire"
)
```

---

## Package Organization

### When to Create a Package

Create a new package when:

1. **Related functionality**: Group similar features
   ```
   pkg/
   ├── storage/      # All storage-related code
   │   ├── sqlite.go
   │   ├── postgres.go
   │   └── storage.go
   ```

2. **Public API**: Code meant for external use
   ```
   pkg/
   ├── client/       # Client library
   │   ├── client.go
   │   └── options.go
   ```

3. **Avoid import cycles**: Separate code with circular dependencies
   ```
   internal/
   ├── service/      # Uses repository
   └── repository/   # Used by service (no cycle)
   ```

### Package Naming

**Rules**:
- Short, concise names (one word if possible)
- Lowercase only
- Avoid generic names like "utils" or "helper"

**Good**:
```
pkg/
├── cache/
├── auth/
├── server/
└── client/
```

**Bad**:
```
pkg/
├── utilities/        # Too generic
├── helpers/          # Too generic
├── common/           # Too generic
└── v1/               # Use versioning in API paths, not package names
```

---

## File Organization

### Naming Conventions

**Test files**: `*_test.go`
```go
// user.go
func GetUser(id string) (*User, error) { ... }

// user_test.go
func TestGetUser(t *testing.T) { ... }
```

**Interface files**: `interface.go` or descriptive name
```go
// storage.go - defines Storage interface
type Storage interface {
    Save(key string, value []byte) error
    Get(key string) ([]byte, error)
}

// sqlite.go - implements Storage
type SQLiteStorage struct { ... }
```

**Implementation files**: Match interface or feature name
```
internal/
├── storage.go       # Interface definition
├── sqlite.go        # SQLite implementation
├── postgres.go      # PostgreSQL implementation
└── storage_test.go  # Tests
```

### Organizing Large Packages

For large packages, split into multiple files:

```
internal/auth/
├── auth.go           # Main interface and public functions
├── jwt.go            # JWT implementation
├── oauth.go          # OAuth2 implementation
├── password.go       # Password utilities
├── session.go        # Session management
└── auth_test.go      # Tests
```

---

## Module Structure

### go.mod

```
module github.com/yourname/myproject

go 1.21

require (
    github.com/some/package v1.2.3
)
```

### Versioning

External packages should support semantic versioning:

```
pkg/
└── v1/
    ├── client.go
    └── types.go
```

This allows multiple versions:

```go
import "myproject/pkg/v1"  // Old API
import "myproject/pkg/v2"  // New API
```

---

## Small vs Large Projects

### Small Project (Single Binary)

```
myapp/
├── main.go
├── handler.go
├── storage.go
├── go.mod
├── README.md
└── Dockerfile
```

### Medium Project (Library + CLI)

```
myproject/
├── cmd/
│   └── cli/main.go
├── pkg/
│   └── mylib/
└── README.md
```

### Large Project (Microservice)

```
myservice/
├── cmd/
│   ├── myservice/main.go
│   └── worker/main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── storage/
│   └── config/
├── pkg/
│   ├── client/
│   └── types/
├── api/
├── docs/
├── scripts/
├── testdata/
└── Dockerfile
```

---

## Best Practices

1. **Use `/internal` for application code**
   - Prevents accidental external dependencies
   - Clear API boundary

2. **Use `/pkg` for libraries**
   - Code meant for external consumption
   - Stable API contract

3. **Keep `/cmd` minimal**
   - Only parse flags and call internal code
   - No business logic in main.go

4. **One responsibility per package**
   - auth, storage, config, etc.
   - Not "utils" or "helpers"

5. **Keep packages small**
   - Easier to understand and test
   - Clear interfaces

6. **Limit internal knowledge**
   - Don't expose internal packages to cmd
   - Use defined interfaces

7. **Test in same package**
   - `mycode.go` + `mycode_test.go`
   - Direct access to private functions

---

## Examples

### Example 1: Small CLI Tool

```
mytool/
├── main.go
├── cmd.go
├── go.mod
└── README.md
```

### Example 2: HTTP API

```
myapi/
├── cmd/
│   └── server/main.go
├── internal/
│   ├── handler/handler.go
│   ├── service/service.go
│   └── storage/storage.go
├── pkg/
│   └── client/client.go
├── api/
│   └── service.proto
├── go.mod
└── Dockerfile
```

### Example 3: Library

```
mylib/
├── pkg/
│   └── mylib/
│       ├── public.go
│       ├── types.go
│       └── public_test.go
├── internal/
│   └── internal_helpers.go
├── examples/
│   └── example.go
├── go.mod
└── README.md
```

---

## Summary

- **`/cmd`**: Application entry points only
- **`/internal`**: Private application code
- **`/pkg`**: Public libraries (if applicable)
- **`/api`**: API definitions
- **`/docs`**: Documentation
- Keep packages focused and small
- Use clear, descriptive names
- Test files alongside implementation

This structure scales from small scripts to large distributed systems.

---

## Next: Chapter 2 - Interface Design

We'll learn how to design effective interfaces for loose coupling and testability.
