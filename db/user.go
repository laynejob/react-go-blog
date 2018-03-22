package db

import (
    "database/sql"
    "github.com/go-sql-driver/mysql"
    "fmt"
    //"github.com/golang/glog"
)

var userTableCreateStatement =
    `CREATE TABLE IF NOT EXISTS user (
        id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
        username VARCHAR(255) NOT NULL,
        password VARCHAR(32) NULL DEFAULT NULL,
        isSuper TINYINT(1) NULL DEFAULT NULL,
        nickname VARCHAR(255) NULL DEFAULT NULL,
        avatar VARCHAR(255) NULL DEFAULT NULL,
        email VARCHAR(255) NULL DEFAULT NULL,
        qq VARCHAR(255) NULL DEFAULT NULL,
        wechat VARCHAR(255) NULL DEFAULT NULL,
        ctime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        ltime DATETIME NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        PRIMARY KEY (id),
        UNIQUE KEY username (username),
        UNIQUE KEY email (email)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

type userStmt struct {
    conn *sql.DB

    insert *sql.Stmt
    selectById *sql.Stmt
    //list   *sql.Stmt
    //listBy *sql.Stmt
    //get    *sql.Stmt
    //update *sql.Stmt
    //delete *sql.Stmt
}

type UserTableInterface interface {
    // ListBooks returns a list of books, ordered by title.
    //ListBooks() ([]*Book, error)
    //
    //// ListBooksCreatedBy returns a list of books, ordered by title, filtered by
    //// the user who created the book entry.
    //ListBooksCreatedBy(userID string) ([]*Book, error)
    //
    // GetBook retrieves a book by its ID.
    SelectById(id uint64) (*User, error)
    SelectByUsername(username string) (*User, error)
    SelectByEmail(email string) (*User, error)

    // AddBook saves a given book, assigning it a new ID.
    Add(b *User) (id int64, err error)
    //
    //// DeleteBook removes a given book by its ID.
    //DeleteBook(id int64) error
    //
    //// UpdateBook updates the entry for a given book.
    //UpdateBook(b *Book) error

    // Close closes the database, freeing up any available resources.
    // TODO(cbro): Close() should return an error.
    Close()
}

// Ensure mysqlDB conforms to the BookDatabasbhj j HY BKIMe interface.
var _ UserTableInterface = &userStmt{}

func (c MySQLConfig) ensureUserTable(conn *sql.DB) error {
    if _, err := conn.Exec("DESCRIBE user"); err != nil {
        // MySQL error 1146 is "table does not exist"
        if mErr, ok := err.(*mysql.MySQLError); ok && mErr.Number == 1146 {
            _, err := conn.Exec(userTableCreateStatement)
            if err != nil {
                return err
            }
        } else {
            // Unknown error.
            return fmt.Errorf("mysql: could not connect to the database: %v", err)
        }
    }
    return nil

}

// Close closes the database, freeing up any resources.
func (db *userStmt) Close() {
    db.conn.Close()
}

func (c MySQLConfig) prepareUser(conn *sql.DB) (UserTableInterface, error) {
    var err error
    db := &userStmt{
        conn: conn,
    }

    if db.insert, err = conn.Prepare(insertStatement); err != nil {
        return nil, fmt.Errorf("mysql: prepare insert: %v", err)
    }
    if db.selectById, err = conn.Prepare(selectByIdStatement); err != nil {
        return nil, fmt.Errorf("mysql: prepare selectById: %v", err)
    }
    return db, nil
}

func scanUser(s rowScanner) (*User, error) {
    var (
        id          uint64
        username     string
        password     sql.NullString
        isSuper     sql.NullBool
        nickname     sql.NullString
        avatar         sql.NullString
        email         string
        qq             sql.NullString
        wechat         sql.NullString
        ctime         NullTime
        ltime         NullTime
    )
    if err := s.Scan(&id, &username, &password, &isSuper, &nickname,
        &avatar, &email, &qq, &wechat, &ctime, &ltime); err != nil {
        return nil, err
    }

    user := &User{
        id,
        username,
        password.String,
        isSuper.Bool,
        nickname.String,
        avatar.String,
        email,
        qq.String,
        wechat.String,
        ctime.Time,
        ltime.Time,
        //ctime,
        //ltime,
    }
    return user, nil
}

const insertStatement = `
  INSERT INTO user (
    Username, Password,    IsSuper, Nickname, Avatar, Email, QQ, WeChat
  ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

// AddBook saves a given book, assigning it a new ID.
func (db *userStmt) Add(b *User) (id int64, err error) {
    r, err := execAffectingOneRow(db.insert, b.Username, b.Password, b.IsSuper,
        b.Nickname, b.Avatar, b.Email, b.QQ, b.WeChat)
    if err != nil {
        return 0, err
    }

    lastInsertID, err := r.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("mysql: could not get last insert ID: %v", err)
    }
    return lastInsertID, nil
}

const selectByIdStatement =    `
    SELECT id,username,password,isSuper,nickname,avatar,email,qq,wechat,UNIX_TIMESTAMP(ctime),UNIX_TIMESTAMP(ltime)
    from user where id=?`

func (db *userStmt) SelectById(id uint64) (*User, error) {
    user, err := scanUser(db.selectById.QueryRow(id))
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("mysql: could not find user with id %d", id)
    }
    if err != nil {
        return nil, fmt.Errorf("mysql: could not get user: %v", err)
    }
    //glog.Info(fmt.Sprintf("%+v", user))
    return user, nil
}

func (db *userStmt) SelectByUsername(username string) (*User, error) {
    return nil, nil
}

func (db *userStmt) SelectByEmail(email string) (*User, error) {
    return nil, nil
}