package db

import (
    "database/sql"
    "fmt"
    "database/sql/driver"
    "github.com/go-sql-driver/mysql"
    "time"
)

type MySQLConfig struct {
    // Optional.
    Username, Password string

    // Host of the MySQL instance.
    //
    // If set, UnixSocket should be unset.
    Host string

    // Port of the MySQL instance.
    //
    // If set, UnixSocket should be unset.
    Port string

    Database string
}

type mysqlDB struct {
    conn *sql.DB
    UserTable UserTableInterface
}

// rowScanner is implemented by sql.Row and sql.Rows
type rowScanner interface {
    Scan(dest ...interface{}) error
}

// dataStoreName returns a connection string suitable for sql.Open.
func (c MySQLConfig) dataStoreName(databaseName string) string {
    var cred string
    // [username[:password]@]
    if c.Username != "" {
        cred = c.Username
        if c.Password != "" {
            cred = cred + ":" + c.Password
        }
        cred = cred + "@"
    }

    //if c.UnixSocket != "" {
    //    return fmt.Sprintf("%sunix(%s)/%s", cred, c.UnixSocket, databaseName)
    //}
    return fmt.Sprintf("%stcp([%s]:%s)/%s", cred, c.Host, c.Port, databaseName)
}

// ensureTableExists checks the table exists. If not, it creates it.
func (c MySQLConfig) ensureTableExists() error {
    conn, err := sql.Open("mysql", c.dataStoreName(""))
    if err != nil {
        return fmt.Errorf("mysql: could not get a connection: %v", err)
    }
    defer conn.Close()

    // Check the connection.
    if conn.Ping() == driver.ErrBadConn {
        return fmt.Errorf("mysql: could not connect to the database. " +
            "could be bad address, or this address is not whitelisted for access")
    }

    if _, err := conn.Exec(fmt.Sprintf("USE %s", c.Database)); err != nil {
        // MySQL error 1049 is "database does not exist"
        if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1049 {
            if err := createDatabase(conn); err != nil {
                return err
            }
        }
    }

    if err := c.ensureUserTable(conn); err != nil {
        return err
    }

    return nil
}

func newMySQLDB(c MySQLConfig) (*mysqlDB, error) {
    // Check database and table exists. If not, create it.
    if err := c.ensureTableExists(); err != nil {
        return nil, err
    }

    conn, err := sql.Open("mysql", c.dataStoreName(c.Database))
    if err != nil {
        return nil, fmt.Errorf("mysql: could not get a connection: %v", err)
    }
    if err := conn.Ping(); err != nil {
        conn.Close()
        return nil, fmt.Errorf("mysql: could not establish a good connection: %v", err)
    }

    db := &mysqlDB{
        conn: conn,
    }

    db.UserTable, err = c.prepareUser(conn)
    return db, nil
}

var createDatabaseStatement = []string{
    `CREATE DATABASE IF NOT EXISTS blog DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
    `USE blog`,
}

func createDatabase(conn *sql.DB) error {
    for _, stmt := range createDatabaseStatement {
        _, err := conn.Exec(stmt)
        if err != nil {
            return err
        }
    }
    return nil
}

// execAffectingOneRow executes a given statement, expecting one row to be affected.
func execAffectingOneRow(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
    r, err := stmt.Exec(args...)
    if err != nil {
        return r, fmt.Errorf("mysql: could not execute statement: %v", err)
    }
    rowsAffected, err := r.RowsAffected()
    if err != nil {
        return r, fmt.Errorf("mysql: could not get rows affected: %v", err)
    } else if rowsAffected != 1 {
        return r, fmt.Errorf("mysql: expected 1 row affected, got %d", rowsAffected)
    }
    return r, nil
}


type NullTime struct{
    Time JsonTime
    Valid  bool // Valid is true if String is not NULL
}
// Scan implements the Scanner interface.
func (t *NullTime) Scan(value interface{}) error {
    if value == nil {
        t.Time, t.Valid = JsonTime{}, false
        return nil
    }
    t.Valid = true
    //tm, err := time.Parse("2006-01-02 15:04:05", string(value.([]uint8)[:]))
    tm := time.Unix(value.(int64), 0)
    //if err != nil {
    //    t.Time, t.Valid = JsonTime{}, false
    //    return nil
    //}
    t.Time = JsonTime(tm)
    return nil
}