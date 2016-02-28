package JX9
import (
	"fmt"
//	"encoding/json"
//	"strings"
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
		}`, database_name,pointer_name,database_name,pointer_name)


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


