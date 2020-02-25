package mmap

import (
	"reflect"
	"syscall"
	"unsafe"
)

type Mmap struct {
	data []byte
	size int
}

func New(size int) (*Mmap, error) {
	data, err := syscall.Mmap(-1, 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		return nil, err
	}
	return &Mmap{data: data, size: size}, nil
}

func (m *Mmap) ToByteSlice() []byte {
	return m.data
}

func (m *Mmap) ToUint64Slice() []uint64 {
	var u []uint64
	dataHeader := (*reflect.SliceHeader)(unsafe.Pointer(&m.data))
	uHeader := (*reflect.SliceHeader)(unsafe.Pointer(&u))
	uHeader.Data = dataHeader.Data
	uHeader.Len = m.size / 8
	uHeader.Cap = m.size / 8
	return u
}

func (m *Mmap) ToFloat32Slice() []float32 {
	var f []float32
	dataHeader := (*reflect.SliceHeader)(unsafe.Pointer(&m.data))
	uHeader := (*reflect.SliceHeader)(unsafe.Pointer(&f))
	uHeader.Data = dataHeader.Data
	uHeader.Len = m.size / 4
	uHeader.Cap = m.size / 4
	return f
}

func (m *Mmap) Close() {
	syscall.Munmap(m.data)
}
