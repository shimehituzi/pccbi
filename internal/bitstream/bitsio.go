package bitstream

import (
	"bufio"
	"fmt"
	"os"
)

const bitsize = 16

type bitbuf struct {
	pos int
	buf uint
}

func NewBitbuf(enc bool) *bitbuf {
	bit := new(bitbuf)
	bit.buf = 0
	if enc {
		bit.pos = 8
	} else {
		bit.pos = 0
	}
	return bit
}

func (bit *bitbuf) putbits(w *bufio.Writer, n int, x uint) int {
	bits := n
	if bits <= 0 {
		return 0
	}
	for n >= bit.pos {
		n -= bit.pos
		if n < 32 {
			bit.buf |= ((x >> n) & (0xff >> (8 - bit.pos)))
		}
		w.WriteByte(byte(bit.buf))
		bit.buf = 0
		bit.pos = 8
	}
	bit.pos -= n
	bit.buf |= ((x & (0xff >> (8 - n))) << bit.pos)
	return bits
}

func (bit *bitbuf) getbits(r *bufio.Reader, n int) uint {
	x := uint(0)
	if n <= 0 {
		return 0
	}
	for n > bit.pos {
		n -= bit.pos
		x = (x << bit.pos) | bit.buf
		byte, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
		bit.buf = uint(byte & 0xff)
		bit.pos = 8
	}
	bit.pos -= n
	x = (x << n) | (bit.buf >> bit.pos)
	bit.buf &= ((1 << bit.pos) - 1)
	return x
}

// ------------Test------------

func testBitsio() {
	filename := "bytes"

	val := []uint{5555, 8941, 12, 3215}
	bits := testBitsioEncode(filename, val)
	fmt.Println(bits)

	length := 4
	valrc := testBitsioDecode(filename, length)
	fmt.Println(valrc)
}

func testBitsioEncode(filename string, val []uint) int {
	fp, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	bitbuf := NewBitbuf(true)

	bits := 0
	for _, elem := range val {
		bits += bitbuf.putbits(w, bitsize, elem)
	}
	w.Flush()

	return bits
}

func testBitsioDecode(filename string, length int) (val []uint) {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	r := bufio.NewReader(fp)

	bitbuf := NewBitbuf(false)

	val = make([]uint, length)
	for i := range val {
		val[i] = bitbuf.getbits(r, bitsize)
	}

	return val
}
