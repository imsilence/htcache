package rocksdb

// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../libs -L${SRCDIR}/../../../rocksdb -lrocksdb -lsnappy -lbz2 -lz -lm -lstdc++ -ldl
// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"errors"
	"htcache/cache"
	"regexp"
	"runtime"
	"strconv"
	"unsafe"
	"time"
)

type Cache struct {
	db *C.rocksdb_t
	ro *C.rocksdb_readoptions_t
	wo *C.rocksdb_writeoptions_t
	e  *C.char
}

func New(ttl time.Duration) (cache.Cache, error) {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char
	// db := C.rocksdb_open(options, C.CString("rocksdb.db"), &e)
	db := C.rocksdb_with_ttl(options, C.Cstring("rocksdb.db"), C.int(ttl), &e)
	if e != nil {
		return nil, errors.New(C.GoString(e)))
	}
	C.rocksdb_options_destroy(options)
	return &Cache{db, C.rocksdb_readoptions_create(), C.rocksdb_writeoptions_create(), e}, nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	var length C.size_t

	v := C.rocksdb_get(c.db, c.ro, k, C.size_t(len(key)), &length, &c.e)
	if c.e != nil {
		return nil, errors.New(C.GoString(c.e))
	}
	defer C.free(unsafe.Pointer(v))
	return C.GoBytes(unsafe.Pointer(v), C.int(length)), nil

}

func (c *Cache) Set(key string, value []byte) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	v := C.CBytes(value)
	defer C.free(v)

	C.rocksdb_put(c.db, c.wo, k, C.size_t(len(key)), (*C.char)(v), C.size_t(len(value)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))
	}
	return nil

}

func (c *Cache) Del(key string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	C.rocksdb_delete(c.db, c.wo, k, C.size_t(len(key)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))
	}
	return nil
}

func (c *Cache) GetStat() cache.Stat {
	var s cache.Stat

	k := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(k))

	v := C.rocksdb_property_value(c.db, k)
	defer C.free(unsafe.Pointer(v))

	value := C.GoString(v)
	reg := regexp.MustCompile(`([^;]+)=([^;]+);`)
	for _, matchers := range reg.FindAllStringSubmatch(value, -1) {
		switch matchers[1] {
		case " # entries":
			s.Count, _ = strconv.ParseInt(matchers[2], 10, 64)
		case " raw key size":
			s.KeySize, _ = strconv.ParseInt(matchers[2], 10, 64)
		case " raw value size":
			s.ValueSize, _ = strconv.ParseInt(matchers[2], 10, 64)

		}
	}
	return s
}