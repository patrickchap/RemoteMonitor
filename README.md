# Project RemoteMonitor

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. 

To run this project, ensure you have the following installed:

- **Go** version 1.23 or higher
- **SQLC** for generating Go code from SQL
- **Templ** for template-based generation

### 1. Clone the repository
```bash
https://github.com/patrickchap/RemoteMonitor.git
cd RemoteMonitor
```

### 2. Install SQLC and TEMPL
```bash
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```


### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Add .env
```bash
echo -e "PORT=8080\nDB_URL=sqlite.db\nACCESS_TOKEN_SECRET=access-token\nREFRESH_TOKEN_SECRET=refresh-token\nJWT_SECRET_KEY=some-secret-key\nJWT_REFRESH_SECRET_KEY=some-refresh-secret-key" > .env
```

### 5. Create sqlite db
```bash
sqlite3 sqlite.db
```

### 6. Create database schema in the sqlite console
```sql
CREATE TABLE host_services (
    id INTEGER PRIMARY KEY,
    host_id INTEGER,
    service_id INTEGER,
    active INTEGER DEFAULT 1,
    schedule_number INTEGER DEFAULT 3,
    schedule_unit VARCHAR(255) DEFAULT 'm',
    last_check TIMESTAMP DEFAULT '0001-01-01 00:00:01',
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(255) DEFAULT 'pending',
    FOREIGN KEY (host_id) REFERENCES hosts(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE hosts (
    id INTEGER PRIMARY KEY,
    host_name VARCHAR(255) NOT NULL,
    canonical_name VARCHAR(255),
    url VARCHAR(255),
    ip VARCHAR(255),
    ipv6 VARCHAR(255),
    location VARCHAR(255),
    os VARCHAR(255),
    active INTEGER DEFAULT 1,
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE services (
    id INTEGER PRIMARY KEY,
    service_name VARCHAR(255),
    active INTEGER DEFAULT 1,
    icon VARCHAR(255),
    last_updated DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    user_active INTEGER DEFAULT 0,
    access_level INTEGER DEFAULT 3,
    email VARCHAR(255),
    password VARCHAR(60),
    deleted_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

```

### 7 Seed database in the sqlite console. Replase hashedpassword with a hashed passowrd
```sql

INSERT INTO users (first_name, last_name, user_active, access_level, email, password)
VALUES ('Admin', 'User', 1, 1, 'admin@example.com', 'hashedpassword');

INSERT INTO services (service_name, active)
VALUES ('https', 1);
```

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
