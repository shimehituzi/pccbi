package codec

import (
	"fmt"
	"math"
	"os"
)

const size = 64
const top = uint64(1) << (size - 8)
const bot = uint64(1) << (size - 16)

func Test() {
	numVals := 2

	val := []uint32{0, 0, 1, 0, 0}

	freq := make([]uint32, numVals)
	for _, v := range val {
		freq[v]++
	}
	cumfreq := make([]uint32, numVals+1)
	cumfreq[0] = 0
	for i := range freq {
		cumfreq[i+1] = cumfreq[i] + freq[i]
	}
	fmt.Println(freq, cumfreq)

	// エンコード

	fp, err := os.Create("compressed")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	var rcrange, low, code uint64

	rcrange = math.MaxUint64
	low = 0
	code = 0

	for _, v := range val {
		rcrange /= uint64(cumfreq[numVals])
		low += uint64(cumfreq[v]) * rcrange
		rcrange *= uint64(freq[v])

		for (low ^ (low + rcrange)) < top {
			stream := byte(low >> (size - 8))
			fp.Write([]byte{stream})
			code += 8
			if code > 1e8 {
				panic("L T")
			}
			rcrange <<= 8
			low <<= 8
		}
		for rcrange < bot {
			stream := byte(low >> (size - 8))
			fp.Write([]byte{stream})
			fmt.Printf("%b\n", stream)
			code += 8
			if code > 1e8 {
				panic("L T")
			}
			rcrange = ((-low) & (bot - 1)) << 8
			low <<= 8
		}
	}

	for i := 0; i < size; i++ {
		stream := byte(low >> (size - 8))
		if stream != 0x00 {
			fp.Write([]byte{stream})
		}
		code += 8
		low <<= 8
	}
}
