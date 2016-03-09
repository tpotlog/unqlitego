/*	The collections package represents a high level approach to use unqlite trough collections
	which allow reading/writing/deleting/updating documents (JSONs) to and from the unqlite database.
 */
package collections

import (
	"../../unqlitego"
	"../JX9"
	"fmt"
)

//A string representing a error related to collections
type UnqliteCollectionError string

//A wrapped to allow use the UnqliteCollectionError with error interface
func (err UnqliteCollectionError) Error() string{
	return string(err)
}

//A struct representing unqlite collection
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


/*	Retrun a new collection related

 */
func NewCollection(name string,database *unqlitego.Database) *UnqliteCollection{
	collection:=&UnqliteCollection{name,database,JX9.NewJX9Script(),true,true,false,nil,"",""}
	collection.Flush()
	return collection
}

func (collection *UnqliteCollection) GetAutoCreate() bool{
	return collection._auto_create
}

func (collection *UnqliteCollection) GetScript() *JX9.JX9_script{
	return collection.script
}

func (collection *UnqliteCollection) SetAutoCreate(auto_create bool) {
	collection._auto_create=auto_create
}


func (collection *UnqliteCollection) GetAutoCommit() bool  {
	return collection._auto_commit
}

func (collection *UnqliteCollection) SetAutoCommit(auto_commit bool){
	collection._auto_commit=auto_commit
}

func (collection *UnqliteCollection) IsCommited() bool  {
	return collection._commited
}

func (collection *UnqliteCollection) GetName() string {
	return collection.name
}

func (collection *UnqliteCollection) GetLastExecutionOutput() string {
	return collection._last_output
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
	if collection._commited!=false {
		var res interface{}
		var err  error
		switch variable_type {
		case "int":
				res,err=collection._last_vm.Extract_variable_as_int(variable_name)
		case "int64":
				res,err=collection._last_vm.Extract_variable_as_int64(variable_name)
		case "string":
				res,err=collection._last_vm.Extract_variable_as_string(variable_name)
		case "bool":
				res,err=collection._last_vm.Extract_variable_as_bool(variable_name)
		case "double":
				res,err=collection._last_vm.Extract_variable_as_double(variable_name)
		default:
			return nil
		}
		if err!=nil{
			return res
		}
	}

	return nil
}

func (collection *UnqliteCollection) SaveRecords(json_code string ) (error,bool){
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.StoreJson(collection.name,json_code,"store_results")
	if collection._auto_commit {
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("store_results","bool")
			if result!=nil {
				return nil, result.(bool)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s",collection._last_error)),false
			}
		} else {
			return res,false
		}

	}
	return nil,false

}

func (collection *UnqliteCollection) GetAll() (error,string){
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.GetAllFromDatatBase(collection.name,"all_obj")
	if collection._auto_commit{
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("all_obj","string")
			if result!=nil {
				return nil, result.(string)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s",collection._last_error)),""
			}
		} else {
			return res,""
		}

	}
	//we should never get here
	return nil,""
}

func (collection *UnqliteCollection) GetTotalNumberOfRecord() (error,int64){
	{
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.GetDatabaseTotalNumberOfRecords(collection.name,"total_number_of_records")
	if collection._auto_commit{
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("total_number_of_records","int64")
			if result!=nil {
				return nil, result.(int64)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s",collection._last_error)),-1
			}
		} else {
			return res,-1
		}

	}
	//we should never get here
	return nil,-1
	}
}

func (collection *UnqliteCollection) DeleteRecordByID(record_id int64) (error,bool){
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.DeleteRecord(collection.name,record_id,"remove_result")
	if collection._auto_commit {
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("remove_result","bool")
			if result!=nil {
				return nil, result.(bool)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s",collection._last_error)),false
			}
		} else {
			return res,false
		}

	}
	return nil,false

}

func (collection *UnqliteCollection) DropCollection() (error,bool){
	var res error
	_auto_create:=collection.GetAutoCreate()
	collection.SetAutoCreate(false)
	if collection._commited{
		collection.Flush()
	}
	collection.SetAutoCreate(_auto_create)
	collection.script.DropCollection(collection.name,"drop_result")
	if collection._auto_commit {
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("drop_result","bool")
			if result!=nil {
				return nil, result.(bool)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s",collection._last_error)),false
			}
		} else {
			return res,false
		}

	}
	return nil,false

}

func (collection *UnqliteCollection) GetRecordByID(record_id int64) (error,string){
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.FetchJsonByID(collection.name,record_id,"record_by_id")
	if collection._auto_commit{
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("record_by_id","string")
			if result!=nil {
				return nil, result.(string)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s",collection._last_error)),""
			}
		} else {
			return res,""
		}

	}
	//we should never get here
	return nil,""
}

func (collection *UnqliteCollection) FetchRecords(path string) (error,string){
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.FetchJsonList(collection.name,"records_by_path",path)
	if collection._auto_commit{
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("records_by_path","string")
			if result!=nil {
				return nil, result.(string)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s",collection._last_error)),""
			}
		} else {
			return res,""
		}

	}
	//we should never get here
	return nil,""
}

func (collection *UnqliteCollection) UpdateRecord(record_id int64,json_code string ) (error,int64){
	var res error
	if collection._commited{
		collection.Flush()
	}
	collection.script.UpdateRecord(collection.name,record_id,json_code,"update_results")
	if collection._auto_commit {
		res=collection.Commit()
		if res==nil{
			result:=collection.getVariableFromLastExecution("update_results","int64")
			if result!=nil {
				return nil, result.(int64)
			} else{
				return UnqliteCollectionError(fmt.Sprintf("Variable could not be extracted ,Last Error:%s\n%s",collection._last_error,collection._last_output)),-1
			}
		} else {
			return res,-1
		}

	}
	return nil,-1

}