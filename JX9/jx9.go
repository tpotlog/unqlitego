package JX9
import (
	"fmt"
	"../../unqlitego"
)

/*
Build and generate JX9 scripts for unqlite database
*/

type JX9_script struct{
	script string
}

type JX9_json struct {
	json string
}
type JsonError string

func (err JsonError) Error() string{
	return string(err)
}

/*jx9_script code*/

func (script *JX9_script) GetScript() string{
	//Return current script of the JX9_struct struct
	return script.script
}

func (script *JX9_script ) UpdateScript(update_code string){
	//Add updated code to the script
	script.script+=update_code
}

func (script *JX9_script) InitScript() {
	//Initialize the scripts, could be also used to flush
	script.script=""
}

func NewScript() JX9_script{
	script:=JX9_script{}
	script.InitScript()
	return script
}

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

func (script *JX9_script) Compile(database unqlitego.Database) (error,string,unqlitego.VM){
	vm:=unqlitego.NewVM()
	err,out:=database.Unqlite_compile(script.GetScript(),vm)
	return err,out,*vm

}


func (script JX9_script) CompileAndExecute(database unqlitego.Database) (error,string,string,unqlitego.VM){
	/*
	Complie and execute a JX9 script returned values are
	*)error - if error occured during coplialtion of the script return the error code
	*)string -if error occured during copliation , return the error message
	*)string - if copliation ended successfully the output of the script is returned
	*)unqlitego.VM - The instance of the vm used to esecute the script is returned
	 */
	err,out,vm:=script.Compile(database)
	vm.VM_exacute()
	return err,out,vm.VM_extract_output(),vm
}

func (script *JX9_script) StoreJson(database_name,json_code string){
	//Store A json or A json list
	f:=fmt.Sprintf(`$j=%s;
	$res=db_store(%s,$j);
	if ( !$res ){
		print db_errlog();
		retrun;
	}
	`,json_code,database_name)
	script.UpdateScript(f)
}

func (script *JX9_script) GetAllFromDatatBase(database_name,variable_name string){
	f:=fmt.Sprintf(`$%s=db_fetch_all('%s');
	`,variable_name,database_name)
	script.UpdateScript(f)
}

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

func (script *JX9_script) FetchJsonByID(database_name string,record_id int64,variable_name string){
	f:=fmt.Sprintf(`$%s=db_fetch_by_id('%s',%d);
	`,variable_name,database_name,record_id)
	script.UpdateScript(f)

}

func (script *JX9_script) DropCollection(database_name string){
	f:=fmt.Sprintf(`db_drop_collection('%s');
	`,database_name)
	script.UpdateScript(f)
}

func (script *JX9_script) DeleteRecord(database_name string,record_id int64,result_variable string){
	f:=fmt.Sprintf(`$%s=db_drop_record('%s',%d);
	`,result_variable,database_name,record_id)
	script.UpdateScript(f)
}

func (script *JX9_script) UpdateRecord(database_name string ,record_id int64,new_json string){
	script.DeleteRecord(database_name,record_id,"delete_record_status")
	script.UpdateScript(`if ( !$delete_record_status ){
	print db_errlog();
	return ;
	}
	`)
	script.StoreJson(database_name,new_json)
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

