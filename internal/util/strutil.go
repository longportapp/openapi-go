package util

import (
	"reflect"
	"unsafe"
)

// UnsafeStringToBytes convert string to bytes without memory allocation
func UnsafeStringToBytes(s string) []byte {
	d := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&s)).Data)
	var b []byte
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	h.Data = uintptr(d)
	h.Cap = len(s)
	h.Len = len(s)
	return b
}

// UnsafeBytesToString convert bytes to string without memory allocation
func UnsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
