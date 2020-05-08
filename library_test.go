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
