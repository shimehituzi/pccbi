package tool

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func Postprocessing(dstPath, etcPath, recPath string) {
	dstData := readDst(dstPath)
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

func readDst(path string) plyData {
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	sccaner := bufio.NewScanner(fp)

	data := plyData{}
	for sccaner.Scan() {
		texts := strings.Split(sccaner.Text(), " ")
		line := make([]int, len(texts))
		for i := range line {
			line[i], err = strconv.Atoi(texts[i])
			if err != nil {
				panic(err)
			}
		}
		data = append(data, line)
	}

	return data
}
