package processing

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

type ply [][3]int

type bitmap [][]byte

type bitCube struct {
	data   []bitmap
	bias   [3]int
	length [3]int
}

type order [3]int

type orderString int

const (
	XYZ orderString = iota
	XZY
	YXZ
	ZXY
	ZYX
	YZX
)

func (o orderString) Order() order {
	switch o {
	case 0:
		return [3]int{0, 1, 2}
	case 1:
		return [3]int{0, 2, 1}
	case 2:
		return [3]int{1, 0, 2}
	case 3:
		return [3]int{2, 0, 1}
	case 4:
		return [3]int{2, 1, 0}
	case 5:
		return [3]int{1, 2, 0}
	default:
		return [3]int{2, 0, 1}
	}
}

func newPly(srcPath string) (ply, error) {
	fp, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	sccaner := bufio.NewScanner(fp)

	ply := ply{}
	for isData := false; sccaner.Scan(); {
		if isData {
			text := sccaner.Text()
			line := strings.Split(text, " ")
			data := [3]int{}
			for i := 0; i < 3; i++ {
				data[i], err = strconv.Atoi(line[i])
				if err != nil {
					return nil, err
				}
			}
			ply = append(ply, data)
		}

		if "end_header" == sccaner.Text() {
			isData = true
		}
	}
	return ply, nil
}

func newBitCube(ply ply, order order) *bitCube {
	bc := new(bitCube)

	bc.length, bc.bias = ply.getLengthAndbias(order)
	bc.data = make([]bitmap, bc.length[0])
	for i := range bc.data {
		bc.data[i] = make(bitmap, bc.length[1])
		for j := range bc.data[i] {
			bc.data[i][j] = make([]byte, bc.length[2])
		}
	}

	for _, point := range ply {
		dim0 := point[order[0]] - bc.bias[0]
		dim1 := point[order[1]] - bc.bias[1]
		dim2 := point[order[2]] - bc.bias[2]
		bc.data[dim0][dim1][dim2] = 1
	}

	return bc
}

func (ply ply) getLengthAndbias(order order) (length, bias [3]int) {
	for d := 0; d < 3; d++ {
		dim := order[d]
		max := math.MinInt32
		min := math.MaxInt32
		for _, point := range ply {
			if max < point[dim] {
				max = point[dim]
			}
			if min > point[dim] {
				min = point[dim]
			}
		}
		length[d] = max - min + 1
		bias[d] = min
	}
	return length, bias
}

func LoadPly(srcPath string, order order) (*bitCube, error) {
	ply, err := newPly(srcPath)
	if err != nil {
		return nil, err
	}

	bc := newBitCube(ply, order)

	return bc, nil
}