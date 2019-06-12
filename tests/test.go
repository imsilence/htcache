package main

// #include "test.h"
// #cgo LDFLAGS: libtest.a
import "C"

func main() {
	C.test()
}
