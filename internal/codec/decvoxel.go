package codec

func WritePly(recPath string, ply Ply) {

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
				// YZX
				coordinates := [3]int{i, j, k}

				// YZX → XYZ としたい
				// X の値は 2 に保存されている → coordinates[2] でアクセスできる
				// Y の値は 0 に保存されている → coordinates[0] でアクセスできる
				// Z の値は 1 に保存されている → coordinates[1] でアクセスできる

				// バイアスを考慮したい
				// X のバイアスは 2 に保存されている → bias[2] でアクセスできる
				// Y のバイアスは 0 に保存されている → bias[0] でアクセスできる
				// Z のバイアスは 1 に保存されている → bias[1] でアクセスできる

				// YZX.getIndex() → [3]int{2, 0, 1}

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
