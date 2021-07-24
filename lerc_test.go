package lerc

import (
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"testing"
)

func TestLerc(t *testing.T) {
	h := 512
	w := 512
	zImg := make([]float32, w*h)
	for i := range zImg {
		zImg[i] = 0
	}

	maskByteImg := make([]byte, w*h)
	for i := range maskByteImg {
		maskByteImg[i] = 0
	}

	k := 0
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			zImg[k] = float32(math.Sqrt((float64)(i*i + j*j)))
			zImg[k] += rand.Float32()

			if j%100 == 0 || i%100 == 0 {
				maskByteImg[k] = 0
			} else {
				maskByteImg[k] = 1
			}
			k++
		}
	}

	maxZErrorWanted := 0.1
	eps := 0.0001
	maxZError := maxZErrorWanted - eps

	size, err := ComputeCompressedSize(zImg, 1, w, h, 1, maskByteImg, maxZError)
	if err != nil {
		t.Error(err)
	}

	if size == 0 {
		t.Error("error")
	}

	buff, err := Encode(zImg, 1, w, h, 1, maskByteImg, maxZError)
	if err != nil {
		t.Error(err)
	}

	if buff == nil {
		t.Error("error")
	}

	binfo, err := GetBlobInfo(buff)
	if err != nil {
		t.Error("error")
	}

	if binfo.Rows() != uint32(h) || binfo.Cols() != uint32(w) {
		t.Error("error")
	}

	f, err := os.Open("./testdata/title_13_3152_6707.atm")
	defer f.Close()

	if err != nil {
		t.Error("error")
	}
	src, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error("error")
	}

	newImg, newMask, err := Decode(src)
	if err != nil {
		t.Error("error")
	}

	zImg3 := newImg.([]float32)

	if len(newMask) != len(zImg3) {
		t.Error("error")
	}
	/**
	maxDelta := 0.0
	k = 0
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if newMask[k] != maskByteImg[k] {
				t.Error("Error in main: decoded valid bytes differ from encoded valid bytes")
			}

			if newMask[k] == 1 {
				delta := math.Abs(float64(zImg3[k] - zImg[k]))
				if delta > maxDelta {
					maxDelta = delta
				}
			}
			k++
		}
	}**/
}
