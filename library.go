package main

import (
  "fmt"
  "io/ioutil"
  //"path/filepath"
  //"reflect"
  "strings"
  //"sync"
  //"time"

  // mysql connector
  _ "github.com/go-sql-driver/mysql"
  sqlx "github.com/jmoiron/sqlx"
)

const (
  User     = "105gun"
  Password = ""
  DBName   = "ass3"

  //StringCreate=
)

type Library struct {
  db *sqlx.DB
}

// From ass2
func mustExecute(db *sqlx.DB, SQLs []string) {
	for _, s := range SQLs {
		_, err := db.Exec(s)
		if err != nil {
			panic(err)
		}
	}
}

func executeSQLsFromFile(filePath string, db *sqlx.DB) error {
	binData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	SQLs := string(binData)
	for _, s := range strings.Split(SQLs, ";") {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		_, err := db.Exec(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// From bolierplate
func (lib *Library) ConnectDB() {
  fmt.Println("--Trying connection.")
  db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
  if err != nil {
    fmt.Println("--Connection fail!")
    panic(err)
  }
  lib.db = db
  fmt.Println("--Connection success!")
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error {
  err := executeSQLsFromFile("CreateTable.sql",lib.db)
  return err
}

// Insert test data
func (lib *Library) TestData() error {
  fmt.Println("--Inserting test data")
  err := executeSQLsFromFile("TestData.sql",lib.db)
  return err
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string) error {
  _, err := lib.db.Exec(fmt.Sprintf("INSERT INTO booktype (ISBN, title, author) VALUES (%s, %s, %s)", ISBN, title, author))
  return err
}

// Init step, including drop same name db & create db & create table & set root user
func initdb(lib *Library) {
  db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
  if err != nil {
    panic(err)
  }

  mustExecute(db, []string{
    "DROP DATABASE IF EXISTS ass3",
    "CREATE DATABASE ass3",
  })
  err=lib.CreateTables()
  if err!=nil {
    panic(err)
  }
  
  db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
  if err != nil {
    fmt.Println("--Connection fail!")
    panic(err)
  }

  // Set root user
  _, err = db.Exec("INSERT INTO user (id, authority) VALUES (0,100)")
  if err != nil {
    panic(err)
  }

  // Insert test data, should be ignored in real application
  err=lib.TestData()
  if err != nil {
    panic(err)
  }

  lib.db = db
}

// Login as student
func stulogin(lib *Library, id int, stat int) {
  fmt.Println("Log success. Welcome student ", id, ".")
  var op int
SL:
  for {
    fmt.Println("1.Borrow\t2.Return\t3.Extend deadline")
    fmt.Println("4.Query books\t5.Query not returned\t6.Check deadline")
    fmt.Println("7.Logout")
    fmt.Scanln(&op)
    switch op {
      case 1:
      case 2:
      case 7:
        fmt.Println("Now logout.")
        break SL
      default:
        fmt.Println("Please type in right operation number.")
    }
  }
}

// Login as admin
func adminlogin(lib *Library, id int) {
  
}

// Login application
func loginapp(lib *Library) {
  fmt.Println("Please Enter your uid!")
  var op int
  fmt.Scanln(&op)
  rows, err := lib.db.Query(fmt.Sprintf("SELECT authority FROM user WHERE id=%d",op))
  if err != nil {
    fmt.Println("Log fail!")
    return
  }
  rows.Next()
  var au int
  rows.Scan(&au)
  switch au {
    case 0:	// banned student
      stulogin(lib, op, 0)
    case 1:	// common student
      stulogin(lib, op, 1)
    default:	// admin
      adminlogin(lib, op)
  }
}

func main() {
  fmt.Println("Welcome to the Library Management System!")
  var op int
  var lib Library
  lib.ConnectDB()
M:
  for {
    fmt.Println("1.Init DB\t2.Login\t3.Exit")
    fmt.Scanln(&op)
    switch op {
      case 1:
        initdb(&lib)
      case 2:
        loginapp(&lib)
      case 3:
        fmt.Println("Bye Bye!")
        break M
      default:
        fmt.Println("Please type in right operation number.")
    }
  }
}
