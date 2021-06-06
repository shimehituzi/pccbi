package plyio

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Ply の座標をソートする順番
// 例えば [1, 2, 0] だと y → z → x の順に安定ソートする
type sortOrders [3]int

// Ply の座標部分を数値としてパースしてソートした型
type Points struct {
	// ソート順
	sortOrders sortOrders
	// 座標データ
	data [][3]int
}

// Points のコンストラクタ x y z のソートする順番を [1, 2, 0] のように与えて初期化
func NewPoints(so sortOrders) (*Points, error) {
	const message = "ソートする順番は 0, 1, 2 の範囲で重複の無いように与えてください"

	for i := range so {
		if so[i] < 0 || 2 < so[i] {
			return nil, errors.New(message)
		}
	}

	if so[0] == so[1] || so[1] == so[2] || so[2] == so[1] {
		return nil, errors.New(message)
	}

	points := new(Points)
	points.sortOrders = so

	return points, nil
}

// ply の構造体からソートされた座標を読み込む
func (points *Points) ReadPly(ply *Ply) error {
	points.data = make([][3]int, ply.NumOfPoints())

	for i := range ply.data {
		line := strings.Split(ply.data[i], " ")
		for j := 0; j < 3; j++ {
			data, err := strconv.Atoi(line[j])
			if err != nil {
				return err
			}
			points.data[i][j] = data
		}
	}

	sort.Slice(points.data, points.generateLessFunc())

	return nil
}

// sort.Slice に渡す型
type lessFunc func(i, j int) bool

// クロージャーを返す高階関数
func (points *Points) generateLessFunc() lessFunc {
	return func(i, j int) bool {
		data := points.data
		so := points.sortOrders
		switch {
		case data[i][so[0]] < data[j][so[0]]:
			return true
		case data[i][so[0]] > data[j][so[0]]:
			return false
		default:
			switch {
			case data[i][so[1]] < data[j][so[1]]:
				return true
			case data[i][so[1]] > data[j][so[1]]:
				return false
			default:
				switch {
				case data[i][so[2]] < data[j][so[2]]:
					return true
				case data[i][so[2]] > data[j][so[2]]:
					return false
				default:
					return false
				}
			}
		}
	}
}

// ある次元の要素の幅とバイアスを返す関数
func (points *Points) getLengthAndBias(sortOrderElem int) (length, bias int) {
	switch sortOrderElem {
	case 0, 1, 2:
		max := math.MinInt32
		min := math.MaxInt32
		data := points.data
		for i := range data {
			if max < data[i][sortOrderElem] {
				max = data[i][sortOrderElem]
			}
			if min > data[i][sortOrderElem] {
				min = data[i][sortOrderElem]
			}
		}
		length = max - min + 1
		bias = min
		return
	default:
		panic("sortOrderElem は 0, 1, 2 のいずれかを与えてください")
	}
}
