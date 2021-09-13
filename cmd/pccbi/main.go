package main

import (
	"github.com/shimehituzi/pccbi/internal/codec"
	"github.com/shimehituzi/pccbi/internal/encoder"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	// 1062090 点のデータ
	srcPath := relativePath[4:]
	voxel, err := encoder.LoadPly(srcPath, encoder.YZX.Order())
	if err != nil {
		panic(err)
	}
	cb := encoder.NewContourBuffer(voxel)

	stream := encoder.NewStream(voxel, cb)
	codec.Encode(stream)

	codec.Decode()
}
