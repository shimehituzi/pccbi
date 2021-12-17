package main

import (
	"time"

	"github.com/shimehituzi/pccbi/internal/codec"
	"github.com/shimehituzi/pccbi/internal/tool"
)

func main() {
	times := [6]time.Time{}
	times[0] = time.Now()

	// Arguments
	axis := codec.YZX

	// Preprocessing
	numPoints := tool.Preprocessing(origPath, srcPlyPath, srcPath, etcPath)

	// Encode
	times[1] = time.Now()
	encPly, eh := codec.ReadPly(srcPath, axis)
	encStream := encPly.ConvertVoxel(eh).ConvertContour(eh).ConvertStream()
	dataBits, headerBits := codec.Encode(dstPath, encStream, eh)
	times[2] = time.Now()

	// Decode
	times[3] = time.Now()
	decStream, dh := codec.Decode(dstPath)
	decPly := decStream.ConvertContour(dh).ConvertVoxel(dh).ConvertPly(dh)
	codec.WritePly(recPath, decPly)
	times[4] = time.Now()

	// Postprocessing
	tool.Postprocessing(recPath, etcPath, recPlyPath)

	// Chack Lossless
	result := tool.TestLossless(srcPath, recPath, srcPlyPath, recPlyPath)
	tool.DeleteTmpFile(result, srcPath, recPath, srcPlyPath, recPlyPath, etcPath)

	// Report
	times[5] = time.Now()
	tool.Report(result, dataBits, headerBits, numPoints, times)
}

// 座標でソートされていないオリジナルの ply
const origPath = "./DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"

// 座標でソートした ply
const srcPlyPath = "./out/src.ply"

// 復元した ply
const recPlyPath = "./out/rec.ply"

// 座標以外の元データ
const etcPath = "./out/attributes_and_header"

// 座標のみの元データ
const srcPath = "./out/coordinates_source"

// 圧縮データ
const dstPath = "./out/coordinates_compressed"

// 座標のみの復元データ
const recPath = "./out/coordinates_reconstructed"
