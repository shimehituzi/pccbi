package tool

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readCoordinates(path string) plyData {
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

func writeCoordinates(path string, data plyData) {
	fp, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	for _, v := range data {
		line := fmt.Sprintf("%d %d %d\n", v[0], v[1], v[2])
		w.WriteString(line)
	}

	w.Flush()
}
