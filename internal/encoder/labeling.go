package encoder

import (
	"math"
	"sync"
)

func NewLabels(voxel *voxel) (labeledVoxel, []int) {
	lv := make(labeledVoxel, voxel.header.length[0])
	numLabels := make([]int, voxel.header.length[0])
	wg := &sync.WaitGroup{}
	for i := range lv {
		wg.Add(1)
		go func(i int) {
			lv[i], numLabels[i] = newLabel(voxel.Data[i], i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	return lv, numLabels
}

func newLabel(bm bitmap, debug int) (label, int) {
	label := make(label, len(bm))
	for y := range label {
		label[y] = make([]int, len(bm[0]))
	}

	lookupTable := []int{}
	counter := 0

	for y := range bm {
		for x := range bm[y] {
			if bm[y][x] == 1 {
				// 周りの 0 以外の値を取得
				al := arroundLavel(x, y, label)
				if len(al) == 0 {
					// 新規ラベル作成
					counter++
					lookupTable = append(lookupTable, counter)
					label[y][x] = counter
				} else {
					// ルックアップテーブルをたどり一番小さい値を取得する
					root := []int{}
					for _, v := range al {
						k := lookupTable[v-1] - 1
						for k != lookupTable[k]-1 {
							k = lookupTable[k] - 1
						}
						root = append(root, k+1)
					}
					// 最小値を算出
					min := math.MaxInt32
					for _, v := range root {
						if min > lookupTable[v-1] {
							min = lookupTable[v-1]
						}
					}
					// ルックアップテーブル更新
					for _, v := range root {
						lookupTable[v-1] = min
					}
					// ラベリング
					label[y][x] = min
				}
			}
		}
	}

	// ルックアップテーブルを更新
	for i := range lookupTable {
		if i == lookupTable[i]-1 {
			continue
		}
		k := lookupTable[i] - 1
		for k != lookupTable[k]-1 {
			k = lookupTable[k] - 1
		}
		lookupTable[i] = k + 1
	}

	// ルックアップテーブルの値を詰める
	updateTable := []int{}
	for _, v := range lookupTable {
		flag := true
		for _, u := range updateTable {
			if u == v {
				flag = false
			}
		}
		if flag {
			updateTable = append(updateTable, v)
		}
	}
	numLabel := len(updateTable)

	for i, u := range updateTable {
		for j, v := range lookupTable {
			if v == u {
				lookupTable[j] = i + 1
			}
		}
	}

	// label を更新
	for y := range label {
		for x, v := range label[y] {
			if v != 0 {
				newLabel := lookupTable[v-1]
				label[y][x] = newLabel
			}
		}
	}

	return label, numLabel
}

func arroundLavel(x, y int, label label) []int {
	mask := [4]point{
		{-1, -1}, {0, -1}, {1, -1}, {-1, 0},
	}

	al := []int{}
	for _, m := range mask {
		p := point{x + m.x, y + m.y}
		if validPointInt(p, label) && label[p.y][p.x] != 0 {
			al = append(al, label[p.y][p.x])
		}
	}

	return al
}

func validPointInt(p point, img label) bool {
	if p.y < 0 || p.x < 0 || len(img) <= p.y || len(img[0]) <= p.x {
		return false
	}
	return true
}
