package plyio

import (
	"bufio"
	"os"
)

// ply ファイルを文字列の配列として取り込んだ構造体
type Ply struct {
	// ヘッダー部分
	Header []string
	// データ部分
	Data []string
}

// Ply のコンストラクタ
func NewPly() *Ply {
	ply := new(Ply)
	return ply
}

// ply ファイルから構造体に読み込む
func (ply *Ply) Read(filename string) error {
	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	sccaner := bufio.NewScanner(fp)

	for isData := false; sccaner.Scan(); {
		if isData {
			ply.Data = append(ply.Data, sccaner.Text())
		} else {
			ply.Header = append(ply.Header, sccaner.Text())
		}

		if "end_header" == sccaner.Text() {
			isData = true
		}
	}

	return nil
}

// 点の総数
func (ply *Ply) Points() int {
	return len(ply.Data)
}
