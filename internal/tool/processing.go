package tool

func Preprocessing(origPath, srcPlyPath, srcPath, etcPath string) int {
	data, header := readPly(origPath)
	data.Sort()

	writePly(srcPlyPath, data, header)

	coodinatesData := make(plyData, len(data))
	attributesData := make(plyData, len(data))
	for i := range data {
		coodinatesData[i] = data[i][:3]
		attributesData[i] = data[i][3:]
	}
	writeCoordinates(srcPath, coodinatesData)
	writePly(etcPath, attributesData, header)

	return len(data)
}

func Postprocessing(recPath, etcPath, recPlyPath string) {
	dstData := readCoordinates(recPath)
	etcData, header := readPly(etcPath)

	if len(dstData) != len(etcData) {
		panic("The plyData length is different")
	}

	data := make(plyData, len(dstData))
	for i := range data {
		data[i] = append(dstData[i], etcData[i]...)
	}

	writePly(recPlyPath, data, header)
}
