package lerc

/*
#include "Lerc_c_api.h"
#cgo CFLAGS: -I ./
*/
import "C"

type LercStatus C.lerc_status

func lercDecode(b []byte, dim, cols, rows, bands int, dataType uint32) ([]byte, error) {
	return nil, nil
}