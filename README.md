# Snippetbox

Snippetbox is a server-rendered web application written in Go. It allows users
to register, log in, create text snippets with an expiry period, and browse
recent snippets. The application uses MySQL for persistent data and for
server-side session storage.

## Features

- Public home page showing the ten most recently created, non-expired snippets.
- Individual snippet pages at `/snippet/view/:id`.
- User registration with:
  - Name, email, and password fields.
  - Required-field validation.
  - Email-format validation.
  - Minimum password length of eight characters.
  - Duplicate-email detection.
- User login and logout.
- Password hashing with bcrypt.
- Authenticated snippet creation.
- Snippet expiry choices:
  - One day.
  - One week.
  - One year.
- Automatic filtering of expired snippets.
- Server-side sessions stored in MySQL.
- Twelve-hour session lifetime.
- CSRF protection on form submissions.
- Authentication middleware for protected pages.
- Session ID renewal during login and logout.
- Flash messages for successful signup, snippet creation, and logout.
- HTML template rendering with reusable base and navigation templates.
- Form values and validation errors are redisplayed after invalid submissions.
- Request logging and panic recovery.
- Security-related HTTP response headers, including:
  - Content Security Policy.
  - Referrer Policy.
  - `X-Content-Type-Options: nosniff`.
  - `X-Frame-Options: deny`.
- Static CSS, JavaScript, logo, and favicon assets.

## Technology stack

- Go
- `net/http`
- `httprouter` for routing
- `alice` for middleware chains
- MySQL
- `database/sql`
- `go-sql-driver/mysql`
- `scs` and `scs/mysqlstore` for sessions
- `nosurf` for CSRF protection
- `go-playground/form` for form decoding
- `bcrypt` from `golang.org/x/crypto` for password hashing
- Go `html/template` for server-side HTML rendering

## Requirements

Install the following before running the application:

1. Go 1.27.1 or a compatible Go version supported by the module.
2. MySQL Server.
3. A MySQL client or another way to execute `sqlFile.sql`.

The application expects to be started from the repository root because it
loads templates and static files using relative paths such as `./ui/html` and
`./ui/static`.

## Database setup

The file [`sqlFile.sql`](./sqlFile.sql) creates the database, tables, sample
snippets, and the application database user.

The script currently uses these defaults:

- Database: `snippetbox`
- MySQL user: `web`
- MySQL password: `pass`
- MySQL host: `localhost`

Run the script from a MySQL client:

```sql
SOURCE P:/letsGo/sqlFile.sql;
```

Or from PowerShell using the MySQL client:

```powershell
mysql -u root -p < sqlFile.sql
```

Use a strong password instead of `pass` for any real deployment. If the
database user or password is changed, pass the matching DSN when starting the
application.

## Run the application

From the project root:

```powershell
go mod download
go run ./cmd/web
```

The default server address is:

```text
http://localhost:4000
```

Open that address in a browser after the server starts.

### Custom address and database connection

The server accepts these command-line flags:

| Flag | Default | Description |
|---|---|---|
| `-addr` | `:4000` | HTTP network address |
| `-dsn` | `web:pass@/snippetbox?parseTime=true` | MySQL data source name |

Examples:

```powershell
go run ./cmd/web -addr ":8080"
```

```powershell
go run ./cmd/web -dsn "web:your-password@tcp(127.0.0.1:3306)/snippetbox?parseTime=true"
```

The DSN must include `parseTime=true` so MySQL datetime values can be scanned
into Go `time.Time` values.

## Application routes

| Method | Path | Access | Purpose |
|---|---|---|---|
| `GET` | `/` | Public | Display the latest ten active snippets |
| `GET` | `/snippet/view/:id` | Public | Display one active snippet |
| `GET` | `/user/signup` | Public | Display the signup form |
| `POST` | `/user/signup` | Public | Validate and create a user |
| `GET` | `/user/login` | Public | Display the login form |
| `POST` | `/user/login` | Public | Authenticate a user |
| `GET` | `/snippet/create` | Authenticated | Display the create-snippet form |
| `POST` | `/snippet/create` | Authenticated | Validate and save a snippet |
| `POST` | `/user/logout` | Authenticated | Log out the current user |
| `GET` | `/static/*filepath` | Public | Serve CSS, JavaScript, and images |

Invalid snippet IDs and expired or missing snippets return a 404 response.
Invalid form submissions return a 422 response and render the form again with
validation errors.

## User flow

1. Open `/user/signup` and create an account.
2. Open `/user/login` and authenticate.
3. Select **Create snippet** from the navigation.
4. Enter a title and content.
5. Choose an expiry period and publish the snippet.
6. The application redirects to the new snippet's public view page.
7. Log out using the navigation form when finished.

## Validation rules

### Signup

- Name must not be blank.
- Email must not be blank and must match the configured email regular
  expression.
- Password must not be blank.
- Password must contain at least eight characters.
- Email addresses must be unique.

### Login

- Email must not be blank and must pass email validation.
- Password must not be blank.
- Invalid credentials are reported as a non-field error without revealing
  whether the email or password was incorrect.

### Snippet creation

- Title must not be blank.
- Title must be no longer than 100 characters.
- Content must not be blank.
- Expiry must be exactly `1`, `7`, or `365` days.

Character limits count Unicode runes rather than raw bytes.

## Project structure

```text
letsGo/
├── README.md
├── go.mod
├── go.sum
├── sqlFile.sql
├── cmd/
│   └── web/
│       ├── context.go       # Request-context keys
│       ├── handlers.go      # HTTP handlers and form validation
│       ├── helpers.go       # Rendering, form decoding, and HTTP helpers
│       ├── main.go          # Application setup and server startup
│       ├── middleware.go    # Authentication, CSRF, logging, and security
│       ├── routes.go        # HTTP route registration
│       └── templates.go     # Template cache and template data
├── internal/
│   ├── models/
│   │   ├── errors.go        # Shared model-layer errors
│   │   ├── snippets.go      # Snippet queries and persistence
│   │   └── users.go         # User creation and authentication
│   └── validator/
│       └── validator.go     # Reusable form-validation helpers
└── ui/
    ├── html/
    │   ├── base.tmpl        # Common HTML document layout
    │   ├── pages/
    │   │   ├── create.tmpl  # Create-snippet page
    │   │   ├── home.tmpl    # Latest-snippets page
    │   │   ├── login.tmpl   # Login page
    │   │   ├── signup.tmpl  # Signup page
    │   │   └── view.tmpl    # Individual-snippet page
    │   └── partials/
    │       └── nav.tmpl     # Shared navigation
    └── static/
        ├── css/
        │   └── main.css     # Application styling
        ├── img/
        │   ├── favicon.ico
        │   └── logo.png
        └── js/
            └── main.js      # Active navigation-link behavior
```

## Database schema

The SQL setup creates four tables:

- `snippets`: snippet title, content, creation time, and expiry time.
- `users`: user name, unique email, bcrypt password hash, and creation time.
- `sessions`: server-side session token, session data, and expiration time.
- The `idx_snippets_created` and `sessions_expiry_idx` indexes support common
  lookup and cleanup operations.

Snippet queries use UTC timestamps and only return rows whose expiry is later
than the current UTC time.

## Security and deployment notes

- Passwords are never stored in plaintext; bcrypt hashes are stored in the
  `users.hashed_password` column.
- SQL statements use placeholders for user and snippet values.
- Form submissions include CSRF tokens.
- Authenticated pages use `Cache-Control: no-store`.
- The session cookie is configured as `HttpOnly` and `Secure`. For production,
  serve the application over HTTPS.
- The default database credentials in `sqlFile.sql` and the default DSN are
  development values. Change them before deployment.
- The current server configuration is intended for a simple local application;
  production deployments should also place it behind HTTPS and a suitable
  reverse proxy.

## Troubleshooting

### Database connection errors

Confirm that:

1. MySQL is running.
2. The `snippetbox` database exists.
3. The `web` user can connect from the configured host.
4. The DSN password and database name match the database setup.
5. `parseTime=true` is present in the DSN.

### Templates or static files are not found

Start the application from the repository root:

```powershell
go run ./cmd/web
```

Running the compiled binary from another working directory can prevent the
relative `./ui/html` and `./ui/static` paths from resolving.

### Running checks

Run the Go test and compile checks with:

```powershell
go test ./...
```

The repository currently does not contain dedicated Go test files, so this
command primarily verifies that all packages compile.
