package plyio

import (
	"bufio"
	"os"
)

// ply ファイルを文字列の配列として取り込んだ構造体
type Ply struct {
	// ヘッダー部分
	header []string
	// データ部分
	data []string
}

// Ply のコンストラクタ
func NewPly(srcfile string) (*Ply, error) {
	fp, err := os.Open(srcfile)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	ply := new(Ply)

	sccaner := bufio.NewScanner(fp)

	for isData := false; sccaner.Scan(); {
		if isData {
			ply.data = append(ply.data, sccaner.Text())
		} else {
			ply.header = append(ply.header, sccaner.Text())
		}

		if "end_header" == sccaner.Text() {
			isData = true
		}
	}
	return ply, nil
}
