package codec

import (
	"fmt"
	"math"
)

const size = 64
const top = uint64(1) << (size - 8)
const bot = uint64(1) << (size - 16)

func Test() {
	numVals := 8

	val := []uint32{0, 1, 2, 3, 0, 1, 0, 0, 0, 0, 4, 5, 0, 7, 0, 0}

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

	// fp, err := os.Create("compressed")
	// if err != nil {
	// 	panic(err)
	// }

	var rcrange, low, code uint64

	rcrange = math.MaxUint64
	low = 0
	code = 0

	for _, v := range val {
		rcrange /= uint64(cumfreq[numVals])
		low += uint64(cumfreq[v]) * rcrange
		rcrange *= uint64(freq[v])

		for (low ^ (low + rcrange)) < top {
			stream := low >> (size - 8)
			fmt.Printf("%b", stream)
			code += 8
			if code > 1e8 {
				panic("L T")
			}
			rcrange <<= 8
			low <<= 8
		}
		for rcrange < bot {
			stream := low >> (size - 8)
			fmt.Printf("%b", stream)
			code += 8
			if code > 1e8 {
				panic("L T")
			}
			rcrange = ((-low) & (bot - 1)) << 8
			low <<= 8
		}
	}

	for i := 0; i < size; i++ {
		stream := low >> (size - 8)
		if stream != 0 {
			fmt.Printf("%b", stream)
		}
		code += 8
		low <<= 8
	}
}
