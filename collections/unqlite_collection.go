package collections

import (
	"../../unqlitego"
	"../JX9"

)

type UnqliteCollectionError string

func (err UnqliteCollectionError) Error() string{
	return string(err)
}

type UnqliteCollection struct {
	name string
	database *unqlitego.Database
	script *JX9.JX9_script
	_auto_commit bool
	_auto_create bool
	_commited bool
	_last_vm *unqlitego.VM
	_last_error string
	_last_output string


}


func NewCollection(name string,database *unqlitego.Database) *UnqliteCollection{
	collection:=&UnqliteCollection{name,database,JX9.NewJX9Script(),true,true,false,nil,"",""}
	collection.Flush()
	return collection
}

func (collection *UnqliteCollection) GetAutoCreate() bool{
	return collection._auto_create
}

func (collection *UnqliteCollection) GetAutoCommit() bool  {
	return collection._auto_commit
}

func (collection *UnqliteCollection) IsCommited() bool  {
	return collection._commited
}

func (collection *UnqliteCollection) GetName() string {
	return collection.name
}

func (collection *UnqliteCollection) Flush(){
	collection._commited=false
	collection._last_vm=nil
	collection._last_error=""
	collection._last_output=""
	collection.script.InitScript()
	if collection._auto_create{
		collection.script.CreateOpenDataBase(collection.name,"ptr")
	}
}

func (collection *UnqliteCollection) Commit() (error) {
	if collection._commited{
		return UnqliteCollectionError("This collection data is marked as commited")
	}
	err,res,out,vm:=collection.script.CompileAndExecute(*collection.database)
	collection._commited=true
	collection._last_vm=&vm
	collection._last_output=out
	collection._last_error=res
	if err!=nil {
		return UnqliteCollectionError(res)
	}
	return nil
}

func (collection *UnqliteCollection) getVariableFromLastExecution(variable_name ,variable_type string) interface{}{
	if collection._commited!=false && collection._last_vm!=nil{
		var res interface{}
		var err  error
		switch variable_name {
		case "int":
				res,err=collection._last_vm.Extract_variable_as_int(variable_name)
		case "int64":
				res,err=collection._last_vm.Extract_variable_as_int64(variable_name)
		case "string":
				res,err=collection._last_vm.Extract_variable_as_int64(variable_name)
		case "bool":
				res,err=collection._last_vm.Extract_variable_as_bool(variable_name)
		case "double":
				res,err=collection._last_vm.Extract_variable_as_bool(variable_name)
		default:
			return nil
		}
		if err!=nil{
			return res
		}
	}

	return nil
}

func (collection *UnqliteCollection) SaveJson(json_code string ) error{
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.StoreJson(collection.name,json_code)
	if collection._auto_commit{
		res=collection.Commit()
		return res
	}
	return nil
}