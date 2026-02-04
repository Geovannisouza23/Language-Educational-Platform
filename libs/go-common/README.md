# Shared Go Libraries

Common libraries and utilities shared across Go microservices.

## Packages

### Logger
Standardized logging with structured fields

### Auth Client
JWT validation and user context extraction

### Database
Common database utilities and helpers

### HTTP Client
Instrumented HTTP client with retry logic

### Error Handling
Standardized error codes and responses

## Usage

Import in your service:

```go
import (
    "github.com/language-platform/libs/go-common/logger"
    "github.com/language-platform/libs/go-common/auth"
)
```
