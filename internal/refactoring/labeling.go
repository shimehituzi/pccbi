package refactoring

import (
	"fmt"
	"math"
	"sync"
)

type label intmap

type labeledVoxel []label

func NewLabels(voxel *voxel) labeledVoxel {
	lv := make(labeledVoxel, voxel.header.length[0])
	wg := &sync.WaitGroup{}
	for i := range lv {
		wg.Add(1)
		go func(i int) {
			lv[i] = newLabel(voxel.Data[i], i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	return lv
}

func newLabel(bm bitmap, debug int) label {
	label := make(label, len(bm))
	for y := range label {
		label[y] = make([]int, len(bm[0]))
	}

	lookupTable := []int{}
	counter := 0

	for y := range bm {
		for x := range bm[y] {
			if bm[y][x] == 1 {
				al := arroundLavel(x, y, label)
				if len(al) == 0 {
					// 新規ラベル作成
					counter++
					lookupTable = append(lookupTable, counter)
					label[y][x] = counter
				} else {
					// 最小値を算出
					min := math.MaxInt32
					for _, v := range al {
						if min > lookupTable[v-1] {
							min = lookupTable[v-1]
						}
					}
					// ルックアップテーブル更新
					for _, v := range al {
						lookupTable[v-1] = min
					}
					// ラベリング
					label[y][x] = min
				}
			}
		}
	}

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
	for i, u := range updateTable {
		for j, v := range lookupTable {
			if v == u {
				lookupTable[j] = i + 1
			}
		}
	}

	if debug == 485 {
		for i := range lookupTable {
			fmt.Println("[", i+1, ",", lookupTable[i], "]")
		}
	}

	for y := range label {
		for x, v := range label[y] {
			if v != 0 {
				newLabel := lookupTable[v-1]
				label[y][x] = newLabel
			}
		}
	}

	return label
}

func arroundLavel(x, y int, label label) []int {
	mask := []point{
		{-1, -1}, {0, -1}, {1, -1}, {-1, 0},
	}

	al := []int{}
	for _, m := range mask {
		p := point{x + m.x, y + m.y}
		if validPoint(p, intmap(label)) && label[p.y][p.x] != 0 {
			al = append(al, label[p.y][p.x])
		}
	}

	return al
}
