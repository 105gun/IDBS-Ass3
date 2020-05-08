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

  DueTime  = 7
)

type Library struct {
  db *sqlx.DB
}

func GetTime() int {
  fmt.Println("--Type in current time")
  var ans int
  fmt.Scanln(&ans)
  return ans
}

func GetNum(name string) int {
  rows, err = lib.db.Query(fmt.Sprintf("SELECT MAX(id) FROM %s",name))
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  return tmp+1
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
func (lib *Library) AddBook(title, author, ISBN string){
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM booktype WHERE ISBN=%s", ISBN)
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  if tmp!=0 {

  } else{
    lib.db.Exec(fmt.Sprintf("INSERT INTO booktype VALUES (%s, %s, %s)", ISBN, title, author))
    lib.db.Exec(fmt.Sprintf("INSERT INTO book VALUES (%d, %s, 1, 0, %s)", ISBN, ""))
  }
}

// Try to borrow a book with a student id
func (lib *Library) BorrowBook(id, bid int) {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT existed FROM book WHERE id=%d",bid))
  if err != nil {
    fmt.Println("Search fail!")
    return
  }

  rows.Next()
  var au int
  rows.Scan(&au)

  if au==0 {
    fmt.Println("Book not available!")
  }
  if au==1 {
    lib.db.Exec(fmt.Sprintf("UPDATE book SET existed=0 WHERE ID=%d",bid))
    lib.db.Exec(fmt.Sprintf("INSERT INTO borrow VALUES	(%d, %d, %d, %d, 0, 0)", GetNum(borrow), bid, id, GetTime(),))
  }
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
  fmt.Println()
  fmt.Println("Log success. Welcome student", id, ".")
  if stat==0 {
    fmt.Println("!!! WARNING, your account is banned because of some reason. !!!")
  }
  var op int
SL:
  for {
    fmt.Println()
    fmt.Println("1.Borrow\t2.Return\t3.Extend deadline")
    fmt.Println("4.Query books\t5.Query not returned\t6.Check deadline")
    fmt.Println("7.Logout")
    fmt.Scanln(&op)
    switch op {
      case 1:
        if stat!=0 {
          var bid int
          fmt.Println("Please enter the book id.")
          fmt.Scanln(&bid)
          lib.BorrowBook(id, bid)
        } else {
          fmt.Println("Your account is banned!")
        }
      case 2:
      case 3:
      case 4:
      case 5:
      case 6:
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
  fmt.Println()
  fmt.Println("Log success. Welcome admin", id, ".")
  var op int
AL:
  for {
    fmt.Println()
    fmt.Println("1.Add book\t2.Remove book\t3.Add user")
    fmt.Println("4.Query books\t5.Query not returned\t6.Query borrow")
    fmt.Println("7.Check deadline\t8.Check overdue\t9.Logout")
    fmt.Scanln(&op)
    switch op {
      case 1:
      case 2:
      case 3:
      case 4:
      case 5:
      case 6:
      case 7:
      case 8:
      case 9:
        fmt.Println("Now logout.")
        break AL
      default:
        fmt.Println("Please type in right operation number.")
    }
  }
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
  fmt.Println()
  fmt.Println("Welcome to the Library Management System!")
  var op int
  var lib Library
  lib.ConnectDB()
M:
  for {
    fmt.Println()
    fmt.Println("1.Init DB\t2.Login  \t3.Exit")
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
