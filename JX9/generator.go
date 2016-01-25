package JX9
import "fmt"

/*
Build and generate JX9 scripts for unqlite database
*/

type JX9_script struct{
	script string
}

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



