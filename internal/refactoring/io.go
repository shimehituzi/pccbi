package refactoring

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
	Data   []bitmap
	Bias   [3]int
	Length [3]int
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

	bc.Length, bc.Bias = ply.getLengthAndbias()
	bc.Data = make([]bitmap, bc.Length[order[0]])
	for i := range bc.Data {
		bc.Data[i] = make(bitmap, bc.Length[order[1]])
		for j := range bc.Data[i] {
			bc.Data[i][j] = make([]byte, bc.Length[order[2]])
		}
	}

	for _, point := range ply {
		dim0 := point[order[0]] - bc.Bias[order[0]]
		dim1 := point[order[1]] - bc.Bias[order[1]]
		dim2 := point[order[2]] - bc.Bias[order[2]]
		bc.Data[dim0][dim1][dim2] = 1
	}

	return bc
}

func (ply ply) getLengthAndbias() (length, bias [3]int) {
	for dim := 0; dim < 3; dim++ {
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
		length[dim] = max - min + 1
		bias[dim] = min
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
