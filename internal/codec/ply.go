package codec

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// =================
//      encode
// =================

func ReadPly(srcPath string, axis Axis) (Ply, *Header) {
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

	ply.Sort()

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
	header := &Header{
		Axis:   axis,
		Length: length,
		Bias:   bias,
	}
	return ply, header
}

// =================
//      decode
// =================

func WritePly(dstPath string, ply Ply) {
	fp, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	for _, v := range ply {
		line := fmt.Sprintf("%d %d %d\n", v[0], v[1], v[2])
		w.WriteString(line)
	}

	w.Flush()
}
