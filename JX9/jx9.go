//The JX9 Package represents a JX9 script code with the ability to compile and execute the code.
//In addtion there are code snippets to allow quick code creation for popular code snippets like
//Creating new collection(database) ,storing JSON/JSON list in a collection , removing JSON , removing collections...
package JX9
import (
	"fmt"
	"../../unqlitego"
)


//The JX9 script struct
type JX9_script struct{
	script string //The script code
}

//Create new JX9_script and initialize it.
func NewJX9Script() *JX9_script{
	s:=&JX9_script{}
	s.InitScript()
	return s
}

//Get the script Code
func (script *JX9_script) GetScript() string{
	//Return current script of the JX9_struct struct
	return script.script
}

//Append new code to the current script
func (script *JX9_script ) UpdateScript(update_code string){
	script.script+=update_code
}

//Initialize the script (flush it to an empty string)
func (script *JX9_script) InitScript() {
	script.script=""
}

/*	Add a code snippet create a database (collection), if the databse does not exists it will be automaticlly created
		database_name:The name of the databse (collection) to create
		pointer_name:The name of the JX9 variable which will hold the databse creation result
*/
func (script *JX9_script) CreateOpenDataBase(database_name,pointer_name string)  {
	//Create JX9 code add an open pointer to the database
	code := fmt.Sprintf(`if (! db_exists('%s')) {
			$%s=db_create('%s');
			if ( !$%s ){
			print db_errlog();
			return;
			}
		}
	`, database_name,pointer_name,database_name,pointer_name)


	script.UpdateScript(code)
}

/*	Complie the jx9 script code to bytecompile code which lkate could be executed
		databse:The databse struct which this code will compile aginst

	The fuction will return the following
		1.error:If compilation error occured it will be returned the error string please check , unqlitego.errString.
		If compilation suceeded than nil is returned.
		2.string:If a compilation failed the error log will be returned, if compilation suceeded this should be ingnored
		3.unqlitego.VM:The VM on which this code was compiled at (if compilation suceeded)
*/
func (script *JX9_script) Compile(database unqlitego.Database) (error,string,unqlitego.VM){
	vm:=unqlitego.NewVM()
	err,out:=database.Unqlite_compile(script.GetScript(),vm)
	return err,out,*vm

}

/*	Complie the JX9 script code and execute it
		databse:The databse struct which this code will compile aginst

	The function will return the following
		1,error:if error occured during coplialtion of the script return the error code
		2.string:if error occured during copliation , return the error message
		3.string:if copliation ended successfully the output of the script is returned
		4.unqlitego.VM:The instance of the vm used to execute the script is returned
*/
func (script JX9_script) CompileAndExecute(database unqlitego.Database) (error,string,string,unqlitego.VM){
	err,out,vm:=script.Compile(database)
	vm.VM_exacute()
	return err,out,vm.VM_extract_output(),vm
}

/*	This will add a JX9 snippet script to the current script which stores a JSON/JSON list to a database.
		database_name:The name of the database.
		json_code:The JSON/JSON list to be added to the databse.
		variable_name:The name of the variable which will hold the result of the db_store function result.
			true: in case that adding a record ended successfully.
			false: in case that adding a record failed.
*/
func (script *JX9_script) StoreJson(database_name,json_code,variable_name string){
	//Store A json or A json list
	f:=fmt.Sprintf(`$json_code=%s;
	$%s=db_store(%s,$json_code);
	if ( !$%s ){
		print db_errlog();
	}
	`,json_code,variable_name,database_name,variable_name)
	script.UpdateScript(f)
}
/*	Add code snippet to the JX9 script which get all avaliable records from the databse.
		database_name:The name of the databse to greb all the record from
		variable_name:The name of the variable which will hold the result of the db_fetch_all function result.
 */
func (script *JX9_script) GetAllFromDatatBase(database_name,variable_name string){
	f:=fmt.Sprintf(`$%s=db_fetch_all('%s');
	`,variable_name,database_name)
	script.UpdateScript(f)
}
/*	Add Code snippet to the JX9 script which will search for specific
		database_name:The name of the databse to greb all the record from
		variable_name:The name of the variable which will hold the result of the db_fetch_all function result.
		search_path:The JX9 json path format to be used for grabbing records:
		 	Ex:
		 		Assumiing that our databse have the following records stored
		 		1. {"x":1,"y":{"a":,2}}
		 		2. {"x":1,"y":{"a":,3}}
		 		3. {"x":1,"y":{"a":,4}}

		 		search_path="y.a>2" will return
		 		[{"x":1,"y":{"a":,3}},{"x":1,"y":{"a":,4}}]
*/
func (script *JX9_script) FetchJsonList(database_name,variable_name,search_path string){
	f:=fmt.Sprintf(`$callback=function($record){
							if ($record.%s){
									return TRUE;
								}
							return FALSE;
					};
	$%s=db_fetch_all('%s',$callback);
	`,search_path,variable_name,database_name)
	script.UpdateScript(f)
}

/*	Add Code snippet to the JX9 script which will search for specific
		database_name:The name of the databse to greb all the record from
		variable_name:The name of the variable which will hold the result of the db_fetch_all function result.
		search_path:The JX9 json path format to be used for grabbing records:
		 	Ex:
		 		Assumiing that our databse have the following records stored
		 		1. {"x":1,"y":{"a":,2}}
		 		2. {"x":1,"y":{"a":,3}}
		 		3. {"x":1,"y":{"a":,4}}

		 		search_path="y.a>2" will return
		 		[{"x":1,"y":{"a":,3}},{"x":1,"y":{"a":,4}}]
*/
func (script *JX9_script) FetchJsonByID(database_name string,record_id int64,variable_name string){
	f:=fmt.Sprintf(`$%s=db_fetch_by_id('%s',%d);
	`,variable_name,database_name,record_id)
	script.UpdateScript(f)

}

func (script *JX9_script) DropCollection(database_name,variable_name string){
	f:=fmt.Sprintf(`$%s=db_drop_collection('%s');
	` ,variable_name,database_name)
	script.UpdateScript(f)
}

func (script *JX9_script) DeleteRecord(database_name string,record_id int64,result_variable string){
	f:=fmt.Sprintf(`$%s=db_drop_record('%s',%d);
	`,result_variable,database_name,record_id)
	script.UpdateScript(f)
}

func (script *JX9_script) UpdateRecord(database_name string ,record_id int64,new_json,variable_name string){
	script.DeleteRecord(database_name,record_id,`delete_record_status`)
	script.UpdateScript(fmt.Sprintf(`$%s=-1;
	if ( !$delete_record_status ){
	print db_errlog();
	return ;
	}
	`, variable_name))
	script.StoreJson(database_name,new_json,"update_record_status")
	script.UpdateScript(fmt.Sprintf(`if ( !$update_record_status ) {
	print db_errlog();
	return ;
	}
	$%s=db_last_record_id('%s');
	`,variable_name,database_name))
}

func (script *JX9_script)  GetDatabaseTotalNumberOfRecords(databse_name ,variable_name string){
	f:=fmt.Sprintf(`$%s=db_total_records(%s)
	`,variable_name,databse_name)
	script.UpdateScript(f)

}

func (script *JX9_script) GetDatabaseVersion(variable_name string){
	f:=fmt.Sprintf(`$%s=db_version();
	`,variable_name)
	script.UpdateScript(f)
}

func (script *JX9_script) GetDatabaseSig(variable_name string){
	f:=fmt.Sprintf(`$%s=db_sig();
	`,variable_name)
	script.UpdateScript(f)
}

func (script *JX9_script) GetDatabaseCopyRight(variable_name string){
	f:=fmt.Sprintf(`$%s=db_copyright();
	`,variable_name)
	script.UpdateScript(f)
}

