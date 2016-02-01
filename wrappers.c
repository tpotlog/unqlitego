/*The wrappers are binding code for cases where cgo can not
execute or call functions directly*/

#include "./unqlite.h"
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





#ifdef __cplusplus
extern  }
#endif