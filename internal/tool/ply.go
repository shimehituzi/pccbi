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

func readCoordinatesFile(path string) plyData {
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

func (data plyData) Sort() { sort.Sort(data) }

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
