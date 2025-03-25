# Go Authentication API

## 📂 Project Structure

```plaintext
/auth-api
│── /cmd               # Entry point commands (optional)
│── /config            # Configuration management (env, config files)
│── /internal          # Internal application logic
│   │── /handlers      # HTTP request handlers (controllers)
│   │── /models        # Database models and schemas
│   │── /repository    # Data persistence layer (queries, DB operations)
│   │── /services      # Business logic and authentication services
│   │── /middleware    # Middleware (Auth, Logging, Rate Limiting)
│── /pkg               # Shared utilities (JWT, password hashing, etc.)
│── /migrations        # Database migration files (SQL scripts)
│── main.go            # Application entry point
│── go.mod             # Go module dependencies
│── go.sum             # Dependency checksums
│── .env               # Environment variables
│── Makefile           # Common commands (run, test, build)
```


## ⚡ Tech Stack

- **Go Web Framework**: [`net/http`] or [`gin-gonic/gin`]
- **Database**: PostgreSQL (via [`gorm`])
- **Authentication**: JWT (`golang-jwt/jwt`)
- **Password Hashing**: bcrypt (`golang.org/x/crypto/bcrypt`)
- **Environment Config**: [`godotenv`]

## 🔥 API Endpoints

| Method | Endpoint  | Description          | Auth Required |
|--------|----------|----------------------|--------------|
| `POST` | `/register` | Register a new user | ❌ |
| `POST` | `/login`    | Authenticate user and return token | ❌ |
| `GET`  | `/profile`  | Fetch user profile  | ✅ |
| `POST` | `/logout`   | Invalidate token    | ✅ |

## 🛠 Setup & Usage
