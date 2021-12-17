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

	srcPath := "../../DATABASE/orig/soldier/soldier/Ply/soldier_vox10_0537.ply"[4:]
	pccPath := "./out/compressed"
	dstPath := "./out/destination"

	sortedPath := "./out/sorted.ply"
	etcPath := "./out/etc"
	recPath := "./out/rec.ply"

	// Preprocessing
	numPoints := tool.Preprocessing(srcPath, sortedPath, etcPath)

	// Encode
	times[1] = time.Now()
	encPly, eh := codec.ReadPly(srcPath, axis)
	encStream := encPly.ConvertVoxel(eh).ConvertContour(eh).ConvertStream()
	dataBits, headerBits := codec.Encode(pccPath, encStream, eh)
	times[2] = time.Now()

	// Decode
	times[3] = time.Now()
	decStream, dh := codec.Decode(pccPath)
	decPly := decStream.ConvertContour(dh).ConvertVoxel(dh).ConvertPly(dh)
	codec.WritePly(dstPath, decPly)
	times[4] = time.Now()

	// Postprocessing
	tool.Postprocessing(dstPath, etcPath, recPath)

	// Chack Lossless
	result := tool.TestLossless(sortedPath, recPath)
	tool.DeleteTmpFile(result, sortedPath, dstPath, etcPath)

	// Report
	times[5] = time.Now()
	tool.Report(result, dataBits, headerBits, numPoints, times)
}
