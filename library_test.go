package main

import (
	"testing"
)

func Test_initdb(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := initdb(&lib)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_QueryBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.QueryBook(0)
	if err != nil {
		t.Errorf("Wrong!")
	}
	err = lib.QueryBook(1)
	if err != nil {
		t.Errorf("Wrong!")
	}
	err = lib.QueryBook(2)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_AddUser(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.AddUser(1)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_AddBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.AddBook("A","B","C")
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_RemoveBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.RemoveBook("TEST!",1)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_QueryNotReturned(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.QueryNotReturned(1)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_QueryBorrow(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.QueryBorrow(1)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_CheckDDL(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.CheckDDL(1)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_BorrowBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.BorrowBook(1,7)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_ReturnBook(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.ReturnBook(3)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

func Test_Extend(t *testing.T) {
	lib := Library{}
	lib.ConnectDB()
	err := lib.Extend(1,1)
	if err != nil {
		t.Errorf("Wrong!")
	}
}

