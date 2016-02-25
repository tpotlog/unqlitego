package JX9
import (
	"fmt"
	"encoding/json"
	"strings"
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

func (script *JX9_script) AddJx9JsonVaraible  (jx9_json JX9_json,variable_name string){
	if ! strings.HasPrefix(variable_name,"$"){
		variable_name="$"+variable_name
	}
	script.script=fmt.Sprintf("%s\n%s=%s;\n",script.script,variable_name,jx9_json.GetJson())
}

func (script *JX9_script) AddJx9JsonListVaraible  (json_list []JX9_json,variable_name string){
	if ! strings.HasPrefix(variable_name,"$"){
		variable_name="$"+variable_name
	}
	script.script=fmt.Sprintf("%s\n%s=%s;\n",script.script,variable_name,JX9_json_list_convertor(json_list))
}

/*JX9_json code*/
func as_jx9_json(j map[string]interface{},jx9_json *string,indent string){

	_as_jx9_json(j,jx9_json,indent)
	*jx9_json+=fmt.Sprintf("%s}",indent)
}

func _as_json_type(t interface{}) string{
	switch t.(type){
		case bool:return fmt.Sprintf("%t",t.(bool))
		case string:return fmt.Sprintf("\"%s\"" ,t.(string))
		case float64:return fmt.Sprintf("%f",t.(float64))
		case nil:return "null"
	}

	return "null"
}
func _as_jx9_json(j map[string]interface{},jx9_json *string,indent string){
	*jx9_json+=fmt.Sprintf(" {%s\n",indent)
	t:=0
	for x,y := range j{
		splitter:=""
		if t<(len(j)-1) {
			splitter = ","
		}else {
			splitter=""
		}
		t+=1
		switch y.(type){
			case []interface{}:
				*jx9_json+=fmt.Sprintf(" %s%s:[\n",indent,x)
				for p,r:= range y.([]interface{}){
					marker:=""

					if p<(len(y.([]interface{}))-1){
						marker=","
					}else{
						marker=""
					}
					switch r.(type){
						case map[string]interface{}:
							tmp:=""
							tmp_indent:="    "+indent
							as_jx9_json(r.(map[string]interface{}),&tmp,tmp_indent)
							*jx9_json+=fmt.Sprintf("   %s%s%s\n",indent,tmp,marker)
						default:*jx9_json+=fmt.Sprintf("   %s%s%s\n",indent,_as_json_type(r),marker)
					}
				}
				*jx9_json+=fmt.Sprintf(" %s]%s\n",indent,splitter)
			case map[string]interface{}:
				*jx9_json+=fmt.Sprintf(" %s%s:",indent,x)
				indent+="      "
				_as_jx9_json(y.(map[string]interface{}),jx9_json,indent)
				*jx9_json+=fmt.Sprintf("%s}%s\n",indent,splitter)
			default:*jx9_json+=fmt.Sprintf(" %s%s:%s%s\n",indent,x,_as_json_type(y),splitter)

		}
	}
}

func (jx9_json JX9_json) GetJson() string{
	return jx9_json.json
}

func (jx9_json *JX9_json) ConvertToJx9Json(json_data interface{},indent string) error {

	//Convert json to jx9 jeson format

	switch json_data.(type){
		case  []byte:
				var js map[string]interface{}
				err:=json.Unmarshal(json_data.([]byte), &js)
				if err!=nil {
					return err
				}
				as_jx9_json(js,&jx9_json.json,indent)
		case map[string]interface{}:
				as_jx9_json(json_data.(map[string]interface{}),&jx9_json.json,"")
		default: return JsonError("json_data should be from the types of []byte or map[string]interface{}")
	}
	return nil
}

func ConvertToJsonList(json_list_data []interface{},indent string) ([]JX9_json,error) {
	json_list := make([]JX9_json, len(json_list_data))
	for idx, elem := range json_list_data {
		j := JX9_json{""}
		err := j.ConvertToJx9Json(elem,indent)
		if err != nil {
			return nil, JsonError(fmt.Sprintf("Element %d is not a valid json [%s]", idx, err))
		}
		json_list[idx] = j

	}
	return json_list,nil
}

func JX9_json_list_convertor(json_list []JX9_json) string{
	out:="%s[\n"
	for _,j := range json_list{
		out+=j.GetJson()

	}
	out+="\n]"
	return out
}


/*automatic JX9 scripting to handle json manipulations over the database*/


func AddJx9JsonToDatabse(database *unqlitego.Database,database_name string,jx9_json JX9_json) (error,string){
	jx9_script:=NewScript()
	jx9_script.CreateOpenDataBase(database_name,"db")
	jx9_script.AddJx9JsonVaraible(jx9_json,"j")
	jx9_script.UpdateScript(fmt.Sprintf(`res = db_store(%s,$j);
	if( !$res ){
    print db_errlog();
  	return;
	}` ,database_name))
	vm:=unqlitego.NewVM()
	vm.Free()
	return database.Unqlite_compile(jx9_script.GetScript(),vm)
}


