package codec

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

func ReadPly(srcPath string, axis Axis) (Ply, *Header) {
	ply := encPly(srcPath)
	header := encHeader(ply, axis)
	return ply, header
}

func EncVoxel(ply Ply, header *Header) Voxel {
	voxel := make(Voxel, header.Length[0])
	for i := range voxel {
		voxel[i] = make(bitmap, header.Length[1])
		for j := range voxel[i] {
			voxel[i][j] = make([]byte, header.Length[2])
		}
	}

	order := header.Axis.getOrder()
	for _, point := range ply {
		dim0 := point[order[0]] - header.Bias[0]
		dim1 := point[order[1]] - header.Bias[1]
		dim2 := point[order[2]] - header.Bias[2]
		voxel[dim0][dim1][dim2] = 1
	}

	return voxel
}

func encPly(srcPath string) Ply {
	fp, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	sccaner := bufio.NewScanner(fp)

	ply := Ply{}
	for isData := false; sccaner.Scan(); {
		if isData {
			text := sccaner.Text()
			line := strings.Split(text, " ")
			data := [3]int{}
			for i := 0; i < 3; i++ {
				data[i], err = strconv.Atoi(line[i])
				if err != nil {
					panic(err)
				}
			}
			ply = append(ply, data)
		}

		if "end_header" == sccaner.Text() {
			isData = true
		}
	}

	// エンコード時とデコード時の比較のために Sort している
	// Voxel を作る上ではしてもしなくても関係ない
	ply.Sort()

	return ply
}

func encHeader(ply Ply, axis Axis) *Header {
	var length, bias [3]int
	order := axis.getOrder()
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
