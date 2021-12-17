package codec

// =================
//      encode
// =================

func (ply Ply) ConvertVoxel(header *Header) Voxel {
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

// =================
//      decode
// =================

func (voxel Voxel) ConvertPly(header *Header) Ply {
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
