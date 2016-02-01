package main


import (
	"fmt"
	"../../unqlitego"
	"../JX9"
)

func main(){
/*	jx:=`if( !db_exists('users') ){
   $rc = db_create('users');
   if ( !$rc ){
     //Handle error
      print db_errlog();
   return;
   }
}`*/
	vm:=unqlitego.NewVM()
	db,err:=unqlitego.NewDatabase("/tmp/unqlite.db")
	jx9script:=JX9.NewScript()
	jx9script.CreateOpenDataBase("x","y")
	res,out:=db.Unqlite_compile(jx9script.GetScript(),vm)
	fmt.Printf("%s:%s" ,res,out)
	if err != nil {
		fmt.Printf("%s", err)

	}

		db.Close()


}