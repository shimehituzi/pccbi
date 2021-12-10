package codec

import (
	"bufio"
	"fmt"
	"os"
)

func WritePly(recPath string, ply Ply) {
	fp, err := os.Create(recPath)
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

func DecVoxcel(voxel Voxel, header *Header) Ply {
	ply := Ply{}
	index := header.Axis.getIndex()
	bias := header.Bias
	for i := range voxel {
		for j := range voxel[i] {
			for k, v := range voxel[i][j] {
				if v == 0 {
					continue
				}
				coordinates := [3]int{i, j, k}

				x := coordinates[index[0]] + bias[index[0]]
				y := coordinates[index[1]] + bias[index[1]]
				z := coordinates[index[2]] + bias[index[2]]

				ply = append(ply, [3]int{x, y, z})
			}
		}
	}

	ply.Sort()

	return ply
}
