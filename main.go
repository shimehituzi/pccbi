package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var srcPath string
	flag.StringVar(&srcPath, "s", "", "入力ファイルのパス")
	flag.Parse()

	fp, err := os.Open(srcPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer fp.Close()

	sccaner := bufio.NewScanner(fp)

	var (
		lines  []string
		isData = false
	)
	for sccaner.Scan() {
		if isData {
			lines = append(lines, sccaner.Text())
		}
		if "end_header" == sccaner.Text() {
			isData = true
		}
	}

	numberOfLines := len(lines)
	data := make(Data, numberOfLines)
	for i := range lines {
		line := strings.Split(lines[i], " ")
		for j := 0; j < 3; j++ {
			data[i][j], err = strconv.Atoi(line[j])
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	sort.Sort(&data)

	for i := range data {
		if i < 30 {
			fmt.Println(data[i])
		}
	}
}

type Data [][3]int

func (data Data) Len() int { return len(data) }

func (data Data) Swap(i, j int) { data[i], data[j] = data[j], data[i] }

// sort.SliceStable(data, func(i, j int) bool { return data[i][0] < data[j][0] })
// sort.SliceStable(data, func(i, j int) bool { return data[i][2] < data[j][2] })
// sort.SliceStable(data, func(i, j int) bool { return data[i][1] < data[j][1] })
func (data Data) Less(i, j int) bool {
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
			switch {
			case data[i][0] < data[j][0]:
				return true
			case data[i][0] > data[j][0]:
				return false
			default:
				return false
			}
		}
	}
}
