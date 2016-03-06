package main


import (
	"fmt"
	"../../unqlitego"
	"../collections"
)

func main(){

	//script:=JX9.NewScript()
	/*script.CreateOpenDataBase("users","ptr")
	script.StoreJson("users","{\"x\":\"y\"}")
	script.StoreJson("users","[{\"VVV\":\"y\"}]")
	script.GetAllFromDatatBase("users","users")
	script.FetchJsonList("users","P","t.f==1")
	script.UpdateRecord("users",6,"{\"b\":1}")
	script.GetAllFromDatatBase("users","users")*/
	//script.GetDatabaseCopyRight("Q")
	db,err:=unqlitego.NewDatabase("/tmp/unqlite.db")
	//err,res,out,vm:=script.CompileAndExecute(*db)
	//fmt.Printf("%s",script.GetScript())
	if err!=nil{
		fmt.Printf("%s", err)

	}else {
		collection := collections.NewCollection("users", db)
		f,t:=collection.UpdateRecord(3,"{\"foo\":\"bar\"}")
		fmt.Printf("%s  -- %s\n" ,f,t)
		//fmt.Print("\n%s\n",collection.GetScript().GetScript())
		_,W:=collection.GetAll()
		fmt.Print("%s  \n" ,W)
	}
/*else{
		fmt.Printf("%s\n\n",out)
		x,y:=vm.Extract_variable_as_string("Q")
		fmt.Printf("\n%s:%s",x,y)

	}*/

	db.Close()


}