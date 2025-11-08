# Install golang-migrate CLI
go install -tags 'sqlserver' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Add the migration library to your project
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/golang-migrate/migrate/v4/database/sqlserver
go get -u github.com/golang-migrate/migrate/v4/source/file
