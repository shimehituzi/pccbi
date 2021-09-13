package main

import (
	"github.com/shimehituzi/pccbi/internal/codec"
	"github.com/shimehituzi/pccbi/internal/processing"
)

func main() {
	relativePath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"
	// 1062090 点のデータ
	srcPath := relativePath[4:]
	voxel, err := processing.LoadPly(srcPath, processing.YZX.Order())
	if err != nil {
		panic(err)
	}
	cb := processing.NewContourBuffer(voxel)

	stream := processing.NewStream(voxel, cb)
	codec.Encode(stream)

	codec.Decode()
}
