package dbs

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Name string `json:"name"`
	Address string `json:"address"`
}

type Content struct {
	Title string `json:"title"`
	ContentPath string `json:"content_path"`
	ContentHash string `json:"content_hash"`
	Address string `json:"address"`
	TokenID string `json:"token_id"`
}

const (
	// TODO: database config
	driverName = "mysql"
	userName = "admin"
	password = "123456"
	ip = "127.0.0.1"
	port = "3306"
	dbName = "receipt"
)

var DBConn *sql.DB

func init() {
	DBConn = InitDB(driverName)
}

func InitDB(driver string) *sql.DB {
	path := strings.Join([]string{userName, ":", password, "@tcp(",ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	db, err := sql.Open(driver, path)
	if err != nil {
		panic(err.Error())
	}
	return db
}


func (u User) AddUser() error {
	_, err := DBConn.Exec("insert into t_user(email, name, password, address) values(?, ?, ?, ?);", u.Email, u.Name, u.Password, u.Address)
	if err != nil {
		fmt.Println("Failed to insert t_user: ", err)
		return err
	}
	return nil
}

func (u *User) QueryUser() (bool, error) {
	rows, err := DBConn.Query("select email, address from t_user where name = ? and password = ?;", u.Name, u.Password)
	if err != nil {
		fmt.Println("failed to select t_user: ", err)
		return false, err
	}
	// process query results
	if rows.Next() {
		err = rows.Scan(&u.Email, &u.Address)
		if err != nil {
			fmt.Println("Failed to scan select t_user", err)
			return false, err
		}
		return true, nil
	}

	return false, err
}

func (c Content) AddContent() error {
	_, err := DBConn.Exec("insert into t_content (title, content, content_hash, address, token_id) values (?, ?, ?, ?, ?);", c.Title, c.ContentPath, c.ContentHash, c.Address, c.TokenID)
	if err != nil {
		fmt.Println("Failed to insert t_content: ", err)
		return err
	}
	return err
}

func QueryContent(address string) ([]Content, error) {
	rows, err := DBConn.Query("select title, content, content_hash, token_id from t_content where address = ?;", address)
	if err != nil {
		fmt.Println("Failed to QueryUser t_content: ", err)
		return nil, err
	}
	// process query results
	var s []Content
	var c Content
	for rows.Next() {
		err = rows.Scan(&c.Title, &c.ContentPath, &c.ContentHash, &c.TokenID)
		if err != nil {
			fmt.Println("Failed to scan select t_content: ", err)
			return s, err
		}
		s = append(s, c)
	}
	return s, nil
}