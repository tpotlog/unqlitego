/*The wrappers are binding code for cases where cgo can not
execute or call functions directly*/

#include "./unqlite.h"
#include <stdio.h>

#define VERSION "1.0"
#ifdef __cplusplus
extern "C" {
#endif

char * extract_unqlite_log_error(unqlite *pDb)
{
    /*This function will return the number of chars at the error message
      If there is no error message than 0 is returned
    */
     int length;
     char *buffer;
    //Extract the errror if exists
    unqlite_config(pDb,UNQLITE_CONFIG_JX9_ERR_LOG,&buffer,&length);
    return buffer;
}

char * extract_vm_output(unqlite_vm *pvm,int *length)
{
    const void *buffer;
    char *t;
    //Extract the VM output
    unqlite_vm_config(pvm,UNQLITE_VM_CONFIG_EXTRACT_OUTPUT,&buffer,length);
    t=(char *)malloc(*length+1);
    memcpy(t,(const char*)buffer,*length);
    return t;

}

char * extract_variable_as_string(unqlite_value *unqlite_value,int *len){
    return unqlite_value_to_string(unqlite_value,len);
}


#ifdef __cplusplus
extern  }
#endif