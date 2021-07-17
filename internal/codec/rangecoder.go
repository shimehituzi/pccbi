package codec

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

const SIZE = 64
const TOP = uint64(1) << (SIZE - 8)
const BOT = uint64(1) << (SIZE - 16)

const filename = "compressed"

func Test() {
	val := []uint32{0, 0, 1, 0, 0}

	max := val[0]
	for i := 1; i < len(val); i++ {
		if val[i] > max {
			max = val[i]
		}
	}

	freq := make([]uint32, max+1)
	for _, v := range val {
		freq[v]++
	}
	cumfreq := make([]uint32, max+2)
	cumfreq[0] = 0
	for i := range freq {
		cumfreq[i+1] = cumfreq[i] + freq[i]
	}

	fmt.Println(freq)
	fmt.Println(cumfreq)

	encode(freq, cumfreq, val)

	decode(freq, cumfreq, len(val))
}

func encode(freq, cumfreq, val []uint32) {
	fp, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	var rcrange, low, code uint64

	rcrange = math.MaxUint64
	low = 0
	code = 0

	for _, v := range val {
		rcrange /= uint64(cumfreq[len(cumfreq)-1])
		low += uint64(cumfreq[v]) * rcrange
		rcrange *= uint64(freq[v])

		for (low ^ (low + rcrange)) < TOP {
			w.WriteByte(byte(low >> (SIZE - 8)))
			code += 8
			if code > 1e8 {
				panic("L T")
			}
			rcrange <<= 8
			low <<= 8
		}
		for rcrange < BOT {
			w.WriteByte(byte(low >> (SIZE - 8)))
			code += 8
			if code > 1e8 {
				panic("L T")
			}
			rcrange = ((-low) & (BOT - 1)) << 8
			low <<= 8
		}
	}

	for i := 0; i < SIZE; i += 8 {
		w.WriteByte(byte(low >> (SIZE - 8)))
		code += 8
		low <<= 8
	}
	w.Flush()
}

func decode(freq, cumfreq []uint32, length int) {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	r := bufio.NewReader(fp)

	var rcrange, low, code uint64

	rcrange = math.MaxUint64
	low = 0
	code = 0

	for i := 0; i < SIZE; i += 8 {
		buf, err := r.ReadByte()
		if err != nil && err != io.EOF {
			panic(err)
		}
		code = (code << 8) | uint64(buf)
	}

	for i := 0; i < length; i++ {
		var offset, totfreq, rfreq uint64
		offset = uint64(cumfreq[0])
		totfreq = uint64(cumfreq[len(cumfreq)-1]) - offset
		rcrange /= uint64(totfreq)
		rfreq = (code - low) / rcrange
		if rfreq >= totfreq {
			panic("C D")
		}

		var value, max, middle int
		value = 0
		max = len(freq) - 1
		for value < max {
			middle = (value + max) / 2
			if uint64(cumfreq[middle+1])-offset <= rfreq {
				value = middle + 1
			} else {
				max = middle
			}
		}

		low += (uint64(cumfreq[value]) - offset) * rcrange
		rcrange *= uint64(freq[value])
		for (low ^ (low + rcrange)) < TOP {
			buf, err := r.ReadByte()
			if err != nil && err != io.EOF {
				panic(err)
			}
			code = (code << 8) | uint64(buf)
			rcrange <<= 8
			low <<= 8
		}
		for rcrange < BOT {
			buf, err := r.ReadByte()
			if err != nil && err != io.EOF {
				panic(err)
			}
			code = (code << 8) | uint64(buf)
			rcrange = ((-low) & (BOT - 1)) << 8
			low <<= 8
		}

		fmt.Print(value)
	}
}
