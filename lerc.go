package lerc

/*
#include "Lerc_c_api.h"
#cgo CFLAGS: -I ./
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	DT_CHAR   = 0
	DT_UCHAR  = 1
	DT_SHORT  = 2
	DT_USHORT = 3
	DT_INT    = 4
	DT_UINT   = 5
	DT_FLOAT  = 6
	DT_DOUBLE = 7
)

type BlobInfo []uint32

func (t BlobInfo) Version() uint32 {
	return t[0]
}

func (t BlobInfo) DataType() uint32 {
	return t[1]
}

func (t BlobInfo) Dim() uint32 {
	return t[2]
}

func (t BlobInfo) Cols() uint32 {
	return t[3]
}

func (t BlobInfo) Rows() uint32 {
	return t[4]
}

func (t BlobInfo) Bands() uint32 {
	return t[5]
}

func (t BlobInfo) ValidPixels() uint32 {
	return t[6]
}

func (t BlobInfo) BlobSize() uint32 {
	return t[7]
}

func NewBlobInfo() BlobInfo {
	return make(BlobInfo, 10)
}

type LercStatus C.lerc_status

func GetBlobInfo(b []byte) (BlobInfo, error) {
	binfo := NewBlobInfo()
	var dataRangeArr [3]float64

	drptr := (*C.double)(unsafe.Pointer(&dataRangeArr[0]))

	biptr := (*C.uint)(unsafe.Pointer(&binfo[0]))
	ptr := (*C.uchar)(C.CBytes(b))
	len := len(b)
	state := C.lerc_getBlobInfo(ptr, C.uint(len), biptr, drptr, 10, 3)
	if uint(state) != 0 {
		return nil, errors.New("lercDecode error")
	}
	return binfo, nil
}

func makeData(size uint32, dataType uint32) (interface{}, unsafe.Pointer) {
	switch dataType {
	case DT_CHAR:
		dt := make([]int8, size)
		return dt, unsafe.Pointer(&dt[0])
	case DT_UCHAR:
		dt := make([]uint8, size)
		return dt, unsafe.Pointer(&dt[0])
	case DT_SHORT:
		dt := make([]int16, size)
		return dt, unsafe.Pointer(&dt[0])
	case DT_USHORT:
		dt := make([]uint16, size)
		return dt, unsafe.Pointer(&dt[0])
	case DT_INT:
		dt := make([]int32, size)
		return dt, unsafe.Pointer(&dt[0])
	case DT_UINT:
		dt := make([]uint32, size)
		return dt, unsafe.Pointer(&dt[0])
	case DT_FLOAT:
		dt := make([]float32, size)
		return dt, unsafe.Pointer(&dt[0])
	case DT_DOUBLE:
		dt := make([]float64, size)
		return dt, unsafe.Pointer(&dt[0])
	}
	return nil, nil
}

func lercDecode(b []byte, info BlobInfo) (interface{}, []byte, error) {
	ptr := (*C.uchar)(C.CBytes(b))
	len := len(b)
	mask := make([]byte, info.Cols()*info.Rows())
	raster, uptr := makeData(info.Cols()*info.Rows(), info.DataType())

	state := C.lerc_decode(ptr, C.uint(len), (*C.uchar)(unsafe.Pointer(&mask[0])),
		C.int(info.Dim()), C.int(info.Cols()), C.int(info.Rows()), C.int(info.Bands()),
		C.uint(info.DataType()), uptr)

	if uint(state) != 0 {
		return nil, nil, errors.New("lercDecode error")
	}
	return raster, mask, nil
}

func Decode(b []byte) (interface{}, []byte, error) {
	binfo, err := GetBlobInfo(b)
	if err != nil {
		return nil, nil, err
	}
	return lercDecode(b, binfo)
}

func getDataType(size interface{}) (uint32, unsafe.Pointer, error) {
	switch t := size.(type) {
	case []int8:
		return DT_CHAR, unsafe.Pointer(&t[0]), nil
	case []uint8:
		return DT_UCHAR, unsafe.Pointer(&t[0]), nil
	case []int16:
		return DT_SHORT, unsafe.Pointer(&t[0]), nil
	case []uint16:
		return DT_USHORT, unsafe.Pointer(&t[0]), nil
	case []int32:
		return DT_INT, unsafe.Pointer(&t[0]), nil
	case []uint32:
		return DT_UINT, unsafe.Pointer(&t[0]), nil
	case []float32:
		return DT_FLOAT, unsafe.Pointer(&t[0]), nil
	case []float64:
		return DT_DOUBLE, unsafe.Pointer(&t[0]), nil
	}
	return 0, nil, nil
}

func ComputeCompressedSize(b interface{}, dim int, cols int, rows int, bands int, mask []byte, maxZErr float64) (uint32, error) {
	dt, ptr, err := getDataType(b)
	if err != nil {
		return 0, err
	}
	var outlen uint32
	state := C.lerc_computeCompressedSize(ptr, C.uint(dt), C.int(dim), C.int(cols),
		C.int(rows), C.int(bands), (*C.uchar)(unsafe.Pointer(&mask[0])),
		C.double(maxZErr), (*C.uint)(unsafe.Pointer(&outlen)))
	if uint(state) != 0 {
		return 0, errors.New("lerc ComputeCompressedSize error")
	}
	return outlen, nil
}

func Encode(b interface{}, dim int, cols int, rows int, bands int, mask []byte, maxZErr float64) ([]byte, error) {
	size, err := ComputeCompressedSize(b, dim, cols, rows, bands, mask, maxZErr)
	if err != nil {
		return nil, err
	}
	dt, ptr, err := getDataType(b)
	buff := make([]byte, size)
	var outlen uint32
	state := C.lerc_encode(ptr, C.uint(dt), C.int(dim), C.int(cols),
		C.int(rows), C.int(bands), (*C.uchar)(unsafe.Pointer(&mask[0])),
		C.double(maxZErr), (*C.uchar)(unsafe.Pointer(&buff[0])),
		C.uint(size), (*C.uint)(unsafe.Pointer(&outlen)))
	if uint(state) != 0 {
		return nil, errors.New("lercEncode error")
	}
	return buff[0:outlen], nil
}
