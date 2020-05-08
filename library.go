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

func GetNum(lib *Library,name string) int {
  rows, _ := lib.db.Query(fmt.Sprintf("SELECT MAX(id) FROM %s",name))
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  tmp++
  return tmp
}

// Daily update
func (lib *Library) GetUpdate() {
  fmt.Println("Update start.")
  //executeSQLsFromFile("Update.sql",lib.db)
  crt:=GetTime()
  lib.db.Exec(fmt.Sprintf("UPDATE user SET authority=0 WHERE id IN (SELECT uid FROM book, booktype, borrow WHERE bid=book.id AND book.ISBN=booktype.ISBN AND removed=0 AND existed=1 AND is_returned=0 AND authority!= 100 AND time+(1 + extend_status)*7<%d)", crt))
  fmt.Println("Update Done.")
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

// Try to connect the data base
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

// Add a new student user
func (lib *Library) AddUser(sid int) error{
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM user WHERE id=%d", sid))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  if tmp!=0 {
    fmt.Println("Already existed!")
  } else{
    _, err=lib.db.Exec(fmt.Sprintf("INSERT INTO user VALUES (%d, 1);", sid))
    fmt.Println("Done!")
  }
  return err
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string) error{
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM booktype WHERE ISBN='%s'", ISBN))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  if tmp!=0 {
    _, err=lib.db.Exec(fmt.Sprintf("INSERT INTO book VALUES (%d, '%s', 1, 0, '');", GetNum(lib,"book"), ISBN))
  } else{
    lib.db.Exec(fmt.Sprintf("INSERT INTO booktype VALUES ('%s', '%s', '%s');", ISBN, title, author))
    _, err=lib.db.Exec(fmt.Sprintf("INSERT INTO book VALUES (%d, '%s', 1, 0, '');", GetNum(lib,"book"), ISBN))
  }
  fmt.Println("Done!")
  return err
}

// Remove book
func (lib *Library) RemoveBook(com string, bid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM book WHERE id=%d", bid))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  if tmp!=0 {
    _, err=lib.db.Exec(fmt.Sprintf("UPDATE book SET existed=0, removed=1, commit='%s' WHERE id=%d", com, bid))
    fmt.Println("Done!")
  } else{
    fmt.Println("Book dont existed!")
  }
  return err
}

// Query book by some order
func (lib *Library) QueryBook(op int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM book"))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  fmt.Println(tmp,"books in all.")
  if tmp!=0 {
    var s string
    switch op {
      case 0:	//title
        s="title"
      case 1:	//anthor
        s="author"
      case 2:	//ISBN
        s="ISBN"
      default:
        fmt.Println("Bad op!")
        return err
    }
    rows, err=lib.db.Query(fmt.Sprintf("SELECT book.id,book.ISBN,title,author,existed,removed FROM book, booktype WHERE book.ISBN=booktype.ISBN ORDER BY %s", s))
    if err!=nil {
      return err
    }
    fmt.Println("ID ISBN Title Author Existed Removed")
    var id,existed,removed int
    var ISBN,title,author string
    for rows.Next() {
      rows.Scan(&id,&ISBN,&title,&author,&existed,&removed)
      fmt.Println(id,ISBN,title,author,existed,removed)
    }
    fmt.Println("Done!")
  } else{
    fmt.Println("Book doesnt existed!")
  }
  return err
}

// Query un returned book
func (lib *Library) QueryNotReturned(uid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM book WHERE removed=0 AND existed=1"))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  fmt.Println(tmp,"unreturned books in all.")
  fmt.Println("Now checking",uid)
  if tmp!=0 {
    rows, err=lib.db.Query(fmt.Sprintf("SELECT book.id,book.ISBN,title,author FROM book, booktype, borrow WHERE bid=book.id AND book.ISBN=booktype.ISBN AND removed=0 AND existed=1 AND is_returned=0 AND uid=%d", uid))
    if err!=nil {
      return err
    }
    fmt.Println("ID ISBN Title Author Userid")
    var id int
    var ISBN,title,author string
    for rows.Next() {
      rows.Scan(&id,&ISBN,&title,&author)
      fmt.Println(id,ISBN,title,author,uid)
    }
    fmt.Println("Done!")
  } else{
    fmt.Println("Book doesnt existed!")
  }
  return err
}

// Query un returned book
func (lib *Library) QueryBorrow(uid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM borrow WHERE uid=%d", uid))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  fmt.Println(tmp,"borrow history in all.")
  fmt.Println("Now checking",uid)
  if tmp!=0 {
    rows, err=lib.db.Query(fmt.Sprintf("SELECT id,bid,uid,time,is_returned,extend_status FROM borrow WHERE uid=%d", uid))
    if err!=nil {
      return err
    }
    fmt.Println("ID BorrowID UserID Time IsReturned ExtendStatus")
    var id,bid,uid,time,is_returned,extend_status int
    for rows.Next() {
      rows.Scan(&id,&bid,&uid,&time,&is_returned,&extend_status)
      fmt.Println(id,bid,uid,time,is_returned,extend_status)
    }
    fmt.Println("Done!")
  } else{
    fmt.Println("History doesnt existed!")
  }
  return err
}

// Query un returned book
func (lib *Library) CheckDDL(uid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM borrow WHERE uid=%d", uid))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp int
  rows.Scan(&tmp)
  fmt.Println(tmp,"borrow history in all.")
  fmt.Println("Now checking",uid)
  if tmp!=0 {
    rows, err=lib.db.Query(fmt.Sprintf("SELECT book.id,title,time+(1 + extend_status)*7 FROM book, booktype, borrow WHERE bid=book.id AND book.ISBN=booktype.ISBN AND removed=0 AND existed=1 AND is_returned=0 AND uid=%d", uid))
    if err!=nil {
      return err
    }
    fmt.Println("BookID Title Deadline")
    var id,ddl int
    var title string
    for rows.Next() {
      rows.Scan(&id,&title,&ddl)
      fmt.Println(id,title,ddl)
    }
    fmt.Println("Done!")
  } else{
    fmt.Println("History doesnt existed!")
  }
  return err
}
// Query un returned book
func (lib *Library) CheckDue(uid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM borrow WHERE uid=%d", uid))
  if err!=nil {
    return err
  }
  rows.Next()
  var tmp,crt int
  crt=GetTime()
  rows.Scan(&tmp)
  fmt.Println(tmp,"borrow history in all.")
  fmt.Println("Now checking",uid)
  if tmp!=0 {
    rows, err=lib.db.Query(fmt.Sprintf("SELECT book.id,title,time+(1 + extend_status)*7 FROM book, booktype, borrow WHERE bid=book.id AND book.ISBN=booktype.ISBN AND removed=0 AND existed=1 AND is_returned=0 AND uid=%d AND time+(1 + extend_status)*7<%d", uid, crt))
    if err!=nil {
      return err
    }
    fmt.Println("BookID Title Deadline")
    var id,ddl int
    var title string
    for rows.Next() {
      rows.Scan(&id,&title,&ddl)
      fmt.Println(id,title,ddl)
    }
    fmt.Println("Done!")
  } else{
    fmt.Println("History doesnt existed!")
  }
  return err
}

// Try to borrow a book with a student id
func (lib *Library) BorrowBook(id, bid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT existed FROM book WHERE id=%d",bid))
  if err != nil {
    fmt.Println("Search fail!")
    return err
  }

  rows.Next()
  var au int
  rows.Scan(&au)

  if au==0 {
    fmt.Println("Book not available!")
  }
  if au==1 {
    lib.db.Exec(fmt.Sprintf("UPDATE book SET existed=0 WHERE id=%d",bid))
    _, err=lib.db.Exec(fmt.Sprintf("INSERT INTO borrow VALUES	(%d, %d, %d, %d, 0, 0)", GetNum(lib, "borrow"), bid, id, GetTime(),))
    fmt.Println("Done!")
  }
  return err
}

// Try to return a book
func (lib *Library) ReturnBook(bid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM book WHERE id=%d",bid))
  if err != nil {
    fmt.Println("Search fail!")
    return err
  }

  rows.Next()
  var au int
  rows.Scan(&au)

  if au==0 {
    fmt.Println("Book not existed!")
  }
  if au==1 {
    lib.db.Exec(fmt.Sprintf("UPDATE book SET existed=1 WHERE id=%d",bid))
    _, err=lib.db.Exec(fmt.Sprintf("UPDATE borrow SET is_returned=1 WHERE bid=%d AND is_returned=0",bid))
    fmt.Println("Done!")
  }
  return err
}

// Extend ddl
func (lib *Library) Extend(oid int,uid int) error {
  rows, err := lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM borrow WHERE id=%d AND is_returned=0 AND uid=%d",oid, uid))
  if err != nil {
    fmt.Println("Search fail!")
    return err
  }

  rows.Next()
  var au, tmp int
  rows.Scan(&au)

  if au==0 {
    fmt.Println("Borrow not existed or already returned!")
  }
  if au==1 {
    rows, err = lib.db.Query(fmt.Sprintf("SELECT extend_status FROM borrow WHERE id=%d AND is_returned=0 AND uid=%d",oid, uid))
    rows.Next()
    rows.Scan(&tmp)
    if tmp>=3 {
      fmt.Println("Can not extend any more!")
      return err
    }
    _, err=lib.db.Exec(fmt.Sprintf("UPDATE borrow SET extend_status=%d WHERE id=%d", tmp+1, oid))
    fmt.Println("Done!")
  }
  return err
}

// Init step, including drop same name db & create db & create table & set root user
func initdb(lib *Library) error {
  db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/", User, Password))
  if err != nil {
    return err
  }

  mustExecute(db, []string{
    "DROP DATABASE IF EXISTS ass3",
    "CREATE DATABASE ass3",
  })
  err=lib.CreateTables()
  if err!=nil {
    return err
  }
  
  db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
  if err != nil {
    fmt.Println("--Connection fail!")
    return err
  }

  // Set root user
  _, err = db.Exec("INSERT INTO user (id, authority) VALUES (0,100)")
  if err != nil {
    return err
  }

  // Insert test data, should be ignored in real application
  err=lib.TestData()
  if err != nil {
    return err
  }

  lib.db = db
  return err
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
        var bid int
        fmt.Println("Please enter the book id.")
        fmt.Scanln(&bid)
        lib.ReturnBook(bid)
      case 3:
        var oid int
        fmt.Println("Please enter the borrow id.")
        fmt.Scanln(&oid)
        lib.Extend(oid, id)
      case 4:
        op:=0
        fmt.Println("0.Title 1.Author 2.ISBN")
        fmt.Scanln(&op)
        lib.QueryBook(op)
      case 5:
        lib.QueryNotReturned(id)
      case 6:
        lib.CheckDDL(id)
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
  lib.GetUpdate()
AL:
  for {
    fmt.Println()
    fmt.Println("1.Add book\t2.Remove book\t3.Add user")
    fmt.Println("4.Query books\t5.Query not returned\t6.Query borrow")
    fmt.Println("7.Check deadline\t8.Check overdue\t9.Logout")
    fmt.Scanln(&op)
    switch op {
      case 1:
        var title, author, ISBN string
        fmt.Println("Type in the title, author, ISBN:")
        fmt.Scanln(&title)
        fmt.Scanln(&author)
        fmt.Scanln(&ISBN)
        lib.AddBook(title, author, ISBN)
      case 2:
        var com string
        bid:=1
        fmt.Println("Type in the commit and the book id:")
        fmt.Scanln(&com)
        fmt.Scanln(&bid)
        lib.RemoveBook(com, bid)
      case 3:
        uid:=0
        fmt.Println("Type in the user id:")
        fmt.Scanln(&uid)
        lib.AddUser(uid)
      case 4:
        op:=0
        fmt.Println("0.Title 1.Author 2.ISBN")
        fmt.Scanln(&op)
        lib.QueryBook(op)
      case 5:
        var uid int
        fmt.Println("Please enter the user id.")
        fmt.Scanln(&uid)
        lib.QueryNotReturned(uid)
      case 6:
        var uid int
        fmt.Println("Please enter the user id.")
        fmt.Scanln(&uid)
        lib.QueryBorrow(uid)
      case 7:
        var uid int
        fmt.Println("Please enter the user id.")
        fmt.Scanln(&uid)
        lib.CheckDDL(uid)
      case 8:
        var uid int
        fmt.Println("Please enter the user id.")
        fmt.Scanln(&uid)
        lib.CheckDue(uid)
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
