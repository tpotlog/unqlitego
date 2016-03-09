package main


import (
	"fmt"
	"../../unqlitego"
	"../collections"
)

func main(){

	db,err:=unqlitego.NewDatabase("/tmp/unqlite.db")
	if err!=nil{
		fmt.Printf("%s", err)

	}else {
		collection := collections.NewCollection("users", db)
		collection2 := collections.NewCollection("books", db)
		collection.SaveRecords("{\"tal\":123}")
		collection2.SaveRecords("[{\"a\":{\"b\":true},\"L\":1.1}]")
		_,W:=collection.GetAll()
		fmt.Printf("%s  \n" ,W)
		_,Z:=collection2.GetRecordByID(0)
		fmt.Printf("%s  \n" ,Z)
	}

	db.Close()


}