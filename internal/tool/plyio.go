package tool

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type plyHeader []string
type plyData [][]int

func readPly(path string) (plyData, plyHeader) {
	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	sccaner := bufio.NewScanner(fp)

	header := plyHeader{}
	data := plyData{}
	for isData := false; sccaner.Scan(); {
		text := sccaner.Text()
		if isData {
			texts := strings.Split(text, " ")
			line := make([]int, len(texts))
			for i := range line {
				line[i], err = strconv.Atoi(texts[i])
				if err != nil {
					panic(err)
				}
			}
			data = append(data, line)
		} else {
			header = append(header, text)
		}

		if text == "end_header" {
			isData = true
		}
	}

	return data, header
}

func writePly(path string, data plyData, header plyHeader) {
	fp, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	w := bufio.NewWriter(fp)

	for _, v := range header {
		w.WriteString(v + "\n")
	}

	for i := range data {
		line := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data[i])), " "), "[]")
		w.WriteString(line + "\n")
	}

	w.Flush()
}
