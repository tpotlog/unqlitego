package main


import (
	"fmt"
	"../../unqlitego"
	"../JX9"
)

func main(){
	/*d:=`$my_data = {
     // Greeting message
     greeting : "Hello world!\n",
     // Dummy field
      __id : 1,
     // Host Operating System
     os_name: uname(), //Built-in OS function
     // Current date
     date : __DATE__,
     // Return the current time using an anonymous function
     time : function(){ return __TIME__; }
 };

$j=[{Brown:1}];
// invoke JSON object members
print $my_data.greeting; //Hello world!
print "Host Operating System: ", $my_data.os_name, JX9_EOL;
print "Current date: ", $my_data.date, JX9_EOL;
print "Current time: ", $my_data.time(); */ // Anonymous function `
	//d:=`$j=1973.1;`
	script:=JX9.NewScript()
	script.CreateOpenDataBase("users","ptr")
	script.StoreJson("users","{\"x\":\"y\"}")
	script.StoreJson("users","[{\"VVV\":\"y\"}]")
	script.GetAllFromDatatBase("users","users")
	script.FetchJsonList("users","P","x==\"y\"")
	script.DeleteRecord("users",0,"Q")
	//vm:=unqlitego.NewVM()
	db,err:=unqlitego.NewDatabase("/tmp/unqlite.db")
	//jx9script:=JX9.NewScript()
	//jx9script.CreateOpenDataBase("x","y")
	//jx9script.UpdateScript("\n"+`$v={ "name" : 'alex', "age" : 19, "mail" : 'alex@example.com'  };`)
	//jx9script.UpdateScript("\n"+`$rc = db_store('users',$v);`)
	//jx9script.UpdateScript("\n"+`$q = "XXX";`)
	//res,out:=db.Unqlite_compile(jx9script.GetScript(),vm)
	//err,res,out,vm:=script.CompileAndExecute(*db)
	err,res,out,vm:=script.CompileAndExecute(*db)
	//err,res:=db.Unqlite_compile(script.GetScript(),vm)
	fmt.Printf("%s",script.GetScript())
	if err!=nil{
		fmt.Printf("%s", res)

	}else{
		fmt.Printf("%s\n\n",out)
		x,y:=vm.Extract_variable_as_string("P")
		fmt.Printf("\n%s:%s",x,y)

	}

	//fmt.Printf("%s",jx9script.GetScript()

		db.Close()


}