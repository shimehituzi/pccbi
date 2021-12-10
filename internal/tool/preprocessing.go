package tool

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type plyHeader []string
type plyData [][]int

func Preprocessing(srcPath, sortedPath, etcPath string) {
	data, header := readSortedPly(srcPath)
	writePly(sortedPath, data, header, true)
	writePly(etcPath, data, header, false)
}

func writePly(path string, data plyData, header plyHeader, all bool) {
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
		var line []int
		if all {
			line = data[i]
		} else {
			line = data[i][3:]
		}
		text := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(line)), " "), "[]")
		w.WriteString(text + "\n")
	}

	w.Flush()
}

func readSortedPly(srcPath string) (plyData, plyHeader) {
	fp, err := os.Open(srcPath)
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

	sort.Sort(data)

	return data, header
}

func (data plyData) Len() int { return len(data) }

func (data plyData) Swap(i, j int) { data[i], data[j] = data[j], data[i] }

func (data plyData) Less(i, j int) bool {
	switch {
	case data[i][0] < data[j][0]:
		return true
	case data[i][0] > data[j][0]:
		return false
	default:
		switch {
		case data[i][1] < data[j][1]:
			return true
		case data[i][1] > data[j][1]:
			return false
		default:
			switch {
			case data[i][2] < data[j][2]:
				return true
			case data[i][2] > data[j][2]:
				return false
			default:
				return false
			}
		}
	}
}
