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
	// データ本体
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

// 点の総数
func (points *Points) NumOfPoints() int {
	return len(points.data)
}

// フレーム数を返す関数
func (points *Points) numFrame() int {
	data := points.data
	so := points.sortOrders
	f := 1
	for i := range points.data {
		if i != 0 && data[i][so[0]] != data[i-1][so[0]] {
			f++
		}
	}
	return f
}

// フレームの高さを返す関数
func (points *Points) frameHeight() (int, int) {
	max := math.MinInt32
	min := math.MaxInt32
	data := points.data
	so := points.sortOrders
	for i := range data {
		if max < data[i][so[1]] {
			max = data[i][so[1]]
		}
		if min > data[i][so[1]] {
			min = data[i][so[1]]
		}
	}
	return max - min + 1, min
}

// フレームの幅を返す関数
func (points *Points) frameWidth() (int, int) {
	max := math.MinInt32
	min := math.MaxInt32
	data := points.data
	so := points.sortOrders
	for i := range data {
		if max < data[i][so[2]] {
			max = data[i][so[2]]
		}
		if min > data[i][so[2]] {
			min = data[i][so[2]]
		}
	}
	return max - min + 1, min
}
