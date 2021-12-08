package encoder

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

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

func newVoxel(ply ply, order order) (Voxel, *Header) {
	header := ply.getHeader(order)
	voxel := make([]bitmap, header.Length[0])
	for i := range voxel {
		voxel[i] = make(bitmap, header.Length[1])
		for j := range voxel[i] {
			voxel[i][j] = make([]byte, header.Length[2])
		}
	}

	for _, point := range ply {
		dim0 := point[order[0]] - header.Bias[0]
		dim1 := point[order[1]] - header.Bias[1]
		dim2 := point[order[2]] - header.Bias[2]
		voxel[dim0][dim1][dim2] = 1
	}

	return voxel, header
}

func (ply ply) getHeader(order order) *Header {
	var length, bias, axis [3]int
	for i := range axis {
		axis[order[i]] = i
	}
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
	return &Header{
		Axis:   axis,
		Length: length,
		Bias:   bias,
	}
}

func LoadPly(srcPath string, order order) (Voxel, *Header) {
	ply, err := newPly(srcPath)
	if err != nil {
		panic(err)
	}

	voxel, header := newVoxel(ply, order)

	return voxel, header
}
