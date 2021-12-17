package tool

func Preprocessing(origPath, srcPlyPath, etcPath string) int {
	data, header := readPly(origPath)
	data.Sort()

	writePly(srcPlyPath, data, header)

	etcData := make(plyData, len(data))
	for i := range etcData {
		etcData[i] = data[i][3:]
	}

	writePly(etcPath, etcData, header)

	return len(data)
}

func Postprocessing(recPath, etcPath, recPlyPath string) {
	dstData := readCoordinatesFile(recPath)
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
