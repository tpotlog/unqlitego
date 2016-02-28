

char * extract_unqlite_log_error(unqlite *pDb,char *buffer);

char * extract_vm_output(unqlite_vm *pvm,int *length);

char * extract_variable_as_string(unqlite_value *unqlite_value,int *len);