# **How to Set Up PostgreSQL in a Golang HTTP Server**

# 1. Connecting to the Database

To connect to a PostgreSQL database, use the `pq` package in combination with `database/sql`.

## Example of Connecting to PostgreSQL:

```go
import (
    "database/sql"
    _ "github.com/lib/pq" // Import PostgreSQL driver
    "log"
)

func main() {
    // Connecting to the PostgreSQL database
    dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Check if the connection is working
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
}
```

Here, we are connecting to the PostgreSQL database using a DSN (Data Source Name). The `Ping` function checks the connection's functionality.

# 2. Executing SQL Queries

The `database/sql` package allows executing SQL queries and returning results.

## Example of Executing an INSERT Query:

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Executing an INSERT query
    result, err := db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "John", 30)
    if err != nil {
        log.Fatal(err)
    }

    // Get the ID of the last inserted row
    lastInsertID, err := result.LastInsertId()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Last inserted ID: %d", lastInsertID)
}
```

In this example, `Exec` is used to insert a new record into the `users` table. The query uses placeholders `$1`, `$2`, which are replaced by parameters.

## Example of Executing a SELECT Query:

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT id, name, age FROM users WHERE age > $1", 25)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var name string
        var age int
        err := rows.Scan(&id, &name, &age)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("User: %d, %s, %d", id, name, age)
    }

    // Check for errors after iteration
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
}
```

This example performs a SELECT query and uses the `Scan` method to extract data from each row of the result.

# 3. Working with Transactions

Transactions allow you to perform multiple SQL operations atomically.

## Example of Using Transactions:

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Start a transaction
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }

    // Perform multiple operations within the transaction
    _, err = tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "Alice", 25)
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }

    _, err = tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = $1", 1)
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }

    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        tx.Rollback()
        log.Fatal(err)
    }

    log.Println("Transaction committed successfully")
}
```

In this example, several operations are performed within a single transaction. If any error occurs, the transaction is rolled back using the `Rollback` method.

# 4. Prepared Statements

Prepared statements allow you to reuse the same SQL query with different parameters.

## Example of Using Prepared Statements:

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Prepare SQL statement
    stmt, err := db.Prepare("INSERT INTO users (name, age) VALUES ($1, $2)")
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()

    // Execute prepared statement with different parameters
    _, err = stmt.Exec("Bob", 29)
    if err != nil {
        log.Fatal(err)
    }

    _, err = stmt.Exec("Charlie", 35)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Users inserted successfully")
}
```

Here, the prepared statement is compiled once and then executed multiple times with different parameters.

# 5. Scanning Results

The `QueryRow` method is used to execute queries that return a single row.

## Example of Using `QueryRow`:

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    dsn := "postgres://user:password@localhost:5432/dbname?sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var name string
    var age int

    // Execute query that returns a single row
    err = db.QueryRow("SELECT name, age FROM users WHERE id = $1", 1).Scan(&name, &age)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("User: %s, %d", name, age)
}
```

The `QueryRow` method is used to execute queries that return a single row. The `Scan` method extracts data from this row.

# 6. Connection Pooling

The `*sql.DB` object in Go represents a connection pool. It automatically manages the number of open connections to the database.

## Configuring the Connection Pool:

```go
db.SetMaxOpenConns(25) // Maximum number of open connections
db.SetMaxIdleConns(25) // Maximum number of idle connections
db.SetConnMaxLifetime(time.Minute * 5) // Connection lifetime
```

These methods allow you to control the maximum number of connections and their lifetime, which helps prevent database resource exhaustion.