# Simple Bank

## Description
This is a simple bank project that allows users to manage their accounts and perform basic banking operations. The project is implemented in Go, using SQLC for database operations.

## Technologies Used
- Go
- PostgreSQL
- Docker
- SQLC

## Installation
1. Clone the repository: `git clone https://github.com/haziqkamel/go-simplebank.git`
2. Install Go dependencies: `go mod download`
3. Set up the PostgreSQL database using Docker: `make postgres`
4. Create the database: `make createdb`
5. Run database migrations: `make migrateup`
6. Generate SQLC code: `make sqlc`

## Mocking

### Required:
1. mockgen
2. mockegen/model

Then run command below:
```
mockgen -package mockdb -destination db/mock/store.go github.com/haziqkamel/simplebank/db/sqlc Store
```

## Contact
For any inquiries or support, please contact us at [haziqkamel@outlook.com](mailto:haziqkamel@outlook.com).