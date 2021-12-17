package tool

func Preprocessing(srcPath, sortedPath, etcPath string) int {
	data, header := readPly(srcPath)
	data.Sort()

	writePly(sortedPath, data, header)

	etcData := make(plyData, len(data))
	for i := range etcData {
		etcData[i] = data[i][3:]
	}

	writePly(etcPath, etcData, header)

	return len(data)
}

func Postprocessing(dstPath, etcPath, recPath string) {
	dstData := readCoordinatesFile(dstPath)
	etcData, header := readPly(etcPath)

	if len(dstData) != len(etcData) {
		panic("The plyData length is different")
	}

	data := make(plyData, len(dstData))
	for i := range data {
		data[i] = append(dstData[i], etcData[i]...)
	}

	writePly(recPath, data, header)
}
