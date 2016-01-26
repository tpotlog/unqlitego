package main


import (
	"fmt"
	"../../unqlitego"
	"../JX9"
)

func main(){
	jx:=`if( !db_exists('users') ){
    /* Try to create it */
   $rc = db_create('users');
   if ( !$rc ){
     //Handle error
      print db_errlog()
   return;
   }
}`
	vm:=unqlitego.NewVM()
	db,err:=unqlitego.NewDatabase("/tmp/unqlite.db")
	jx9script:=JX9.NewScript()
	jx9script.CreateOpenDataBase("x","y")
//	res:=db.Unqlite_compile(jx9script.GetScript(),vm)
	res:=db.Unqlite_compile(jx,vm)

	fmt.Printf("%s" ,res)
	if err != nil {
		fmt.Printf("%s", err)

	} else {
		a:=[]byte("x")
		b:=[]byte("y")
		db.Store(a,b)
		v,err:=db.Fetch([]byte("z"))
		if err != nil {
			fmt.Printf("%s", err)
		}else{
			fmt.Printf("%s" ,v)
		}

		db.Close()
	}

}