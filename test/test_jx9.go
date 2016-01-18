package main

import (
	"fmt"
	"../JX9"
)

func main(){
	jx9script:=new(JX9.JX9_script)
	jx9script.InitScript()
	jx9script.CreateOpenDataBase("rc","bc")
	fmt.Printf("%s" ,jx9script.GetScript())
}