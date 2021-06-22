package plyio

import (
	"image"

	"github.com/shimehituzi/pccbi/internal/bitmap"
)

// 2値画像の集合として点群の位置情報を表す構造体
type BitMaps struct {
	// ある次元の要素の値の幅
	Length [3]int
	// ある次元の要素の値のバイアス
	Bias [3]int
	// 2値画像の集合
	Data []bitmap.BitMap
}

// BitMaps のコンストラクタ
func NewBitMaps(points *Points) *BitMaps {
	bms := new(BitMaps)
	so := points.sortOrders
	for i := 0; i < 3; i++ {
		bms.Length[i], bms.Bias[i] = points.getLengthAndBias(so[i])
	}

	bms.Data = make([]bitmap.BitMap, bms.Length[0])
	for i := range bms.Data {
		bms.Data[i] = make([][]byte, bms.Length[1])
		for j := range bms.Data[i] {
			bms.Data[i][j] = make([]byte, bms.Length[2])
		}
	}

	pdata := points.data
	for i := range points.data {
		dim0 := pdata[i][so[0]] - bms.Bias[0]
		dim1 := pdata[i][so[1]] - bms.Bias[1]
		dim2 := pdata[i][so[2]] - bms.Bias[2]
		bms.Data[dim0][dim1][dim2] = 1
	}
	return bms
}

func (bms *BitMaps) GetImage(f int) image.Image {
	return bms.Data[f]
}

func (bms *BitMaps) GetLength() bitmap.DimensionLength {
	return bitmap.DimensionLength{
		D0: bms.Length[0],
		D1: bms.Length[1],
		D2: bms.Length[2],
	}
}
