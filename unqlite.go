package unqlitego

/*
#cgo linux CFLAGS: -DUNQLITE_ENABLE_THREADS=1 -Wno-unused-but-set-variable
#cgo darwin CFLAGS: -DUNQLITE_ENABLE_THREADS=1
#cgo windows CFLAGS: -DUNQLITE_ENABLE_THREADS=1
#include "./unqlite.h"
#include "./wrappers.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

// UnQLiteError ... standard error for this module

type GlobaLError string

func (s GlobaLError) Error() string {
	return string(s)
}

type UnQLiteError int

func (e UnQLiteError) Error() string {
	s := errString[e]
	if s == "" {
		return fmt.Sprintf("errno %d", int(e))
	}
	return s
}

var errString = map[UnQLiteError]string{
	C.UNQLITE_LOCKERR:        "Locking protocol error",
	C.UNQLITE_READ_ONLY:      "Read only Key/Value storage engine",
	C.UNQLITE_CANTOPEN:       "Unable to open the database file",
	C.UNQLITE_FULL:           "Full database",
	C.UNQLITE_VM_ERR:         "Virtual machine error",
	C.UNQLITE_COMPILE_ERR:    "Compilation error",
	C.UNQLITE_DONE:           "Operation done", // Not an error.
	C.UNQLITE_CORRUPT:        "Corrupt pointer",
	C.UNQLITE_NOOP:           "No such method",
	C.UNQLITE_PERM:           "Permission error",
	C.UNQLITE_EOF:            "End Of Input",
	C.UNQLITE_NOTIMPLEMENTED: "Method not implemented by the underlying Key/Value storage engine",
	C.UNQLITE_BUSY:           "The database file is locked",
	C.UNQLITE_UNKNOWN:        "Unknown configuration option",
	C.UNQLITE_EXISTS:         "Record exists",
	C.UNQLITE_ABORT:          "Another thread have released this instance",
	C.UNQLITE_INVALID:        "Invalid parameter",
	C.UNQLITE_LIMIT:          "Database limit reached",
	C.UNQLITE_NOTFOUND:       "No such record",
	C.UNQLITE_LOCKED:         "Forbidden Operation",
	C.UNQLITE_EMPTY:          "Empty record",
	C.UNQLITE_IOERR:          "IO error",
	C.UNQLITE_NOMEM:          "Out of memory",
}


// Database ...
type Database struct {
	handle *C.unqlite
}

// Cursor ...
type Cursor struct {
	parent *Database
	handle *C.unqlite_kv_cursor
}

type VM struct {
	vm *C.unqlite_vm
}

type Unqlite_value struct {
	unqlite_value *C.unqlite_value
}

func init() {
	C.unqlite_lib_init()
	if !IsThreadSafe() {
		panic("unqlite library was not compiled for thread-safe option UNQLITE_ENABLE_THREADS=1")
	}
}



// NewDatabase ...
func NewDatabase(filename string) (db *Database, err error) {
	db = &Database{}
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))
	res := C.unqlite_open(&db.handle, name, C.UNQLITE_OPEN_CREATE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	if db.handle != nil {
		runtime.SetFinalizer(db, (*Database).Close)
	}
	return
}

func NewVM() (vm *VM) {
	vm = &VM{}
	return
}

func (db *Database,) Unqlite_compile(jx9_script string,vm *VM) (error,string ){
	res := C.unqlite_compile(db.handle,C.CString(jx9_script),C.int(len(jx9_script)),&vm.vm)
	if res != C.UNQLITE_OK{
		if res == C.UNQLITE_COMPILE_ERR{
			err:=UnQLiteError(res)
			error_log:=new(C.char)
			err_msg:=C.extract_unqlite_log_error(db.handle,error_log)
			g_err_msg:=C.GoString(err_msg)
			C.free(unsafe.Pointer(err_msg))
			return err,g_err_msg
		}
	}
	return UnQLiteError(res),""
}

func (vm *VM)Unqlite_vm_extract_variable(variable_name string) (*Unqlite_value){
	/*This function must be used with extra causion since it might return
	a variable from the type of *C.unqlite_value ,be sure to free this pointer
	In case of no such variable of out-of-memory issue NULL is returned

	In case where the VM have not been executed the return value will be C.int(0)
	In case of unqlite is compiled with threads support and the vm.vm instance have been releases
	by a different thread 0 will be returned

	For summary:
	-----------
	*) 0 or NULL = Bad
	*) *C.unqlite_value = Good */
	c_variable_name := C.CString(variable_name)
	defer C.free(unsafe.Pointer(c_variable_name))
	unqlite_value:=C.unqlite_vm_extract_variable(vm.vm,c_variable_name)
	return &Unqlite_value{unqlite_value}
}


func unqlite_value_ok(unqlite_value *Unqlite_value)(bool) {
	switch unqlite_value.unqlite_value{
	case nil:return false//User Data is wrtong
	default:return true
	}
}

func (vm *VM)Extract_variable_as_int(variable_name string) (int,error){
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	 */
	var unqlite_value *Unqlite_value
	unqlite_value=vm.Unqlite_vm_extract_variable(variable_name)
	if ! unqlite_value_ok(unqlite_value){
		return 0,nil
	}

	defer C.free(unsafe.Pointer(unqlite_value.unqlite_value))
	res:=int(C.unqlite_value_to_int(unqlite_value.unqlite_value))
	return res,GlobaLError("OK")
}

func (vm *VM)Extract_variable_as_string(variable_name string) (string,error){
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	 */
	var unqlite_value *Unqlite_value
	unqlite_value=vm.Unqlite_vm_extract_variable(variable_name)
	if ! unqlite_value_ok(unqlite_value){
		return "",nil
	}

	defer C.free(unsafe.Pointer(unqlite_value.unqlite_value))
	var plen *C.int
	c_res:=C.unqlite_value_to_string(unqlite_value.unqlite_value,plen)
	res:=C.GoStringN(c_res,*plen)
	defer C.free(unsafe.Pointer(c_res))
	defer C.free(unsafe.Pointer(plen))
	return res,GlobaLError("OK")
}

func (vm *VM)Extract_variable_as_bool(variable_name string) (bool,error){
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	 */
	var unqlite_value *Unqlite_value
	unqlite_value=vm.Unqlite_vm_extract_variable(variable_name)
	if ! unqlite_value_ok(unqlite_value){
		return false,nil
	}

	defer C.free(unsafe.Pointer(unqlite_value.unqlite_value))
	res:=int(C.unqlite_value_to_bool(unqlite_value.unqlite_value))
	return res!=0,GlobaLError("OK")
}

func (vm *VM)Extract_variable_as_int64(variable_name string) (int64,error){
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	 */
	var unqlite_value *Unqlite_value
	unqlite_value=vm.Unqlite_vm_extract_variable(variable_name)
	if ! unqlite_value_ok(unqlite_value){
		return 0,nil
	}

	defer C.free(unsafe.Pointer(unqlite_value.unqlite_value))
	res:=int64(C.unqlite_value_to_int64(unqlite_value.unqlite_value))
	return res,GlobaLError("OK")
}

func (vm *VM)Extract_variable_as_double(variable_name string) (float64,error){
	/*Extract a variable from the VM after if have been executed
	If something went wrong return nil
	 */
	var unqlite_value *Unqlite_value
	unqlite_value=vm.Unqlite_vm_extract_variable(variable_name)
	if ! unqlite_value_ok(unqlite_value){
		return 0.0,nil
	}

	defer C.free(unsafe.Pointer(unqlite_value.unqlite_value))
	res:=float64(C.unqlite_value_to_double(unqlite_value.unqlite_value))
	return res,GlobaLError("OK")
}
// Close ...
func (db *Database) Close() (err error) {
	if db.handle != nil {
		res := C.unqlite_close(db.handle)
		if res != C.UNQLITE_OK {
			err = UnQLiteError(res)
		}
		db.handle = nil
	}
	return
}

// Store ...
func (db *Database) Store(key, value []byte) (err error) {
	var k, v unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	if len(value) > 0 {
		v = unsafe.Pointer(&value[0])
	}

	res := C.unqlite_kv_store(db.handle,
		k, C.int(len(key)),
		v, C.unqlite_int64(len(value)))
	if res == C.UNQLITE_OK {
		return nil
	}
	return UnQLiteError(res)
}

// Append ...
func (db *Database) Append(key, value []byte) (err error) {
	var k, v unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	if len(value) > 0 {
		v = unsafe.Pointer(&value[0])
	}

	res := C.unqlite_kv_append(db.handle,
		k, C.int(len(key)),
		v, C.unqlite_int64(len(value)))
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Fetch ...
func (db *Database) Fetch(key []byte) (value []byte, err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	var n C.unqlite_int64
	res := C.unqlite_kv_fetch(db.handle, k, C.int(len(key)), nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}
	value = make([]byte, int(n))
	res = C.unqlite_kv_fetch(db.handle, k, C.int(len(key)), unsafe.Pointer(&value[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Delete ...
func (db *Database) Delete(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_delete(db.handle, k, C.int(len(key)))
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Begin ...
func (db *Database) Begin() (err error) {
	res := C.unqlite_begin(db.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Commit ...
func (db *Database) Commit() (err error) {
	res := C.unqlite_commit(db.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Rollback ...
func (db *Database) Rollback() (err error) {
	res := C.unqlite_rollback(db.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// NewCursor ...
func (db *Database) NewCursor() (cursor *Cursor, err error) {
	cursor = &Cursor{parent: db}
	res := C.unqlite_kv_cursor_init(db.handle, &cursor.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	runtime.SetFinalizer(cursor, (*Cursor).Close)
	return
}

// Close ...
func (curs *Cursor) Close() (err error) {
	if curs.parent.handle != nil && curs.handle != nil {
		res := C.unqlite_kv_cursor_release(curs.parent.handle, curs.handle)
		if res != C.UNQLITE_OK {
			err = UnQLiteError(res)
		}
		curs.handle = nil
	}
	return
}

// Seek ...
func (curs *Cursor) Seek(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(curs.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_EXACT)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// SeekLE ...
func (curs *Cursor) SeekLE(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(curs.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_LE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// SeekGE ...
func (curs *Cursor) SeekGE(key []byte) (err error) {
	var k unsafe.Pointer

	if len(key) > 0 {
		k = unsafe.Pointer(&key[0])
	}

	res := C.unqlite_kv_cursor_seek(curs.handle, k, C.int(len(key)), C.UNQLITE_CURSOR_MATCH_GE)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// First ...
func (curs *Cursor) First() (err error) {
	res := C.unqlite_kv_cursor_first_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Last ...
func (curs *Cursor) Last() (err error) {
	res := C.unqlite_kv_cursor_last_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// IsValid ...
func (curs *Cursor) IsValid() (ok bool) {
	return C.unqlite_kv_cursor_valid_entry(curs.handle) == 1
}

// Next ...
func (curs *Cursor) Next() (err error) {
	res := C.unqlite_kv_cursor_next_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Prev ...
func (curs *Cursor) Prev() (err error) {
	res := C.unqlite_kv_cursor_prev_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Delete ...
func (curs *Cursor) Delete() (err error) {
	res := C.unqlite_kv_cursor_delete_entry(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Reset ...
func (curs *Cursor) Reset() (err error) {
	res := C.unqlite_kv_cursor_reset(curs.handle)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Key ...
func (curs *Cursor) Key() (key []byte, err error) {
	var n C.int
	res := C.unqlite_kv_cursor_key(curs.handle, nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}
	key = make([]byte, int(n))
	res = C.unqlite_kv_cursor_key(curs.handle, unsafe.Pointer(&key[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Value ...
func (curs *Cursor) Value() (value []byte, err error) {
	var n C.unqlite_int64
	res := C.unqlite_kv_cursor_data(curs.handle, nil, &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
		return
	}
	value = make([]byte, int(n))
	res = C.unqlite_kv_cursor_data(curs.handle, unsafe.Pointer(&value[0]), &n)
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// Shutdown ...
func Shutdown() (err error) {
	res := C.unqlite_lib_shutdown()
	if res != C.UNQLITE_OK {
		err = UnQLiteError(res)
	}
	return
}

// IsThreadSafe ...
func IsThreadSafe() bool {
	return C.unqlite_lib_is_threadsafe() == 1
}

// Version ...
func Version() string {
	return C.GoString(C.unqlite_lib_version())
}

// Signature ...
func Signature() string {
	return C.GoString(C.unqlite_lib_signature())
}

// Ident ...
func Ident() string {
	return C.GoString(C.unqlite_lib_ident())
}

// Copyright ...
func Copyright() string {
	return C.GoString(C.unqlite_lib_copyright())
}

/* TODO: implement

// Database Engine Handle
int unqlite_config(unqlite *pDb,int nOp,...);

// Key/Value (KV) Store Interfaces
int unqlite_kv_fetch_callback(unqlite *pDb,const void *pKey,
	                    int nKeyLen,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);
int unqlite_kv_config(unqlite *pDb,int iOp,...);

//  Cursor Iterator Interfaces
int unqlite_kv_cursor_key_callback(unqlite_kv_cursor *pCursor,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);
int unqlite_kv_cursor_data_callback(unqlite_kv_cursor *pCursor,int (*xConsumer)(const void *,unsigned int,void *),void *pUserData);

// Utility interfaces
int unqlite_util_load_mmaped_file(const char *zFile,void **ppMap,unqlite_int64 *pFileSize);
int unqlite_util_release_mmaped_file(void *pMap,unqlite_int64 iFileSize);

// Global Library Management Interfaces
int unqlite_lib_config(int nConfigOp,...);
*/
